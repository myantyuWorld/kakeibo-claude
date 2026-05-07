export type LayoutType = 'public' | 'private'

declare module 'vue-router' {
  // eslint-disable-next-line @typescript-eslint/consistent-type-definitions
  interface RouteMeta {
    layout?: LayoutType
    skipAuth?: boolean
  }
}
