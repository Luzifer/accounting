<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col">
        <!-- TODO: Add time-selector -->
      </div>
      <div class="col d-flex align-items-center justify-content-end">
        <div :class="unallocatedMoneyClass">
          <span class="fs-4">{{ unallocatedMoney.toFixed(2) }} €</span>
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
              <td>&nbsp;</td>
              <td :class="{'text-end': true, 'text-danger': cat.balance < 0}">
                {{ cat.balance.toFixed(2) }} €
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { unallocatedMoneyAcc } from '../constants'

export default {
  computed: {
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

  name: 'AccountingAppBudgetDashboard',

  props: {
    accounts: {
      required: true,
      type: Array,
    },
  },
}
</script>
