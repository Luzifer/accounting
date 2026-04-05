<template>
  <div
    id="transferBudgetCategoryMoneyModal"
    ref="transferBudgetCategoryMoneyModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="transferBudgetCategoryMoneyModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h1
            id="transferBudgetCategoryMoneyModalLabel"
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
              for="transferBudgetCategoryMoneyModalFrom"
              class="form-label"
            >From</label>
            <select
              id="transferBudgetCategoryMoneyModalFrom"
              v-model="form.from"
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
              for="transferBudgetCategoryMoneyModalTo"
              class="form-label"
            >To</label>
            <select
              id="transferBudgetCategoryMoneyModalTo"
              v-model="form.to"
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
              for="transferBudgetCategoryMoneyModalAmount"
              class="form-label"
            >Amount</label>
            <div class="input-group">
              <input
                id="transferBudgetCategoryMoneyModalAmount"
                v-model.number="form.amount"
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
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { Modal } from 'bootstrap'

import type { Account } from '../types'

interface TransferBudgetCategoryMoneyForm {
  amount: number
  from: string
  to: string
}

export default defineComponent({
  computed: {
    transferModalValid(): boolean {
      if (!this.form.from || !this.form.to) {
        return false
      }

      if (this.form.from === this.form.to) {
        return false
      }

      return true
    },

    transferableCategories(): Account[] {
      const accounts = this.accounts
        .filter(acc => acc.type === 'category')
      accounts.sort((a, b) => a.name.localeCompare(b.name))
      return accounts
    },
  },

  created() {
    this.resetForm()
  },

  data() {
    return {
      closeReason: 'reject' as 'reject' | 'resolve',
      form: {
        amount: 0,
        from: '',
        to: '',
      } as TransferBudgetCategoryMoneyForm,
      modal: null as Modal | null,
    }
  },

  emits: ['reject', 'resolve'],

  methods: {
    resetForm() {
      this.form = {
        amount: 0,
        from: this.initialFrom,
        to: this.initialTo,
      }
    },

    async transferMoney() {
      const params = new URLSearchParams()
      params.set('amount', this.form.amount.toFixed(2))

      await fetch(`/api/accounts/${this.form.from}/transfer/${this.form.to}?${params.toString()}`, {
        method: 'PUT',
      })

      this.closeReason = 'resolve'
      this.modal?.hide()
    },
  },

  mounted() {
    const modalElement = this.$refs.transferBudgetCategoryMoneyModal as HTMLElement
    this.modal = Modal.getOrCreateInstance(modalElement)
    modalElement.addEventListener('hidden.bs.modal', () => this.$emit(this.closeReason))
    this.modal.show()
  },

  name: 'AccountingAppTransferBudgetCategoryMoneyModal',

  props: {
    accounts: {
      required: true,
      type: Array as PropType<Account[]>,
    },

    initialFrom: {
      required: true,
      type: String,
    },

    initialTo: {
      required: true,
      type: String,
    },
  },
})
</script>
