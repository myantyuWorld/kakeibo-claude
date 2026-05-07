# 開発ガイド

このプロジェクトは Claude Code による AI 駆動開発を前提としています。設計から実装・レビュー・PR 作成までをスキル（スラッシュコマンド）で自動化し、開発者は意思決定とレビューに集中します。

## 開発フロー

プロジェクトの開発は以下の工程で進みます。各工程に対応するスキルを順に実行することで、一貫した品質の成果物が得られます。

```
設計 → 仕様生成 → 実装 → PR
```

| 工程 | 成果物 | 使用するスキル |
|------|--------|---------------|
| 設計 | drawio 図、業務ルール.md | `/drawio` `/review-design` `/review-business-rules` |
| 仕様生成 | OpenAPI YAML、ER 図（DBML） | `/generate-api-spec` `/generate-er` |
| 実装（バックエンド） | Go コード、テスト、Draft PR | `/impl-api` |
| 実装（フロントエンド） | Vue コンポーネント、テスト、Story、Draft PR | `/impl-front` |

## 設計工程

### /drawio

図やダイアグラムを `.drawio` ファイルとして生成します。コンテキストマップ、ドメインモデル図、ユースケース図など、sudo モデリングの成果物作成に使用します。PNG / SVG / PDF へのエクスポートにも対応しています。

### /review-design

sudo モデリングの drawio 成果物をレビューします。図の曖昧さや不整合を検出し、CRITICAL / HIGH の指摘は付箋として drawio ファイルに直接追記します。複数の図を横断したクロスチェックも行います。

### /review-business-rules

業務ルール.md の完全性・明確性をレビューします。5 つの観点（網羅性、明確性、エッジケース、状態遷移、ドメイン間整合性）で検査し、曖昧な点は対話で確認しながら受入基準（Given/When/Then）まで整理します。

## 仕様生成工程

### /generate-api-spec

設計ドキュメント（ドメインモデル.md、業務ルール.md）から OpenAPI YAML を生成します。ドメイン名を引数に指定します。

```
/generate-api-spec customer
/generate-api-spec customer 顧客登録    # ユースケース指定
```

### /generate-er

ドメインモデルから ER 図（DBML）を生成します。DDD のモデルを RDB スキーマにマッピングするルール（値オブジェクトの埋め込み、Enum のカラム化など）に従って変換します。

```
/generate-er customer
```

## 実装工程

### /impl-api

バックエンドの機能実装を一気通貫で実行します。設計ドキュメントの確認からドメインモデル実装、API 実装、テスト生成・実行、AI レビュー、Draft PR 作成、PR へのインラインコメント投稿までを自動で行います。

```
/impl-api customer 顧客登録
```

実行される工程:

1. 設計ドキュメントの確認
2. ドメインモデルの実装（未実装の場合）
3. API の実装（handler → service → repository）
4. テスト生成・実行（失敗時は最大 3 回修正）
5. AI コードレビュー（go-reviewer → security-reviewer → go-test-reviewer）
6. Draft PR の作成
7. 実装意図のインラインコメント投稿

### /impl-front

フロントエンドの画面実装を一気通貫で実行します。設計・API 仕様の確認からルート定義、composable 実装、Vue コンポーネント実装、テスト・Storybook 生成、ブラウザ QA、AI レビュー、Draft PR 作成までを自動で行います。

```
/impl-front 顧客一覧
```

実行される工程:

1. 設計・API 仕様の確認
2. ルート定義の追加
3. composable の実装（API 呼び出し・バリデーション）
4. Vue ページコンポーネントの実装
5. テスト・Storybook Story の生成・実行（失敗時は最大 3 回修正）
6. ブラウザ QA 確認（`/browser-qa`）
7. AI コードレビュー（vue-reviewer → a11y-reviewer → vue-test-reviewer）
8. Draft PR の作成
9. 実装意図のインラインコメント投稿

## 補助スキル

### /generate-api-tests

受入基準（業務ルール.md の Given/When/Then）からバックエンドのテストコードを生成します。単体テスト・統合テストの両方を生成し、テストと受入基準の対応表も出力します。`/impl-api` の工程内で自動実行されるため、通常は単体で使用する必要はありません。

### /generate-frontend-tests

既存のフロントエンドコードに対するテストを生成します。composable、コンポーネント、フォームバリデーションなど、対象に応じた適切なテストパターンを適用します。`/impl-front` の工程内で自動実行されます。

### /browser-qa

Playwright MCP を使って開発サーバーや Storybook の画面をブラウザで確認します。スクリーンショット取得、コンソールエラー検出、DOM 構造チェックを行います。`/impl-front` の工程内で自動実行されます。

## AI レビュー

実装スキル（`/impl-api`、`/impl-front`）は、PR 作成前に複数の AI レビュアーを自動実行します。手動でレビューを依頼する必要はありません。

### バックエンド

| レビュアー | 観点 |
|-----------|------|
| go-reviewer | DDD アーキテクチャ・設計ドキュメントとの整合性 |
| security-reviewer | SQLi / XSS / 認可漏れ / 情報漏洩 |
| go-test-reviewer | テスト網羅性・受入基準との対応 |

### フロントエンド

| レビュアー | 観点 |
|-----------|------|
| vue-reviewer | FSD アーキテクチャ・Vue 3 コーディング規約 |
| a11y-reviewer | セマンティック HTML・ARIA・アクセシビリティ |
| vue-test-reviewer | テスト網羅性・テストパターンの適切さ |

## Rules

`.claude/rules/` 配下に、コーディング規約やアーキテクチャルールが定義されています。スキル実行時に自動で読み込まれるため、開発者が意識する必要はありません。変更したい場合は以下のディレクトリを確認してください。

```
.claude/rules/
├── common/       # 共通（Git ワークフロー、コーディング規約）
├── backend/      # バックエンド（Go / DDD / テスト / DB 設計）
└── frontend/     # フロントエンド（Vue / TypeScript / FSD / Storybook）
```

## クイックスタート

環境構築は [README.md](README.md) を参照してください。構築後、以下の流れで開発を始められます。

```bash
# 1. 設計ドキュメントを確認
#    docs/design/{domain}/ 配下にドメインモデル.md と業務ルール.md があることを確認

# 2. API 仕様書を生成（未作成の場合）
/generate-api-spec {domain}

# 3. バックエンドを実装
/impl-api {domain} {usecase}

# 4. フロントエンドを実装
/impl-front {page}
```
