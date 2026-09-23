# notify-jev-router

Go library that resolves **notification channels, targets, and fan-out** using [TypeSafe AI Jev](https://www.jevtypesafeai.com/) with typed policies, confidence gates, and deterministic fallbacks.

Hosts pass **non-PII** `RoutingContext` (template, category, severity, prefs, capability flags). The library returns a `RoutingPlan` for your existing mail / push / SMS workers.

## Status

| Item | State |
|------|--------|
| Repository | Public |
| License | Apache-2.0 |
| Version | [v0.1.0](https://github.com/kyrenzk/notify-jev-router/releases/tag/v0.1.0) |
| Go | 1.22+ |

## Quick start (no API key)

Official mock demo using recorded Jev JSON:

```bash
go run ./examples/mockdemo
```

This prints an indented `RoutingPlan` JSON for `security.new_login`.

## Quick start (live Jev)

```go
router, err := notifyjev.New(notifyjev.Config{
	JevClient: notifyjev.NewJevGoAdapter(jc, os.Getenv("JEV_MODEL")),
	Model:     os.Getenv("JEV_MODEL"), // pin jev-x.y.z in production
})
plan, err := router.Resolve(ctx, notifyjev.ResolveInput{Context: rc})
```

Run unit tests (fixtures only):

```bash
go test ./...
```

## Production readiness

- **Shadow first**: Call `ShadowResolve` and log with `LogShadowObservation` while keeping legacy delivery. See [docs/shadow-pilot.md](docs/shadow-pilot.md).
- **Pin the model**: Use an explicit `jev-x.y.z` model in production. Reserve `jev-latest` for staging only.
- **Delivery is yours**: This library does not send email or push; it only returns a plan.
- **PII**: Do not put contact fields or message bodies in context; validation rejects common mistakes. See [docs/constraints.md](docs/constraints.md) and [SECURITY.md](SECURITY.md).

## Documentation

| Doc | Description |
|-----|-------------|
| [Overview](docs/overview.md) | Goals, scope, OSS checklist |
| [Architecture](docs/architecture.md) | Boundaries and data flow |
| [Routing model](docs/routing-model.md) | Questions, confidence, fallbacks |
| [Constraints](docs/constraints.md) | PII, pinning, availability |
| [Roadmap](docs/roadmap.md) | Phases and milestones |
| [ADRs](docs/decisions/001-v0.1.0-contract.md) | Design decisions 001–004 |
| [Shadow pilot](docs/shadow-pilot.md) | Internal shadow rollout summary |
| [Jev fixtures](docs/jev-fixtures.md) | Refreshing `testdata/jev/` |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). CI runs `go test ./...` on every PR.

---

## 日本語

**notify-jev-router** は、個人情報を含まない通知メタデータから、送信チャネル・宛先種別・fan-out を型安全に決める Go ライブラリです。

- **キー不要の試行**: リポジトリルートで `go run ./examples/mockdemo`
- **本番**: Jev モデルを `jev-x.y.z` にピン留めし、先に [shadow モード](docs/shadow-pilot.md) で既存ルールと Plan を比較してください
- **配送**: メール / Push の実送信はホストサービスが担当します

詳細は上記英語ドキュメントリンク（[全体像](docs/overview.md)、[制約](docs/constraints.md)）を参照してください。

## License

Apache-2.0 — see [LICENSE](LICENSE).
