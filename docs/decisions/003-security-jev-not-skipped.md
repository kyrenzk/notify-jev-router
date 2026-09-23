# ADR 003: Security category still calls Jev

## Status

Accepted (2026-09)

## Context

Security notifications could skip Jev and rely on `PolicyEngine` alone for lower latency and cost. The open item in [routing-model.md](../routing-model.md) asked whether to skip Jev for security.

## Decision

v0.1.0 **does not skip Jev** for `category=security` when channels are deliverable:

- `PolicyEngine` still seeds required destinations (e.g. recovery email).
- Jev runs for primary strategy, push fan-out, and recovery copy confidence.
- Merge keeps policy `Required` flags; Jev cannot clear them.

Hosts may add their own fast path outside this library; the default router always calls Jev when deliverable.

## Consequences

- Security routes pay Jev latency; shadow mode is recommended before cutover.
- Explainability improves via `answers_snapshot` on security plans.
