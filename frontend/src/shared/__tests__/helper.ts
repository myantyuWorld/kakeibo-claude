import { createApp, defineComponent, h } from 'vue'

export const withSetup = <T>(composable: () => T): T => {
  let result!: T
  const app = createApp(
    defineComponent({
      setup() {
        result = composable()
        return () => h('div')
      },
    }),
  )
  app.mount(document.createElement('div'))
  return result
}
