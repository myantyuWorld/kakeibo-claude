import type { RouteRecordRaw } from 'vue-router'

export const authRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    meta: { layout: 'public' },
    component: () => import('@/pages/login'),
  },
]
