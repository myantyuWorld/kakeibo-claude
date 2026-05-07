# ADR-0001: フロントエンドアーキテクチャ

## コンテキスト

Vue + TypeScript による SPA のフロントエンド開発において、コードの構造・レイヤー間の依存方向・機能の分割方針を統一する必要がある。

Feature-Sliced Design（FSD）をベースとし、プロジェクト規模に合わせて 4 レイヤー構成に簡略化する。

## 決定

### ディレクトリ構成

```
frontend/src/
├── app/                    # App レイヤー（スライスなし）
│   ├── routes/             # ルーティング定義
│   ├── assets/             # グローバルアセット
│   │   ├── css/            #   スタイル
│   │   ├── fonts/          #   フォント（必要な場合のみ）
│   │   └── img/            #   画像（必要な場合のみ）
│   ├── app.vue
│   └── main.ts
│
├── pages/                  # Pages レイヤー
│   └── {page-name}/
│       ├── {page-name}-page.vue  # ページコンポーネント
│       ├── model/                # ページ固有のロジック・状態（必要な場合のみ）
│       └── index.ts              # Public API
│
├── features/               # Features レイヤー（ドメイン単位でスライス）
│   └── {domain}/
│       ├── ui/             # ドメイン固有の共通コンポーネント（必要な場合のみ）
│       ├── model/          # composable（use-create-{domain}.ts 等）
│       ├── api/            # API 呼び出し
│       ├── lib/            # feature 内ユーティリティ（必要な場合のみ）
│       └── index.ts        # Public API
│
└── shared/                 # Shared レイヤー（スライスなし）
    ├── ui/                 # 汎用 UI コンポーネント
    ├── api/                # API クライアント・インターセプター
    ├── lib/                # 汎用ユーティリティ
    ├── config/             # 環境変数・定数
    └── types/              # 共通型定義
```

### レイヤーの責務

| レイヤー | 責務 | スライス |
|---------|------|---------|
| app | アプリケーション初期化。ルーティング、グローバルプロバイダー、グローバルスタイル | なし |
| pages | 1 画面の構成。features の組み合わせとページ固有のレイアウト | あり（画面単位） |
| features | ドメインごとの API 呼び出しと composable をカプセル化。UI コンポーネントは持たない | あり（ドメイン単位） |
| shared | ビジネスロジックを含まない再利用可能な基盤。UI キット、API クライアント、ユーティリティ | なし（セグメントで整理） |

### 省略したレイヤーとその理由

| レイヤー | 理由 |
|---------|------|
| entities | ドメインモデルはバックエンドが管理する。フロントではAPI レスポンスの型定義（shared/types）と features 内のロジックで十分 |
| widgets | ページ間で再利用する大規模コンポーネントは features レイヤーで対応可能。必要になった時点で追加を検討する |

### レイヤー間の依存ルール

依存は上から下への一方向のみ許可する。

```
app → pages → features → shared
```

- 上位レイヤーは下位レイヤーにのみ依存できる
- 同一レイヤー内のスライス間の依存は禁止する
- shared は他のレイヤーに依存しない

```
例）
  features/accounts → shared/api          OK（下位レイヤー）
  features/accounts → features/orders     NG（同一レイヤー）
  shared/api → features/accounts          NG（上位レイヤー）
```

### Public API（index.ts）

- pages / features の各スライスは必ず `index.ts` を持つ
- 外部から参照できるのは `index.ts` で export されたものだけ
- スライス内部のファイルを直接 import することは禁止する

```typescript
// features/accounts/index.ts
export { useCreateAccount } from './model/use-create-account'
export { useListAccounts } from './model/use-list-accounts'

// pages/account-list/index.ts
export { default as AccountListPage } from './account-list-page.vue'
```

```typescript
// OK
import { useCreateAccount } from '@/features/accounts'

// NG — 内部パスへの直接アクセス
import { useCreateAccount } from '@/features/accounts/model/use-create-account'
```

### セグメントの役割

| セグメント | 配置するもの | 配置しないもの |
|-----------|------------|-------------|
| ui/ | Vue コンポーネント、スタイル | API 呼び出し、ビジネスロジック |
| model/ | composable、Pinia ストア、バリデーション | UI コンポーネント、API 実装 |
| api/ | API リクエスト関数、レスポンス型 | UI、状態管理 |
| lib/ | ヘルパー関数、定数 | 状態を持つもの |
| config/ | 環境変数、設定値（shared で使用） | ロジック |

すべてのセグメントが必須ではない。必要なものだけ作成する。

### 命名規則

| 対象 | 規則 | 例 |
|------|------|-----|
| ファイル名（すべて） | kebab-case | `login-form.vue`, `auth-store.ts` |
| スライス名（ディレクトリ） | kebab-case | `features/user-profile/` |
| composable | ファイル名は kebab-case、関数名は use プレフィックス | `auth-store.ts` → `useAuthStore()` |

## 結果

### 良い影響

- 4 レイヤーのシンプルな構成で学習コストが低い
- FSD の依存ルールにより、機能間の不要な結合を構造的に防止できる
- Public API により、リファクタリング時の影響範囲がスライス内に閉じる
- バックエンドの集約単位パッケージ（ADR-0001）と考え方が一致し、チーム内で設計の語彙を共有できる

### 悪い影響

- entities レイヤーを省略しているため、複数 features で同じ型やロジックを共有する場合に shared に寄せるか判断が必要になる。頻発するなら entities レイヤーの追加を検討する
- Public API（index.ts）の維持コストがかかる。export の追加漏れに注意が必要
