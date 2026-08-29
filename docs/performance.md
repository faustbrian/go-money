# Performance

Correctness and bounds take precedence over allocation count. Exact
arbitrary-precision arithmetic is expected to cost more than fixed `int64`
money implementations.

The benchmark suite gates identical EUR addition results before timing
`money`, `govalues/money`, and `Rhymond/money`. On an Apple M4 Max with Go
1.26.5, a 100-iteration local sample measured approximately 587 ns/op and 21
allocations for `money`, 12 ns/op and zero allocations for `govalues/money`,
and 337 ns/op and two allocations for `Rhymond/money`. These figures are not
portable performance promises.

The ordinary test suite separately compares signed fixed-scale addition,
subtraction, equal splits, and currency-mismatch rejection with both mature
packages. Weighted allocation and cash-rounding policies use independent
conservation and hand-calculated matrices because those policies are not
identical across packages.

The `benchmark` gate runs the benchmark suite on the target architecture after
its correctness preflight. Benchmark results are valid only when that
preflight passes. Prefer minor-unit fixed operations for hot paths and retain
rational results only until the required boundary.

Run `golib check --module .` for the complete contract, or `make check` as the
repository shorthand. The shared tool reports benchmark output and keeps its
Go caches and temporary files task-owned.
