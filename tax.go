package money

import (
	"context"
	"fmt"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

// MaxTaxRate bounds a tax multiplier to 1000 percent.
const MaxTaxRate int64 = 10

// TaxRate is an exact validated nonnegative tax fraction: 0.24 means 24%.
// Its zero value is invalid; a parsed zero rate is valid.
type TaxRate struct{ rate Rate }

// ParseTaxRate parses an exact decimal or fraction from zero through MaxTaxRate
// without binary floating point. Invalid input returns ErrInvalidRate; the
// caller is responsible for selecting a jurisdictionally applicable rate.
func ParseTaxRate(input string) (TaxRate, error) {
	rate, err := ParseRate(input)
	if err != nil {
		return TaxRate{}, err
	}
	maximum, _ := rational.New(MaxTaxRate, 1)
	if rate.value.Cmp(maximum) > 0 {
		return TaxRate{}, ErrInvalidRate
	}

	return TaxRate{rate: rate}, nil
}

// Rate returns the generic exact multiplier.
func (rate TaxRate) Rate() Rate { return rate.rate }

// TaxResult contains conserved net, tax, and gross components.
type TaxResult struct {
	net      Money
	tax      Money
	gross    Money
	rounding RoundingResult
}

// Net returns the amount excluding tax.
func (result TaxResult) Net() Money { return result.net }

// Tax returns the conserved tax component. AddTax rounds this component;
// ExtractTax derives it as gross minus rounded net, without separately rounding tax.
func (result TaxResult) Tax() Money { return result.tax }

// Gross returns the amount including tax.
func (result TaxResult) Gross() Money { return result.gross }

// Rounding returns conditions raised at the explicit monetary boundary.
func (result TaxResult) Rounding() RoundingResult { return result.rounding }

// AddTax calculates tax from a net amount, rounds the tax once, then derives
// gross by addition so net + tax always equals gross on success. It uses the
// net context and requires valid Money, a parsed rate, non-nil ctx and a valid
// mode. Automatic contexts can fail if the rounded component scale differs
// from net. Signed amounts follow arithmetic, not jurisdictional tax policy.
func AddTax(ctx context.Context, net Money, rate TaxRate, mode gomath.RoundingMode) (TaxResult, error) {
	if !rate.rate.Valid() {
		return TaxResult{}, ErrInvalidRate
	}
	exact, err := net.Mul(ctx, rate.rate)
	if err != nil {
		return TaxResult{}, err
	}
	tax, rounding, err := exact.Round(net.context, mode)
	if err != nil {
		return TaxResult{}, err
	}
	gross, err := net.Add(tax)
	if err != nil {
		return TaxResult{}, err
	}

	return TaxResult{net: net, tax: tax, gross: gross, rounding: rounding}, nil
}

// ExtractTax derives exact net from gross/(1+rate), rounds net once, then
// subtracts it from gross so the returned components conserve gross on success.
// Rounding describes the net boundary, not an independently rounded tax. It
// requires valid Money, a parsed rate, non-nil ctx and a valid mode. Automatic
// contexts whose resolved net scale differs from gross trigger an
// internal-invariant panic during subtraction; use a fixed context.
func ExtractTax(ctx context.Context, gross Money, rate TaxRate, mode gomath.RoundingMode) (TaxResult, error) {
	if !gross.Valid() {
		return TaxResult{}, ErrInvalidMoney
	}
	if !rate.rate.Valid() {
		return TaxResult{}, ErrInvalidRate
	}
	one, _ := rational.New(1, 1)
	factor, err := one.Add(ctx, rate.rate.value, arithmeticLimits())
	if err != nil {
		return TaxResult{}, fmt.Errorf("money: tax factor: %w", err)
	}
	grossRational, err := rationalFromAmount(gross.amount)
	if err != nil {
		return TaxResult{}, err
	}
	netRational, err := grossRational.Quo(ctx, factor, arithmeticLimits())
	if err != nil {
		return TaxResult{}, fmt.Errorf("money: extract tax: %w", err)
	}
	exact := RationalMoney{amount: netRational, currency: gross.currency, source: gross.context}
	net, rounding, err := exact.Round(gross.context, mode)
	if err != nil {
		return TaxResult{}, err
	}
	tax := mustInvariant(gross.Sub(net))

	return TaxResult{net: net, tax: tax, gross: gross, rounding: rounding}, nil
}
