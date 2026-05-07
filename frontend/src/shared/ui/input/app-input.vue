<script setup lang="ts">
  defineProps<{
    label: string
    error?: string
  }>()

  const model = defineModel<string>()
</script>

<template>
  <div>
    <label>
      <span>{{ label }}</span>
      <input v-model="model" v-bind="$attrs" />
    </label>
    <p v-if="error">{{ error }}</p>
  </div>
</template>
