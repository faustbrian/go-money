# Security model

The released v1 line remains supported. This document also describes the
stricter fail-closed ceilings in planned, unpublished v2 source.

## Assets and trust boundaries

The protected assets are exact monetary values, currency and context identity,
allocation conservation, deterministic persistence, and caller availability.
Untrusted data crosses the boundary through decimal, rate, currency, locale,
JSON, text, SQL, and PostgreSQL numeric inputs. Caller-created `go-math`
decimal and integer values, allocation ratios, and money-bag slices are also
untrusted because their magnitude or count can exceed monetary policy even
when their Go types are valid.

The module does not authenticate callers, authorize database access, fetch
exchange rates, open files or sockets, start goroutines, or process payments.
Those controls remain with the application.

## Bounds and fail-closed controls

Amounts are limited to 256 digits and 18 fractional places. Allocation ratios
are limited to 64 digits, allocation outputs to 10,000 parts, total money-bag
inputs and entries to 1,000, exchange-rate source attribution to 128 bytes,
encoded and formatted values to 2,048 bytes, monetary intermediate arithmetic
to 8,192 bits, and errors to fixed category strings. Rate magnitude defaults to
1,000,000, tax rates to 10, and cash steps to 1,000,000,000,000,000,000.
Planned v2 checks scanner byte lengths before copying and admits
typed decimal and integer operands under bounded arithmetic before converting
them to text. JSON traversal is bounded by the encoded-byte ceiling, rejects
duplicate keys and trailing data, and accepts only the versioned closed schema.

The package rejects absent or unknown currencies, unavailable default minor
units, mismatched currency or context, negative or excessive rates, zero or
negative allocation weights, excessive part or bag-input counts,
non-terminating automatic precision, malformed versions, and
precision-losing construction. Context-aware arithmetic and allocation return
cancellation or deadline errors from the supplied operation context.

## Ownership, concurrency, and privacy

Public values are immutable. Mutable `math/big` internals are owned by
`go-math`; accessors return immutable values or defensive copies. Bag and
allocation slices are copied. There is no global mutable monetary policy and
no background lifecycle. Shared values and formatters are exercised under the
race detector.

Errors contain bounded operation categories, not source monetary records,
credentials, or customer identifiers. The package has no logger. Applications
should avoid attaching raw payment or customer data when wrapping errors.

## Verification and accepted risk

Hostile boundary tests cover exact limits, oversized encodings, excessive
collection counts, huge typed numeric operands, cancellation, malformed and
nested JSON, and diagnostic redaction. Fuzzing covers money and rate parsing,
allocation conservation, versioned JSON, PostgreSQL numeric scanning, and
locale formatting; benchmarks exercise bounded work and output.

### MONEY-RISK-001: application-level resource policy

- **Severity:** low.
- **Disposition:** accepted.
- **Owner:** go-money maintainers.
- **Rationale:** package ceilings protect library availability but cannot encode
  business, tenant, request, or database limits without application identity
  and policy context.
- **Mitigation:** applications must impose narrower amount, rate, request,
  database, and tenant quotas appropriate to their domain before calling the
  package.
- **Review condition:** reconsider this disposition if the module begins owning
  I/O, authentication, persistence, tenancy, or background work, or if a
  package ceiling proves too large for a supported deployment class.
