---
paths:
  - "frontend/src/**/*.test.ts"
---

# フロントエンドテスト規約

## 配置

- テストファイルは対象スライスの `__tests__/` ディレクトリに置く
- ファイル名はテスト対象に合わせる（`use-auth.ts` → `use-auth.test.ts`）

```
features/auth/
├── model/
│   ├── use-auth.ts
│   └── use-login-form.ts
├── __tests__/
│   ├── use-auth.test.ts
│   └── use-login-form.test.ts
└── index.ts
```

## 命名

- `describe` / `it` の説明は日本語で振る舞いを記述する

```typescript
describe('useAuth', () => {
  it('初期状態では未認証である', () => { ... })
  it('checkAuth が 401 を返すと未認証のままになる', async () => { ... })
})
```

## composable のテスト

- ALWAYS `withSetup` ヘルパーを使って Vue アプリケーションコンテキスト内で実行する
- NEVER composable を直接呼び出さない（`ref` や `inject` 等の Vue API がコンテキスト外で動作しない）

```typescript
import { withSetup } from '@/shared/__tests__/helper'
import { useLoginForm } from '../model/use-login-form'

it('空のまま送信するとエラーが出る', async () => {
  const { errors, onSubmit } = withSetup(() => useLoginForm())
  await onSubmit()
  expect(errors.value.email).toBeDefined()
})
```

## API モック

- ALWAYS MSW + `mockApi` ヘルパーを使う
- NEVER `vi.mock` 等で手動モックしない
- デフォルトハンドラは `shared/__tests__/handlers.ts` に定義済み。テスト固有のレスポンスだけ上書きする

```typescript
import { mockApi } from '@/shared/__tests__/server'
import { HttpResponse } from 'msw'

it('401 を返すと未認証のままになる', async () => {
  mockApi.get('/auth/me', () => new HttpResponse(null, { status: 401 }))
  // ...
})
```

## 状態リセット

- composable がモジュールスコープの状態（シングルトン）を持つ場合、`beforeEach` でリセットする

```typescript
beforeEach(() => {
  const { logout } = useAuth()
  logout()
})
```

## 関連

- **vue-test-reviewer** agent -- テストレビュー時
- skill: `generate-frontend-tests` -- テストコード生成時
