import createClient, { type Middleware } from 'openapi-fetch'
import type { paths } from './schema'
import { getAccessToken } from '@/shared/auth'

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    const token = getAccessToken()
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`)
    }
    return request
  },
}

const errorMiddleware: Middleware = {
  async onResponse({ response }) {
    if (response.status >= 500) {
      window.location.href = '/error/500'
    }
    return response
  },
  async onError() {
    window.location.href = '/error/500'
  },
}

export const client = createClient<paths>({
  baseUrl: import.meta.env.VITE_API_BASE_URL ?? '',
})

client.use(authMiddleware, errorMiddleware)
