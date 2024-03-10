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
        @keyup.enter="$refs.payee.focus()"
      >
    </td>
    <td v-if="accountIsCategory">
      {{ accountIdToName[edit.account] }}
    </td>
    <td>
      <input
        ref="payee"
        v-model="form.payee"
        class="form-control form-control-sm"
        type="text"
        @keyup.esc="sendCancel"
        @keyup.enter="$refs.category.focus()"
      >
    </td>
    <td v-if="account.type !== 'tracking'">
      <select
        ref="category"
        v-model="form.category"
        class="form-select form-select-sm"
        @keyup.esc="sendCancel"
        @keyup.enter="$refs.description.focus()"
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
        @keyup.enter="$refs.amount.focus()"
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

<script>
import { responseToJSON } from '../helpers'

export default {
  computed: {
    accountIdToName() {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.name]))
    },

    accountIsCategory() {
      return this.account.type === 'category'
    },

    categories() {
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
      },
    }
  },

  emits: ['editCancelled', 'editSaved'],

  methods: {
    saveTransaction() {
      const body = JSON.stringify({
        ...this.form,
        account: this.account.id,
        category: this.form.category || null,
        time: new Date(this.form.date),
      })

      const headers = {
        'Content-Type': 'application/json',
      }

      let prom = null
      if (this.edit?.id) {
        prom = fetch(`/api/transactions/${this.edit.id}`, { body, headers, method: 'PUT' })
      } else {
        prom = fetch('/api/transactions', { body, headers, method: 'POST' })
      }

      return prom
        .then(responseToJSON)
        .then(() => {
          this.$emit('editSaved')
        })
        .then(() => {
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
            this.$refs.date.focus()
          }
        })
    },

    sendCancel() {
      this.$emit('editCancelled')
    },
  },

  mounted() {
    if (this.edit) {
      this.form = {
        ...this.edit,
        date: new Date(this.edit.time).toISOString()
          .split('T')[0],
      }
    }

    this.$refs.date.focus()
  },

  name: 'AccountingAppTXEditor',

  props: {
    account: {
      required: true,
      type: Object,
    },

    accounts: {
      required: true,
      type: Array,
    },

    edit: {
      default: null,
      type: Object,
    },
  },
}
</script>
