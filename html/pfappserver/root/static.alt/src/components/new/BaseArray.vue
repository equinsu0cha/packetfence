<template>
  <div class="w-100">
    <template v-for="(_, index) in inputValue">
      <template v-for="slotName in slotsList">
        <slot :name="slotName" v-bind="{ index }"></slot>
      </template>
    </template>
  </div>
</template>
<script>
import { useArray, useArrayProps } from '@/composables/useArray'

export const props = {
  ...useArrayProps
}

// @vue/component
export default {
  name: 'base-array',
  inheritAttrs: false,
  props,
  setup(props, context) {
    const {
      value,
      slotsList
    } = useArray(props, context)

    return {
      inputValue: value,
      slotsList
    }
  }
}
</script>
