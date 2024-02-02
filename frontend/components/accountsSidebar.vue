<template>
  <div>
    <router-link
      to="/"
      class="d-flex align-items-center fs-5 mb-3 mb-md-0 me-md-auto text-start text-white text-decoration-none"
    >
      <i class="fas fa-fw fa-hand-holding-dollar me-2" />
      My Budget
    </router-link>

    <hr>

    <ul class="list-unstyled lh-lg ps-0">
      <acc-list
        :accounts="budgetAccounts"
        header="Budget"
      />
      <acc-list
        :accounts="trackingAccounts"
        header="Tracking"
      />
    </ul>

    <button
      class="btn btn-sm w-100"
      data-bs-toggle="modal"
      data-bs-target="#createAccountModal"
    >
      <i class="fas fa-fw fa-circle-plus mr-1" />
      Add Account
    </button>

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
                v-model.trim="modals.addAccount.name"
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
                v-model="modals.addAccount.type"
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
                  v-model.number="modals.addAccount.startingBalance"
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
  </div>
</template>

<script>
/* eslint-disable sort-imports */
import { Modal } from 'bootstrap'

import accList from './accountsSidebarAccList.vue'

import { formatNumber } from '../helpers'

export default {
  components: { accList },

  computed: {
    budgetAccounts() {
      const accs = (this.accounts || []).filter(acc => acc.type === 'budget')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },

    budgetSum() {
      return this.budgetAccounts.reduce((sum, acc) => sum + acc.balance, 0)
    },

    trackingAccounts() {
      const accs = (this.accounts || []).filter(acc => acc.type === 'tracking')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },

    trackingSum() {
      return this.trackingAccounts.reduce((sum, acc) => sum + acc.balance, 0)
    },
  },

  data() {
    return {
      modals: {
        addAccount: {
          name: '',
          startingBalance: 0,
          type: 'budget',
        },
      },
    }
  },

  methods: {
    addAccount() {
      return fetch('/api/accounts', {
        body: JSON.stringify({
          name: this.modals.addAccount.name,
          startingBalance: this.modals.addAccount.startingBalance,
          type: this.modals.addAccount.type,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then(() => this.$emit('update-accounts'))
        .then(() => {
          Modal.getInstance(this.$refs.createAccountModal).toggle()

          this.modals.addAccount = {
            name: '',
            startingBalance: 0,
            type: 'budget',
          }
        })
    },

    formatNumber,
  },

  name: 'AccountingAppSidebar',

  props: {
    accounts: {
      required: true,
      type: Array,
    },
  },
}
</script>
