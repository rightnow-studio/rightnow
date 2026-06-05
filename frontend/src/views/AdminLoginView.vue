<template>
  <div class="admin-login">
    <div class="login-card">
      <div class="logo">
        <span class="logo-icon">⚙</span>
        <h1>此刻 Admin</h1>
      </div>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="field">
          <label>用户名</label>
          <input
            v-model="username"
            type="text"
            placeholder="admin"
            autocomplete="username"
            required
          />
        </div>
        <div class="field">
          <label>密码</label>
          <input
            v-model="password"
            type="password"
            placeholder="••••••"
            autocomplete="current-password"
            required
          />
        </div>
        <p v-if="error" class="error-msg">{{ error }}</p>
        <button type="submit" :disabled="loading" class="login-btn">
          <span v-if="loading">登录中...</span>
          <span v-else>登录</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminAuthStore } from '@/stores/adminAuth'

const router = useRouter()
const adminAuth = useAdminAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await adminAuth.login(username.value, password.value)
    router.push('/admin')
  } catch {
    error.value = '用户名或密码错误'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login {
  min-height: 100vh;
  background: #0a0a0a;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.login-card {
  width: 100%;
  max-width: 360px;
  background: #141414;
  border: 1px solid #222;
  border-radius: 12px;
  padding: 2.5rem 2rem;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
}

.logo-icon {
  font-size: 1.5rem;
  opacity: 0.7;
}

.logo h1 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #e0e0e0;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.field label {
  font-size: 0.8rem;
  color: #888;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.field input {
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 0.7rem 0.9rem;
  color: #e0e0e0;
  font-size: 0.95rem;
  outline: none;
  transition: border-color 0.15s;
}

.field input:focus {
  border-color: #555;
}

.error-msg {
  color: #f87171;
  font-size: 0.85rem;
  margin: -0.5rem 0 0;
}

.login-btn {
  background: #e0e0e0;
  color: #0a0a0a;
  border: none;
  border-radius: 8px;
  padding: 0.75rem;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s;
  margin-top: 0.5rem;
}

.login-btn:hover:not(:disabled) {
  opacity: 0.88;
}

.login-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
