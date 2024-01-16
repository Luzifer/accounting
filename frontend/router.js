import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { component: null, name: 'budget', path: '/' },
  { component: null, name: 'account-transactions', path: '/accounts/:id' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
