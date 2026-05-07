# Frontend

Vue 3 + TypeScript + Vite によるフロントエンド。FSD（Feature-Sliced Design）アーキテクチャを採用。

## ディレクトリ構成

```
frontend/src/
├── app/                  # アプリ初期化・ルーティング・レイアウト
│   ├── layouts/          # レイアウトコンポーネント（public / private）
│   └── routes/           # ルート定義
├── pages/                # ページコンポーネント
│   ├── home/
│   ├── login/
│   ├── not-found/
│   └── internal-error/
├── features/             # 機能モジュール
│   └── auth/             # 認証（composable / schema / types）
└── shared/               # 共通モジュール
    ├── api/              # openapi-fetch クライアント・型定義
    ├── auth/             # トークン管理
    ├── ui/               # 共通 UI コンポーネント（button / input）
    ├── lib/              # ライブラリラッパー（zod）
    └── __tests__/        # テストヘルパー・MSW 設定
```

## 使用ライブラリ

### プロダクション

| ライブラリ | 用途 |
|-----------|------|
| [Vue 3](https://vuejs.org/) | UI フレームワーク |
| [vue-router](https://router.vuejs.org/) | ルーティング |
| [openapi-fetch](https://openapi-ts.dev/openapi-fetch/) | 型安全な API クライアント |
| [vee-validate](https://vee-validate.logaretm.com/v4/) | フォームバリデーション |
| [zod](https://zod.dev/) | スキーマバリデーション |

### 開発

| ライブラリ | 用途 |
|-----------|------|
| [Vite](https://vite.dev/) | ビルドツール |
| [TypeScript](https://www.typescriptlang.org/) | 型システム |
| [Vitest](https://vitest.dev/) | テストフレームワーク |
| [happy-dom](https://github.com/nicedayfor/happy-dom) | DOM エミュレーション |
| [MSW](https://mswjs.io/) | API モックサーバー |
| [Storybook](https://storybook.js.org/) | コンポーネントカタログ |
| [openapi-typescript](https://openapi-ts.dev/) | OpenAPI → TypeScript 型生成 |
| [ESLint](https://eslint.org/) | Linter |
| [Prettier](https://prettier.io/) | Formatter |

## mise タスク

`frontend/` ディレクトリで実行する。

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
