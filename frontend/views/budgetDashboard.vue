<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col d-flex align-items-center">
        <range-selector
          v-model="timeRange"
          :multi-month="false"
        />
      </div>
      <div class="col d-flex align-items-center justify-content-center">
        <div :class="unallocatedMoneyClass">
          <span class="fs-4">{{ formatNumber(unallocatedMoney) }} €</span>
          <span class="small">Unallocated</span>
        </div>
      </div>
      <div class="col d-flex align-items-center justify-content-end">
        <div class="btn-group btn-group-sm">
          <button
            class="btn btn-primary"
            @click="() => openTransferModal()"
          >
            <i class="fas fa-fw fa-arrow-right-arrow-left mr-1" />
            Add Transfer
          </button>
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
              <td>
                <a
                  class="text-white text-decoration-none"
                  href="#"
                  title="Click to Edit"
                  @click.prevent="editAccount(cat.id)"
                >
                  {{ cat.name }}
                </a>
              </td>
              <td :class="classFromNumber(allocatedByCategory[cat.id] || 0, ['text-end'])">
                {{ formatNumber(allocatedByCategory[cat.id] || 0) }} €
              </td>
              <td :class="classFromNumber(activityByCategory[cat.id] || 0, ['text-end'])">
                <router-link
                  class="text-white text-decoration-none"
                  :to="{ name: 'account-transactions', params: { accountId: cat.id }}"
                  title="List Transactions"
                >
                  {{ formatNumber(activityByCategory[cat.id] || 0) }} €
                </router-link>
              </td>
              <td class="text-end">
                <a
                  :class="classFromNumber(cat.balance, ['text-decoration-none'], 'text-white')"
                  href="#"
                  title="Transfer Money using this Category"
                  @click.prevent="initTransfer(cat.id)"
                >
                  {{ formatNumber(cat.balance) }} €
                </a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script lang="ts">

import { defineComponent, type PropType } from 'vue'

import accountEditor from '../components/accountEditor.vue'
import modalHost from '../components/modal.vue'
import transferBudgetCategoryMoneyModal from '../components/transferBudgetCategoryMoneyModal.vue'
import { classFromNumber, formatNumber, responseToJSON } from '../helpers'
import rangeSelector from '../components/rangeSelector.vue'
import { unallocatedMoneyAcc } from '../constants'
import type { Account, DateRange, Transaction } from '../types'

export default defineComponent({
  components: { rangeSelector },

  computed: {
    activityByCategory(): Record<string, number> {
      return this.transactions
        .filter(tx => tx.account)
        .reduce((alloc, tx) => {
          if (!tx.category) {
            return alloc
          }

          alloc[tx.category] = (alloc[tx.category] || 0) + tx.amount
          return alloc
        }, {} as Record<string, number>)
    },

    allocatedByCategory(): Record<string, number> {
      return this.transactions
        .filter(tx => !tx.account)
        .reduce((alloc, tx) => {
          if (!tx.category) {
            return alloc
          }

          alloc[tx.category] = (alloc[tx.category] || 0) + tx.amount
          return alloc
        }, {} as Record<string, number>)
    },

    categories(): Account[] {
      const accounts = this.accounts
        .filter(acc => acc.type === 'category')
        .filter(acc => acc.id !== unallocatedMoneyAcc)
      accounts.sort((a, b) => a.name.localeCompare(b.name))
      return accounts
    },

    unallocatedMoney(): number {
      const acc = this.accounts.filter(acc => acc.id === unallocatedMoneyAcc)[0] || null
      if (acc === null) {
        return 0
      }
      return acc.balance
    },

    unallocatedMoneyClass(): string {
      const classes = ['d-inline-flex', 'flex-column', 'text-center', 'p-2', 'rounded']
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
      timeRange: {} as DateRange,
      transactions: [] as Transaction[],
    }
  },

  emits: ['update-accounts'],

  methods: {
    classFromNumber,

    async editAccount(accountId: string) {
      const account = this.accounts.filter(acc => acc.id === accountId)[0] || null
      if (!account) {
        return
      }

      try {
        await modalHost.openModal(accountEditor, { account })
        this.$emit('update-accounts')
      } catch {
        // Dismissed by user.
      }
    },

    async fetchTransactions() {
      if (!this.timeRange.start || !this.timeRange.end) {
        return
      }

      const since = this.timeRange.start.toISOString()
      const until = this.timeRange.end.toISOString()

      const resp = await fetch(`/api/transactions?since=${since}&until=${until}`)
      const txs = await responseToJSON<Transaction[]>(resp)
      this.transactions = txs ?? []
    },

    formatNumber,

    async initTransfer(catId: string) {
      const cat = this.categories.filter(cat => cat.id === catId)[0] || null
      if (!cat) {
        console.error(`init transfer from non existent category ${catId}`)
        return
      }

      let initialFrom = ''
      let initialTo = ''

      if (cat.balance <= 0) {
        // Most likely scenario: Transfer in
        initialFrom = unallocatedMoneyAcc
        initialTo = cat.id
      } else {
        // Most likely scenario: Transfer out
        initialFrom = cat.id
      }

      await this.openTransferModal(initialFrom, initialTo)
    },

    async openTransferModal(initialFrom = unallocatedMoneyAcc, initialTo = '') {
      try {
        await modalHost.openModal(transferBudgetCategoryMoneyModal, {
          accounts: this.accounts,
          initialFrom,
          initialTo,
        })
        this.$emit('update-accounts')
        await this.fetchTransactions()
      } catch {
        // Dismissed by user.
      }
    },
  },

  mounted() {
    void this.fetchTransactions()
  },

  name: 'AccountingAppBudgetDashboard',

  props: {
    accounts: {
      required: true,
      type: Array as PropType<Account[]>,
    },
  },

  watch: {
    timeRange() {
      this.fetchTransactions()
    },
  },
})
</script>
