<template>
  <div class="login-screen">
    <div class="login-card">
      <h1>{{ t('login.title') }}</h1>
      <p class="subtitle">{{ t('login.subtitle') }}</p>
      <input v-model="email" type="email" :placeholder="t('login.email')" />
      <input v-model="password" type="password" :placeholder="t('login.password')" />
      <button class="btn-primary" @click="handleSubmit">
        {{ isRegister ? t('login.register') : t('login.login') }}
      </button>
      <p class="switch">
        {{ isRegister ? t('login.hasAccount') : t('login.noAccount') }}
        <span @click="isRegister = !isRegister">
          {{ isRegister ? t('login.login') : t('login.register') }}
        </span>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const isRegister = ref(false)

async function handleSubmit() {
  try {
    if (isRegister.value) {
      await auth.register(email.value, password.value)
    } else {
      await auth.login(email.value, password.value)
    }
    router.push('/')
  } catch (e) {
    const msg = (e as any).response?.data?.error || t('common.unknownError')
    alert(t('common.error', { msg }))
  }
}
</script>

<style scoped>
.login-screen {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0a0a0a;
  padding: 1.5rem;
}
.login-card {
  width: 100%;
  max-width: 340px;
}
h1 {
  font-size: 32px;
  font-weight: 700;
  color: #f0f0f0;
  margin-bottom: 4px;
  letter-spacing: -0.5px;
}
.subtitle {
  font-size: 14px;
  color: #888;
  margin-bottom: 2.5rem;
}
input {
  width: 100%;
  padding: 0.85rem 1rem;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 10px;
  color: #f0f0f0;
  font-size: 15px;
  margin-bottom: 0.75rem;
  outline: none;
}
input:focus {
  border-color: #2dd4bf;
}
.btn-primary {
  width: 100%;
  padding: 0.9rem;
  background: #2dd4bf;
  color: #0a0a0a;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  margin-top: 0.5rem;
}
.switch {
  font-size: 13px;
  color: #666;
  margin-top: 1.25rem;
  text-align: center;
}
.switch span {
  color: #2dd4bf;
  cursor: pointer;
}
</style>
