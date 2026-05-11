# 家計簿アプリ (kakeibo)

複数人（夫婦・パートナーなど）で共有する家計簿アプリ。スプレッドシート運用をシステム化し、入力の手間削減・年単位の比較・視覚化を実現する。

## 配信形態

**LIFF ミニアプリ（LINE 内 Web アプリ）** として提供する。詳細は [要件定義書](../docs/要件定義書.md) §5.5 を参照。

- 認証: LINE Login API
- 通知: LINE Messaging API（公式アカウント）経由の Push
- 外部ブラウザ単独での利用は v1 では想定しない

## 主要ドメイン（暫定）

実装着手前に sudo モデリング（`/review-design`）でドメインを確定させる。現時点での想定:

| ドメイン | 概要 |
|---------|------|
| ユーザー (User) | アプリ利用者。LINE ログインで認証。プロフィール（ニックネーム・サムネ）を持つ |
| 家計簿 (Household) | 収支を共有管理する単位。1 ユーザーは 1 家計簿のみに所属 |
| メンバーシップ (Membership) | ユーザーと家計簿の所属関係。ロール（管理者 / メンバー）を持つ。管理者は 1 家計簿に 1 人（作成者） |
| カテゴリ (Category) | 収支の分類。収入／支出の 2 系統。カテゴリ自体は家計簿レベルで管理（管理者のみ編集可） |
| 個人カテゴリ設定 (UserCategoryPreference) | ユーザーごとのカテゴリ表示・並び順・**固定費フラグ**の設定。固定費フラグはメンバーが自身の使い方に応じて支出カテゴリ単位で個別に設定する |
| 収支明細 (Entry) | 1 件の収入または支出。日付・区分・カテゴリ・金額・メモ・**所有者**（誰の収支か）を持つ。所有者は登録時に選択可能（代理入力対応） |
| 予算 (Budget) | 支出カテゴリごとの月次予算（管理者が設定） |
| 招待 (Invitation) | 招待 URL。有効期限 24 時間・1 回限り |
| チャット (ChatMessage) | 家計簿に紐づくグループチャットのメッセージ。既読・オンライン状態を扱う |

> NOTE: 詳細なエンティティ・値オブジェクト・業務ルールは `docs/design/` 配下のドメイン設計ドキュメントで管理する。

## 重要な業務ルール

- **暦月集計**: 月次集計・固定費自動コピー・予算判定はすべて暦月（1 日〜末日）固定
- **固定費の月初自動コピー**: 各メンバーが固定費フラグを立てた支出カテゴリについて、そのメンバー自身の前月明細（金額・所有者・メモ）を月初に自動コピー
- **代理入力**: 収支明細の所有者は登録時に他メンバーへ変更可能（編集履歴は v1 では保持しない）
- **メンバーの権限**: 所属家計簿の収支明細は所有者を問わず閲覧・編集・削除可能
- **退会時のデータ**: 退会してもユーザーの収支明細は家計簿に残る（家計簿削除時にのみ全データ削除）

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

- [要件定義書](../docs/要件定義書.md) - 背景・目的・スコープ・機能要件・非機能要件・制約事項
- [設計ドキュメント](../docs/design/README.md) - ドメインごとの設計ドキュメント入口
- [フロントエンド ADR](../docs/adr/frontend/0001-architecture.md) - FSD アーキテクチャの決定記録
- [API 仕様書](../docs/api/openapi.yml) - OpenAPI 定義
