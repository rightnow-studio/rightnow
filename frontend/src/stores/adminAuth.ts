import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useAdminAuthStore = defineStore('adminAuth', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const username = ref(localStorage.getItem('admin_username') || '')

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(t: string, u: string) {
    token.value = t
    username.value = u
    localStorage.setItem('admin_token', t)
    localStorage.setItem('admin_username', u)
  }

  function clearAuth() {
    token.value = ''
    username.value = ''
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_username')
  }

  async function login(u: string, password: string) {
    const { data } = await axios.post('/admin/auth/login', { username: u, password })
    setAuth(data.token, data.username)
  }

  function logout() {
    clearAuth()
  }

  return { token, username, isLoggedIn, login, logout, clearAuth }
})
