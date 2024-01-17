<template>
  <div>
    <router-link
      to="/"
      class="d-flex align-items-center fs-5 mb-3 mb-md-0 me-md-auto text-start text-white text-decoration-none"
    >
      <i class="fas fa-fw fa-hand-holding-dollar me-2" />
      My Budget
    </router-link>
    <!--
    <router-link
      to="/"
      class="btn fs-5 text-start w-100"
    >
      <i class="fas fa-fw fa-chart-line" />
      Reports
    </router-link>
    -->
    <hr>
    <ul class="list-unstyled lh-lg ps-0">
      <li class="mb-1 fw-semibold">
        <div class="d-flex align-items-center">
          <i class="fas fa-fw fa-credit-card me-1" /> Budget
          <span :class="{'ms-auto': true, 'text-danger': budgetSum < 0}">
            {{ formatNumber(budgetSum) }} €
          </span>
        </div>
        <ul class="btn-toggle-nav list-unstyled fw-normal pb-1 ps-3 small">
          <li>
            <router-link
              v-for="acc in budgetAccounts"
              :key="acc.id"
              class="d-flex align-items-center text-white text-decoration-none"
              :to="{ name: 'account-transactions', params: { accountId: acc.id }}"
            >
              {{ acc.name }}
              <span :class="{'ms-auto': true, 'text-danger': acc.balance < 0}">
                {{ formatNumber(acc.balance) }} €
              </span>
            </router-link>
          </li>
        </ul>
      </li>
      <li class="mb-1 fw-semibold">
        <div class="d-flex align-items-center">
          <i class="fas fa-fw fa-coin me-1" /> Tracking
          <span :class="{'ms-auto': true, 'text-danger': trackingSum < 0}">
            {{ formatNumber(trackingSum) }} €
          </span>
        </div>
        <ul class="btn-toggle-nav list-unstyled fw-normal pb-1 ps-3 small">
          <li>
            <router-link
              v-for="acc in trackingAccounts"
              :key="acc.id"
              class="d-flex align-items-center text-white text-decoration-none"
              :to="{ name: 'account-transactions', params: { accountId: acc.id }}"
            >
              {{ acc.name }}
              <span :class="{'ms-auto': true, 'text-danger': acc.balance < 0}">
                {{ formatNumber(acc.balance) }} €
              </span>
            </router-link>
          </li>
        </ul>
      </li>
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

import { formatNumber } from '../helpers'
import { unallocatedMoneyAcc } from '../constants'

export default {
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
          type: this.modals.addAccount.type,
        }),
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      })
        .then(resp => resp.json())
        .then(account => {
          if (account.type === 'budget') {
            return fetch('/api/transactions', {
              body: JSON.stringify({
                account: account.id,
                amount: this.modals.addAccount.startingBalance,
                category: unallocatedMoneyAcc,
                cleared: true,
                description: 'Starting Balance',
                time: new Date(),
              }),
              headers: {
                'Content-Type': 'application/json',
              },
              method: 'POST',
            })
          } else if (account.type === 'tracking') {
            return fetch('/api/transactions', {
              body: JSON.stringify({
                account: account.id,
                amount: this.modals.addAccount.startingBalance,
                cleared: true,
                description: 'Starting Balance',
                time: new Date(),
              }),
              headers: {
                'Content-Type': 'application/json',
              },
              method: 'POST',
            })
          } else if (account.type === 'category') {
            return fetch(`/api/accounts/${unallocatedMoneyAcc}/transfer/${account.id}?amount=${this.modals.addAccount.startingBalance}`, {
              method: 'PUT',
            })
          }
          throw new Error('invalid account type detected')
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
