<template>
  <div>
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
              data-bs-toggle="modal"
              data-bs-target="#transferMoneyModal"
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
                <td>{{ cat.name }}</td>
                <td :class="{'text-end': true, 'text-danger': (allocatedByCategory[cat.id] || 0) < 0}">
                  {{ formatNumber(allocatedByCategory[cat.id] || 0) }} €
                </td>
                <td :class="{'text-end': true, 'text-danger': (activityByCategory[cat.id] || 0) < 0}">
                  {{ formatNumber(activityByCategory[cat.id] || 0) }} €
                </td>
                <td class="text-end">
                  <a
                    :class="{'text-decoration-none': true, 'text-danger': cat.balance < 0, 'text-white': cat.balance >= 0}"
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

    <div
      id="transferMoneyModal"
      ref="transferMoneyModal"
      class="modal fade"
      tabindex="-1"
      aria-labelledby="transferMoneyModalLabel"
      aria-hidden="true"
    >
      <div class="modal-dialog modal-sm">
        <div class="modal-content">
          <div class="modal-header">
            <h1
              id="transferMoneyModalLabel"
              class="modal-title fs-5"
            >
              Transfer Money
            </h1>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              aria-label="Close"
            />
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label
                for="transferMoneyModalFrom"
                class="form-label"
              >From</label>
              <select
                id="transferMoneyModalFrom"
                v-model="modals.createTransfer.from"
                class="form-select"
              >
                <option
                  v-for="cat in transferableCategories"
                  :key="cat.id"
                  :value="cat.id"
                >
                  {{ cat.name }}
                </option>
              </select>
            </div>

            <div class="mb-3">
              <label
                for="transferMoneyModalTo"
                class="form-label"
              >To</label>
              <select
                id="transferMoneyModalTo"
                v-model="modals.createTransfer.to"
                class="form-select"
              >
                <option
                  v-for="cat in transferableCategories"
                  :key="cat.id"
                  :value="cat.id"
                >
                  {{ cat.name }}
                </option>
              </select>
            </div>

            <div class="mb-3">
              <label
                for="transferMoneyModalAmount"
                class="form-label"
              >Amount</label>
              <div class="input-group">
                <input
                  id="transferMoneyModalAmount"
                  v-model.number="modals.createTransfer.amount"
                  type="number"
                  min="0.01"
                  step="0.01"
                  class="form-control text-end"
                >
                <span class="input-group-text">€</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-primary"
              :disabled="!transferModalValid"
              @click="transferMoney"
            >
              <i class="fas fa-fw fa-arrow-right-arrow-left mr-1" />
              Transfer
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
/* eslint-disable sort-imports */
import { Modal } from 'bootstrap'

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

    transferModalValid() {
      if (!this.modals.createTransfer.from || !this.modals.createTransfer.to) {
        return false
      }

      if (this.modals.createTransfer.from === this.modals.createTransfer.to) {
        return false
      }

      return true
    },

    transferableCategories() {
      const accounts = this.accounts
        .filter(acc => acc.type === 'category')
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
      modals: {
        createTransfer: {
          amount: 0,
          from: unallocatedMoneyAcc,
          to: '',
        },
      },

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

    initTransfer(catId) {
      const cat = this.categories.filter(cat => cat.id === catId)[0] || null
      if (!cat) {
        console.error(`init transfer from non existent category ${catId}`)
        return
      }

      this.modals.createTransfer = {
        amount: 0,
        from: '',
        to: '',
      }

      if (cat.balance <= 0) {
        // Most likely scenario: Transfer in
        this.modals.createTransfer.from = unallocatedMoneyAcc
        this.modals.createTransfer.to = cat.id
      } else {
        // Most likely scenario: Transfer out
        this.modals.createTransfer.from = cat.id
      }

      Modal.getOrCreateInstance(this.$refs.transferMoneyModal).toggle()
    },

    resetTransfer() {
      this.modals.createTransfer = {
        amount: 0,
        from: unallocatedMoneyAcc,
        to: '',
      }
    },

    transferMoney() {
      const params = new URLSearchParams()
      params.set('amount', this.modals.createTransfer.amount.toFixed(2))

      return fetch(`/api/accounts/${this.modals.createTransfer.from}/transfer/${this.modals.createTransfer.to}?${params.toString()}`, {
        method: 'PUT',
      })
        .then(() => {
          this.$emit('update-accounts')
          this.fetchTransactions()
          Modal.getInstance(this.$refs.transferMoneyModal).toggle()
        })
    },
  },

  mounted() {
    this.$refs.transferMoneyModal
      .addEventListener('hidden.bs.modal', () => this.resetTransfer())
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
