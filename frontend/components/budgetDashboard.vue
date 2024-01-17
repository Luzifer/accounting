<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col d-flex align-items-center">
        <range-selector
          v-model="timeRange"
          :multi-month="false"
        />
      </div>
      <div class="col d-flex align-items-center justify-content-end">
        <div :class="unallocatedMoneyClass">
          <span class="fs-4">{{ formatNumber(unallocatedMoney) }} €</span>
          <span class="small">Unallocated</span>
        </div>
      </div>
    </div>
    <div class="row mt-3">
      <div class="col">
        <table class="table table-striped small">
          <thead>
            <tr>
              <th>Category</th>
              <th class="text-end">
                Allocated
              </th>
              <th class="text-end">
                Activity
              </th>
              <th class="text-end">
                Available
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="cat in categories"
              :key="cat.id"
            >
              <td>{{ cat.name }}</td>
              <td :class="{'text-end': true, 'text-danger': (allocatedByCategory[cat.id] || 0) < 0}">
                {{ formatNumber(allocatedByCategory[cat.id] || 0) }} €
              </td>
              <td :class="{'text-end': true, 'text-danger': (activityByCategory[cat.id] || 0) < 0}">
                {{ formatNumber(activityByCategory[cat.id] || 0) }} €
              </td>
              <td :class="{'text-end': true, 'text-danger': cat.balance < 0}">
                {{ formatNumber(cat.balance) }} €
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { formatNumber } from '../helpers'
import rangeSelector from './rangeSelector.vue'
import { unallocatedMoneyAcc } from '../constants'

export default {
  components: { rangeSelector },

  computed: {
    activityByCategory() {
      return this.transactions
        .filter(tx => tx.account)
        .reduce((alloc, tx) => {
          alloc[tx.category] = (alloc[tx.category] || 0) + tx.amount
          return alloc
        }, {})
    },

    allocatedByCategory() {
      return this.transactions
        .filter(tx => !tx.account)
        .reduce((alloc, tx) => {
          alloc[tx.category] = (alloc[tx.category] || 0) + tx.amount
          return alloc
        }, {})
    },

    categories() {
      const accounts = this.accounts
        .filter(acc => acc.type === 'category')
        .filter(acc => acc.id !== unallocatedMoneyAcc)
      accounts.sort((a, b) => a.name.localeCompare(b.name))
      return accounts
    },

    unallocatedMoney() {
      const acc = this.accounts.filter(acc => acc.id === unallocatedMoneyAcc)[0] || null
      if (acc === null) {
        return 0
      }
      return acc.balance
    },

    unallocatedMoneyClass() {
      const classes = ['d-inline-flex', 'flex-column', 'text-center', 'ms-auto', 'p-2', 'rounded']
      if (this.unallocatedMoney < 0) {
        classes.push('bg-danger')
      } else if (this.unallocatedMoney === 0) {
        classes.push('bg-warning')
      } else {
        classes.push('bg-success')
      }

      return classes.join(' ')
    },
  },

  data() {
    return {
      timeRange: {},
      transactions: [],
    }
  },

  methods: {
    fetchTransactions() {
      const since = this.timeRange.start.toISOString()
      const until = this.timeRange.end.toISOString()

      return fetch(`/api/transactions?since=${since}&until=${until}`)
        .then(resp => resp.json())
        .then(txs => {
          this.transactions = txs
        })
    },

    formatNumber,
  },

  name: 'AccountingAppBudgetDashboard',

  props: {
    accounts: {
      required: true,
      type: Array,
    },
  },

  watch: {
    timeRange() {
      this.fetchTransactions()
    },
  },
}
</script>
