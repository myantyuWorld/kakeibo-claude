import type { RouteRecordRaw } from 'vue-router'

export const errorRoutes: RouteRecordRaw[] = [
  {
    path: '/error/404',
    name: 'not-found',
    meta: { layout: 'public', skipAuth: true },
    component: () => import('@/pages/not-found'),
  },
  {
    path: '/error/500',
    name: 'internal-error',
    meta: { layout: 'public', skipAuth: true },
    component: () => import('@/pages/internal-error'),
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/error/404',
  },
]
