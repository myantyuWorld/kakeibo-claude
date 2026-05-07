import type { Meta, StoryObj } from '@storybook/vue3-vite'
import { expect, userEvent, within } from 'storybook/test'
import LoginPage from './login-page.vue'

const meta = {
  component: LoginPage,
} satisfies Meta<typeof LoginPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

export const ValidationError: Story = {
  name: '未入力のまま送信すると必須エラーが表示される',
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement)

    await userEvent.click(canvas.getByRole('button', { name: 'ログイン' }))

    await expect(canvas.getByText('必須項目です')).toBeVisible()
  },
}

export const InvalidEmail: Story = {
  name: '不正なメールアドレスを入力した場合、送信すると形式エラーが表示される',
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement)

    await userEvent.type(canvas.getByLabelText('メールアドレス'), 'invalid')
    await userEvent.type(canvas.getByLabelText('パスワード'), 'password123')
    await userEvent.click(canvas.getByRole('button', { name: 'ログイン' }))

    await expect(
      canvas.getByText('メールアドレスの形式が正しくありません'),
    ).toBeVisible()
  },
}

export const ValidInput: Story = {
  name: '正しい値を入力した場合、送信してもエラーが表示されない',
  play: async ({ canvasElement }) => {
    const canvas = within(canvasElement)

    await userEvent.type(
      canvas.getByLabelText('メールアドレス'),
      'test@example.com',
    )
    await userEvent.type(canvas.getByLabelText('パスワード'), 'password123')
    await userEvent.click(canvas.getByRole('button', { name: 'ログイン' }))

    const errors = canvasElement.querySelectorAll('p')
    for (const error of errors) {
      await expect(error).not.toBeVisible()
    }
  },
}
