# money

[![CI](https://github.com/faustbrian/go-money/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-money/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-money/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-money.svg)](https://pkg.go.dev/github.com/faustbrian/go-money)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-money?sort=semver)](https://github.com/faustbrian/go-money/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`money` is an immutable exact monetary value package. It delegates decimal,
integer, and rational arithmetic to `math` and delegates ISO 4217 identity
and metadata to `international/currency`.

The package never accepts or emits `float32` or `float64`. Fixed `Money`
operations preserve currency and resolved context identity. Multiplication and
division return `RationalMoney`; callers choose the context and rounding mode at
the boundary where a fixed amount is required.

## Status, lifecycle, and platform

The module is stable at v1 and requires Go 1.26.6. It is portable Go and does
not require an operating-system service or a network connection. Operations do
not start background work or acquire resources that callers must close. A
supplied context bounds cancellation only for the operation that receives it.

## When to use it

Use `money` when monetary values must remain exact and currency, precision,
rounding, allocation, or conversion policy must be explicit. Use the package
for value calculations and persistence boundaries, with exchange rates and
business rules supplied by the caller.

## When not to use it

Do not use it as an exchange-rate client, pricing or accounting system, payment
processor, currency registry, or process-wide monetary configuration. Do not
use it when binary floating-point inputs or implicit rounding are required.

## Install

```sh
go get github.com/faustbrian/go-money
```

## Quick start

```go
euro, _ := currency.Parse("EUR")
ctx, _ := money.DefaultContext(euro)
left, _ := money.Parse("12.30", euro, ctx)
right, _ := money.Parse("0.45", euro, ctx)
total, _ := left.Add(right) // 12.75 EUR
```

Exact multiplication remains rational until rounding is explicit:

```go
rate, _ := money.ParseRate("1/3")
exact, _ := total.Mul(context.Background(), rate)
fixed, result, _ := exact.Round(ctx, gomath.RoundHalfEven)
_ = result.Inexact()
```

## Guarantees

- Cross-currency or cross-context arithmetic returns an error.
- Default scales come from authoritative ISO metadata; historic currencies are
  accepted only through the explicit historic parse policy.
- Amount, scale, rate, ratio, allocation, output, and diagnostic work is
  bounded.
- Equal and weighted allocations distribute minor-unit remainders
  deterministically and conserve the source total.
- Tax and discount results derive one component by subtraction, so their
  documented totals are conserved.
- Conversion uses only a caller-supplied directed exact rate with timestamp and
  source metadata. There is no live FX client.
- Versioned JSON and SQL representations encode amounts as strings.

## Packages

- `money`: values, contexts, arithmetic, allocation, tax, discount, and FX.
- `money/format`: exact locale display using `international` locale tags.
- `money/encoding`: versioned JSON, text, SQL, and PostgreSQL numeric adapters.
- `money/moneytest`: official edge fixtures and conservation assertions.

## Verification

`make check` runs formatting, analysis, tests, meaningful 100% production
coverage, race checks, docs, compatibility, dependency, float-contamination,
vulnerability, fuzzing, mutation, and correctness-gated comparative benchmark
gates through the released `go-library-tools` contract.

See [docs/api.md](docs/api.md),
[docs/contexts-and-rounding.md](docs/contexts-and-rounding.md), and
[docs/cookbook.md](docs/cookbook.md).

## Documentation

Use the [documentation index](docs/README.md) for package-owned guides,
operational contracts, [executable examples](examples_test.go), and maintainer
references. See the [FAQ and troubleshooting guide](docs/faq.md),
[support policy](SUPPORT.md), and [security policy](SECURITY.md) for help and
reporting routes.

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).
