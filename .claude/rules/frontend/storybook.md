---
paths:
  - "frontend/src/**/*.stories.ts"
---

# Storybook 規約

## 配置

- Story ファイルは対象コンポーネントと同じディレクトリに置く
- ファイル名はコンポーネントに合わせる（`login-page.vue` → `login-page.stories.ts`）

## CSF 3.0 形式

- `meta` を `satisfies Meta<typeof Component>` で定義し、`export default` する
- 各 Story は `StoryObj<typeof meta>` 型の named export とする

```typescript
import type { Meta, StoryObj } from '@storybook/vue3-vite'
import LoginPage from './login-page.vue'

const meta = {
  component: LoginPage,
} satisfies Meta<typeof LoginPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}
```

## Story の命名

- `Default`（または `Primary`）以外の Story には必ず `name` プロパティを設定する
- `name` は Given-When-Then 形式で、日本語で振る舞いを記述する

```typescript
// OK — Given-When-Then 形式
export const ValidationError: Story = {
  name: '未入力の場合、送信すると必須エラーが表示される',
}

export const SubmitSuccess: Story = {
  name: '正しい値を入力した場合、送信すると次の画面に遷移する',
}

// NG — name なし
export const Disabled: Story = { args: { disabled: true } }

// NG — 曖昧・動作のみ
export const Validation: Story = {
  name: 'バリデーション',
}
```

## play function によるインタラクションテスト

- ユーザー操作を `play` 関数で記述し、状態遷移を Story として表現する
- `within(canvasElement)` でスコープを限定し、`userEvent` で操作する

```typescript
import { expect, userEvent, within } from 'storybook/test'

export const ValidationError: Story = {
  name: '未入力の場合、送信すると必須エラーが表示される',
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement)
    await userEvent.click(canvas.getByRole('button', { name: 'ログイン' }))
    await expect(canvas.getByText('必須項目です')).toBeVisible()
  },
}
```

## バリアント

- ALWAYS `Default` Story を用意する
- 状態の切り替わりがあるコンポーネントは、各状態を Story として網羅する（エラー状態、空状態、ローディング等）

## 関連

- **vue-reviewer** agent -- コンポーネントレビュー時
- skill: `impl-front` -- Story作成時
