<script setup lang="ts">
  import { useLoginForm } from '@/features/auth'
  import { AppInput } from '@/shared/ui/input'
  import { PrimaryButton } from '@/shared/ui/button'

  const { defineField, errors, onSubmit } = useLoginForm()

  const [email, emailAttrs] = defineField('email')
  const [password, passwordAttrs] = defineField('password')
</script>

<template>
  <div>
    <h1>ログイン</h1>
    <form @submit.prevent="onSubmit">
      <AppInput
        v-model="email"
        v-bind="emailAttrs"
        label="メールアドレス"
        type="email"
        autocomplete="email"
        :error="errors.email"
      />

      <AppInput
        v-model="password"
        v-bind="passwordAttrs"
        label="パスワード"
        type="password"
        autocomplete="current-password"
        :error="errors.password"
      />

      <PrimaryButton type="submit">ログイン</PrimaryButton>
    </form>
  </div>
</template>
