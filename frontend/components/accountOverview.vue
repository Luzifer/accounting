<template>
  <div>
    <div class="container-fluid">
      <div class="row">
        <div
          class="col d-flex fs-3 align-items-center text-semibold"
        >
          <a
            class="text-white text-decoration-none"
            href="#"
            title="Click to Edit"
            @click.prevent="editedAccId = accountId"
          >
            {{ accountIdToName[accountId] }}
          </a>
        </div>
        <div class="col d-flex align-items-center justify-content-end">
          <range-selector
            v-model="timeRange"
            store-key="account-overview"
          />
        </div>
      </div>

      <div class="row mt-3">
        <div class="col d-flex align-items-center">
          <div class="d-flex align-items-start">
            <div class="d-inline-flex text-center flex-column me-4">
              <span :class="classFromNumber(account.balance - balanceUncleared)">
                {{ formatNumber(account.balance - balanceUncleared) }} €
              </span>
              <span class="form-text mt-0">
                <i class="fas fa-fw fa-copyright mr-1 text-success" />
                Cleared Balance
              </span>
            </div>
            +
            <div class="d-inline-flex text-center flex-column mx-4">
              <span :class="classFromNumber(balanceUncleared)">
                {{ formatNumber(balanceUncleared) }} €
              </span>
              <span class="form-text mt-0">
                <i class="fas fa-fw fa-copyright mr-1" />
                Uncleared Balance
              </span>
            </div>
            =
            <div class="d-inline-flex text-center flex-column ms-4">
              <span :class="classFromNumber(account.balance)">
                {{ formatNumber(account.balance) }} €
              </span>
              <span class="form-text mt-0">Working Balance</span>
            </div>
          </div>
        </div>
        <div class="col d-flex align-items-center justify-content-end">
          <div class="btn-group btn-group-sm">
            <button
              v-shortkey.once="['ctrl', 'alt', 'n']"
              class="btn btn-primary"
              @click="showAddTransaction = !showAddTransaction"
              @shortkey="showAddTransaction = !showAddTransaction"
            >
              <i class="fas fa-fw fa-plus-circle mr-1" />
              Add Transaction
            </button>
            <button
              class="btn btn-primary"
              data-bs-toggle="modal"
              data-bs-target="#transferMoneyModal"
            >
              <i class="fas fa-fw fa-arrow-right-arrow-left mr-1" />
              Add Transfer
            </button>

            <button
              class="btn btn-success"
              @click="markAccountReconciled"
            >
              <i class="fas fa-fw fa-square-check mr-1" />
              Mark Reconciled
            </button>

            <button
              class="btn btn-secondary"
              :disabled="selectedTx.length !== 1"
              @click="editSelected"
            >
              <i class="fas fa-fw fa-pencil mr-1" />
              Edit Transaction
            </button>

            <button
              type="button"
              class="btn btn-secondary dropdown-toggle dropdown-toggle-split"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <span class="visually-hidden">Toggle Dropdown</span>
            </button>
            <ul class="dropdown-menu">
              <li>
                <button
                  class="dropdown-item"
                  :disabled="selectedTx.length < 1"
                  @click="moveToToday"
                >
                  <i class="fas fa-fw fa-calendar-days mr-1" />
                  Move to Today
                </button>
              </li>
              <li>
                <button
                  v-shortkey="['del']"
                  class="dropdown-item text-danger"
                  :disabled="selectedTx.length < 1"
                  @click="deleteSelected"
                  @shortkey="deleteSelected"
                >
                  <i class="fas fa-fw fa-trash mr-1" />
                  Delete Selected
                </button>
              </li>
            </ul>
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
                <th v-if="account.type !== 'tracking'">
                  Category
                </th>
                <th>Description</th>
                <th class="minimized-amount text-end">
                  Amount
                </th>
                <th class="minimized-column">
                  <i class="fas fa-copyright" />
                </th>
              </tr>
            </thead>
            <tbody>
              <tx-editor
                v-if="showAddTransaction"
                :account="account"
                :accounts="accounts"
                @editCancelled="showAddTransaction = false"
                @editSaved="txSaved"
              />

              <template
                v-for="tx in sortedTransactions"
              >
                <tr
                  v-if="tx.id !== editedTxId"
                  :key="tx.id"
                  @dblclick="editTx(tx.id)"
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
                  <td v-if="account.type !== 'tracking'">
                    {{ accountIdToName[tx.category] }}
                  </td>
                  <td>{{ tx.description }}</td>
                  <td :class="classFromNumber(tx.amount, ['minimized-amount', 'text-end'])">
                    {{ formatNumber(tx.amount) }} €
                  </td>
                  <td>
                    <span
                      v-if="tx.reconciled"
                      class="text-success"
                    >
                      <i class="fas fa-lock" />
                    </span>
                    <a
                      v-else
                      href="#"
                      :class="{'text-decoration-none':true, 'text-muted': !tx.cleared, 'text-success': tx.cleared}"
                      @click.prevent="markCleared(tx.id, !tx.cleared)"
                    >
                      <i class="fas fa-copyright" />
                    </a>
                  </td>
                </tr>
                <tx-editor
                  v-else
                  :key="tx.id"
                  :account="account"
                  :accounts="accounts"
                  :edit="tx"
                  @editCancelled="editedTxId = null"
                  @editSaved="txSaved"
                />
              </template>
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

    <account-editor
      :account="editedAcc"
      @editClosed="editedAccId = null"
      @editComplete="$emit('update-accounts'); editedAccId = null"
    />
  </div>
</template>

<script>
/* eslint-disable sort-imports */
import { Modal } from 'bootstrap'

import accountEditor from './accountEditor.vue'
import { classFromNumber, formatNumber } from '../helpers'
import rangeSelector from './rangeSelector.vue'
import txEditor from './txEditor.vue'

export default {
  components: { accountEditor, rangeSelector, txEditor },

  computed: {
    account() {
      return this.accounts.filter(acc => acc.id === this.accountId)[0] || {}
    },

    accountIdToName() {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.name]))
    },

    accountTypes() {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.type]))
    },

    balanceUncleared() {
      return this.transactions
        .filter(tx => !tx.cleared)
        .reduce((sum, tx) => sum + tx.amount, 0)
    },

    categories() {
      const cats = this.accounts.filter(acc => acc.type === 'category')
      cats.sort((a, b) => a.name.localeCompare(b.name))
      return cats
    },

    editedAcc() {
      if (!this.editedAccId) {
        return null
      }
      return this.accounts.filter(acc => acc.id === this.editedAccId)[0] || null
    },

    selectedTx() {
      return Object.entries(this.selectedTxRaw)
        .filter(e => e[1])
        .map(e => e[0])
    },

    sortedTransactions() {
      const tx = [...this.transactions]
      tx.sort((b, a) => {
        if (a.cleared && !b.cleared) {
          return -1
        } else if (!a.cleared && b.cleared) {
          return 1
        }

        return new Date(a.time).getTime() - new Date(b.time).getTime()
      })
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
      editedAccId: null,
      editedTxId: null,

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

  emits: ['update-accounts'],

  methods: {
    classFromNumber,

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

    editSelected() {
      this.editTx(this.selectedTx[0])
    },

    editTx(txId) {
      this.editedTxId = txId
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

    markAccountReconciled() {
      return fetch(`/api/accounts/${this.accountId}/reconcile`, {
        method: 'PUT',
      })
        .then(() => this.fetchTransactions())
    },

    markCleared(txId, cleared) {
      return fetch(`/api/transactions/${txId}?cleared=${cleared}`, {
        method: 'PATCH',
      })
        .then(() => this.fetchTransactions())
    },

    moveToToday() {
      const actions = []
      for (const id of this.selectedTx) {
        actions.push(fetch(`/api/transactions/${id}`, {
          body: JSON.stringify([
            {
              op: 'replace',
              path: '/time',
              value: new Date(new Date().toISOString()
                .split('T')[0]),
            },
          ]),
          method: 'PATCH',
        }))
      }

      Promise.all(actions)
        .then(() => {
          this.$emit('update-accounts')
          this.fetchTransactions()
        })
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
          this.fetchTransactions()
          Modal.getInstance(this.$refs.transferMoneyModal).toggle()

          this.modals.createTransfer = {
            amount: 0,
            category: '',
            from: '',
            to: this.accountId,
          }
        })
    },

    txSaved() {
      this.editedTxId = null
      this.$emit('update-accounts')
      this.fetchTransactions()
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
