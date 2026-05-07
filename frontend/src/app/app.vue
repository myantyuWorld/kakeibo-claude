<script setup lang="ts">
  import { computed, type Component } from 'vue'
  import { useRoute } from 'vue-router'
  import type { LayoutType } from './routes/types'
  import PublicLayout from './layouts/public-layout.vue'
  import PrivateLayout from './layouts/private-layout.vue'

  const layouts: Record<LayoutType, Component> = {
    public: PublicLayout,
    private: PrivateLayout,
  }

  const route = useRoute()
  const layout = computed(() => layouts[route.meta.layout ?? 'private'])
</script>

<template>
  <component :is="layout">
    <RouterView />
  </component>
</template>
