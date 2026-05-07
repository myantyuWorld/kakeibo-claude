import type { paths } from '@/shared/api/schema'

export const apiUrl = <P extends keyof paths>(path: P) =>
  `${import.meta.env.VITE_API_BASE_URL ?? ''}${path as string}`
