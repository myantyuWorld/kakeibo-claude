---
paths:
  - "backend/internal/**/*.go"
---

# ドメインモデル実装ガイドライン

## 参照すべきドキュメント

1. `docs/design/{ドメイン名}/index.md` — ドメインの概要とユースケース
2. `docs/design/{ドメイン名}/ドメインモデル.md` — 集約・エンティティ・値オブジェクトの構造、対訳表
3. `docs/design/{ドメイン名}/業務ルール.md` — ライフサイクル、バリデーション、ビジネスルール

## 実装先

- `internal/{ドメイン名}/{ドメイン名}.go` に以下を実装する
  - エンティティ（unexported struct）
  - 値オブジェクト（unexported struct、値レシーバ）
  - ステータス列挙（unexported type + const）
  - コンストラクタ関数（バリデーション込み）
  - ドメインロジック（エンティティのメソッド）

## 生成しないもの

- handler / service / repository
- テストコード
- マイグレーション

## 命名

- ALWAYS 対訳表（ユビキタス言語）の英語名をそのままコード上の命名に使う
- `.claude/rules/backend/go-coding.md` の命名規約に従う
- パッケージ名と型名の重複を避ける（`order.Status` であって `order.OrderStatus` ではない）

## エラー

- `.claude/rules/backend/error-handling.md` に従い、エラーは `pkg/errors` を使う

## 関連

- **go-reviewer** agent -- ドメインモデルの変更レビュー時
- skill: `impl-api` -- ドメインモデル実装時
