<template>
  <div
    id="transferAccountMoneyModal"
    ref="transferAccountMoneyModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="transferAccountMoneyModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h1
            id="transferAccountMoneyModalLabel"
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
              for="transferAccountMoneyModalFrom"
              class="form-label"
            >From</label>
            <select
              id="transferAccountMoneyModalFrom"
              v-model="form.from"
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
              for="transferAccountMoneyModalTo"
              class="form-label"
            >To</label>
            <select
              id="transferAccountMoneyModalTo"
              v-model="form.to"
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
              for="transferAccountMoneyModalCategory"
              class="form-label"
            >Category</label>
            <select
              id="transferAccountMoneyModalCategory"
              v-model="form.category"
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
              for="transferAccountMoneyModalDescription"
              class="form-label"
            >Description</label>
            <input
              id="transferAccountMoneyModalDescription"
              v-model="form.description"
              type="text"
              class="form-control"
            >
          </div>

          <div class="mb-3">
            <label
              for="transferAccountMoneyModalAmount"
              class="form-label"
            >Amount</label>
            <div class="input-group">
              <input
                id="transferAccountMoneyModalAmount"
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

interface TransferAccountMoneyForm {
  amount: number
  category: string
  description: string
  from: string
  to: string
}

export default defineComponent({
  computed: {
    accountTypes(): Record<string, string> {
      return Object.fromEntries(this.accounts.map(acc => [acc.id, acc.type]))
    },

    categories(): Account[] {
      const cats = this.accounts.filter(acc => acc.type === 'category')
      cats.sort((a, b) => a.name.localeCompare(b.name))
      return cats
    },

    transferCategoryActive(): boolean {
      return Boolean(this.form.from
        && this.form.to
        && this.accountTypes[this.form.from] !== this.accountTypes[this.form.to])
    },

    transferModalValid(): boolean {
      if (!this.form.from || !this.form.to) {
        return false
      }

      if (this.form.from === this.form.to) {
        return false
      }

      const tFrom = this.accountTypes[this.form.from]
      const tTo = this.accountTypes[this.form.to]
      if (tFrom !== tTo && !this.form.category) {
        return false
      }

      if (tFrom === tTo && this.form.category) {
        return false
      }

      return true
    },

    transferableAccounts(): Account[] {
      const accs = this.accounts.filter(acc => acc.type !== 'category')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
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
        category: '',
        description: '',
        from: '',
        to: '',
      } as TransferAccountMoneyForm,

      modal: null as Modal | null,
    }
  },

  emits: ['reject', 'resolve'],

  methods: {
    resetForm() {
      this.form = {
        amount: 0,
        category: '',
        description: '',
        from: '',
        to: this.accountId,
      }
    },

    async transferMoney() {
      const params = new URLSearchParams()
      params.set('amount', this.form.amount.toFixed(2))
      if (this.form.category) {
        params.set('category', this.form.category)
      }
      if (this.form.description) {
        params.set('description', this.form.description)
      }

      await fetch(`/api/accounts/${this.form.from}/transfer/${this.form.to}?${params.toString()}`, {
        method: 'PUT',
      })

      this.closeReason = 'resolve'
      this.modal?.hide()
    },
  },

  mounted() {
    const modalElement = this.$refs.transferAccountMoneyModal as HTMLElement
    this.modal = Modal.getOrCreateInstance(modalElement)
    modalElement.addEventListener('hidden.bs.modal', () => this.$emit(this.closeReason))
    this.modal.show()
  },

  name: 'AccountingAppTransferAccountMoneyModal',

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
})
</script>
