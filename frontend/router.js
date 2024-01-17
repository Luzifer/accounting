import { createRouter, createWebHistory } from 'vue-router'

import accountOverview from './components/accountOverview.vue'
import budgetDashboard from './components/budgetDashboard.vue'

const routes = [
  { component: budgetDashboard, name: 'budget', path: '/' },
  { component: accountOverview, name: 'account-transactions', path: '/accounts/:accountId', props: true },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
