# Backend

Go + Echo v4 による API サーバー。DDD ベースのモジュール構成を採用。

## ディレクトリ構成

```
backend/
├── cmd/
│   ├── api/              # API サーバーのエントリポイント
│   └── migrate/          # マイグレーション実行
├── internal/             # ドメインモジュール（非公開）
│   ├── _template/        # 新規ドメイン作成用テンプレート
│   └── auth/             # 認証
├── pkg/                  # 横断的パッケージ（公開）
│   ├── api/              # Echo サーバー設定・バリデーター・エラーハンドラー
│   ├── config/           # 環境変数の読み込み
│   ├── errors/           # アプリケーション共通エラー
│   ├── postgres/         # DB 接続・トランザクション実装
│   └── transaction/      # トランザクションインターフェース
├── migrations/           # SQL マイグレーションファイル
├── tests/                # 統合テスト
└── dbconfig.yml          # sql-migrate 設定
```

## 使用ライブラリ

| ライブラリ | 用途 |
|-----------|------|
| [Echo v4](https://echo.labstack.com/) | HTTP フレームワーク |
| [sqlx](https://github.com/jmoiron/sqlx) | SQL クライアント |
| [sql-migrate](https://github.com/rubenv/sql-migrate) | マイグレーション |
| [validator/v10](https://github.com/go-playground/validator) | バリデーション |
| [testify](https://github.com/stretchr/testify) | テストアサーション・mock |

## mise タスク

`backend/` ディレクトリで実行する。

| コマンド | 用途 |
|---------|------|
| `mise run dev` | Air によるホットリロード開発サーバー起動 |
| `mise run lint` | golangci-lint 実行 |
| `mise run test` | 単体テスト実行（`./internal/...`） |
| `mise run test-integration` | 統合テスト実行（`./tests/...`） |
| `mise run test-all` | 全テスト実行 |
