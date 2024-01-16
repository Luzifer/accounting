import { createRouter, createWebHistory } from 'vue-router'

import budgetDashboard from './components/budgetDashboard.vue'

const routes = [
  { component: budgetDashboard, name: 'budget', path: '/' },
  { component: null, name: 'account-transactions', path: '/accounts/:id' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
