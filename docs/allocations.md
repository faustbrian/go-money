# Allocations

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

`EqualSplit` divides exact minor units by a positive bounded part count. Any
signed remainder is assigned one unit at a time from the first part onward.

`Allocate` accepts positive arbitrary-precision integer ratios. It calculates
each absolute share, orders fractional remainders from largest to smallest,
uses original ratio order as the stable tie-break, and reapplies the source
sign. Zero and negative ratios are rejected.

Since v2, the module admits each typed ratio under the monetary arithmetic ceiling before
formatting it, and rejects more than `MaxAllocationParts` ratios. Money-bag
construction similarly rejects more than `MaxMoneyBagEntries` total inputs,
including duplicate identities.

Both algorithms are deterministic and conserve the original total for positive,
negative, and zero amounts. `AllocationResult.Parts` returns an independent
slice, and `Sum` verifies the conserved identity.
