---
name: generate-frontend-tests
description: フロントエンドの既存コードに対するテストを生成する。「フロントのテストを書いて」「コンポーネントのテストを生成して」など、フロントエンドのテスト作成を依頼されたときに使用する。
---

# generate-frontend-tests

既存のフロントエンドコードに対してテストコードを生成する。

## When to Activate

- 「フロントのテストを書いて」「コンポーネントのテストを生成して」と依頼されたとき
- `impl-front` の Phase 5 から呼び出されたとき
- composable やページコンポーネントの実装後にテストを追加するとき

## 入力

- 対象ファイルまたはスライス（例: `features/auth`、`pages/login`）

## 手順

### 1. テスト対象の分析

対象のコードを読み、テストすべき振る舞いを洗い出す。

- composable: 公開している ref / computed / メソッドの振る舞い
- コンポーネント: ユーザー操作に対する表示の変化

### 2. テストコードの生成

`.claude/rules/frontend/testing.md` に従いテストコードを生成する。

#### composable のテスト

- `withSetup` ヘルパーを使い Vue コンテキスト内で実行する
- API 呼び出しがある場合は MSW + `mockApi` でモックする
- シングルトン状態を持つ composable は `beforeEach` でリセットする

```typescript
// PASS: withSetup + mockApi を使い、状態リセットも行う
import { withSetup } from '@/shared/__tests__/helper'
import { mockApi } from '@/shared/__tests__/server'
import { HttpResponse } from 'msw'
import { useAuth } from '../model/use-auth'

describe('useAuth', () => {
  beforeEach(() => {
    const { logout } = useAuth()
    logout()
  })

  it('checkAuth 成功時に認証済みになる', async () => {
    const { isAuthenticated, checkAuth } = withSetup(() => useAuth())
    await checkAuth()
    expect(isAuthenticated.value).toBe(true)
  })

  it('checkAuth が 401 を返すと未認証のままになる', async () => {
    mockApi.get('/auth/me', () => new HttpResponse(null, { status: 401 }))
    const { isAuthenticated, checkAuth } = withSetup(() => useAuth())
    await checkAuth()
    expect(isAuthenticated.value).toBe(false)
  })
})

// FAIL: withSetup を使わず直接呼び出す
import { useAuth } from '../model/use-auth'

it('認証チェック', () => {
  const { isAuthenticated } = useAuth()  // NG: Vue コンテキスト外で呼び出し
  expect(isAuthenticated.value).toBe(false)
})

// FAIL: vi.mock で手動モックする
vi.mock('@/shared/api/client', () => ({  // NG: MSW を使うべき
  client: { GET: vi.fn() }
}))
```

#### フォームバリデーションのテスト

vee-validate + zod のフォームは、バリデーションの境界値を重点的にテストする。

```typescript
// PASS: 空入力・不正形式・正常入力の各パターンをテスト
describe('useLoginForm', () => {
  it('空のまま送信するとメールとパスワードにエラーが出る', async () => {
    const { errors, onSubmit } = withSetup(() => useLoginForm())
    await onSubmit()
    expect(errors.value.email).toBeDefined()
    expect(errors.value.password).toBeDefined()
  })

  it('不正なメール形式ではエラーが出る', async () => {
    const { email, errors, onSubmit } = withSetup(() => useLoginForm())
    email.value = 'invalid-email'
    await onSubmit()
    expect(errors.value.email).toBeDefined()
  })

  it('正しい入力ではエラーが出ない', async () => {
    const { email, password, errors, onSubmit } = withSetup(() => useLoginForm())
    email.value = 'user@example.com'
    password.value = 'password123'
    await onSubmit()
    expect(errors.value.email).toBeUndefined()
    expect(errors.value.password).toBeUndefined()
  })
})

// FAIL: バリデーションの異常系がない
describe('useLoginForm', () => {
  it('送信できる', async () => {  // NG: 正常系しかテストしていない
    // ...
  })
})
```

#### Vue コンポーネントのテスト

- 条件付き表示: 表示/非表示の切り替え
- イベント: ボタンクリック、フォーム送信
- スロット・props の反映

```typescript
// PASS: ユーザー操作と表示の変化をテスト
import { mount } from '@vue/test-utils'

it('削除ボタンをクリックすると確認ダイアログが表示される', async () => {
  const wrapper = mount(CustomerDetail)
  await wrapper.find('[data-testid="delete-button"]').trigger('click')
  expect(wrapper.find('[data-testid="confirm-dialog"]').exists()).toBe(true)
})

// FAIL: 実装の内部構造に依存する
it('削除', async () => {
  const wrapper = mount(CustomerDetail)
  expect(wrapper.vm.showDialog).toBe(false)  // NG: 内部の ref を直接参照
  wrapper.vm.showDialog = true                // NG: 内部状態を直接操作
})
```

### 3. テストケースの網羅性確認

生成後、以下をチェックする。

- 正常系と異常系の両方があるか
- 境界値（空文字、最大文字数等）がカバーされているか
- 条件分岐の各パスがテストされているか
- API エラー時の振る舞いがテストされているか

## テストケース命名パターン

| パターン | テストケース名の例 |
|---------|-----------------|
| 正常系 | `checkAuth 成功時に認証済みになる` |
| APIエラー | `checkAuth が 401 を返すと未認証のままになる` |
| バリデーション（空） | `空のまま送信するとエラーが出る` |
| バリデーション（不正） | `不正なメール形式ではエラーが出る` |
| 条件付き表示 | `未認証の場合ログインリンクが表示される` |
| ユーザー操作 | `削除ボタンをクリックすると確認ダイアログが表示される` |
| 初期状態 | `初期状態では未認証である` |

---

**Remember**: composable は withSetup、API モックは MSW。内部実装ではなく振る舞いをテストする。
