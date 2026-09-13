# Compatibility

The minimum supported toolchain is Go 1.27.0. The released v1 public API remains
captured byte-for-byte in `api/v1.txt`. Planned v2 source uses the semantic
import path `github.com/faustbrian/go-money/v2`; its active compatibility
baseline is `api/v2.txt`.

The v2 boundary is required because rejecting more than
`MaxMoneyBagEntries` total constructor inputs changes v1 acceptance behavior,
even when duplicate identities would combine to fewer output entries. Existing
consumers must remain on v1 until v2 is published, then opt in by updating
imports and reviewing the documented ceilings.

The directly owned `go-knapsack/objective/gomoney` and
`go-knapsack/objective/money` modules and the `go-library-tools` release
compatibility consumer remain pinned to `github.com/faustbrian/go-money` v1.
Their `/v2` migration is blocked on publication of the first v2 release and is
not part of this source-only transition.

Persistence compatibility is independent of Go API compatibility. Version-1
decoders reject unknown versions and fields. Future representations require a
new version and migration documentation rather than changing version-1 meaning.

Currency metadata compatibility belongs to `international`. A persisted
historic code remains that code; this package does not map it to a successor.
Decimal and rational semantics belong to `math`; this package does not carry
a competing implementation.
