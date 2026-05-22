import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')

  const isLoggedIn = computed(() => !!token.value)

  function setTokens(access: string, refresh: string) {
    token.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  function clearTokens() {
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  async function login(email: string, password: string) {
    const { data } = await client.post('/auth/login', { email, password })
    setTokens(data.access_token, data.refresh_token)
  }

  async function register(email: string, password: string) {
    const { data } = await client.post('/auth/register', { email, password })
    setTokens(data.access_token, data.refresh_token)
  }

  function logout() {
    clearTokens()
  }

  return { token, refreshToken, isLoggedIn, setTokens, clearTokens, login, register, logout }
})
