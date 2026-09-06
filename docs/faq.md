# FAQ

## Why no float constructor?

Binary floats cannot represent many decimal fractions exactly. Callers must
provide exact text, minor units, `math` values, or validated persistence.

## Why does `1.230` fail in a scale-2 context?

The extra represented scale is a context difference. Rejecting it avoids silent
normalization. Parse `1.23` or select a scale-3 context explicitly.

## Why does multiplication return `RationalMoney`?

Rates can create repeating or higher-scale results. Returning a rational keeps
the value exact until the caller selects the accounting boundary.

## Does the package fetch exchange rates?

No. Conversion accepts only an injected exact directed rate with attribution.

## Can a `MoneyBag` total different currencies?

No single scalar total exists across currencies without an injected conversion
policy. The bag preserves each currency/context entry independently.

## How should I classify an error?

Use `errors.Is`; do not compare error strings. Root-package validation and
identity failures match `money.ErrUnknownCurrency` and the other exported
root-package sentinels. Malformed versioned representations and unsupported
scanner source types match `encoding.ErrInvalidEncoding`. Invalid values passed
for encoding match `money.ErrInvalidMoney`, while `encoding.ScanNumeric`
preserves the applicable `money.Parse` classification. Display failures match
`format.ErrInvalidFormat`. Wrapped arithmetic and cancellation errors preserve
their underlying classification.

## Should I retry a failed operation?

Not with the same invalid inputs. Unknown currencies, invalid contexts or
values, currency/context mismatches, precision loss, invalid rates or
allocations, and configured-limit failures are deterministic local rejections.
Correct the input or explicitly select a different policy first. Retry
`context.Canceled` or `context.DeadlineExceeded` only as a new bounded operation
when the surrounding workflow allows it; the package has no remote service or
internal backoff to wait for.

## What should I inspect when parsing or arithmetic fails?

Keep diagnosis bounded to the failed value and operation:

1. classify the error with `errors.Is`;
2. confirm the currency code and the resolved context's kind, scale, and cash
   step;
3. confirm decimal scale and digit counts fit the selected context and package
   limits; and
4. for binary operations, compare both operands' currency and full resolved
   context before changing or retrying anything.

Do not log full monetary payloads merely to diagnose a classification error.

## Why was persisted or formatted input rejected?

For `encoding.ErrInvalidEncoding`, check only the bounded representation:
version 1, a string amount, a recognized currency, a valid context, no unknown
or duplicate fields, and at most `encoding.MaxEncodedBytes`. For
`format.ErrInvalidFormat`, validate both the `Money` value and locale and ensure
the result fits `format.MaxFormattedBytes`.
