/* eslint-disable sort-imports */
import { createApp } from 'vue'
import VueShortkey from 'vue3-shortkey'

import './style.scss'

import appMain from './components/app.vue'
import router from './router'

const app = createApp(appMain)

app.use(router)
app.use(VueShortkey)
app.mount('#app')
