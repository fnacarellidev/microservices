<script setup lang="ts">
import { z } from 'zod'
import type { FormSubmitEvent } from '#ui/types'

const state = reactive({
  username: undefined,
  password: undefined
})

const registrationSchema = z.object({
  username: z.string(),
  password: z.string()
})

type RegistrationSchema = z.output<typeof registrationSchema>

async function logUserIn(event: FormSubmitEvent<RegistrationSchema>) {
  const response = await $fetch('http://localhost:8080/auth/login', {
    method: 'POST',
    credentials: 'include',
    body: {
      username: event.data.username,
      password: event.data.password
    }
  })
}

async function registerUser(e: FormSubmitEvent<RegistrationSchema>) {
  const response = await $fetch('http://localhost:8080/auth/register', {
    method: 'POST',
    body: {
      username: e.data.username,
      password: e.data.password
    }
  })
}

onMounted(() => {
  console.log(document.cookie)
})

</script>
<template>
  <div>
    <h2 class="ml-2 mt-2">Login</h2>
    <UForm :schema="registrationSchema" :state="state" class="space-y-4 w-1/4 m-2" @submit="logUserIn">
      <UFormGroup label="Username" name="username">
        <UInput v-model="state.username" />
      </UFormGroup>

      <UFormGroup label="Password" name="password">
        <UInput v-model="state.password" type="password" />
      </UFormGroup>

      <UButton type="submit">
        Login
      </UButton>
    </UForm>

    <h2 class="ml-2 mt-2">Register</h2>
    <UForm :schema="registrationSchema" :state="state" class="space-y-4 w-1/4 m-2" @submit="registerUser">
      <UFormGroup label="Username" name="username">
        <UInput v-model="state.username" />
      </UFormGroup>

      <UFormGroup label="Password" name="password">
        <UInput v-model="state.password" type="password" />
      </UFormGroup>

      <UButton type="submit">
        Register
      </UButton>
    </UForm>
  </div>
</template>
