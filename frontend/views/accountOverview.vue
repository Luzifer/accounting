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
            @click.prevent="editAccount(accountId)"
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
              :disabled="accountIsCategory"
              @click="showAddTransaction = !showAddTransaction"
              @shortkey="showAddTransaction = !showAddTransaction"
            >
              <i class="fas fa-fw fa-plus-circle mr-1" />
              Add Transaction
            </button>
            <button
              class="btn btn-primary"
              :disabled="accountIsCategory"
              data-bs-toggle="modal"
              data-bs-target="#transferMoneyModal"
            >
              <i class="fas fa-fw fa-arrow-right-arrow-left mr-1" />
              Add Transfer
            </button>

            <button
              class="btn btn-success"
              :disabled="accountIsCategory"
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
                  class="dropdown-item"
                  :disabled="selectedTx.length < 1"
                  @click="clearToday"
                >
                  <i class="fas fa-fw fa-copyright mr-1" />
                  Clear Today
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
                <th v-if="accountIsCategory">
                  Account
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
                @edit-cancelled="showAddTransaction = false"
                @edit-saved="txSaved"
              />

              <template
                v-for="tx in sortedTransactions"
              >
                <tr
                  v-if="tx.id !== editedTxId"
                  :key="`display-${tx.id}`"
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
                  <td v-if="accountIsCategory">
                    {{ accountIdToName[tx.account!] }}
                  </td>
                  <td>{{ tx.payee }}</td>
                  <td v-if="account.type !== 'tracking'">
                    {{ accountIdToName[tx.category!] }}
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
                  :key="`editor-${tx.id}`"
                  :account="account"
                  :accounts="accounts"
                  :edit="tx"
                  @edit-cancelled="editedTxId = null"
                  @edit-saved="txSaved"
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
      <div class="modal-dialog">
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
                for="transferMoneyModalDescription"
                class="form-label"
              >Description</label>
              <input
                id="transferMoneyModalDescription"
                v-model.number="modals.createTransfer.description"
                type="text"
                class="form-control"
              >
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

<script lang="ts">

import { defineComponent, type PropType } from 'vue'
import { Modal } from 'bootstrap'

import accountEditor from '../components/accountEditor.vue'
import modalHost from '../components/modal.vue'
import { classFromNumber, formatNumber, responseToJSON } from '../helpers'
import rangeSelector from '../components/rangeSelector.vue'
import txEditor from '../components/txEditor.vue'
import type { Account, DateRange, JsonPatchOperation, Transaction } from '../types'

interface OverviewTransferForm {
  amount: number
  category: string
  description: number | string
  from: string
  to: string
}

const emptyAccount: Account = {
  balance: 0,
  hidden: false,
  id: '',
  name: '',
  type: 'budget',
}

export default defineComponent({
  components: { rangeSelector, txEditor },

  computed: {
    account(): Account {
      return this.accounts.filter(acc => acc.id === this.accountId)[0] || emptyAccount
    },

    accountIdToName(): Record<string, string> {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.name]))
    },

    accountIsCategory(): boolean {
      return this.account.type === 'category'
    },

    accountTypes(): Record<string, string> {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.type]))
    },

    balanceUncleared(): number {
      return this.transactions
        .filter(tx => !tx.cleared)
        .reduce((sum, tx) => sum + tx.amount, 0)
    },

    categories(): Account[] {
      const cats = this.accounts.filter(acc => acc.type === 'category')
      cats.sort((a, b) => a.name.localeCompare(b.name))
      return cats
    },

    selectedTx(): string[] {
      return Object.entries(this.selectedTxRaw)
        .filter(e => e[1])
        .map(e => e[0])
    },

    sortedTransactions(): Transaction[] {
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

    transferCategoryActive(): boolean {
      return Boolean(this.modals.createTransfer.from
        && this.modals.createTransfer.to
        && this.accountTypes[this.modals.createTransfer.from] !== this.accountTypes[this.modals.createTransfer.to])
    },

    transferModalValid(): boolean {
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

    transferableAccounts(): Account[] {
      const accs = (this.accounts || []).filter(acc => acc.type !== 'category')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },
  },

  data() {
    return {
      editedTxId: null as null | string,

      modals: {
        createTransfer: {
          amount: 0,
          category: '',
          description: '',
          from: '',
          to: '',
        } as OverviewTransferForm,
      },

      selectedTxRaw: {} as Record<string, boolean>,
      showAddTransaction: false,
      timeRange: {} as DateRange,
      transactions: [] as Transaction[],
    }
  },

  emits: ['update-accounts'],

  methods: {
    classFromNumber,

    clearToday() {
      return this.patchSelected([
        {
          op: 'replace',
          path: '/time',
          value: new Date(new Date().toISOString()
            .split('T')[0]),
        },
        {
          op: 'replace',
          path: '/cleared',
          value: true,
        },
      ])
    },

    async deleteSelected() {
      const actions = []
      for (const id of this.selectedTx) {
        actions.push(fetch(`/api/transactions/${id}`, {
          method: 'DELETE',
        }))
      }

      await Promise.all(actions)
      this.$emit('update-accounts')
      await this.fetchTransactions()
    },

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

    editSelected() {
      this.editTx(this.selectedTx[0])
    },

    editTx(txId: string) {
      this.editedTxId = txId
    },

    async fetchTransactions() {
      if (!this.timeRange.start || !this.timeRange.end) {
        return
      }

      const since = this.timeRange.start.toISOString()
      const until = this.timeRange.end.toISOString()

      const resp = await fetch(`/api/accounts/${this.accountId}/transactions?since=${since}&until=${until}`)
      const txs = await responseToJSON<Transaction[]>(resp)
      this.transactions = txs ?? []
      this.selectedTxRaw = {}
    },

    formatNumber,

    async markAccountReconciled() {
      await fetch(`/api/accounts/${this.accountId}/reconcile`, {
        method: 'PUT',
      })
      await this.fetchTransactions()
    },

    async markCleared(txId: string, cleared: boolean) {
      await fetch(`/api/transactions/${txId}?cleared=${cleared}`, {
        method: 'PATCH',
      })
      await this.fetchTransactions()
    },

    moveToToday() {
      return this.patchSelected([
        {
          op: 'replace',
          path: '/time',
          value: new Date(new Date().toISOString()
            .split('T')[0]),
        },
      ])
    },

    async patchSelected(patchset: JsonPatchOperation[] = []) {
      const actions = []
      for (const id of this.selectedTx) {
        actions.push(fetch(`/api/transactions/${id}`, {
          body: JSON.stringify(patchset),
          method: 'PATCH',
        }))
      }

      await Promise.all(actions)
      this.$emit('update-accounts')
      await this.fetchTransactions()
    },

    async transferMoney() {
      const params = new URLSearchParams()
      params.set('amount', this.modals.createTransfer.amount.toFixed(2))
      if (this.modals.createTransfer.category) {
        params.set('category', this.modals.createTransfer.category)
      }
      if (this.modals.createTransfer.description) {
        params.set('description', String(this.modals.createTransfer.description))
      }

      await fetch(`/api/accounts/${this.modals.createTransfer.from}/transfer/${this.modals.createTransfer.to}?${params.toString()}`, {
        method: 'PUT',
      })

      this.$emit('update-accounts')
      await this.fetchTransactions()
      const modal = Modal.getInstance(this.$refs.transferMoneyModal as Element)
      modal?.toggle()

      this.modals.createTransfer = {
        amount: 0,
        category: '',
        description: '',
        from: '',
        to: this.accountId,
      }
    },

    async txSaved() {
      this.editedTxId = null
      this.$emit('update-accounts')
      await this.fetchTransactions()
    },

    updateSelectAll(evt: Event) {
      const checked = (evt.target as HTMLInputElement).checked
      this.selectedTxRaw = Object.fromEntries(this.transactions.map(tx => [tx.id, checked]))
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
      type: Array as PropType<Account[]>,
    },
  },

  watch: {
    accountId(to: string) {
      this.modals.createTransfer.to = to
      void this.fetchTransactions()
    },

    timeRange() {
      void this.fetchTransactions()
    },
  },

})
</script>
