<template>
  <div>
    <div class="container-fluid">
      <div class="row">
        <div class="col">
          <range-selector
            v-model="timeRange"
          />
        </div>
        <div class="col fs-4 text-center">
          {{ accountIdToName[accountId] }}
        </div>
        <div class="col d-flex align-items-center justify-content-end">
          <div class="btn-group btn-group-sm">
            <button
              class="btn"
              @click="showAddTransaction = !showAddTransaction"
            >
              <i class="fas fa-fw fa-plus-circle mr-1" />
              Add Transaction
            </button>
            <button
              class="btn"
              data-bs-toggle="modal"
              data-bs-target="#transferMoneyModal"
            >
              <i class="fas fa-fw fa-arrow-right-arrow-left mr-1" />
              Add Transfer
            </button>
            <button
              class="btn text-danger"
              :disabled="selectedTx.length < 1"
              @click="deleteSelected"
            >
              <i class="fas fa-fw fa-trash mr-1" />
              Delete Selected
            </button>
          </div>
        </div>
      </div>
      <div class="row mt-3">
        <div class="col">
          <table class="table table-striped small">
            <thead>
              <tr>
                <th class="minimized-column">
                  <input
                    type="checkbox"
                    @change="updateSelectAll"
                  >
                </th>
                <th class="minimized-column">
                  Date
                </th>
                <th>Payee</th>
                <th>Category</th>
                <th>Description</th>
                <th class="minimized-column text-end">
                  Amount
                </th>
                <th class="minimized-column">
                  <i class="fas fa-copyright" />
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="showAddTransaction">
                <td />
                <td>
                  <input
                    v-model="form.date"
                    class="form-control form-control-sm"
                    type="date"
                  >
                </td>
                <td>
                  <input
                    v-model="form.payee"
                    class="form-control form-control-sm"
                    type="text"
                  >
                </td>
                <td>
                  <select
                    v-model="form.category"
                    class="form-select form-select-sm"
                  >
                    <option
                      v-for="cat in categories"
                      :key="cat.id"
                      :value="cat.id"
                    >
                      {{ cat.name }}
                    </option>
                  </select>
                </td>
                <td>
                  <input
                    v-model="form.description"
                    class="form-control form-control-sm"
                    type="text"
                  >
                </td>
                <td>
                  <input
                    v-model="form.amount"
                    class="form-control form-control-sm"
                    type="number"
                    step="0.01"
                    @keypress.enter="createTransaction"
                  >
                </td>
                <td class="align-middle">
                  <input
                    v-model="form.cleared"
                    type="checkbox"
                    @keypress.enter="createTransaction"
                  >
                </td>
              </tr>
              <tr
                v-for="tx in sortedTransactions"
                :key="tx.id"
              >
                <td>
                  <input
                    v-model="selectedTxRaw[tx.id]"
                    type="checkbox"
                  >
                </td>
                <td class="minimized-column">
                  {{ new Date(tx.time).toLocaleDateString() }}
                </td>
                <td>{{ tx.payee }}</td>
                <td>{{ accountIdToName[tx.category] }}</td>
                <td>{{ tx.description }}</td>
                <td :class="{'minimized-column text-end': true, 'text-danger': tx.amount < 0}">
                  {{ formatNumber(tx.amount) }} €
                </td>
                <td>
                  <a
                    href="#"
                    :class="{'text-decoration-none':true, 'text-muted': !tx.cleared, 'text-success': tx.cleared}"
                    @click.prevent="markCleared(tx.id, !tx.cleared)"
                  >
                    <i class="fas fa-copyright" />
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
                  v-for="acc in transferableAccounts"
                  :key="acc.id"
                  :value="acc.id"
                >
                  {{ acc.name }}
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
                  v-for="acc in transferableAccounts"
                  :key="acc.id"
                  :value="acc.id"
                >
                  {{ acc.name }}
                </option>
              </select>
            </div>

            <div class="mb-3">
              <label
                for="transferMoneyModalCategory"
                class="form-label"
              >Category</label>
              <select
                id="transferMoneyModalCategory"
                v-model="modals.createTransfer.category"
                class="form-select"
                :disabled="!transferCategoryActive"
              >
                <option
                  v-for="cat in categories"
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

import { formatNumber, responseToJSON } from '../helpers'
import rangeSelector from './rangeSelector.vue'

export default {
  components: { rangeSelector },

  computed: {
    accountIdToName() {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.name]))
    },

    accountTypes() {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.type]))
    },

    categories() {
      const cats = this.accounts.filter(acc => acc.type === 'category')
      cats.sort((a, b) => a.name.localeCompare(b.name))
      return cats
    },

    selectedTx() {
      return Object.entries(this.selectedTxRaw)
        .filter(e => e[1])
        .map(e => e[0])
    },

    sortedTransactions() {
      const tx = [...this.transactions]
      tx.sort((b, a) => new Date(a.time).getTime() - new Date(b.time).getTime())
      return tx
    },

    transferCategoryActive() {
      return this.modals.createTransfer.from && this.modals.createTransfer.to && this.accountTypes[this.modals.createTransfer.from] !== this.accountTypes[this.modals.createTransfer.to]
    },

    transferModalValid() {
      if (!this.modals.createTransfer.from || !this.modals.createTransfer.to) {
        return false
      }

      if (this.modals.createTransfer.from === this.modals.createTransfer.to) {
        return false
      }

      const tFrom = this.accountTypes[this.modals.createTransfer.from]
      const tTo = this.accountTypes[this.modals.createTransfer.to]
      if (tFrom !== tTo && !this.modals.createTransfer.category) {
        return false
      }

      if (tFrom === tTo && this.modals.createTransfer.category) {
        return false
      }

      return true
    },

    transferableAccounts() {
      const accs = (this.accounts || []).filter(acc => acc.type !== 'category')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },
  },

  data() {
    return {
      form: {
        amount: 0,
        category: '',
        cleared: false,
        date: new Date().toISOString()
          .split('T')[0],

        description: '',
        payee: '',
      },

      modals: {
        createTransfer: {
          amount: 0,
          category: '',
          from: '',
          to: '',
        },
      },

      selectedTxRaw: {},
      showAddTransaction: false,
      timeRange: {},
      transactions: [],
    }
  },

  methods: {
    createTransaction() {
      return fetch('/api/transactions', {
        body: JSON.stringify({
          ...this.form,
          account: this.accountId,
          time: new Date(this.form.date),
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then(responseToJSON)
        .then(() => {
          this.$emit('update-accounts')
          this.fetchTransactions()

          this.showAddTransaction = false
          this.form = {
            amount: 0,
            category: '',
            cleared: false,
            date: this.form.date,
            description: '',
            payee: '',
          }
        })
    },

    deleteSelected() {
      const actions = []
      for (const id of this.selectedTx) {
        actions.push(fetch(`/api/transactions/${id}`, {
          method: 'DELETE',
        }))
      }

      Promise.all(actions)
        .then(() => {
          this.$emit('update-accounts')
          this.fetchTransactions()
        })
    },

    fetchTransactions() {
      const since = this.timeRange.start.toISOString()
      const until = this.timeRange.end.toISOString()

      return fetch(`/api/accounts/${this.accountId}/transactions?since=${since}&until=${until}`)
        .then(resp => resp.json())
        .then(txs => {
          this.transactions = txs
          this.selectedTxRaw = {}
        })
    },

    formatNumber,

    markCleared(txId, cleared) {
      return fetch(`/api/transactions/${txId}?cleared=${cleared}`, {
        method: 'PATCH',
      })
        .then(() => this.fetchTransactions())
    },

    transferMoney() {
      const params = new URLSearchParams()
      params.set('amount', this.modals.createTransfer.amount.toFixed(2))
      if (this.modals.createTransfer.category) {
        params.set('category', this.modals.createTransfer.category)
      }

      return fetch(`/api/accounts/${this.modals.createTransfer.from}/transfer/${this.modals.createTransfer.to}?${params.toString()}`, {
        method: 'PUT',
      })
        .then(() => {
          this.$emit('update-accounts')
          Modal.getInstance(this.$refs.transferMoneyModal).toggle()

          this.modals.createTransfer = {
            amount: 0,
            category: '',
            from: '',
            to: this.accountId,
          }
        })
    },

    updateSelectAll(evt) {
      this.selectedTxRaw = Object.fromEntries(this.transactions.map(tx => [tx.id, evt.target.checked]))
    },
  },

  mounted() {
    this.modals.createTransfer.to = this.accountId
    this.fetchTransactions()
  },

  name: 'AccountingAppAccountOverview',

  props: {
    accountId: {
      required: true,
      type: String,
    },

    accounts: {
      required: true,
      type: Array,
    },
  },

  watch: {
    accountId(to) {
      this.modals.createTransfer.to = to
      this.fetchTransactions()
    },

    timeRange() {
      this.fetchTransactions()
    },
  },

}
</script>
