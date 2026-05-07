# プロジェクト名

プロジェクト概要を記載する

## 技術スタック

### バックエンド

- Go 1.26.2 / Echo v4
- PostgreSQL / sqlx
- マイグレーション: sql-migrate
- テスト: testify（テーブル駆動テスト + mock）

### フロントエンド

- Vue 3 / TypeScript / Vite
- vue-router（レイアウト: meta 方式、lazy loading）
- openapi-fetch + openapi-typescript（API クライアント + 型生成）
- vee-validate + zod（フォームバリデーション）
- テスト: Vitest + happy-dom + MSW
- Storybook 10（コンポーネントカタログ + インタラクションテスト）
- アーキテクチャ: FSD（app → pages → features → shared）

## コマンド

直接 `go test` / `pnpm run` 等を実行せず、必ず `mise run` 経由で実行すること。

### 共通（ルートで実行）

| コマンド | 用途 |
|---------|------|
| `mise run setup` | backend/.env と frontend/.env を作成 |
| `mise run up` | Docker コンテナ起動 |
| `mise run down` | Docker コンテナ停止 |

### バックエンド（`backend/` で実行）

| コマンド | 用途 |
|---------|------|
| `mise run dev` | Air によるホットリロード開発サーバー起動 |
| `mise run lint` | golangci-lint 実行 |
| `mise run test` | 単体テスト実行（`./internal/...`） |
| `mise run test-integration` | 統合テスト実行（`./tests/...`） |
| `mise run test-all` | 全テスト実行 |

### フロントエンド（`frontend/` で実行）

| コマンド | 用途 |
|---------|------|
| `mise run dev` | Vite 開発サーバー起動 |
| `mise run build` | プロダクションビルド |
| `mise run lint` | ESLint 実行 |
| `mise run format` | Prettier フォーマット |
| `mise run test` | Vitest ウォッチモード |
| `mise run test-run` | Vitest 単発実行 |
| `mise run generate-api` | OpenAPI 仕様書から型定義を生成 |
| `mise run storybook` | Storybook 開発サーバー起動（port 6006） |

## Git ルール

`.claude/rules/common/git-workflow.md` を参照すること。

## ドキュメント

実装時は以下のドキュメントを参照すること:

- [設計ドキュメント](../docs/design/README.md) - ドメインごとの設計ドキュメント入口
- [フロントエンド ADR](../docs/adr/frontend/0001-architecture.md) - FSD アーキテクチャの決定記録
- [API 仕様書](../docs/api/openapi.yml) - OpenAPI 定義
