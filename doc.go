// Package money provides immutable, bounded, exact monetary values backed by
// go-math and identified by go-international currency codes. Start with Parse
// and DefaultContext, or select an explicit custom, cash or automatic context.
//
// Fixed Money never rounds implicitly. Multiplication and division produce
// RationalMoney, which requires an explicit context and rounding mode before it
// can return Money. Currency and resolved context mismatches are errors, even
// when the numeric amounts or scales agree. Cash construction preserves the
// amount; only RationalMoney.Round enforces the cash increment.
//
// Values can be shared across goroutines. Operations create new values, and
// collection accessors copy their slices. Context-aware arithmetic requires a
// non-nil context and observes cancellation between bounded arithmetic stages;
// it does not preempt individual big-integer primitives. Round is synchronous
// and has no cancellation context. Resource-limit failures can wrap go-math
// errors as well as this package's sentinels; classify errors with errors.Is.
//
// This package does not fetch exchange rates, choose tax policy, process payments
// or own database transactions. The encoding, format and moneytest subpackages
// provide persistence, display and test support. No package requires cgo,
// external services or a background lifecycle.
package money
