import { ref } from 'vue'
import { client } from '@/shared/api'
import { clearAccessToken } from '@/shared/auth'
import type { User } from './types'

// module スコープでシングルトン化。全 useAuth() 呼び出しで同じ状態を共有する
const currentUser = ref<User>()

export const useAuth = () => {
  const isAuthenticated = () => currentUser.value !== undefined

  const checkAuth = async (): Promise<boolean> => {
    if (currentUser.value) return true
    const { data } = await client.GET('/auth/me')
    if (!data) {
      currentUser.value = undefined
      return false
    }
    currentUser.value = data
    return true
  }

  const logout = () => {
    currentUser.value = undefined
    clearAccessToken()
  }

  return { currentUser, isAuthenticated, checkAuth, logout }
}
