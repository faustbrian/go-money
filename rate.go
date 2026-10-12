package money

import (
	"fmt"
	"strings"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

// MaxRateMagnitude is the inclusive default ceiling used by ParseRate.
// ParseRateWithMaximum can select a different monetary-policy ceiling.
const MaxRateMagnitude int64 = 1_000_000

// Rate is an immutable, exact, nonnegative rational multiplier. Its zero value
// is invalid; a parsed zero is valid. Parsing bounds do not guarantee that a
// later multiplication fits the monetary output budget.
type Rate struct {
	value rational.Rational
	valid bool
}

// ParseRate parses either exact decimal text or a numerator/denominator pair,
// accepting values from zero through MaxRateMagnitude. Malformed, negative or
// excessive input returns an error matching ErrInvalidRate.
func ParseRate(input string) (Rate, error) {
	return ParseRateWithMaximum(input, fmt.Sprint(MaxRateMagnitude))
}

// ParseRateWithMaximum parses a rate while applying an explicit inclusive
// caller-defined maximum. The maximum accepts the same exact decimal or
// numerator/denominator syntax as input and must be positive. The generic
// MaxRateMagnitude does not cap this override, but digit, scale and intermediate
// budgets still apply. Invalid input or bounds return ErrInvalidRate; underlying
// parse errors are diagnostic text, not additionally wrapped error identities.
func ParseRateWithMaximum(input, maximum string) (Rate, error) {
	value, err := parseRateValue(input)
	if err != nil {
		return Rate{}, fmt.Errorf("%w: %v", ErrInvalidRate, err)
	}
	bound, err := parseRateValue(maximum)
	if err != nil || bound.Sign() <= 0 {
		return Rate{}, ErrInvalidRate
	}
	if value.Sign() < 0 || value.Cmp(bound) > 0 {
		return Rate{}, ErrInvalidRate
	}

	return Rate{value: value, valid: true}, nil
}

func parseRateValue(input string) (rational.Rational, error) {
	limits := arithmeticLimits()
	var value rational.Rational
	var err error
	if strings.Contains(input, "/") {
		value, err = rational.Parse(input, limits)
	} else {
		amount, amountErr := ParseAmount(input)
		if amountErr != nil {
			err = amountErr
		} else {
			value, err = rationalFromAmount(amount)
		}
	}
	if err != nil {
		return rational.Rational{}, err
	}

	return value, nil
}

// Rational returns the immutable exact math value.
func (rate Rate) Rational() rational.Rational { return rate.value }

// String returns the normalized exact fraction. It does not indicate validity;
// use Valid to distinguish the absent Rate from a parsed zero.
func (rate Rate) String() string { return rate.value.String() }

// IsZero reports whether the multiplier is zero.
func (rate Rate) IsZero() bool { return rate.valid && rate.value.Sign() == 0 }

// Valid reports whether rate was constructed through a validating API.
func (rate Rate) Valid() bool { return rate.valid }

func rationalFromAmount(amount Amount) (rational.Rational, error) {
	return rational.FromBig(amount.Decimal().BigRat(), arithmeticLimits())
}

func arithmeticLimits() gomath.Limits {
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = MaxAmountDigits
	limits.MaxOutputDigits = MaxAmountDigits
	limits.MaxDecimalExpansion = int(MaxScale)
	limits.MaxIntermediateBits = 8192
	limits.MaxDiagnosticBytes = 256

	return limits
}
