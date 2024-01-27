<template>
  <li class="mb-1 fw-semibold">
    <div class="d-flex align-items-center">
      <i class="fas fa-fw fa-credit-card me-1" /> {{ header }}
      <span :class="classFromNumber(sum, ['ms-auto'])">
        {{ formatNumber(sum) }} €
      </span>
    </div>
    <ul class="btn-toggle-nav list-unstyled fw-normal pb-1 ps-3 small">
      <li>
        <router-link
          v-for="acc in accounts"
          :key="acc.id"
          class="d-flex align-items-center text-white text-decoration-none"
          :to="{ name: 'account-transactions', params: { accountId: acc.id }}"
        >
          {{ acc.name }}
          <span :class="classFromNumber(acc.balance, ['ms-auto'])">
            {{ formatNumber(acc.balance) }} €
          </span>
        </router-link>
      </li>
    </ul>
  </li>
</template>

<script>
import { classFromNumber, formatNumber } from '../helpers'

export default {
  computed: {
    sum() {
      return this.accounts.reduce((sum, acc) => sum + acc.balance, 0)
    },
  },

  methods: {
    classFromNumber,
    formatNumber,
  },

  name: 'AccountingAppAccountSidebarAccountList',

  props: {
    accounts: {
      required: true,
      type: Array,
    },

    header: {
      required: true,
      type: String,
    },
  },
}
</script>
