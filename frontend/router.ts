/* eslint-disable sort-imports */
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import accountOverview from './views/accountOverview.vue'
import budgetDashboard from './views/budgetDashboard.vue'

const routes: RouteRecordRaw[] = [
  { component: budgetDashboard, name: 'budget', path: '/' },
  { component: accountOverview, name: 'account-transactions', path: '/accounts/:accountId', props: true },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
