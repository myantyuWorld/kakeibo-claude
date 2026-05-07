<script setup lang="ts">
  import { useRouter } from 'vue-router'

  const router = useRouter()
</script>

<template>
  <div>
    <p>404</p>
    <h1>ページが見つかりません</h1>
    <p>お探しのページは存在しないか、移動した可能性があります。</p>
    <button @click="router.push('/')">トップページに戻る</button>
  </div>
</template>
