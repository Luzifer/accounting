<template>
  <div class="d-inline-flex align-items-center">
    <span
      v-if="multiMonth"
      class="me-1"
    >From:</span>
    <div class="input-group">
      <select
        v-model="dateComponents.fromMonth"
        class="form-select"
      >
        <option
          v-for="(name, idx) in monthNames"
          :key="idx"
          :value="idx"
          :disabled="dateComponents.fromYear === dateComponents.toYear && idx > dateComponents.toMonth"
        >
          {{ name }}
        </option>
      </select>
      <select
        v-model="dateComponents.fromYear"
        class="form-select"
      >
        <option
          v-for="year in years"
          :key="year"
          :value="year"
          :disabled="year > dateComponents.toYear"
        >
          {{ year }}
        </option>
      </select>
    </div>
    <template v-if="multiMonth">
      <span class="me-1 ms-2">To:</span>
      <div class="input-group">
        <select
          v-model="dateComponents.toMonth"
          class="form-select"
        >
          <option
            v-for="(name, idx) in monthNames"
            :key="idx"
            :value="idx"
            :disabled="dateComponents.fromYear === dateComponents.toYear && idx < dateComponents.fromMonth"
          >
            {{ name }}
          </option>
        </select>
        <select
          v-model="dateComponents.toYear"
          class="form-select"
        >
          <option
            v-for="year in years"
            :key="year"
            :value="year"
            :disabled="year < dateComponents.fromYear"
          >
            {{ year }}
          </option>
        </select>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue'

import type { DateRange, DateRangeStorage } from '../types'

interface DateComponents {
  fromMonth: number
  fromYear: number
  toMonth: number
  toYear: number
}

export default defineComponent({
  computed: {
    monthNames(): string[] {
      const format = new Intl
        .DateTimeFormat(undefined, { month: 'long' }).format
      return [...Array(12).keys()]
        .map(m => format(new Date(Date.UTC(2021, m))))
    },

    years(): number[] {
      return [...Array(this.endYear - this.startYear).keys()]
        .map(y => y + this.startYear)
    },
  },

  created() {
    const date = new Date()
    let start = this.modelValue.start || new Date(date.getFullYear(), date.getMonth(), 1)
    let end = this.modelValue.end || new Date(date.getFullYear(), date.getMonth() + 1, 0)
    const storageKey = this.storeKey ? `accounting:range-select:${this.storeKey}` : null

    if (storageKey && localStorage.getItem(storageKey)) {
      const stored = JSON.parse(localStorage.getItem(storageKey) as string) as DateRangeStorage
      start = new Date(stored.start)
      end = new Date(stored.end)
    }

    this.dateComponents = {
      fromMonth: start.getMonth(),
      fromYear: start.getFullYear(),
      toMonth: end.getMonth(),
      toYear: end.getFullYear(),
    }
  },

  data() {
    return {
      dateComponents: {
        fromMonth: 0,
        fromYear: 0,
        toMonth: 0,
        toYear: 0,
      } as DateComponents,
    }
  },

  emits: ['update:modelValue'],

  methods: {
    emitRange() {
      const start = new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth, 1)
      const end = this.multiMonth
        ? new Date(this.dateComponents.toYear, this.dateComponents.toMonth + 1, 1, 0)
        : new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth + 1, 1, 0)

      if (this.storeKey) {
        localStorage.setItem(
          `accounting:range-select:${this.storeKey}`,
          JSON.stringify({
            end: this.multiMonth
              ? new Date(this.dateComponents.toYear, this.dateComponents.toMonth, 1)
              : new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth, 1),
            start,
          }),
        )
      }

      this.$emit('update:modelValue', { end, start } as DateRange)
    },
  },

  mounted() {
    this.emitRange()
  },

  name: 'AccountingAppRangeSelector',

  props: {
    endYear: {
      default: 2034,
      type: Number,
    },

    modelValue: {
      default: (): DateRange => ({}),
      type: Object as PropType<DateRange>,
    },

    multiMonth: {
      default: true,
      type: Boolean,
    },

    startYear: {
      default: 2020,
      type: Number,
    },

    storeKey: {
      default: null,
      type: String as PropType<null | string>,
    },
  },

  watch: {
    'dateComponents.fromMonth'() {
      this.emitRange()
    },

    'dateComponents.fromYear'() {
      this.emitRange()
    },

    'dateComponents.toMonth'() {
      this.emitRange()
    },

    'dateComponents.toYear'() {
      this.emitRange()
    },
  },
})
</script>
