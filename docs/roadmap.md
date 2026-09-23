# ロードマップ

## フェーズ 0 — 全体設定

- [x] 背景・目的・ゴールの文書化
- [x] アーキテクチャと Routing モデルの草案
- [x] オープン設計判断 3 件の決定（[ADR 002–004](decisions/002-partial-fallback-on-low-confidence.md)）
- [x] QuestionSet v0.1.0（Go embed、`questions/v0_1_0.go`）

## フェーズ 1 — コアライブラリ（MVP / v0.1.0）

- [x] `go.mod` / モジュール `github.com/kyrenzk/notify-jev-router`
- [x] `RoutingContext`, `RoutingPlan`, validate（PII 拒否）
- [x] `PolicyEngine`（security recovery must 等）
- [x] `JevClient` インターフェース + jev-go アダプタ
- [x] `Interpreter` + 部分フォールバック
- [x] ユニットテスト（録画 JSON fixture）

**MVP の完了条件（v0.1.0）**: 次の 3 シナリオがテストで green。

| # | シナリオ |
|---|----------|
| 1 | security.new_login → primary + recovery (required) + push all_devices |
| 2 | Jev API 失敗 → billing critical フォールバック（`FallbackMode=full`） |
| 3 | Policy must — Jev が recovery を付けなくても recovery Required が残る |

## フェーズ 2 — 運用品質

- [x] Shadow mode（`ShadowResolve`, metrics, [shadow-pilot.md](shadow-pilot.md)）
- [x] 例: `examples/mockdemo`（キー不要）
- [x] Optional live Jev integration（`integration` build tag + workflow）
- [ ] QuestionSet embed の semver 運用自動化
- [ ] `slog` / OTel の拡張サンプル
- [ ] ベンチ（Resolve p99）

## フェーズ 3 — OSS 公開

- [x] LICENSE（Apache-2.0）
- [x] CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md
- [x] README 英語 + Production readiness、日本語要約
- [ ] Private → Public リポジトリ切替（merge 後 maintainer — [RELEASE_CHECKLIST.md](../.github/RELEASE_CHECKLIST.md)）
- [ ] GitHub Release `v0.1.0`（tag 付与後）

## フェーズ 4 — エコシステム（任意）

- [ ] gRPC sidecar
- [ ] QuestionSet ジェネレータ CLI
- [ ] Langfuse 等トレース連携（jev SDK tracer）

---

## マイルストーン目安

| マイルストーン | 内容 |
|----------------|------|
| **M0** | ドキュメント合意 |
| **M1** | MVP テスト 3 本 |
| **M2** | 社内 1 サービスで shadow mode（[shadow-pilot.md](shadow-pilot.md)） |
| **M3** | OSS 公開 + Release v0.1.0 |

Shadow mode: 既存ルータと Plan を比較し、差分メトリクスを取ってから切替することを推奨します。
