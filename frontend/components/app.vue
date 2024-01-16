<template>
  <div class="d-flex h-100">
    <div
      class="d-flex flex-column flex-shrink-0 p-3"
      style="width: 300px"
    >
      <accounts-sidebar
        :accounts="accounts"
        @updateAccounts="fetchAccounts"
      />
    </div>
    <div class="d-flex flex-column flex-grow-1 py-3">
      <router-view :accounts="accounts" />
    </div>
  </div>
</template>

<script>
import accountsSidebar from './accountsSidebar.vue'

export default {
  components: { accountsSidebar },

  created() {
    this.fetchAccounts()
  },

  data() {
    return {
      accounts: [],
    }
  },

  methods: {
    fetchAccounts() {
      return fetch('/api/accounts?with-balances')
        .then(resp => resp.json())
        .then(data => {
          this.accounts = data
        })
    },
  },

  name: 'AccountingAppMain',
}
</script>
