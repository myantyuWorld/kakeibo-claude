import { describe, it, expect } from 'vitest'
import { withSetup } from '@/shared/__tests__/helper'
import { useLoginForm } from '../model/use-login-form'

describe('useLoginForm', () => {
  it('空のまま送信するとメールアドレスとパスワードの両方にエラーが出る', async () => {
    const { errors, onSubmit } = withSetup(() => useLoginForm())

    await onSubmit()

    expect(errors.value.email).toBeDefined()
    expect(errors.value.password).toBeDefined()
  })

  it('不正なメールアドレスではメール形式エラーが出る', async () => {
    const { defineField, errors, onSubmit } = withSetup(() => useLoginForm())
    const [email] = defineField('email')
    const [password] = defineField('password')

    email.value = 'invalid'
    password.value = 'password123'
    await onSubmit()

    expect(errors.value.email).toBeDefined()
    expect(errors.value.password).toBeUndefined()
  })

  it('正しい入力ではエラーなしで送信できる', async () => {
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
