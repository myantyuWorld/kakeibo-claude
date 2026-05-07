---
paths:
  - "frontend/src/**/*.vue"
---

# Vue コーディング規約

## script

- ALWAYS `<script setup lang="ts">` を使う
- NEVER Options API や通常の `<script>` を使わない

## props / emits / model

- props は型引数で定義する

```vue
<script setup lang="ts">
  defineProps<{
    label: string
    error?: string
  }>()
</script>
```

- emits も型引数で定義する

```vue
<script setup lang="ts">
  defineEmits<{
    update: [value: string]
    submit: []
  }>()
</script>
```

- ALWAYS v-model には `defineModel()` を使う
- NEVER 手動で prop + emit を組み合わせない

```vue
<script setup lang="ts">
  const model = defineModel<string>()
</script>

<template>
  <input v-model="model" />
</template>
```

## コンポーネント設計

- `$attrs` でネイティブ属性を転送する。ラッパーコンポーネントは受け取った属性を握りつぶさない

```vue
<template>
  <BaseButton :size="size" v-bind="$attrs">
    <slot />
  </BaseButton>
</template>
```

## 命名

- ファイル名は kebab-case（`login-page.vue`、`app-input.vue`）
- テンプレート内のコンポーネントは PascalCase（`<AppInput />`、`<PrimaryButton />`）

## 関連

- **vue-reviewer** agent -- Vue規約のレビュー時
- skill: `impl-front` -- コンポーネント実装時
