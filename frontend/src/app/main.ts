import { createApp } from 'vue'
import { z } from 'zod'
import { createAppRouter } from './routes'
import { zodErrorMapJa } from '@/shared/lib/zod'
import './assets/css/global.css'
import App from './app.vue'

z.setErrorMap(zodErrorMapJa)

createApp(App).use(createAppRouter()).mount('#app')
