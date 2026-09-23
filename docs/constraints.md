# 制約と安全

## PII と Jev state

### 原則

- **RoutingContext に PII を入れない** — 型定義と validate で機械的に防ぐ
- ホストが誤って渡した場合: **Resolve 前にエラー**（サイレント strip はデバッグを難しくするため v0 は拒否推奨）
- ログにも **state 全文を本番で DEBUG 出力しない** — サンプリング + フィールド allowlist

### 許可フィールド（allowlist 思考）

| 種別 | 許可例 |
|------|--------|
| 識別子 | `template_id`, `event_id`（ユーザー ID は **ハッシュまたは内部 UUID のみ**、メール等と結びつく外部 ID は避ける） |
| 列挙 | `category`, `severity`, `locale` |
| 集計 | デバイス数、bool フラグ、チャネル ON/OFF |
| 禁止 | 連絡先、トークン、本文、IP、位置情報、決済カード下 4 桁以外の情報 |

`user_id` を state に含めるかは組織次第。**Jev は第三者 API** であることを踏まえ、可能なら `subject_ref` のような opaque ID に限定。

---

## 信頼度とモデルバージョン

| 項目 | 推奨 |
|------|------|
| 本番モデル | `jev-x.y.z` に **ピン留め** |
| 実験 | `jev-latest` は staging のみ |
| 閾値変更 | QuestionSet の semver を bump |

Jev の `jev-latest` は挙動が変わり得るため、confidence 閾値を本番で依存させる場合は **必ずピン**（[Jev ドキュメント](https://www.jevtypesafeai.com/how-to-use) 参照）。

---

## 可用性

- Jev 障害時は **フォールバック Plan** で送信継続（安全系は過剰配送を許容、マーケは under-delivery を許容、などカテゴリ方針）
- タイムアウトはホスト SLA に合わせ **短め**（例: 800ms）+ コンテキストキャンセル
- 同一イベントの **重複 Resolve** は idempotent に（同 `event_id` + 同 QuestionSet → 同 Plan を期待。キャッシュはホスト任意）

---

## セキュリティ

| リスク | 対策 |
|--------|------|
| API キー漏洩 | 環境変数 / シークレットマネージャ。リポジトリにコミットしない |
| プロンプトインジェクション | state は構造化 JSON のみ。自由テキストフィールドを持たない |
| 過剰配送 | PolicyEngine + カテゴリ別 ceiling |
| 配送不足（アカウント侵害） | security フォールバックは over-delivery 寄り |

---

## コンプライアンス（マーケティング）

- マーケ系は **オプトイン bool** を Capabilities で明示し、Jev は補助に留める案を推奨
- 地域別規制（TCPA, GDPR 等）は **ホスト Legal** が最終判断。本ライブラリは「ルート候補」を返すだけ

---

## OSS 公開前のセキュリティチェック

- [ ] `gosec` または同等の静的解析
- [ ] 依存 SDK のライセンス確認（jev-go 等）
- [ ] SECURITY.md（連絡先、対応方針）
- [ ] サンプルにダミーキー・ダミー state のみ
