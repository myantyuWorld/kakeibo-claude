import { http, HttpResponse } from 'msw'
import { apiUrl } from './api-url'

export const handlers = [
  http.get(apiUrl('/auth/me'), () => {
    return HttpResponse.json({
      id: '00000000-0000-0000-0000-000000000000',
      email: 'test@example.com',
    })
  }),
]
