import { describe, it, expect, beforeEach } from 'vitest'
import { HttpResponse } from 'msw'
import { mockApi } from '@/shared/__tests__/server'
import { useAuth } from '../model/use-auth'

describe('useAuth', () => {
  beforeEach(() => {
    const { logout } = useAuth()
    logout()
  })

  it('初期状態では未認証である', () => {
    const { isAuthenticated } = useAuth()
    expect(isAuthenticated()).toBe(false)
  })

  it('checkAuth が成功するとユーザー情報を保持する', async () => {
    const { checkAuth, currentUser, isAuthenticated } = useAuth()
    const ok = await checkAuth()
    expect(ok).toBe(true)
    expect(isAuthenticated()).toBe(true)
    expect(currentUser.value?.email).toBe('test@example.com')
  })

  it('checkAuth が 401 を返すと未認証のままになる', async () => {
    mockApi.get('/auth/me', () => new HttpResponse(null, { status: 401 }))
    const { checkAuth, isAuthenticated } = useAuth()
    const ok = await checkAuth()
    expect(ok).toBe(false)
    expect(isAuthenticated()).toBe(false)
  })
})
