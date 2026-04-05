<template>
  <div
    id="createAccountModal"
    ref="createAccountModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="createAccountModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h1
            id="createAccountModalLabel"
            class="modal-title fs-5"
          >
            Add Account
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
              for="createAccountModalName"
              class="form-label"
            >Account Name</label>
            <input
              id="createAccountModalName"
              v-model.trim="form.name"
              type="text"
              class="form-control"
              required
            >
          </div>
          <div class="mb-3">
            <label
              for="createAccountModalType"
              class="form-label"
            >Account Type</label>
            <select
              id="createAccountModalType"
              v-model="form.type"
              class="form-select"
            >
              <option value="budget">
                Budget
              </option>
              <option value="tracking">
                Tracking
              </option>
              <option value="category">
                Category
              </option>
            </select>
          </div>
          <div class="mb-3">
            <label
              for="createAccountModalBalance"
              class="form-label"
            >Starting Balance</label>
            <div class="input-group">
              <input
                id="createAccountModalBalance"
                v-model.number="form.startingBalance"
                type="number"
                step="0.01"
                class="form-control"
              >
              <span class="input-group-text">€</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-primary"
            @click="addAccount"
          >
            <i class="fas fa-fw fa-circle-plus mr-1" />
            Add
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { Modal } from 'bootstrap'
import { requestAPI } from '../helpers'

interface AddAccountForm {
  name: string
  startingBalance: number
  type: 'budget' | 'category' | 'tracking'
}

export default defineComponent({
  created() {
    this.resetForm()
  },

  data() {
    return {
      closeReason: 'reject' as 'reject' | 'resolve',
      form: {
        name: '',
        startingBalance: 0,
        type: 'budget',
      } as AddAccountForm,

      modal: null as Modal | null,
    }
  },

  emits: ['reject', 'resolve'],

  methods: {
    async addAccount() {
      await requestAPI('POST', '/api/accounts', {
        name: this.form.name,
        startingBalance: this.form.startingBalance,
        type: this.form.type,
      })

      this.closeReason = 'resolve'
      this.modal?.hide()
    },

    resetForm() {
      this.form = {
        name: '',
        startingBalance: 0,
        type: 'budget',
      }
    },
  },

  mounted() {
    const modalElement = this.$refs.createAccountModal as HTMLElement
    this.modal = Modal.getOrCreateInstance(modalElement)
    modalElement.addEventListener('hidden.bs.modal', () => this.$emit(this.closeReason))
    this.modal.show()
  },

  name: 'AccountingAppCreateAccountModal',
})
</script>
