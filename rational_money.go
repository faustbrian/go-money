package money

import (
	"context"
	"fmt"
	"math/big"

	"github.com/faustbrian/go-international/v5/currency"
	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

// RationalMoney is an immutable exact fractional monetary result. It must be
// rounded explicitly before persistence as Money. Its zero value is absent;
// String returns empty text and Round returns ErrInvalidMoney.
type RationalMoney struct {
	amount   rational.Rational
	currency currency.Code
	source   Context
}

// Ratio is an immutable exact signed relationship between compatible monetary
// values. The zero value is absent: String returns empty text, while Rational
// returns numeric zero without conveying validity.
type Ratio struct {
	value rational.Rational
	valid bool
}

// Rational returns the exact math fraction.
func (ratio Ratio) Rational() rational.Rational { return ratio.value }

// String returns the normalized fraction or empty text for an absent ratio.
func (ratio Ratio) String() string {
	if !ratio.valid {
		return ""
	}

	return ratio.value.String()
}

// RoundingResult describes whether an explicit boundary discarded value.
type RoundingResult struct{ conditions gomath.Condition }

// Conditions returns the math arithmetic conditions.
func (result RoundingResult) Conditions() gomath.Condition { return result.conditions }

// Inexact reports whether nonzero value was discarded.
func (result RoundingResult) Inexact() bool {
	return result.conditions.Has(gomath.ConditionInexact)
}

// String returns the normalized exact fraction and currency code.
func (money RationalMoney) String() string {
	if money.currency.IsZero() {
		return ""
	}

	return money.amount.String() + " " + money.currency.String()
}

// Rational returns the immutable exact math value.
func (money RationalMoney) Rational() rational.Rational { return money.amount }

// Currency returns the monetary identity.
func (money RationalMoney) Currency() currency.Code { return money.currency }

// Round creates Money at an explicit rounding boundary. Default targets must
// belong to the same currency; custom targets use their scale and cash targets
// round to their step. Automatic targets require a terminating exact decimal
// within MaxScale places or return ErrPrecisionLoss. Absent money/target or an
// invalid rounding mode returns ErrInvalidMoney. Failure returns invalid Money
// and an empty RoundingResult; no partial rounded value is returned.
func (money RationalMoney) Round(target Context, mode gomath.RoundingMode) (Money, RoundingResult, error) {
	if money.currency.IsZero() || target.IsZero() || !mode.Valid() {
		return Money{}, RoundingResult{}, ErrInvalidMoney
	}
	if target.kind == ContextAutomatic {
		scale, err := terminatingScale(money.amount)
		if err != nil {
			return Money{}, RoundingResult{}, err
		}
		target.scale = scale
	}
	if target.kind == ContextDefault && target.currency != money.currency {
		return Money{}, RoundingResult{}, ErrContextMismatch
	}
	text, conditions, err := roundedDecimal(money.amount, target, mode)
	if err != nil {
		return Money{}, RoundingResult{}, fmt.Errorf("money: round rational: %w", err)
	}
	value, err := Parse(text, money.currency, target)
	if err != nil {
		return Money{}, RoundingResult{}, err
	}

	return value, RoundingResult{conditions: conditions}, nil
}

func terminatingScale(value rational.Rational) (uint8, error) {
	// A reduced fraction terminates in base ten iff its denominator contains
	// only factors 2 and 5. The larger exponent is the required decimal scale.
	denominator := value.Denominator()
	denominator, twos := removeFactor(denominator, 2)
	denominator, fives := removeFactor(denominator, 5)
	if denominator.Cmp(big.NewInt(1)) != 0 || twos > MaxScale || fives > MaxScale {
		return 0, ErrPrecisionLoss
	}
	return max(twos, fives), nil
}

func removeFactor(denominator *big.Int, factor int64) (*big.Int, uint8) {
	remainder := new(big.Int)
	divisor := big.NewInt(factor)
	count := uint8(0)
	// One factor beyond the admitted scale is enough to reject the expansion;
	// do not scan arbitrarily large powers supplied through rational arithmetic.
	for range MaxScale + 1 {
		quotient, _ := new(big.Int).QuoRem(denominator, divisor, remainder)
		if remainder.Sign() != 0 {
			return denominator, count
		}
		denominator = quotient
		count++
	}

	return denominator, count
}

func roundedDecimal(amount rational.Rational, target Context, mode gomath.RoundingMode) (string, gomath.Condition, error) {
	if target.kind != ContextCash {
		return amount.Decimal(int(target.scale), mode, arithmeticLimits())
	}

	// Round an integer count of cash increments once, then restore decimal
	// units. Decimal-place rounding alone would not enforce the cash step.
	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(target.scale)), nil)
	step := mustInvariant(rational.NewChecked(
		new(big.Int).SetUint64(target.cashStep), denominator, arithmeticLimits(),
	))
	var err error
	units, err := amount.Quo(context.Background(), step, arithmeticLimits())
	if err != nil {
		return "", 0, err
	}
	integerText, conditions, err := units.Decimal(0, mode, arithmeticLimits())
	if err != nil {
		return "", 0, err
	}
	integer := mustInvariant(rational.Parse(integerText, arithmeticLimits()))
	rounded := mustInvariant(integer.Mul(context.Background(), step, arithmeticLimits()))
	text, finalConditions, err := rounded.Decimal(int(target.scale), mode, arithmeticLimits())

	return text, conditions | finalConditions, err
}

// Mul returns an exact rational result without applying the Money context's
// fixed scale. The rate must be valid (zero is accepted); nil ctx and bounded
// arithmetic failures return errors. The input Money remains unchanged.
func (money Money) Mul(ctx context.Context, rate Rate) (RationalMoney, error) {
	if !money.Valid() {
		return RationalMoney{}, ErrInvalidMoney
	}
	if !rate.Valid() {
		return RationalMoney{}, ErrInvalidRate
	}
	amount, err := rationalFromAmount(money.amount)
	if err != nil {
		return RationalMoney{}, err
	}
	product, err := amount.Mul(ctx, rate.value, arithmeticLimits())
	if err != nil {
		return RationalMoney{}, fmt.Errorf("money: multiply: %w", err)
	}

	return RationalMoney{amount: product, currency: money.currency, source: money.context}, nil
}

// Quo divides by an exact nonzero rate without rounding. An invalid or zero
// rate returns ErrInvalidRate; nil ctx and bounded arithmetic failures return
// errors. The input Money remains unchanged.
func (money Money) Quo(ctx context.Context, rate Rate) (RationalMoney, error) {
	if !money.Valid() {
		return RationalMoney{}, ErrInvalidMoney
	}
	if !rate.Valid() || rate.IsZero() {
		return RationalMoney{}, ErrInvalidRate
	}
	amount, err := rationalFromAmount(money.amount)
	if err != nil {
		return RationalMoney{}, err
	}
	quotient, err := amount.Quo(ctx, rate.value, arithmeticLimits())
	if err != nil {
		return RationalMoney{}, fmt.Errorf("money: divide: %w", err)
	}

	return RationalMoney{amount: quotient, currency: money.currency, source: money.context}, nil
}

// Ratio returns money/other exactly after verifying currency and resolved
// context. Division by numeric zero wraps go-math.ErrDivisionByZero. Nil ctx
// and bounded arithmetic failures return errors; no partial ratio is returned.
func (money Money) Ratio(ctx context.Context, other Money) (Ratio, error) {
	if err := money.compatible(other); err != nil {
		return Ratio{}, err
	}
	left, err := rationalFromAmount(money.amount)
	if err != nil {
		return Ratio{}, err
	}
	right, err := rationalFromAmount(other.amount)
	if err != nil {
		return Ratio{}, err
	}
	value, err := left.Quo(ctx, right, arithmeticLimits())
	if err != nil {
		return Ratio{}, fmt.Errorf("money: ratio: %w", err)
	}

	return Ratio{value: value, valid: true}, nil
}
