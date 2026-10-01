package money

import (
	"fmt"
	"strings"

	"github.com/faustbrian/go-international/currency"
	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/integer"
)

// FromMinorUnits constructs Money from an arbitrary-precision integer count of
// units at the supplied fixed context scale.
func FromMinorUnits(units integer.Integer, code currency.Code, context Context) (Money, error) {
	if context.IsZero() || context.kind == ContextAutomatic {
		return Money{}, ErrInvalidContext
	}
	if err := validateCurrencyContext(code, context); err != nil {
		return Money{}, err
	}

	// Compare against a fixed monetary boundary before formatting a caller's
	// arbitrary-precision coefficient. Negation copies only the fixed bound.
	limits := arithmeticLimits()
	limits.MaxInputDigits = MaxAmountDigits + 1
	maximum := mustInvariant(integer.Parse("1"+strings.Repeat("0", MaxAmountDigits), integer.ParseOptions{
		Base: 10, Limits: limits,
	}))
	if units.Cmp(maximum) >= 0 || units.Cmp(maximum.Neg()) <= 0 {
		return Money{}, fmt.Errorf("money: parse amount: %w", gomath.ErrLimitExceeded)
	}

	return Parse(decimalTextFromMinor(units.String(), context.scale), code, context)
}

// MinorUnits returns the exact arbitrary-precision integer coefficient at the
// Money context's resolved scale.
func (money Money) MinorUnits() (integer.Integer, error) {
	if !money.Valid() {
		return integer.Integer{}, ErrInvalidMoney
	}

	text := strings.ReplaceAll(money.amount.String(), ".", "")
	return integer.Parse(text, integer.ParseOptions{
		Base:              10,
		AllowLeadingZeros: true,
		Limits:            arithmeticLimits(),
	})
}

func decimalTextFromMinor(text string, scale uint8) string {
	if scale == 0 {
		return text
	}

	sign := ""
	if strings.HasPrefix(text, "-") {
		sign = "-"
		text = strings.TrimPrefix(text, "-")
	}
	places := int(scale)
	if len(text) <= places {
		text = strings.Repeat("0", places-len(text)+1) + text
	}
	point := len(text) - places

	return sign + text[:point] + "." + text[point:]
}
