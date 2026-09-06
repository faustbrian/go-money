# Security policy

## Supported versions

The stable v1 release line is supported. Users should run the latest published
v1 patch because security fixes are not promised for superseded patch releases.
The `main` branch is unreleased development and is not a supported deployment
target.

## Reporting

Report suspected vulnerabilities through a private
[GitHub security advisory](https://github.com/faustbrian/go-money/security/advisories/new).
Include a minimal reproducer, affected version, expected invariant, and
observed impact. Do not attach customer records, account identifiers, or
production monetary payloads.

## Threat model

Untrusted inputs may attempt excessive digits, scales, ratios, allocation
counts, JSON nesting, locale expansion, or diagnostic amplification. Public
parsers and adapters apply fixed bounds before or during expensive work.

The package performs no network access, loads no ambient currency rates, uses
no unsafe code, and logs no monetary source records. Callers remain responsible
for database authorization, transport authentication, rate provenance, and
business-level amount limits narrower than the package maximums.

See [docs/security.md](docs/security.md) for the full boundary inventory.
