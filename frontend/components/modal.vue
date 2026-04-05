<template>
  <component
    :is="currentModal.component"
    v-if="currentModal"
    v-bind="currentModal.props"
    @reject="rejectModal"
    @resolve="resolveModal"
  />
</template>

<script lang="ts">
import { defineComponent, shallowRef, type Component } from 'vue'

interface ActiveModal {
  component: Component
  props: Record<string, unknown>
  reject: (reason?: unknown) => void // eslint-disable-line no-unused-vars
  resolve: (value?: unknown) => void // eslint-disable-line no-unused-vars
}

const currentModal = shallowRef<ActiveModal | null>(null)

function clearModal() {
  currentModal.value = null
}

function openModal(component: Component, props: Record<string, unknown> = {}) {
  if (currentModal.value !== null) {
    throw new Error('cannot open modal while another modal is active')
  }

  return new Promise<unknown>((resolve, reject) => {
    currentModal.value = {
      component,
      props,
      reject,
      resolve,
    }
  })
}

const ModalHost = defineComponent({
  methods: {
    rejectModal(reason?: unknown) {
      currentModal.value?.reject(reason)
      clearModal()
    },

    resolveModal(value?: unknown) {
      currentModal.value?.resolve(value)
      clearModal()
    },
  },

  name: 'AccountingAppModalHost',

  setup() {
    return {
      currentModal,
    }
  },
})

type ModalHostComponent = typeof ModalHost & {
  openModal: typeof openModal
}

export default Object.assign(ModalHost, { openModal }) as ModalHostComponent
</script>
