import { afterAll, afterEach, beforeAll } from 'vitest'
import { server } from './server'

// happy-dom の localStorage が不完全な環境向けのフォールバック
if (typeof localStorage.removeItem !== 'function') {
  const store = new Map<string, string>()
  Object.defineProperty(globalThis, 'localStorage', {
    value: {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, String(value)),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
      get length() {
        return store.size
      },
      key: (index: number) => [...store.keys()][index] ?? null,
    },
  })
}

beforeAll(() => server.listen({ onUnhandledRequest: 'error' }))
afterEach(() => server.resetHandlers())
afterAll(() => server.close())
