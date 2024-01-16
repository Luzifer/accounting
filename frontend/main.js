/* eslint-disable sort-imports */
import { createApp, h } from 'vue'

import './style.scss'

import appMain from './components/app.vue'
import router from './router.js'

const app = createApp({
  name: 'AccountingApp',

  render() {
    return h(appMain)
  },

  router,
})

app.use(router)
app.mount('#app')
