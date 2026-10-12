# Source documentation audit

## Scope and plan

Baseline: `14c26b349c6d27886b5da2717d38ae7efd1d3f3d` (main).
Owner: go-money source-documentation pass. Tier A: comments and audit evidence
only; preserve all declarations, executable tokens, examples and test behavior.
The public module is `github.com/faustbrian/go-money/v2` in `go.mod`; there are
four packages: root money, encoding, format and moneytest. Consumers are users
of these APIs and the package-owned examples. No consumer API or dependency
change is planned. Named risks are misleading precision/rounding, cancellation,
input admission, ownership and persistence documentation.

The smallest sufficient checks are formatting, objective source-documentation
structure, unchanged executable-token comparison, local links and marker/
generator/secret inspection, representative rendered Go documentation and one
ordinary complete-diff review. No runtime campaign or release is warranted.
Shared executable documentation enforcement belongs to go-library-tools; its
coordinator-owned integration is not implemented here.

## Reproducible inventory

Run `git ls-files '*go.mod' '*.go'` for module/file coverage; inspect package
clauses and every declaration in those tracked files. Run
`rg -n 'TO[D]O|FI[X]ME|HA[C]K|NO[T]E|WAR[N]ING|Deprecated:|Code generated|go:generate|go:build' --glob '*.go'`
for special markers. The initial inventory has 45 Go files, one module, no build
variants, cgo, generated Go files, generation directives, deprecated declarations
or debt markers. No vendor files are present. Each row below was read in full,
including all unexported helpers and struct/interface members; test declarations
are listed as coverage evidence, not public library APIs.

| File | Package | Role | Exported declarations/members inspected |
| --- | --- | --- | --- |
| `allocation.go` | `money` | API and implementation | `MaxAllocationParts`, `MaxRatioDigits`, `AllocationResult`, `Parts`, `Sum`, `EqualSplit`, `Allocate` |
| `allocation_test.go` | `money` | tests | `TestEqualSplitConservesPositiveAndNegativeTotalsDeterministically`, `TestWeightedAllocationUsesStableLargestRemainders` |
| `amount.go` | `money` | API and implementation | `MaxAmountDigits`, `Amount`, `ParseAmount`, `AmountFromDecimal`, `Decimal`, `String`, `Scale`, `Sign`, `IsZero`, `Equal` |
| `amount_test.go` | `money` | tests | `TestAmountIsAnExactBoundedDecimalDomainValue` |
| `bag.go` | `money` | API and implementation | `MaxMoneyBagEntries`, `MoneyBag`, `NewMoneyBag`, `Add`, `Values`, `Get` |
| `bag_test.go` | `money` | tests | `TestNewMoneyBagRejectsExcessiveInputValues`, `TestNewMoneyBagAcceptsExactInputLimitAndConservesTotal`, `TestMoneyBagCombinesOnlyIdenticalCurrencyAndContexts` |
| `benchmark_test.go` | `money` | tests | `BenchmarkExactAdd`, `BenchmarkGovaluesAdd`, `BenchmarkRhymondAdd`, `BenchmarkRejectOversizedTypedInputs` |
| `concurrency_test.go` | `money_test` | tests | `TestSharedValuesAndFormattersAreRaceSafe` |
| `context.go` | `money` | API and implementation | `MaxScale`, `MaxCashStep`, `ContextDefault`, `ContextCustom`, `ContextCash`, `ContextAutomatic`, `ContextKind`, `Context`, `DefaultContext`, `CustomContext`, `CashContext`, `AutomaticContext`, `Scale`, `Kind`, `CashStep`, `IsZero` |
| `context_test.go` | `money` | tests | `TestDefaultContextUsesAuthoritativeCurrencyMinorUnits`, `TestExplicitContextsValidateTheirBounds` |
| `contracts_internal_test.go` | `money` | tests | `TestMoneyZeroValuesAndOperationErrorsAreExplicit`, `TestAmountAndMoneyArithmeticEnforceOutputBounds`, `TestAllocationRejectsInvalidWorkAndKeepsBagAliasesPrivate`, `TestRateRationalAndConversionErrorContracts`, `TestTaxAndDiscountAccessorsAndInvalidInputs`, `TestCancellationPropagatesFromEveryBoundedArithmeticStage`, `TestAllocatePreservesCancellationDuringTypedRatioAdmission`, `TestCorruptInternalValuesFailClosedAtEveryNumericBoundary`, `TestInternalInvariantAssertionsFailClosed`, `Deadline`, `Done`, `Value`, `Err` |
| `conversion.go` | `money` | API and implementation | `MaxRateSourceBytes`, `ExchangeRate`, `ConversionResult`, `NewExchangeRate`, `Base`, `Quote`, `Exact`, `ObservedAt`, `Source`, `Source`, `Converted`, `Rate`, `Rounding`, `Convert` |
| `conversion_test.go` | `money` | tests | `TestConversionRequiresInjectedAttributedExactRate` |
| `differential_test.go` | `money` | tests | `TestFixedScaleArithmeticMatchesMatureImplementations`, `TestEqualSplitMatchesMatureImplementationsAcrossSigns`, `TestCurrencyMismatchMatchesMatureImplementations` |
| `discount.go` | `money` | API and implementation | `DiscountRate`, `DiscountResult`, `ParseDiscountRate`, `Rate`, `Original`, `Discount`, `Final`, `Rounding`, `ApplyDiscount` |
| `discount_test.go` | `money` | tests | `TestDiscountRoundsTheComponentAndConservesTheOriginal` |
| `doc.go` | `money` | API and implementation |  |
| `encoding/encoding.go` | `encoding` | API and implementation | `Version`, `MaxEncodedBytes`, `ErrInvalidEncoding`, `SQLMoney.Money`, `SQLMoney`, `MarshalJSON`, `UnmarshalJSON`, `MarshalText`, `UnmarshalText`, `Value`, `Scan`, `NumericValue`, `ScanNumeric` |
| `encoding/encoding_test.go` | `encoding_test` | tests | `TestVersionedJSONAndSQLRoundTripHistoricMoneyExactly`, `TestVersionedJSONRejectsDuplicateKeysAtEveryLevel` |
| `encoding/fuzz_test.go` | `encoding` | tests | `FuzzVersionedJSON`, `FuzzPostgreSQLNumeric` |
| `encoding/internal_test.go` | `encoding` | tests | `TestEveryContextEncodingAndPersistenceBoundary`, `TestEncodingInternalInvariantsAndDelimiterValidation`, `TestDecoderBoundsAndNestedTraversalAreExact`, `TestEncodingRejectsEveryInvalidShape`, `TestSQLMoneyScanRejectsOversizedInputsWithoutCopying`, `TestScanNumericRejectsOversizedInputsWithoutConverting` |
| `encoding/security_test.go` | `encoding_test` | tests | `TestSQLScanRejectsOversizedInputBeforeCopying`, `TestNumericScanRetainsInclusivePreconversionBoundary` |
| `errors.go` | `money` | API and implementation | `ErrUnknownCurrency`, `ErrMinorUnitsUnavailable`, `ErrInvalidContext`, `ErrPrecisionLoss`, `ErrCurrencyMismatch`, `ErrContextMismatch`, `ErrInvalidMoney`, `ErrAmountLimit`, `ErrInvalidRate`, `ErrInvalidAllocation`, `ErrMoneyBagLimit` |
| `examples_test.go` | `money_test` | tests | `ExampleMoney_Add`, `ExampleMoney_Allocate`, `ExampleConvert` |
| `format/format.go` | `format` | API and implementation | `MaxFormattedBytes`, `ErrInvalidFormat`, `Options.Symbol`, `Options`, `Exact`, `Locale` |
| `format/format_test.go` | `format_test` | tests | `TestLocaleFormattingKeepsExactValueAndSeparateIdentity`, `TestLocaleFormattingUsesTheLocalesGroupingPattern` |
| `format/fuzz_test.go` | `format_test` | tests | `FuzzLocale` |
| `format/internal_test.go` | `format` | tests | `TestFormattingBoundariesAndNumberShapes`, `TestFormattingFallbacksAndOutputBound`, `TestFormattingComparisonAndGroupingBoundaries` |
| `fuzz_test.go` | `money` | tests | `FuzzParseMoney`, `FuzzAllocationConservation`, `FuzzWeightedAllocationConservation`, `FuzzRate` |
| `invariant.go` | `money` | API and implementation |  |
| `minor_units.go` | `money` | API and implementation | `FromMinorUnits`, `MinorUnits` |
| `minor_units_test.go` | `money` | tests | `TestMinorUnitConstructionRoundTripsWithoutNarrowing` |
| `money.go` | `money` | API and implementation | `Money`, `Parse`, `Amount`, `Currency`, `Context`, `Valid`, `String`, `Add`, `Sub`, `Compare`, `Equal`, `Neg`, `Abs`, `Sign`, `IsZero` |
| `money_test.go` | `money` | tests | `TestParseProducesExactContextBoundMoney`, `TestArithmeticPreservesIdentityAndRejectsMismatches`, `TestComparisonAndUnaryOperationsAreExactAndImmutable` |
| `moneytest/internal_test.go` | `moneytest` | tests | `TestAssertionsReportEveryFailedLaw`, `Helper`, `Fatalf` |
| `moneytest/moneytest.go` | `moneytest` | API and implementation | `TestingT.Helper`, `TestingT.Fatalf`, `CurrencyFixture.Code`, `CurrencyFixture.Context`, `TestingT`, `CurrencyFixture`, `CurrencyFixtures`, `MustParse`, `AssertEqual`, `AssertAllocationConserved`, `AssertTaxConserved`, `AssertDiscountConserved` |
| `moneytest/moneytest_test.go` | `moneytest_test` | tests | `TestFixturesAndConservationAssertionsCoverCurrencyEdges` |
| `mutation_contract_test.go` | `money` | tests | `TestInclusivePublicBoundsRemainUsable`, `TestTypedNumericInputsAreBoundedBeforeFormatting`, `TestMoneyBagOrderingUsesEveryIdentityField`, `TestIdentityValidationRejectsEachMissingField`, `TestAllocationPreservesTieOrderAndZeroSign`, `TestTerminatingScaleAcceptsExactLimitsAndRejectsBeyondThem` |
| `properties_test.go` | `money` | tests | `TestEveryAuthoritativeCurrencyCanBeRepresentedExplicitly`, `TestAllocationConservationPropertyAcrossSignsAndCounts`, `TestTaxDiscountAndCashMatricesConserveForEveryRoundingMode`, `TestCashRoundingMatrixAcrossStepsSignsAndModes`, `TestResourceAndIdentityBoundariesRejectHostileInputs`, `TestParseRateWithMaximumSupportsCallerDefinedMagnitude` |
| `rate.go` | `money` | API and implementation | `MaxRateMagnitude`, `Rate`, `ParseRate`, `ParseRateWithMaximum`, `Rational`, `String`, `IsZero`, `Valid` |
| `rational_money.go` | `money` | API and implementation | `RationalMoney`, `Ratio`, `RoundingResult`, `Rational`, `String`, `Conditions`, `Inexact`, `String`, `Rational`, `Currency`, `Round`, `Mul`, `Quo`, `Ratio` |
| `rational_money_test.go` | `money` | tests | `TestMoneyRateArithmeticStaysRationalUntilExplicitRounding`, `TestCashRoundingUsesTheConfiguredIncrement`, `TestDivisionAndMoneyRatiosRemainExact`, `TestAutomaticRoundingOnlyAcceptsSafeTerminatingPrecision` |
| `security_budget_test.go` | `money_test` | tests | `TestExternalNumericValuesAreBoundedBeforeFormatting`, `TestExternalNumericValuesPreserveInclusiveMonetaryLimits`, `TestMinorUnitAdmissionPreservesMetadataErrorPrecedence`, `TestAmountAdmissionPreservesDecimalRepresentation`, `TestAmountAdmissionRetainsPrintedDigitAndScaleBoundaries`, `TestAutomaticContextResolvesOnlyAdmittedAmountScales`, `TestMinorUnitPowerBoundaryPreservesRejectionAndMetadata`, `TestMinorUnitsAcceptGreatestAdmittedCoefficient` |
| `tax.go` | `money` | API and implementation | `MaxTaxRate`, `TaxRate`, `TaxResult`, `ParseTaxRate`, `Rate`, `Net`, `Tax`, `Gross`, `Rounding`, `AddTax`, `ExtractTax` |
| `tax_test.go` | `money` | tests | `TestExclusiveAndInclusiveTaxPreserveTheirDocumentedTotals` |

## Initial gaps and disposition

- Canonical package comments already exist, but need entry points, ownership
  exclusions, infrastructure and lifecycle boundaries. There is one comment
  per package; external test package clauses are not additional public packages.
- Missing exported comments: encoding.ErrInvalidEncoding, format.ErrInvalidFormat,
  SQLMoney.Money, CurrencyFixture.Code/Context and TestingT.Helper/Fatalf.
  Other exported declaration comments exist; audit their semantics rather than
  adding boilerplate.
- Cash construction enforces represented scale, not divisibility by cash step;
  only RationalMoney.Round applies the cash increment. Allocation uses minor
  units rather than cash-step units and rejects automatic contexts.
- AllocationResult's earlier-part remainder description is too broad for weighted
  allocation: largest fractional remainder wins, with original order for ties.
- Zero values, nil context admission, output-limit errors, rate maxima, automatic
  precision, scanner mutation/NULL rejection and locale symbol fallback need
  caller-visible clarification.
- TaxResult.Tax is misleading for ExtractTax: net is rounded and tax is the
  conserved residual, not a separately rounded component.
- Add adjacent rationale for typed-ratio admission before formatting, absolute
  weighted allocation and sign restoration, terminating-decimal factorization,
  cash-step rounding and bounded recursive JSON preflight. Keep existing typed
  decimal/minor-unit admission rationale.
- Fuzz targets and benchmark workloads need concise boundary/setup statements.
  Existing example outputs and regression/property assertions are self-contained;
  complex staged-cancellation and fatal-test doubles need a short explanation.

## Investigated rationale and constraints

Code plus contracts_internal_test.go, security_budget_test.go,
mutation_contract_test.go and encoding/security_test.go establish inclusive
limits, metadata-error precedence, bounded input conversion and slice ownership.
allocation_test.go/properties_test.go protect conservation, sign and stable ties;
rational_money_test.go/properties_test.go protect rounding mode and cash policy;
encoding tests protect the closed version-1 schema and transactional receiver
replacement. format tests protect exact digits, grouping and CLDR fallbacks.
No performance promise is inferred from the comparative benchmarks. Currency
identity and ISO metadata belong to the pinned International dependency; no
new specification-compliance or monetary business-policy claim is introduced.

Cancellation is checked at bounded arithmetic stages, not forced preemption;
Round has no operation context. AddTax returns a context mismatch when
automatic result scale differs; ApplyDiscount and ExtractTax instead reach the
internal-invariant panic during residual subtraction. Their source comments
warn consumers to use fixed contexts. This pre-existing runtime hazard is
recorded, not repaired in a prose-only pass. A concrete source trace is
`Parse("1.00", EUR, AutomaticContext())` with `ParseDiscountRate("0")`:
ApplyDiscount rounds zero to automatic scale 0, then subtracts it from automatic
scale 2 and panics. The same original amount with `ParseTaxRate("0")` reaches
ExtractTax's scale-0 net versus scale-2 gross subtraction. Fixed contexts keep
these scales aligned. This trace follows Round -> terminatingScale -> Parse ->
Money.compatible -> mustInvariant; it is not a newly executed runtime test.

## Unresolved questions

No unknown rationale was filled with a guessed guarantee. A follow-up runtime
fix is needed if automatic contexts are to be safely supported by ApplyDiscount
and ExtractTax when result scale differs; current behavior is documented without
a runtime change. Existing resource and cancellation tradeoffs are maintained in [the security model](security.md).
This audit does not claim new runtime, race, fuzz, benchmark, conformance,
public-consumer or release qualification.

## Completion evidence

The four canonical package comments now describe entry points, owner boundaries
and lifecycle; exported sentinels, public fields and interface obligations have
comments. Existing good accessor and regression-test comments were retained.
All 45 files were inspected, including 7 fuzz targets, 4 benchmark functions,
3 examples and test doubles. No generated or platform-specific source was skipped.

Selected local evidence:

- `gofmt -l` on all 23 changed Go files: no output (formatted).
- `git diff --check`: passed.
- `make -f verification/package.mk no-float`: passed.
- Comment-aware lexical comparison of baseline/current Go files: executable
  content and all string literals unchanged in all 23 changed files.
- Source marker/generator/build-variant scan: none.
- Changed audit-local link validation: passed; changed-content private-key,
  cloud/access-token and chat-token indicators: absent.
- Ordinary complete-final-diff inspection: no documentation finding remaining;
  the automatic-context runtime hazard above is explicitly retained as follow-up.

- Go 1.27.2 AST source-comment check: all 154 exported declarations/members
  have correctly prefixed comments; exactly one canonical comment exists for
  each of the four packages.
- `make -f verification/package.mk docs`: passed (native required-document,
  package listing, rendered Go documentation and unfinished-marker checks).
- `go doc .`, `go doc ./encoding`, `go doc ./format`, `go doc ./moneytest`,
  `go doc ApplyDiscount` and `go doc ExtractTax`: passed; rendered output was
  inspected for useful entry points, boundaries and the automatic-context caveat.
- Go commands used `.go-version`'s Go 1.27.2 with bounded private build/module/
  temporary caches; caches and temporary AST checker were removed after capture.

Runtime tests, race/fuzz/mutation/benchmark
campaigns and public-consumer/release qualification are not selected for this
comment-only batch. Shared CI tooling enforcement and remote-main delivery are
coordinator-owned; this record does not certify those later boundaries.

## Objective enforcement adoption

The local `source-docs-check` Make target runs the released source-only strict
checker and is a prerequisite of ordinary/full checks and CI. Shared CI selects
`source_documentation: check` using the immutable Tools v2.4.0 workflow and
matching authenticated binary/checksum declaration. This verifies canonical
package/exported/member comments, marker/generator availability, supported local
links and narrow secret indicators; it does not certify the semantic truth of
contracts or replace the maintained audit and native example/doc checks.
