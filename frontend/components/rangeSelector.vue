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

<script>
export default {
  computed: {
    monthNames() {
      const format = new Intl
        .DateTimeFormat(undefined, { month: 'long' }).format
      return [...Array(12).keys()]
        .map(m => format(new Date(Date.UTC(2021, m))))
    },

    years() {
      return [...Array(this.endYear - this.startYear).keys()]
        .map(y => y + this.startYear)
    },
  },

  created() {
    const date = new Date()
    let start = this.modelValue.start || new Date(date.getFullYear(), date.getMonth(), 1)
    let end = this.modelValue.end || new Date(date.getFullYear(), date.getMonth() + 1, 0)

    if (this.storeKey && localStorage.getItem(`accounting:range-select:${this.storeKey}`)) {
      const stored = JSON.parse(localStorage.getItem(`accounting:range-select:${this.storeKey}`))
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
      },
    }
  },

  emits: ['update:modelValue'],

  methods: {
    emitRange() {
      this.start = new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth, 1)
      if (this.multiMonth) {
        this.end = new Date(this.dateComponents.toYear, this.dateComponents.toMonth + 1, 1, 0)
      } else {
        this.end = new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth + 1, 1, 0)
      }

      if (this.storeKey) {
        localStorage.setItem(
          `accounting:range-select:${this.storeKey}`,
          JSON.stringify({
            end: new Date(this.dateComponents.toYear, this.dateComponents.toMonth, 1),
            start: new Date(this.dateComponents.fromYear, this.dateComponents.fromMonth, 1),
          }),
        )
      }

      this.$emit('update:modelValue', { end: this.end, start: this.start })
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
      default: () => ({}),
      type: Object,
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
      type: String,
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
}
</script>
