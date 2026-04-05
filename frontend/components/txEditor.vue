<template>
  <tr>
    <td />
    <td>
      <input
        ref="date"
        v-model="form.date"
        class="form-control form-control-sm"
        type="date"
        @keyup.esc="sendCancel"
        @keyup.enter="focusRef('payee')"
      >
    </td>
    <td v-if="accountIsCategory">
      {{ accountIdToName[edit?.account!] }}
    </td>
    <td>
      <input
        ref="payee"
        v-model="form.payee"
        class="form-control form-control-sm"
        type="text"
        @keyup.esc="sendCancel"
        @keyup.enter="focusRef('category')"
      >
    </td>
    <td v-if="account.type !== 'tracking'">
      <select
        ref="category"
        v-model="form.category"
        class="form-select form-select-sm"
        @keyup.esc="sendCancel"
        @keydown.enter.prevent="() => {}"
        @keyup.enter.prevent="focusRef('description')"
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
        ref="description"
        v-model="form.description"
        class="form-control form-control-sm"
        type="text"
        @keyup.esc="sendCancel"
        @keyup.enter="focusRef('amount')"
      >
    </td>
    <td>
      <input
        ref="amount"
        v-model="form.amount"
        class="form-control form-control-sm"
        type="number"
        step="0.01"
        @keypress.enter="saveTransaction"
        @keyup.esc="sendCancel"
      >
    </td>
    <td class="align-middle">
      <input
        v-model="form.cleared"
        type="checkbox"
        @keypress.enter="saveTransaction"
        @keyup.esc="sendCancel"
      >
    </td>
  </tr>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'

import type { Account, Transaction } from '../types'
import { requestAPI } from '../helpers'

interface TransactionForm {
  amount: number | string
  category: string
  cleared: boolean
  date: string
  description: string
  payee: string
}

export default defineComponent({
  computed: {
    accountIdToName(): Record<string, string> {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.name]))
    },

    accountIsCategory(): boolean {
      return this.account.type === 'category'
    },

    categories(): Account[] {
      const cats = this.accounts.filter(acc => acc.type === 'category')
      cats.sort((a, b) => a.name.localeCompare(b.name))
      return cats
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
      } as TransactionForm,
    }
  },

  emits: ['editCancelled', 'editSaved'],

  methods: {
    focusRef(refName: 'amount' | 'category' | 'date' | 'description' | 'payee') {
      (this.$refs[refName] as HTMLInputElement | HTMLSelectElement | undefined)?.focus()
    },

    async saveTransaction() {
      const body = {
        ...this.form,
        account: this.account.id,
        category: this.form.category || null,
        time: new Date(this.form.date),
      }

      if (this.edit?.id) {
        await requestAPI('PUT', `/api/transactions/${this.edit.id}`, body)
      } else {
        await requestAPI('POST', '/api/transactions', body)
      }
      this.$emit('editSaved')

      if (!this.edit) {
        this.form = {
          amount: 0,
          category: '',
          cleared: false,
          date: new Date().toISOString()
            .split('T')[0],

          description: '',
          payee: '',
        }
        this.focusRef('date')
      }
    },

    sendCancel() {
      this.$emit('editCancelled')
    },
  },

  mounted() {
    if (this.edit) {
      this.form = {
        ...this.edit,
        category: this.edit.category ?? '',
        date: new Date(this.edit.time).toISOString()
          .split('T')[0],
      }
    }

    this.focusRef('date')
  },

  name: 'AccountingAppTXEditor',

  props: {
    account: {
      required: true,
      type: Object as PropType<Account>,
    },

    accounts: {
      required: true,
      type: Array as PropType<Account[]>,
    },

    edit: {
      default: null,
      type: Object as PropType<Transaction | null>,
    },
  },
})
</script>
