package money_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/faustbrian/go-international/currency"
	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/decimal"
	"github.com/faustbrian/go-math/integer"
	"github.com/faustbrian/go-money/v2"
)

func TestExternalNumericValuesAreBoundedBeforeFormatting(t *testing.T) {
	euro, err := currency.Parse("EUR")
	if err != nil {
		t.Fatal(err)
	}
	monetaryContext, err := money.DefaultContext(euro)
	if err != nil {
		t.Fatal(err)
	}

	// Construction belongs to the caller, outside the measured Money boundary.
	oversized := strings.Repeat("9", 64*1024)
	units, err := integer.Parse(oversized, integer.ParseOptions{Base: 10, Limits: gomath.DefaultLimits()})
	if err != nil {
		t.Fatal(err)
	}
	amount, err := decimal.Parse(oversized)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := money.AmountFromDecimal(amount); !errors.Is(err, money.ErrAmountLimit) {
		t.Fatalf("AmountFromDecimal(oversized) error = %v", err)
	}
	for _, test := range []struct {
		name   string
		reject func() error
	}{
		{
			name: "minor units",
			reject: func() error {
				_, err := money.FromMinorUnits(units, euro, monetaryContext)
				return err
			},
		},
		{
			name: "decimal amount",
			reject: func() error {
				_, err := money.AmountFromDecimal(amount)
				return err
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := test.reject(); err == nil {
				t.Fatal("accepted oversized numeric value")
			}
			result := testing.Benchmark(func(b *testing.B) {
				for range b.N {
					if err := test.reject(); err == nil {
						b.Fatal("accepted oversized numeric value")
					}
				}
			})
			// A generous fixed ceiling admits monetary-limit bookkeeping but not
			// formatting an externally owned 64-KiB coefficient before rejection.
			if result.AllocedBytesPerOp() > 16*1024 {
				t.Errorf("oversized numeric rejection allocated %d bytes per call", result.AllocedBytesPerOp())
			}
		})
	}
}

func TestExternalNumericValuesPreserveInclusiveMonetaryLimits(t *testing.T) {
	euro, err := currency.Parse("EUR")
	if err != nil {
		t.Fatal(err)
	}
	for _, scale := range []uint8{0, money.MaxScale} {
		monetaryContext, err := money.CustomContext(scale)
		if err != nil {
			t.Fatal(err)
		}
		for _, sign := range []string{"", "-"} {
			for _, digits := range []int{money.MaxAmountDigits, money.MaxAmountDigits + 1} {
				text := sign + strings.Repeat("9", digits)
				units, err := integer.Parse(text, integer.ParseOptions{Base: 10, Limits: gomath.DefaultLimits()})
				if err != nil {
					t.Fatal(err)
				}
				value, err := money.FromMinorUnits(units, euro, monetaryContext)
				if digits > money.MaxAmountDigits {
					if err == nil {
						t.Fatalf("FromMinorUnits accepted %d digits at scale %d", digits, scale)
					}
					continue
				}
				if err != nil {
					t.Fatalf("FromMinorUnits(exact digit bound, scale %d) error = %v", scale, err)
				}
				roundTrip, err := value.MinorUnits()
				if err != nil || !roundTrip.Equal(units) {
					t.Fatalf("MinorUnits(exact digit bound, scale %d) did not preserve signed coefficient: %v", scale, err)
				}
				if value.Context().Scale() != scale {
					t.Fatalf("FromMinorUnits changed context scale %d", scale)
				}
			}
		}
	}
}

func TestMinorUnitAdmissionPreservesMetadataErrorPrecedence(t *testing.T) {
	euro, err := currency.Parse("EUR")
	if err != nil {
		t.Fatal(err)
	}
	dollar, err := currency.Parse("USD")
	if err != nil {
		t.Fatal(err)
	}
	fixed, err := money.DefaultContext(euro)
	if err != nil {
		t.Fatal(err)
	}
	units, err := integer.Parse(strings.Repeat("9", money.MaxAmountDigits+1), integer.ParseOptions{
		Base: 10, Limits: gomath.DefaultLimits(),
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name    string
		code    currency.Code
		context money.Context
		want    error
	}{
		{name: "unknown currency", context: fixed, want: money.ErrUnknownCurrency},
		{name: "mismatched context", code: dollar, context: fixed, want: money.ErrContextMismatch},
		{name: "missing context first", want: money.ErrInvalidContext},
		{name: "automatic context", code: euro, context: money.AutomaticContext(), want: money.ErrInvalidContext},
		{name: "valid metadata", code: euro, context: fixed, want: gomath.ErrLimitExceeded},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := money.FromMinorUnits(units, test.code, test.context); !errors.Is(err, test.want) {
				t.Fatalf("FromMinorUnits metadata admission error=%v, want %v", err, test.want)
			}
		})
	}
}

func TestAmountAdmissionPreservesDecimalRepresentation(t *testing.T) {
	for _, text := range []string{"1.2300", "-1.2300", "0.000000000000000000"} {
		value, err := decimal.Parse(text)
		if err != nil {
			t.Fatal(err)
		}
		amount, err := money.AmountFromDecimal(value)
		if err != nil || amount.String() != value.String() || amount.Scale() != value.Scale() {
			t.Fatalf("AmountFromDecimal changed represented amount %q: %v", text, err)
		}
	}
}
