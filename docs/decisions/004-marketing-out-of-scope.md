# ADR 004: Marketing routes are out of v0.1.0 scope

## Status

Accepted (2026-09)

## Context

Marketing notifications involve opt-in regulation and often do not need probabilistic routing. [routing-model.md](../routing-model.md) asked whether marketing should use Jev.

## Decision

**Marketing is out of scope for v0.1.0**:

- `CategoryMarketing` is not part of the MVP contract (see ADR 001).
- Hosts should route marketing with static user prefs and org policy only.
- The fallback table documents marketing behavior for future versions but is not exercised by MVP tests.

A dedicated marketing router may be added in a later release with its own QuestionSet and compliance review.

## Consequences

- OSS consumers must not expect `Resolve` for marketing templates in v0.1.0.
- Reduces regulatory and model-risk surface for the first public release.
