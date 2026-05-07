<script setup lang="ts">
  defineProps<{
    variant: 'primary' | 'secondary'
    size?: 'sm' | 'md' | 'lg'
    disabled?: boolean
  }>()
</script>

<template>
  <button
    :data-variant="variant"
    :data-size="size"
    :disabled="disabled"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>
