# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's docs.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and ordinary CI; the public minimum remains
Go 1.27.0. An ordinary standalone checkout is sufficient for tests and a
focused active fuzz smoke, without sibling repositories or goal ledgers:

```bash
go test ./...
go test . -run '^$' -fuzz '^FuzzAllocationConservation$' -fuzztime=2s
```

Select the native fuzz target and package affected by your change. Ordinary
tests run deterministic seeds and saved corpora. For shared checks, install
the public `golib` v2.5.1 release declared in `.golib.yaml` and verify its
checksums:

```bash
make check
make ci
```

The local `make source-docs-check` target checks all tracked Go documentation
without compilation. `make check`, `make check-full` and `make ci` require it;
shared CI runs the same strict check before its required result can succeed.

Both ordinary Make routes retain the cheap no-float source check and bounded
local assurance. Ordinary CI runs the same script in a checkout-only, two-minute job
and requires both that job and shared CI to succeed. Base-free Make checks run
deterministic seeds/corpora; changed parser inputs receive causal active smoke
in CI. The sole root module declares no services or test tags; PostgreSQL
numeric adapter tests do not require a database service.

Optional mutation analysis can select one affected package:

```bash
golib mutation --module . --package ./encoding
```

The no-binary-float policy check is also independently reproducible:

```bash
make -f verification/package.mk no-float
```

Use `make check-full` for explicit full assurance with the unchanged
10,000-iteration native fuzz campaigns, no-float policy, exact/conservation
and maintained-peer comparisons, API baseline, docs and comparative benchmarks.
Ordinary checks are not full coverage/mutation/race/benchmark qualification.
Select directly affected public consumers and relevant release qualification,
including `golib release check` and `golib release dry-run --all`, when their
contract changes. Full campaigns remain available for selected scheduled or
release risks; ordinary scheduled CI stays bounded. Report unavailable or
failing checks without claiming release readiness. CI/configuration-only
changes do not require a library release.

## Selected Monthly Assurance

[`assurance.yml`](.github/workflows/assurance.yml) runs on main at 07:11 UTC on
the 28th day of each month. It selects one of the 7 declared root-module fuzz
targets using the shared UTC-week selector; it does not run every target or
promise a complete monthly rotation. Explicit main-only dispatch accepts one
exact declared selector, for example `^FuzzAllocationConservation$`.

The selected risks are parse/encoding round trips, allocation conservation,
rate acceptance and bounded locale formatting. The PostgreSQL numeric target
checks in-process conversion, not a live database. This is not complete exact
arithmetic, no-float, independent-peer or public-consumer qualification; retain
those affected checks and the separate test-support package contract.

The immutable shared runner bounds active fuzzing to 30 seconds, two workers
and two CPUs, with a five-minute combined native session, 2 GiB memory and
1 MiB output limits inside a ten-minute job. It installs the declared public
CLI without a source build, uses private disposable caches and keeps failures,
timeouts and exhausted budgets visible. Ordinary deterministic tests still run
seeds/corpus; ordinary CI and its Required contract are unchanged. Full/manual
operations and their original budgets remain in `.golib.yaml`.

Treat findings as owner-attributed diagnostics: promptly repair viable defects
with focused regressions; retain reviewed narrow equivalent/invalid dispositions
rather than weakening tests. Review cadence against useful findings and actual
cost; do not interpret a selected pass as full or release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
