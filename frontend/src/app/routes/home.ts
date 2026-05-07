import type { RouteRecordRaw } from 'vue-router'

export const homeRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    meta: { layout: 'private' },
    component: () => import('@/pages/home'),
  },
]
