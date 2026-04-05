<template>
  <div class="d-flex h-100">
    <div
      class="d-flex flex-column flex-shrink-0 p-3"
      style="width: 300px"
    >
      <accounts-sidebar
        :accounts="accounts"
        @update-accounts="fetchAccounts"
      />
    </div>

    <div class="d-flex flex-column flex-grow-1 py-3">
      <router-view
        :accounts="accounts"
        @update-accounts="fetchAccounts"
      />
    </div>

    <modal-host />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import accountsSidebar from 'components/accountsSidebar.vue'
import modalHost from 'components/modal.vue'
import { responseToJSON } from './helpers'
import type { Account } from './types'

export default defineComponent({
  components: { accountsSidebar, modalHost },

  created() {
    void this.fetchAccounts()
  },

  data() {
    return {
      accounts: [] as Account[],
    }
  },

  methods: {
    async fetchAccounts() {
      const resp = await fetch('/api/accounts?with-balances')
      const data = await responseToJSON<Account[]>(resp)
      this.accounts = data ?? []
    },
  },

  name: 'AccountingAppMain',
})
</script>
