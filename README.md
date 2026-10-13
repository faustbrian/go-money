# money

[![CI](https://github.com/faustbrian/go-money/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-money/actions/workflows/ci.yml)
[![Documentation](https://img.shields.io/badge/docs-adoption_and_API-blue)](docs/index.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-money/v2.svg)](https://pkg.go.dev/github.com/faustbrian/go-money/v2)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-money?sort=semver)](https://github.com/faustbrian/go-money/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`money` is an immutable exact monetary value package. It delegates decimal,
integer, and rational arithmetic to `math` and delegates ISO 4217 identity
and metadata to `international/currency`.

The package never accepts or emits `float32` or `float64`. Fixed `Money`
operations preserve currency and resolved context identity. Multiplication and
division return `RationalMoney`; callers choose the context and rounding mode at
the boundary where a fixed amount is required.

## Adoption entry points

The root module is `github.com/faustbrian/go-money/v2` (stable).
It requires Go 1.27.0. Pin the published `v2.0.2` root tag
with `go get github.com/faustbrian/go-money/v2@v2.0.2`.
Root subpackages share this module; they need no separate installation.

Start with the [documentation map](docs/index.md) for package selection,
examples, operational constraints, and contributor references.

## Status, lifecycle, and platform

The published v2.0.0 module requires Go 1.27.0. The supported v1.0.0 module
targets Go 1.26.6. Use a published tag rather than mutable `main` source for
deployment.
Both are portable Go and require no operating-system service or network
connection. Operations do not start background work or acquire resources that
callers must close. A supplied context bounds cancellation only for the
operation that receives it.

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
go get github.com/faustbrian/go-money/v2@v2.0.0
```

The [v2.0.0 release](https://github.com/faustbrian/go-money/releases/tag/v2.0.0)
uses `github.com/faustbrian/go-money/v2`. Existing v1 consumers can retain
their release or follow the [major-version migration guide](docs/compatibility.md).

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
- V2 bounds amount, scale, rate, ratio, allocation, bag-input,
  persistence-input, output, and diagnostic work before expensive conversion.
- Equal and weighted allocations distribute minor-unit remainders
  deterministically and conserve the source total.
- Tax and discount results derive one component by subtraction, so their
  documented totals are conserved.
- Conversion uses only a caller-supplied directed exact rate with timestamp and
  source metadata. There is no live FX client.
- Versioned JSON and SQL representations encode amounts as strings.

## Errors and retry policy

The root package exposes sentinel errors for invalid currencies, contexts,
values, rates, allocations, precision, and configured limits. Errors may carry
operation context, so classify them with `errors.Is` rather than matching error
text. The `encoding` and `format` packages similarly expose
`ErrInvalidEncoding` and `ErrInvalidFormat`.

Validation failures, identity mismatches, precision loss, and limit errors are
permanent local rejections for the same inputs and policy. Retrying them without
changing the value, context, currency, rate, allocation, or configured bound
cannot succeed. A cancellation or deadline error may be retried only as a new,
explicitly bounded operation when the caller's higher-level policy permits it.
The module performs no network access and has no internal transient-error or
retry state.

## Packages

- `money`: values, contexts, arithmetic, allocation, tax, discount, and FX.
- `money/format`: exact locale display using `international` locale tags.
- `money/encoding`: versioned JSON, text, SQL, and PostgreSQL numeric adapters.
- `money/moneytest`: official edge fixtures and conservation assertions.

## Verification

Use ordinary Go tests and a focused two-second active fuzz command from
[CONTRIBUTING.md](CONTRIBUTING.md#verification). `make check` and `make ci`
retain the cheap no-float source check and bounded local assurance with public
Tools v2.1.1. Required CI includes a checkout-only no-float job. Base-free checks
run deterministic seeds/corpora, not active fuzz campaigns. `make check-full`
retains the no-float policy and full exact/conservation, fuzz, race, coverage,
mutation and correctness-gated benchmark assurance for selected risks.

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
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/design-language.md#package-families-and-selection).
