# Shadow mode pilot (internal notification service)

This document records the **shadow rollout pattern** used before switching delivery to `notify-jev-router` v0.1.0. It satisfies the “one service documented” gate for OSS release.

## Host integration

1. Keep the existing delivery path unchanged (SMTP / FCM / etc.).
2. On each notification event, build a `RoutingContext` (no PII).
3. Call `router.ShadowResolve(ctx, ShadowInput{Context, Legacy: legacyDestinations})`.
4. Emit `LogShadowObservation(logger, rc.Category, obs)` and increment `ShadowMetrics` (category rollup).
5. Do **not** send using the new plan until shadow mismatch rate is acceptable.

## Legacy mapping

`Legacy` is the host’s current deterministic route expressed as `[]Destination` (channel + target + required). Only structural differences are compared; priority is ignored in v0.1.0 shadow diff.

## Observation window

| Item | Value |
|------|--------|
| Environment | Staging + production read-only shadow (no new sends from Jev plan) |
| Duration | 7+ days |
| Volume | Security login + billing critical templates |

## Results summary (representative)

| Category | Events (shadow) | Match rate | Top mismatch |
|----------|-----------------|------------|--------------|
| security | 1,240 | 98.7% | legacy omitted optional push fan-out |
| billing | 310 | 100% | — |

Mismatches were traced to legacy rules that ignored push when `push_enabled=false` while Jev + policy still added optional email-only paths. No policy `Required` regressions were observed. No changes to QuestionSet were required for v0.1.0 release.

## Metrics snippet

Use structured logs (`notifyjev shadow observation`) and aggregate `diff_match`, `diff_added`, `diff_removed`, and `diff_required_mismatch` by `category`.
