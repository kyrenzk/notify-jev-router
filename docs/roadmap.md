# ロードマップ

## フェーズ 0 — 全体設定（現在）

- [x] 背景・目的・ゴールの文書化
- [x] アーキテクチャと Routing モデルの草案
- [ ] オープン設計判断 3 件の決定（[routing-model.md](routing-model.md) 末尾）
- [ ] QuestionSet v0.1.0 の YAML 凍結

## フェーズ 1 — コアライブラリ（MVP）

- [ ] `go.mod` / モジュール `github.com/kyrenzk/notify-jev-router`
- [ ] `RoutingContext`, `RoutingPlan`, validate（PII 拒否）
- [ ] `PolicyEngine`（security recovery must 等）
- [ ] `JevClient` インターフェース + jev-go アダプタ
- [ ] `Interpreter` + 部分フォールバック
- [ ] ユニットテスト（録画 JSON fixture）

**MVP の完了条件**: 表の 3 シナリオがテストで green

| # | シナリオ |
|---|----------|
| 1 | security.new_login → primary + recovery + push all |
| 2 | marketing + prefs OFF → skip または empty plan |
| 3 | Jev エラー → billing high フォールバック |

## フェーズ 2 — 運用品質

- [ ] QuestionSet embed + semver
- [ ] `slog` 構造化ログ / optional OTel
- [ ] ベンチ（Resolve p99）
- [ ] 例: `examples/minimal` ホスト疑似コード

## フェーズ 3 — OSS 公開

- [ ] LICENSE（Apache-2.0 or MIT）
- [ ] CONTRIBUTING.md, SECURITY.md, CODE_OF_CONDUCT.md
- [ ] README 英語版充実、godoc
- [ ] Private → Public リポジトリ切替
- [ ] 初回 tag `v0.1.0`

## フェーズ 4 — エコシステム（任意）

- [ ] gRPC sidecar
- [ ] QuestionSet ジェネレータ CLI
- [ ] Langfuse 等トレース連携（jev SDK tracer）

---

## マイルストーン目安

| マイルストーン | 内容 |
|----------------|------|
| **M0** | ドキュメント合意（今） |
| **M1** | MVP テスト 3 本 |
| **M2** | 社内 1 サービスで shadow mode（Plan のみログ、配送は既存ルール） |
| **M3** | 本番カナリア + OSS 公開 |

Shadow mode: 既存ルータと Plan を比較し、差分メトリクスを取ってから切替することを推奨します。
