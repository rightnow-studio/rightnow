<template>
  <div class="task-list-screen">
    <div class="task-list-header">
      <h2>{{ t('taskList.title') }}</h2>
    </div>

    <div class="task-list-scroll">
      <div v-if="loading" class="empty">{{ t('profile.loading') }}</div>
      <div v-else-if="sortedTasks.length === 0" class="empty">{{ t('taskList.empty') }}</div>
      <div v-else class="task-list">
        <div
          v-for="task in sortedTasks"
          :key="task.id"
          class="task-item"
          :class="{ inactive: !task.is_active }"
        >
          <div class="task-content">
            <div class="task-text" :class="{ 'is-completed': !task.is_active }">{{ task.text_key }}</div>
            <div class="task-meta">
              <span class="category">{{ t(`categories.${task.category}`) }}</span>
              <span class="duration">{{ formatDuration(task.duration_seconds) }}</span>
            </div>
          </div>
          <div class="task-actions">
            <button
              class="btn-toggle"
              :class="{ active: task.is_active }"
              @click="toggleTask(task)"
              :title="task.is_active ? t('taskList.deactivate') : t('taskList.activate')"
            >
              <svg v-if="task.is_active" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
                <polyline points="22 4 12 14.01 9 11.01" />
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10" />
              </svg>
            </button>
            <button class="btn-delete" @click="deleteTask(task.id)" :title="t('taskList.delete')">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <button class="btn-back" @click="router.replace('/')">{{ t('common.back') }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import client from '@/api/client'

const router = useRouter()
const { t } = useI18n()

interface Task {
  id: string
  text_key: string
  duration_seconds: number
  category: string
  is_preset: boolean
  created_by: string
  is_active: boolean
  created_at: string
}

const tasks = ref<Task[]>([])
const loading = ref(true)

const myTasks = computed(() => tasks.value.filter(t => !t.is_preset))

const sortedTasks = computed(() => {
  return [...myTasks.value].sort((a, b) => {
    // 活跃的在前，停用的在后
    if (a.is_active === b.is_active) {
      return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    }
    return a.is_active ? -1 : 1
  })
})

function formatDuration(s: number) {
  if (s < 60) return `${s}${t('common.second')}`
  return `${Math.floor(s / 60)}${t('common.minute')}`
}

async function fetchTasks() {
  loading.value = true
  try {
    const { data } = await client.get('/tasks')
    tasks.value = data.tasks ?? []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function toggleTask(task: Task) {
  try {
    await client.patch(`/tasks/${task.id}`, { is_active: !task.is_active })
    task.is_active = !task.is_active
  } catch (e) {
    alert(t('common.error', { msg: t('common.unknownError') }))
  }
}

async function deleteTask(id: string) {
  if (!confirm(t('taskList.confirmDelete'))) return
  try {
    await client.delete(`/tasks/${id}`)
    tasks.value = tasks.value.filter(t => t.id !== id)
  } catch (e) {
    alert(t('common.error', { msg: t('common.unknownError') }))
  }
}

onMounted(fetchTasks)
</script>

<style scoped>
.task-list-screen {
  width: 100%;
  height: 100%;
  background: #121212;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #f0f0f0;
  overflow: hidden;
}

.task-list-header {
  width: 100%;
  max-width: 600px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 0.5rem;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.task-list-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #f0f0f0;
}

.task-list-scroll {
  width: 100%;
  max-width: 600px;
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.empty {
  font-size: 13px;
  color: #555;
  text-align: center;
  padding: 2rem 0;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.task-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  padding: 0.75rem 1rem;
}

.task-item.inactive {
  opacity: 0.7;
  background: #161616;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-text {
  font-size: 14px;
  font-weight: 500;
  color: #e0e0e0;
  line-height: 1.4;
  margin-bottom: 0.25rem;
  word-break: break-all;
}

.task-text.is-completed {
  text-decoration: line-through;
  color: #666;
}

.task-meta {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.category {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 8px;
  background: rgba(45, 212, 191, 0.12);
  color: #2dd4bf;
}

.duration {
  font-size: 10px;
  color: #666;
}

.task-actions {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  flex-shrink: 0;
}

.btn-toggle,
.btn-delete {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: transparent;
  border: 1px solid #333;
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-toggle:hover {
  color: #2dd4bf;
  border-color: #2dd4bf;
}

.btn-toggle.active {
  color: #2dd4bf;
  border-color: #2dd4bf;
}

.btn-delete:hover {
  color: #ef4444;
  border-color: #ef4444;
}

.btn-toggle svg,
.btn-delete svg {
  width: 16px;
  height: 16px;
}

.btn-back {
  width: 100%;
  max-width: 600px;
  padding: 0.75rem;
  background: transparent;
  color: #888;
  border: 1px solid #333;
  border-radius: 12px;
  font-size: 15px;
  margin-top: 0.75rem;
  margin-bottom: 0.5rem;
  cursor: pointer;
  flex-shrink: 0;
}

.btn-back:hover {
  border-color: #444;
  color: #aaa;
}
</style>
