<script setup lang="ts">
  import BaseButton from './base-button.vue'

  defineProps<{
    size?: 'sm' | 'md' | 'lg'
    disabled?: boolean
  }>()
</script>

<template>
  <BaseButton
    variant="secondary"
    :size="size"
    :disabled="disabled"
    v-bind="$attrs"
  >
    <slot />
  </BaseButton>
</template>
