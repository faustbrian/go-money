// Package encoding provides bounded, versioned money wire and persistence
// adapters. Amounts are always exact decimal strings, never JSON numbers.
// Start with MarshalJSON/UnmarshalJSON for whole values or SQLMoney for
// database/sql. NumericValue/ScanNumeric require separate currency/context data.
//
// The adapters perform no database I/O, migrations, authentication or transaction
// management, and require neither a database driver nor cgo. Stateless functions
// can be called concurrently; a mutable SQLMoney receiver needs caller-owned
// synchronization. Version 1 is a package-owned schema, not a payment protocol.
package encoding

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/faustbrian/go-international/currency"
	"github.com/faustbrian/go-money/v2"
)

const (
	// Version is the current persistence representation.
	Version uint8 = 1
	// MaxEncodedBytes bounds JSON, text, and SQL input before decoding.
	MaxEncodedBytes = 2_048
)

// ErrInvalidEncoding classifies malformed, unsupported or excessive wire/SQL
// representations. Decode errors match it with errors.Is without exposing input.
var ErrInvalidEncoding = errors.New("money encoding: invalid representation")

type wireContext struct {
	Kind     string `json:"kind"`
	Scale    uint8  `json:"scale"`
	CashStep uint64 `json:"cash_step,omitempty"`
}

type wireMoney struct {
	Version  uint8       `json:"version"`
	Amount   string      `json:"amount"`
	Currency string      `json:"currency"`
	Context  wireContext `json:"context"`
}

// MarshalJSON returns deterministic versioned JSON with a string amount,
// retaining currency, context kind, resolved scale and cash step. Invalid Money
// returns money.ErrInvalidMoney and nil output. Returned bytes are caller-owned.
func MarshalJSON(value money.Money) ([]byte, error) {
	if !value.Valid() {
		return nil, money.ErrInvalidMoney
	}
	kind := mustContextKind(value.Context().Kind())

	return json.Marshal(wireMoney{
		Version:  Version,
		Amount:   value.Amount().String(),
		Currency: value.Currency().String(),
		Context: wireContext{
			Kind:     kind,
			Scale:    value.Context().Scale(),
			CashStep: value.Context().CashStep(),
		},
	})
}

// UnmarshalJSON validates version-1 JSON and reconstructs Money without rounding.
// Empty/nil or excessive input, duplicate keys at any depth, unknown fields,
// trailing data and inconsistent metadata return ErrInvalidEncoding or an error
// wrapping it. Historic currency codes are accepted. Input is not retained;
// decoding failure returns invalid Money. This validates representation, not
// authenticity or business authorization.
func UnmarshalJSON(data []byte) (money.Money, error) {
	if len(data) == 0 || len(data) > MaxEncodedBytes {
		return money.Money{}, ErrInvalidEncoding
	}
	// encoding/json normally accepts later duplicate keys. Preflight the whole
	// bounded document so nested duplicates cannot override persisted identity.
	if err := rejectDuplicateKeys(data); err != nil {
		return money.Money{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var wire wireMoney
	if err := decoder.Decode(&wire); err != nil {
		return money.Money{}, fmt.Errorf("%w: malformed JSON", ErrInvalidEncoding)
	}
	if wire.Version != Version || wire.Amount == "" || wire.Currency == "" {
		return money.Money{}, ErrInvalidEncoding
	}
	code, err := currency.ParseWithOptions(wire.Currency, currency.ParseOptions{AllowHistoric: true})
	if err != nil {
		return money.Money{}, fmt.Errorf("%w: currency", ErrInvalidEncoding)
	}
	context, err := decodeContext(wire.Context, code)
	if err != nil {
		return money.Money{}, err
	}
	value, err := money.Parse(wire.Amount, code, context)
	if err != nil {
		return money.Money{}, fmt.Errorf("%w: amount", ErrInvalidEncoding)
	}
	if value.Context().Scale() != wire.Context.Scale {
		return money.Money{}, ErrInvalidEncoding
	}

	return value, nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := scanJSONValue(decoder); err != nil {
		return fmt.Errorf("%w: malformed or duplicate JSON", ErrInvalidEncoding)
	}
	if err := requireEOF(decoder); err != nil {
		return err
	}

	return nil
}

// The production decoder calls scanJSONValue after checking the byte ceiling.
// That ceiling
// bounds recursive nesting and token work; this is not a streaming/unbounded
// JSON traversal API.
func scanJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return nil
	}
	switch delimiter {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, keyErr := decoder.Token()
			if keyErr != nil {
				return keyErr
			}
			key := keyToken.(string)
			if _, duplicate := seen[key]; duplicate {
				return ErrInvalidEncoding
			}
			seen[key] = struct{}{}
			valueErr := scanJSONValue(decoder)
			switch valueErr {
			case nil:
			default:
				return valueErr
			}
		}
	case '[':
		for decoder.More() {
			valueErr := scanJSONValue(decoder)
			switch valueErr {
			case nil:
			default:
				return valueErr
			}
		}
	}
	closing, err := decoder.Token()

	return errors.Join(err, validateClosingDelimiter(delimiter, closing))
}

func validateClosingDelimiter(open json.Delim, closing json.Token) error {
	if closing != matchingDelimiter(open) {
		return ErrInvalidEncoding
	}

	return nil
}

func matchingDelimiter(open json.Delim) json.Delim {
	if open == '{' {
		return '}'
	}

	return ']'
}

// MarshalText uses the same canonical versioned representation as JSON.
func MarshalText(value money.Money) ([]byte, error) { return MarshalJSON(value) }

// UnmarshalText validates the canonical versioned representation.
func UnmarshalText(data []byte) (money.Money, error) { return UnmarshalJSON(data) }

// SQLMoney adapts versioned Money to database/sql Scanner and Valuer. The zero
// value can receive Scan but cannot produce Value until populated. A receiver
// is mutable and must not be scanned or read concurrently without synchronization.
type SQLMoney struct {
	// Money is the current value, replaced only after Scan successfully decodes.
	Money money.Money
}

// Value returns canonical versioned text suitable for text or JSON columns.
// Invalid Money returns nil and money.ErrInvalidMoney; SQL NULL is not emitted.
func (value SQLMoney) Value() (driver.Value, error) {
	data, err := MarshalText(value.Money)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

// Scan accepts only string or byte database values and replaces Money after a
// complete successful decode. Nil receivers, SQL NULL, unsupported source types
// and out-of-budget inputs return ErrInvalidEncoding. Byte inputs are copied
// after admission; no input buffer is retained. Every failure preserves Money.
func (value *SQLMoney) Scan(source any) error {
	if value == nil {
		return ErrInvalidEncoding
	}
	var data []byte
	switch source := source.(type) {
	case string:
		if len(source) == 0 || len(source) > MaxEncodedBytes {
			return ErrInvalidEncoding
		}
		data = []byte(source)
	case []byte:
		if len(source) == 0 || len(source) > MaxEncodedBytes {
			return ErrInvalidEncoding
		}
		data = append([]byte(nil), source...)
	default:
		return ErrInvalidEncoding
	}
	decoded, err := UnmarshalText(data)
	if err != nil {
		return err
	}
	value.Money = decoded

	return nil
}

// NumericValue returns only exact decimal text for PostgreSQL numeric columns.
// Currency and context must be persisted separately by the caller.
func NumericValue(value money.Money) (driver.Value, error) {
	if !value.Valid() {
		return nil, money.ErrInvalidMoney
	}

	return value.Amount().String(), nil
}

// ScanNumeric reconstructs Money from PostgreSQL numeric text plus explicit
// currency and context columns. Only string and []byte are accepted; SQL NULL,
// binary floats and inputs over money.MaxAmountDigits+2 bytes return
// ErrInvalidEncoding before conversion. Admitted text passes directly to
// money.Parse, preserving its errors and precision policy. No input is retained.
func ScanNumeric(source any, code currency.Code, context money.Context) (money.Money, error) {
	var text string
	switch source := source.(type) {
	case string:
		if len(source) > money.MaxAmountDigits+2 {
			return money.Money{}, ErrInvalidEncoding
		}
		text = source
	case []byte:
		if len(source) > money.MaxAmountDigits+2 {
			return money.Money{}, ErrInvalidEncoding
		}
		text = string(source)
	default:
		return money.Money{}, ErrInvalidEncoding
	}
	return money.Parse(text, code, context)
}

func decodeContext(wire wireContext, code currency.Code) (money.Context, error) {
	switch wire.Kind {
	case "default":
		if wire.CashStep != 0 {
			return money.Context{}, ErrInvalidEncoding
		}
		context, err := money.DefaultContext(code)
		if err != nil || context.Scale() != wire.Scale {
			return money.Context{}, ErrInvalidEncoding
		}
		return context, nil
	case "custom":
		if wire.CashStep != 0 {
			return money.Context{}, ErrInvalidEncoding
		}
		context, err := money.CustomContext(wire.Scale)
		return context, wrapContextError(err)
	case "cash":
		context, err := money.CashContext(wire.Scale, wire.CashStep)
		return context, wrapContextError(err)
	case "automatic":
		if wire.CashStep != 0 {
			return money.Context{}, ErrInvalidEncoding
		}
		return money.AutomaticContext(), nil
	default:
		return money.Context{}, ErrInvalidEncoding
	}
}

func contextKind(kind money.ContextKind) (string, error) {
	switch kind {
	case money.ContextDefault:
		return "default", nil
	case money.ContextCustom:
		return "custom", nil
	case money.ContextCash:
		return "cash", nil
	case money.ContextAutomatic:
		return "automatic", nil
	default:
		return "", money.ErrInvalidContext
	}
}

func mustContextKind(kind money.ContextKind) string {
	encoded, err := contextKind(kind)
	if err != nil {
		panic("money encoding: internal context invariant violated")
	}

	return encoded
}

func wrapContextError(err error) error {
	if err != nil {
		return ErrInvalidEncoding
	}

	return nil
}

func requireEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return ErrInvalidEncoding
	}

	return nil
}
