import { http, type HttpResponseResolver } from 'msw'
import { setupServer } from 'msw/node'
import type { paths } from '@/shared/api/schema'
import { handlers } from './handlers'
import { apiUrl } from './api-url'

export const server = setupServer(...handlers)

type PathsWith<M extends string> = {
  [P in keyof paths]: paths[P] extends { [K in M]: unknown } ? P : never
}[keyof paths]

const define = <M extends keyof typeof http>(method: M) => {
  return <P extends PathsWith<M>>(path: P, resolver: HttpResponseResolver) => {
    server.use(http[method](apiUrl(path as keyof paths), resolver))
  }
}

export const mockApi = {
  get: define('get'),
  post: define('post'),
  put: define('put'),
  patch: define('patch'),
  delete: define('delete'),
}
