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
      @click="openCreateAccountModal"
    >
      <i class="fas fa-fw fa-circle-plus mr-1" />
      Add Account
    </button>
  </div>
</template>

<script lang="ts">

import { defineComponent, type PropType } from 'vue'

import accList from './accountsSidebarAccList.vue'
import createAccountModal from './createAccountModal.vue'
import modalHost from './modal.vue'

import type { Account } from '../types'

export default defineComponent({
  components: { accList },

  computed: {
    budgetAccounts(): Account[] {
      const accs = (this.accounts || []).filter(acc => acc.type === 'budget')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },

    budgetSum(): number {
      return this.budgetAccounts.reduce((sum, acc) => sum + acc.balance, 0)
    },

    trackingAccounts(): Account[] {
      const accs = (this.accounts || []).filter(acc => acc.type === 'tracking')
      accs.sort((a, b) => a.name.localeCompare(b.name))
      return accs
    },

    trackingSum(): number {
      return this.trackingAccounts.reduce((sum, acc) => sum + acc.balance, 0)
    },
  },

  emits: ['update-accounts'],

  methods: {
    async openCreateAccountModal() {
      try {
        await modalHost.openModal(createAccountModal)
        this.$emit('update-accounts')
      } catch {
        // Dismissed by user.
      }
    },
  },

  name: 'AccountingAppSidebar',

  props: {
    accounts: {
      required: true,
      type: Array as PropType<Account[]>,
    },
  },
})
</script>
