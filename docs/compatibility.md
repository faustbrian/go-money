# Compatibility

Current v2 source requires Go 1.27.0; published v1.0.0 targets Go 1.26.6.
The released v1 public API remains
captured byte-for-byte in `api/v1.txt`. Planned v2 source uses the semantic
import path `github.com/faustbrian/go-money/v2`; its active compatibility
baseline is `api/v2.txt`.

The v2 boundary is required because rejecting more than
`MaxMoneyBagEntries` total constructor inputs changes v1 acceptance behavior,
even when duplicate identities would combine to fewer output entries. Existing
consumers must remain on v1 until v2 is published, then opt in by updating
imports and reviewing the documented ceilings.

The maintained `go-knapsack/objective/gomoney/v2` and
`go-knapsack/objective/money/v2` modules and the `go-library-tools` release
compatibility consumer still select `github.com/faustbrian/go-money` v1.0.0.
Their adoption of Money v2 requires publication and separate consumer review;
the released adapters expose nominal Money types, so changing those signatures
requires their own major-version decision. Frozen historical compatibility
cohorts remain unchanged.

Update all Money package imports together when adopting v2. V1 and v2 Money,
Context, Amount, and sentinel values are distinct; do not assume cross-major
`errors.Is` matches a Money sentinel from the other module. Shared Math error
classifications still belong to the selected Math module. Minor-unit rejection
retains the original wrapped Math limit classification and validates currency
and context metadata before formatting a coefficient.

Persistence compatibility is independent of Go API compatibility. Version-1
decoders reject unknown versions and fields. Future representations require a
new version and migration documentation rather than changing version-1 meaning.

Currency metadata compatibility belongs to `international`. A persisted
historic code remains that code; this package does not map it to a successor.
Decimal and rational semantics belong to `math`; this package does not carry
a competing implementation.
