package money

import (
	"context"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

// DiscountRate is an exact validated fraction in the inclusive range [0, 1].
// Its zero value is invalid, unlike a parsed zero discount.
type DiscountRate struct{ rate Rate }

// ParseDiscountRate parses an exact decimal or fraction from zero through one
// without floats. Malformed or out-of-range values return ErrInvalidRate.
func ParseDiscountRate(input string) (DiscountRate, error) {
	rate, err := ParseRate(input)
	if err != nil {
		return DiscountRate{}, err
	}
	one, _ := rational.New(1, 1)
	if rate.value.Cmp(one) > 0 {
		return DiscountRate{}, ErrInvalidRate
	}

	return DiscountRate{rate: rate}, nil
}

// Rate returns the generic exact multiplier.
func (rate DiscountRate) Rate() Rate { return rate.rate }

// DiscountResult contains conserved original, discount, and final amounts.
type DiscountResult struct {
	original Money
	discount Money
	final    Money
	rounding RoundingResult
}

// Original returns the amount before discount.
func (result DiscountResult) Original() Money { return result.original }

// Discount returns the rounded discount component.
func (result DiscountResult) Discount() Money { return result.discount }

// Final returns original minus discount.
func (result DiscountResult) Final() Money { return result.final }

// Rounding returns conditions raised at the monetary boundary.
func (result DiscountResult) Rounding() RoundingResult { return result.rounding }

// ApplyDiscount rounds the discount once and derives final by subtraction so
// final + discount always equals original on success. It uses the original
// context and requires valid Money, a parsed rate, non-nil ctx and a valid mode.
// With an automatic context, a rounded scale differing from the original
// resolved scale triggers an internal-invariant panic during subtraction. Use
// a fixed context for this operation. Signed values retain their arithmetic
// sign; business rules about negative prices or discounts remain caller policy.
func ApplyDiscount(ctx context.Context, original Money, rate DiscountRate, mode gomath.RoundingMode) (DiscountResult, error) {
	if !rate.rate.Valid() {
		return DiscountResult{}, ErrInvalidRate
	}
	exact, err := original.Mul(ctx, rate.rate)
	if err != nil {
		return DiscountResult{}, err
	}
	discount, rounding, err := exact.Round(original.context, mode)
	if err != nil {
		return DiscountResult{}, err
	}
	final := mustInvariant(original.Sub(discount))

	return DiscountResult{
		original: original,
		discount: discount,
		final:    final,
		rounding: rounding,
	}, nil
}
