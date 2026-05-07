import { describe, it, expect } from 'vitest'
import { withSetup } from '@/shared/__tests__/helper'
import { useLoginForm } from '@/features/auth'

describe('LoginPage', () => {
  it('空のまま送信するとバリデーションエラーが表示される', async () => {
    const { errors, onSubmit } = withSetup(() => useLoginForm())

    await onSubmit()

    expect(errors.value.email).toBe('必須項目です')
    expect(errors.value.password).toBe('必須項目です')
  })

  it('不正なメールアドレスで送信するとエラーが表示される', async () => {
    const { defineField, errors, onSubmit } = withSetup(() => useLoginForm())

    const [email] = defineField('email')
    const [password] = defineField('password')
    email.value = 'invalid'
    password.value = 'password123'

    await onSubmit()

    expect(errors.value.email).toBe('メールアドレスの形式が正しくありません')
  })

  it('正しい入力ではエラーが表示されない', async () => {
    const { defineField, errors, onSubmit } = withSetup(() => useLoginForm())

    const [email] = defineField('email')
    const [password] = defineField('password')
    email.value = 'test@example.com'
    password.value = 'password123'

    await onSubmit()

    expect(errors.value.email).toBeUndefined()
    expect(errors.value.password).toBeUndefined()
  })
})
