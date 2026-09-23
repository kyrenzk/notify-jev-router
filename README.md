# notify-jev-router

[TypeSafe AI Jev](https://www.jevtypesafeai.com/) を使い、**通知イベントの宛先・チャネル・配信方針**を型安全に解決する Go ライブラリ（予定）です。

マイクロサービス内に組み込み、テンプレート ID・重要度・ユーザー設定の要約など **個人情報を除いたコンテキスト** から、「どのチャネルに送るか」「再設定用メールにも送るか」「Push は全デバイスか」などを判断します。

## ステータス

| 項目 | 状態 |
|------|------|
| リポジトリ | Private（動作・設計が固まり次第 OSS 公開予定） |
| 実装 | v0.1.0 MVP（ルーティングコア） |
| 言語 | Go 1.22+ |

## ドキュメント

| ドキュメント | 内容 |
|--------------|------|
| [全体像（背景・目的・ゴール）](docs/overview.md) | なぜ作るか、何を達成するか、スコープ |
| [アーキテクチャ](docs/architecture.md) | 境界、データフロー、ホストサービスとの関係 |
| [ルーティング判断モデル](docs/routing-model.md) | Jev への入力・質問・出力（RoutingPlan）の考え方 |
| [制約と安全](docs/constraints.md) | PII 禁止、信頼度、フォールバック、監査 |
| [ロードマップ](docs/roadmap.md) | フェーズと公開条件 |

## クイックスタート（ホスト疑似コード）

```go
package main

import (
	"context"
	"log"
	"os"

	jev "github.com/havlan/jev-go"
	notifyjev "github.com/kyrenzk/notify-jev-router"
)

func main() {
	jc, err := jev.New(jev.Config{
		APIKey: os.Getenv("TYPESAFE_API_KEY"),
		Model:  os.Getenv("JEV_MODEL"), // 本番は jev-x.y.z にピン留め
	})
	if err != nil {
		log.Fatal(err)
	}

	router, err := notifyjev.New(notifyjev.Config{
		JevClient: notifyjev.NewJevGoAdapter(jc, os.Getenv("JEV_MODEL")),
		Model:     os.Getenv("JEV_MODEL"),
	})
	if err != nil {
		log.Fatal(err)
	}

	plan, err := router.Resolve(context.Background(), notifyjev.ResolveInput{
		Context: notifyjev.RoutingContext{
			EventID:    "evt-123",
			TemplateID: "security.new_login",
			Category:   notifyjev.CategorySecurity,
			Severity:   notifyjev.SeverityCritical,
			Locale:     "ja-JP",
			UserPrefs: notifyjev.UserPrefs{PushEnabled: true},
			Capabilities: notifyjev.Capabilities{
				EmailVerified:    true,
				PushDeviceCount:  2,
				HasRecoveryEmail: true,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	_ = plan // → ホストがチャネルワーカーへ配送
}
```

`go test ./...` で fixture ベースのテストが実行できます（live Jev API 不要）。

## 位置づけ（一行）

**通知の「中身」を作るライブラリではなく、送る先と送り方を決めるルーティングエンジン** — 実際の配送は各サービス既存のメール / Push / SMS クライアントが担当します。

## ライセンス

OSS 公開時に決定（Apache-2.0 または MIT を想定）。公開前は All rights reserved。
