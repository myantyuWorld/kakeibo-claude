---
paths:
  - "backend/**/*.go"
  - "backend/internal/**"
---

# アーキテクチャ

詳細は `docs/adr/backend/0001-architecture.md` を参照すること。

## ディレクトリ構成

```
internal/
├── _template/         # 実装テンプレート（ビルド対象外）
└── {aggregate}/       # 集約単位でパッケージを配置
    ├── {domain}.go    # ドメインモデル・値オブジェクト
    ├── handler.go     # HTTPハンドラ & Service インターフェース
    ├── service.go     # アプリケーションサービス & Repository インターフェース
    ├── repository.go  # リポジトリ実装
    ├── module.go      # DI・ルーティング登録
    ├── resolver.go    # 他ドメインへの読み取り公開（必要な場合のみ）
    └── creator.go     # 他ドメインへの書き込み公開（必要な場合のみ）
```

## レイヤー間の依存方向

```
handler → service → repository
```

ドメインモデル（{domain}.go）はどのレイヤーにも依存しない。

## パッケージ間の依存方向

コンテキストマップに従う。NEVER コンテキストマップに明記されていない依存は作成しない。

## 実装テンプレート

新規ドメイン追加時は `internal/_template/` を参照すること。各ファイルのコメントに設計判断を記載している。

## 関連

- **go-reviewer** agent -- コードレビュー時
- skill: `impl-api` -- API実装時
