import axios from 'axios'
import { useAdminAuthStore } from '@/stores/adminAuth'

const adminClient = axios.create({
  baseURL: '/admin',
  headers: { 'Content-Type': 'application/json' },
})

adminClient.interceptors.request.use((config) => {
  const store = useAdminAuthStore()
  if (store.token) {
    config.headers.Authorization = `Bearer ${store.token}`
  }
  return config
})

adminClient.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401 || err.response?.status === 403) {
      const store = useAdminAuthStore()
      store.clearAuth()
      window.location.href = '/admin/login'
    }
    return Promise.reject(err)
  },
)

export interface AdminTask {
  id: string
  text_key: string
  duration_seconds: number
  category: string
  is_preset: boolean
  is_active: boolean
  created_at: string
}

export interface AdminCategory {
  id: string
  name: string
  label_key: string
  is_active: boolean
}

export const adminApi = {
  // 任务管理
  listTasks(): Promise<AdminTask[]> {
    return adminClient.get('/tasks').then((r) => r.data)
  },
  createTask(data: { text_key: string; duration_seconds: number; category: string }): Promise<{ id: string }> {
    return adminClient.post('/tasks', data).then((r) => r.data)
  },
  updateTask(
    id: string,
    data: Partial<{ text_key: string; duration_seconds: number; category: string; is_active: boolean }>,
  ): Promise<void> {
    return adminClient.put(`/tasks/${id}`, data).then(() => undefined)
  },
  deleteTask(id: string): Promise<void> {
    return adminClient.delete(`/tasks/${id}`).then(() => undefined)
  },

  // 分类管理
  listCategories(): Promise<AdminCategory[]> {
    return adminClient.get('/categories').then((r) => r.data)
  },
  createCategory(data: { name: string; label_key: string }): Promise<{ id: string }> {
    return adminClient.post('/categories', data).then((r) => r.data)
  },
  deleteCategory(id: string): Promise<void> {
    return adminClient.delete(`/categories/${id}`).then(() => undefined)
  },
}
