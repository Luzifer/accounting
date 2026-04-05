<template>
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
            @click="openTransferModal"
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
</template>

<script lang="ts">

import { defineComponent, type PropType } from 'vue'

import accountEditor from '../components/accountEditor.vue'
import modalHost from '../components/modal.vue'
import transferAccountMoneyModal from '../components/transferAccountMoneyModal.vue'
import { classFromNumber, formatNumber, requestAPI } from '../helpers'
import rangeSelector from '../components/rangeSelector.vue'
import txEditor from '../components/txEditor.vue'
import type { Account, DateRange, JsonPatchOperation, Transaction } from '../types'

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

    balanceUncleared(): number {
      return this.transactions
        .filter(tx => !tx.cleared)
        .reduce((sum, tx) => sum + tx.amount, 0)
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
  },

  data() {
    return {
      editedTxId: null as null | string,

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
        actions.push(requestAPI('DELETE', `/api/transactions/${id}`))
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

      const txs = await requestAPI<Transaction[]>('GET', `/api/accounts/${this.accountId}/transactions?since=${since}&until=${until}`)
      this.transactions = txs ?? []
      this.selectedTxRaw = {}
    },

    formatNumber,

    async markAccountReconciled() {
      await requestAPI('PUT', `/api/accounts/${this.accountId}/reconcile`)
      await this.fetchTransactions()
    },

    async markCleared(txId: string, cleared: boolean) {
      await requestAPI('PATCH', `/api/transactions/${txId}?cleared=${cleared}`)
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

    async openTransferModal() {
      try {
        await modalHost.openModal(transferAccountMoneyModal, {
          accountId: this.accountId,
          accounts: this.accounts,
        })
        this.$emit('update-accounts')
        await this.fetchTransactions()
      } catch {
        // Dismissed by user.
      }
    },

    async patchSelected(patchset: JsonPatchOperation[] = []) {
      const actions = []
      for (const id of this.selectedTx) {
        actions.push(requestAPI('PATCH', `/api/transactions/${id}`, patchset))
      }

      await Promise.all(actions)
      this.$emit('update-accounts')
      await this.fetchTransactions()
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
    accountId() {
      void this.fetchTransactions()
    },

    timeRange() {
      void this.fetchTransactions()
    },
  },

})
</script>
