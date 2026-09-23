# ADR 002: Per-question partial fallback on low confidence

## Status

Accepted (2026-09)

## Context

[routing-model.md](../routing-model.md) lists an open question: when Jev confidence is low, should the router fall back entirely or only for failing questions?

## Decision

v0.1.0 uses **per-question partial fallback**:

- Each asked question is gated independently (`choice`, `noul` thresholds in `ConfidenceConfig`).
- Failed questions are filled from the question-level fallback table; passing questions keep Jev-derived destinations.
- `PlanMetadata.FallbackMode` is `"partial"` when any question failed the gate, `""` when all passed, and `"full"` only when the Jev API call fails.

## Consequences

- Hosts may receive a mix of Jev and deterministic destinations in one plan; audit fields (`fallback_questions`, `answers_snapshot`) must be logged.
- Tuning confidence is per question, not a single global switch.
