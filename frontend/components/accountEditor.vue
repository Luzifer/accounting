<template>
  <div
    id="editAccountModal"
    ref="editAccountModal"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="editAccountModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog modal-sm">
      <div class="modal-content">
        <div class="modal-header">
          <h1
            id="editAccountModalLabel"
            class="modal-title fs-5"
          >
            Edit Account
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
              for="editAccountModalName"
              class="form-label"
            >Name</label>
            <input
              id="editAccountModalName"
              v-model="form.name"
              class="form-control"
              type="text"
              @keypress.enter="updateAccount"
            >
          </div>

          <div class="mb-3">
            <div class="form-check">
              <input
                id="editAccountModalHidden"
                v-model="form.hidden"
                class="form-check-input"
                type="checkbox"
                :disabled="form.balance !== 0"
              >
              <label
                for="editAccountModalHidden"
                class="form-check-label"
              >Hide account</label>
            </div>
            <div class="form-text">
              Note: You cannot un-hide the account using UI, you need to unhide it via API.
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-primary"
            @click="updateAccount"
          >
            <i class="fas fa-fw fa-pencil mr-1" />
            Update
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'
import { Modal } from 'bootstrap'

import { responseToJSON } from '../helpers'
import type { Account } from '../types'

interface AccountEditForm {
  balance: number
  hidden: boolean
  name: string
}

export default defineComponent({
  data() {
    return {
      form: {
        balance: 0,
        hidden: false,
        name: '',
      } as AccountEditForm,
    }
  },

  emits: ['editClosed', 'editComplete'],

  methods: {
    async updateAccount() {
      if (this.account === null) {
        return
      }

      const update = new URLSearchParams()

      if (this.form.name !== this.account.name) {
        update.set('name', this.form.name)
      }

      if (this.form.hidden !== this.account.hidden) {
        update.set('hidden', String(this.form.hidden))
      }

      if (update.toString().length === 0) {
        // No updates, why are we here?
        return
      }

      const resp = await fetch(`/api/accounts/${this.account.id}?${update.toString()}`, {
        method: 'PATCH',
      })

      await responseToJSON(resp)
      this.$emit('editComplete')
    },
  },

  mounted() {
    const modalElement = this.$refs.editAccountModal as HTMLElement
    modalElement.addEventListener('hidden.bs.modal', () => this.$emit('editClosed'))
  },

  name: 'AccountingAppAccountEditor',

  props: {
    account: {
      default: null,
      required: false,
      type: Object as PropType<Account | null>,
    },
  },

  watch: {
    account(to: Account | null) {
      const modal = Modal.getOrCreateInstance(this.$refs.editAccountModal as Element)
      if (!to) {
        modal.hide()
        return
      }

      this.form = {
        balance: to.balance,
        hidden: to.hidden,
        name: to.name,
      }

      modal.show()
    },
  },
})
</script>
