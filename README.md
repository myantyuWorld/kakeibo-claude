# 家計簿アプリ (kakeibo)

2人（パートナー / 夫婦など）で共有する家計簿アプリ。スプレッドシートでの運用をシステム化することを目的とする。

## 主な機能（想定）

- メンバーごとの収入・支出の入力
- カテゴリ別の分類と集計
- 月次・カテゴリ別の支出推移の可視化

詳細なドメイン設計は [docs/design/](docs/design/) を参照。

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| バックエンド | Go 1.26 / Echo v4 / PostgreSQL / sqlx / sql-migrate |
| フロントエンド | Vue 3 / TypeScript / Vite / vue-router |
| タスクランナー | mise |
| CI | GitHub Actions |

## ディレクトリ構成

```
.
├── .claude/              # Claude Code 設定（rules / agents / skills）
├── .github/              # GitHub Actions / PR テンプレート
├── backend/              # バックエンド実装
├── frontend/             # フロントエンド実装
├── backlog.md            # タスク・要望の回収キュー（開発ループの入口）
└── docs/
    ├── design/           # ドメイン設計ドキュメント
    ├── api/              # OpenAPI 仕様書
    └── adr/              # Architecture Decision Records
```

各ディレクトリの詳細（構成・使用ライブラリ・コマンド等）は個別の README を参照してください。

- [backend/README.md](backend/README.md)
- [frontend/README.md](frontend/README.md)

## 前提条件

- [mise](https://mise.jdx.dev/) がインストールされていること
- Docker / Docker Compose が利用可能であること

## セットアップ

```bash
# 1. リポジトリをクローン
git clone <repository-url>
cd kakeibo-claude

# 2. mise でツールをインストール
mise install

# 3. Git hooks を有効化（lefthook + gitleaks）
mise exec -- lefthook install

# 4. 環境変数ファイルを作成
mise run setup

# 5. Docker コンテナ（PostgreSQL）を起動
mise run up

# 6. マイグレーションを実行
cd backend && go run cmd/migrate/main.go

# 7. 開発サーバーを起動（バックエンド / フロントエンドそれぞれ）
cd backend && mise run dev      # localhost:8080
cd frontend && mise run dev     # localhost:5173
```

## コマンド一覧

プロジェクトルートで実行する共通コマンド。バックエンド・フロントエンド固有のコマンドは各 README を参照。

| コマンド | 用途 |
|---------|------|
| `mise run setup` | `.env.example` から `.env` を作成 |
| `mise run up` | Docker コンテナ起動 |
| `mise run down` | Docker コンテナ停止 |

## Claude Code を使った開発

このリポジトリは Claude Code による AI 駆動開発を前提としています。

### 主要スキル

| スキル | 用途 |
|-------|------|
| `/impl-api` | バックエンドの機能実装（設計確認 → 実装 → テスト → PR） |
| `/impl-front` | フロントエンドの画面実装（設計確認 → 実装 → テスト → PR） |
| `/review-design` | sudo モデリング成果物のレビュー |
| `/generate-api-spec` | 設計ドキュメントから OpenAPI 仕様書を生成 |
| `/generate-er` | ドメインモデルから ER 図（DBML）を生成 |
| `/loop-review` | 要望の仕分け → ルール反映 → Loop Contract 更新 |
| `/pr-watch` | merge 駆動のタスク送り（PR close → 次タスク着手） |
| `/loop-diagnosis` | 操作ログの診断（繰り返し・過介入・催促の洗い出し） |

開発ループの入口は [backlog.md](backlog.md) です。人間が握るのは **優先度の決定** と **PR の merge** の2点だけ。

詳細は [DEVELOPMENT.md](DEVELOPMENT.md) を参照してください。
