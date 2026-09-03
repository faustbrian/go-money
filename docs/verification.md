# Verification

The repository contract is executed by `golib` through `make check` or `make
ci`. The enabled gates are:

- `golib check`: unit, matrix, property, official-currency, SQL, JSON, and
  differential tests against `govalues/money` and `Rhymond/money`.
- exact 100% production statement coverage with no threshold rounding or
  aggregate masking.
- race checks for shared immutable values and formatter concurrency.
- fuzzing of decimal, rate, allocation, JSON, PostgreSQL numeric, and locale
  formatting fuzzing.
- mutation testing with exact 100% efficacy and mutant coverage, including
  mismatch, rounding,
  remainder, sign, rate, tax, discount, and conversion mutants.
- correctness-gated benchmark comparisons against maintained packages.
- formatting, vet, static analysis, lint, advisory NilAway, docs, API,
  dependencies, float contamination, and vulnerability checks.

`make ci` runs repository validation and the complete enabled contract. CI uses
the same released `go-library-tools` workflow and configuration. A skipped,
unavailable, or advisory-only substitution for a required gate is not release
evidence.
