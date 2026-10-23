package money

import (
	"github.com/faustbrian/go-international/v5/currency"
)

// MaxScale bounds decimal work and serialized representation sizes.
const MaxScale uint8 = 18

// MaxCashStep bounds cash increments expressed at their context scale.
const MaxCashStep uint64 = 1_000_000_000_000_000_000

func validateCurrencyContext(code currency.Code, context Context) error {
	if code.IsZero() {
		return ErrUnknownCurrency
	}
	if context.IsZero() || context.scale > MaxScale {
		return ErrInvalidContext
	}
	if context.kind == ContextDefault && context.currency != code {
		return ErrContextMismatch
	}
	return nil
}

// ContextKind identifies the monetary precision policy without carrying
// currency metadata into the arithmetic layer.
type ContextKind uint8

const (
	contextInvalid ContextKind = iota
	// ContextDefault uses the authoritative ISO minor-unit exponent captured
	// when the context is constructed.
	ContextDefault
	// ContextCustom uses an application-selected decimal scale.
	ContextCustom
	// ContextCash applies a positive cash increment at a selected scale.
	ContextCash
	// ContextAutomatic preserves exact input scale until an operation requires
	// an explicit rounding boundary.
	ContextAutomatic
)

// Context is an immutable, comparable monetary precision policy. Currency
// identity and metadata are owned by go-international; default contexts retain
// the currency binding needed to validate their captured scale. The
// zero value is invalid. Equality includes kind, scale, cash step and binding,
// not just the number of fractional places.
type Context struct {
	kind     ContextKind
	scale    uint8
	cashStep uint64
	currency currency.Code
}

// DefaultContext constructs a fixed-scale context from authoritative ISO 4217
// metadata. Currencies without an applicable minor-unit exponent require an
// explicit custom context. A zero code returns ErrUnknownCurrency; absent minor
// units return ErrMinorUnitsUnavailable. Metadata comes from the pinned
// go-international dataset, not a live currency service.
func DefaultContext(code currency.Code) (Context, error) {
	if code.IsZero() {
		return Context{}, ErrUnknownCurrency
	}

	scale, ok := code.MinorUnits()
	if !ok {
		return Context{}, ErrMinorUnitsUnavailable
	}
	return Context{kind: ContextDefault, scale: scale, currency: code}, nil
}

// CustomContext constructs an application-selected fixed decimal scale from
// zero through MaxScale. Larger scales return ErrInvalidContext.
func CustomContext(scale uint8) (Context, error) {
	if scale > MaxScale {
		return Context{}, ErrInvalidContext
	}

	return Context{kind: ContextCustom, scale: scale}, nil
}

// CashContext constructs a fixed scale with a positive cash increment. Step is
// expressed as an integer count of units at scale: scale 2 and step 5 means
// 0.05. Scale must not exceed MaxScale and step must be in [1, MaxCashStep],
// otherwise ErrInvalidContext is returned. Parse does not enforce divisibility
// by step; RationalMoney.Round is the cash-rounding boundary.
func CashContext(scale uint8, step uint64) (Context, error) {
	if scale > MaxScale || step == 0 || step > MaxCashStep {
		return Context{}, ErrInvalidContext
	}

	return Context{kind: ContextCash, scale: scale, cashStep: step}, nil
}

// AutomaticContext constructs a policy that captures exact input scale. The
// resolved Money context carries that scale, making representation differences
// explicit during later arithmetic. RationalMoney.Round accepts this policy
// only for exact decimal expansions terminating within MaxScale places.
func AutomaticContext() Context { return Context{kind: ContextAutomatic} }

// Scale returns the number of decimal fractional digits. An automatic context
// starts at zero and a Money value carries its resolved input or result scale.
func (context Context) Scale() uint8 { return context.scale }

// Kind returns the context's precision policy.
func (context Context) Kind() ContextKind { return context.kind }

// CashStep returns the integer number of scale-sized units in one cash
// increment. It is zero for non-cash contexts.
func (context Context) CashStep() uint64 { return context.cashStep }

// IsZero reports whether context is absent or invalid.
func (context Context) IsZero() bool { return context.kind == contextInvalid }
