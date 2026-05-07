import { createRouter, createWebHistory, createMemoryHistory } from 'vue-router'
import { useAuth } from '@/features/auth'
import { homeRoutes } from './home'
import { authRoutes } from './auth'
import { errorRoutes } from './error'

const routes = [...homeRoutes, ...authRoutes, ...errorRoutes]

export const createAppRouter = (mode: 'web' | 'memory' = 'web') => {
  const router = createRouter({
    history: mode === 'web' ? createWebHistory() : createMemoryHistory(),
    routes,
  })

  router.beforeEach(async (to) => {
    if (to.meta.skipAuth) return true

    const { isAuthenticated, checkAuth } = useAuth()

    if (to.meta.layout === 'public') {
      if (isAuthenticated()) return { name: 'home' }
      return true
    }

    if (await checkAuth()) return true
    return { name: 'login' }
  })

  return router
}
