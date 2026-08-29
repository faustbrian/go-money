# Changelog

All notable changes follow Keep a Changelog conventions. Compatibility follows
semantic versioning after the first tagged release.

## Unreleased

### Changed

- Replace copied repository-local verification tooling with the released
  `go-library-tools` v1.0.5 contract while preserving package-owned evidence,
  API baselines, and source-specific verification operations.

### Documentation

- Replace archived monorepo links and completed execution artifacts with a
  standalone, human-oriented documentation structure.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Fixed

- Fixed-context parsing now returns a bounded arithmetic error instead of
  panicking when scale alignment would exceed the maximum amount digits.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-money` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Require exact 100% mutant coverage in the package mutation command while
  retaining the package's explicit behavioral mutants.
- Reuse the repository-provided API compatibility tool during isolated checks
  so mutable local dependency archives never conflict with source checksums.
- Pin unpublished owned modules to exact resolvable revisions so clean
  consumers no longer depend on nonexistent `v0.1.0` tags.
- Refresh owned-module checksums against the final consolidated archives.
- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.
- Replaced coverage-only defensive exclusions with tested internal invariant
  assertions, while retaining returned errors for every caller-reachable
  invalid input and arithmetic failure.
- Validate the complete pinned currency dataset against the supported money
  scale instead of carrying an unreachable per-call scale branch.
- Make encoding invariants and CLDR fallback parsing directly testable while
  preserving strict duplicate-key, delimiter, and formatted-output bounds.
- Removed the annotation-based defensive-statement coverage exclusion and
  require exact deduplicated coverage for every production package.

### Added

- Caller-defined exact maximum magnitudes for rates while preserving the
  generic one-million default.
- Immutable exact `Money`, `Amount`, `RationalMoney`, `MoneyBag`, rate, ratio,
  context, allocation, tax, discount, conversion, and result values.
- ISO default, custom, cash, and safe automatic precision contexts.
- Deterministic equal and weighted allocation with signed conservation.
- Explicit tax, discount, cash-rounding, and attributed FX boundaries.
- Locale formatting, versioned JSON/text/SQL persistence, PostgreSQL numeric
  support, and reusable `moneytest` laws.
- Property, exhaustive currency, fuzz, race, mutation, coverage, benchmark, and
  compatibility gates.
