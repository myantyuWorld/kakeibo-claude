import type { Preview } from '@storybook/vue3'
import { z } from 'zod'
import { zodErrorMapJa } from '../src/shared/lib/zod/error-map'

z.setErrorMap(zodErrorMapJa)

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /date$/i,
      },
    },
  },
}

export default preview
