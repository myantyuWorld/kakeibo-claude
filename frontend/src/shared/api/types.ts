import type { components } from './schema'

export type ApiSchema<T extends keyof components['schemas']> =
  components['schemas'][T]
