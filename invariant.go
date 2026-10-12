package money

const invariantViolation = "money: internal invariant violated"

// mustInvariant is reserved for results established by prior admission or
// arithmetic invariants. Its panic message is fixed rather than echoing a
// dependency diagnostic; caller-reachable validation must return errors instead.
func mustInvariant[T any](value T, err error) T {
	if err != nil {
		panic(invariantViolation)
	}

	return value
}
