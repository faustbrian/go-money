# Verification

Ordinary `make check` and `make ci` execute the bounded local contract through
the public `golib` release declared in `.golib.yaml`: deterministic unit,
matrix, property, official-currency, SQL and JSON tests, maintained-peer
comparisons, formatting, vet, static analysis, lint, docs, API, dependency and
applicable security checks. Base-free Make checks run fuzz seeds/corpora;
changed parser inputs receive causal active smoke in CI. A focused local
two-second active fuzz command is documented in
[CONTRIBUTING.md](../CONTRIBUTING.md#verification).

Ordinary Make checks also retain the cheap no-float source policy. CI runs the
same existing script in a checkout-only job bounded to two minutes, without
Go or tooling bootstrap. Required fails unless both it and shared CI succeed.

`make check-full` retains the original full gates for selected material risks
or release qualification:

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

The no-float check is separately reproducible with
`make -f verification/package.mk no-float`; it also remains in the full contract.
`make ci` retains repository and cohesion validation without selecting the full
campaign. CI uses the same released verifier and pinned shared workflow.
A skipped, unavailable, or advisory-only substitution for a required gate is
not release evidence.
