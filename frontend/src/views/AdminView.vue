<template>
  <div class="admin-page">
    <!-- 顶部导航栏 -->
    <header class="admin-header">
      <div class="header-left">
        <span class="header-icon">⚙</span>
        <span class="header-title">此刻 Admin</span>
      </div>
      <div class="header-right">
        <span class="header-user">{{ adminAuth.username }}</span>
        <button class="logout-btn" @click="handleLogout">退出</button>
      </div>
    </header>

    <!-- Tab 切换 -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        :class="['tab-btn', { active: activeTab === tab.key }]"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 预设任务管理 -->
    <div v-if="activeTab === 'tasks'" class="tab-content">
      <div class="section-header">
        <h2>预设任务 ({{ tasks.length }})</h2>
        <button class="btn-primary" @click="showCreateTask = true">+ 新增任务</button>
      </div>

      <!-- 新增任务表单 -->
      <div v-if="showCreateTask" class="form-card">
        <h3>新增预设任务</h3>
        <div class="form-row">
          <div class="field">
            <label>任务 Key（text_key）</label>
            <input v-model="newTask.text_key" placeholder="task.my_task_name" />
          </div>
          <div class="field">
            <label>时长（秒，10-300）</label>
            <input v-model.number="newTask.duration_seconds" type="number" min="10" max="300" placeholder="60" />
          </div>
          <div class="field">
            <label>分类</label>
            <select v-model="newTask.category">
              <option v-for="cat in activeCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
            </select>
          </div>
        </div>
        <p v-if="taskFormError" class="error-msg">{{ taskFormError }}</p>
        <div class="form-actions">
          <button class="btn-primary" @click="handleCreateTask" :disabled="taskFormLoading">
            {{ taskFormLoading ? '创建中...' : '创建' }}
          </button>
          <button class="btn-ghost" @click="cancelCreateTask">取消</button>
        </div>
      </div>

      <!-- 任务列表 -->
      <div v-if="tasksLoading" class="loading">加载中...</div>
      <div v-else-if="tasks.length === 0" class="empty">暂无预设任务</div>
      <div v-else class="task-table-wrap">
        <table class="task-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Text Key</th>
              <th>时长(秒)</th>
              <th>分类</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in tasks" :key="task.id" :class="{ inactive: !task.is_active }">
              <td class="id-cell">{{ task.id }}</td>
              <td>
                <template v-if="editingTask?.id === task.id">
                  <input v-model="editingTask.text_key" class="inline-input" />
                </template>
                <template v-else>{{ task.text_key }}</template>
              </td>
              <td>
                <template v-if="editingTask?.id === task.id">
                  <input v-model.number="editingTask.duration_seconds" type="number" min="10" max="300" class="inline-input small" />
                </template>
                <template v-else>{{ task.duration_seconds }}</template>
              </td>
              <td>
                <template v-if="editingTask?.id === task.id">
                  <select v-model="editingTask.category" class="inline-select">
                    <option v-for="cat in activeCategories" :key="cat.id" :value="cat.name">{{ cat.name }}</option>
                  </select>
                </template>
                <template v-else>
                  <span class="cat-badge">{{ task.category }}</span>
                </template>
              </td>
              <td>
                <span :class="['status-dot', task.is_active ? 'active' : 'inactive']">
                  {{ task.is_active ? '启用' : '禁用' }}
                </span>
              </td>
              <td class="action-cell">
                <template v-if="editingTask?.id === task.id">
                  <button class="btn-sm btn-primary" @click="handleSaveTask">保存</button>
                  <button class="btn-sm btn-ghost" @click="editingTask = null">取消</button>
                </template>
                <template v-else>
                  <button class="btn-sm btn-ghost" @click="startEditTask(task)">编辑</button>
                  <button
                    class="btn-sm"
                    :class="task.is_active ? 'btn-warn' : 'btn-primary'"
                    @click="handleToggleTask(task)"
                  >
                    {{ task.is_active ? '禁用' : '启用' }}
                  </button>
                  <button class="btn-sm btn-danger" @click="handleDeleteTask(task.id)">删除</button>
                </template>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 任务分类管理 -->
    <div v-if="activeTab === 'categories'" class="tab-content">
      <div class="section-header">
        <h2>任务分类 ({{ categories.length }})</h2>
      </div>

      <!-- 分类列表 -->
      <div v-if="catsLoading" class="loading">加载中...</div>
      <div v-else class="cat-list">
        <div v-for="cat in categories" :key="cat.id" class="cat-item">
          <div class="cat-info">
            <span class="cat-name">{{ cat.name }}</span>
            <span class="cat-key">{{ cat.label_key }}</span>
            <span :class="['status-dot', cat.is_active ? 'active' : 'inactive']">
              {{ cat.is_active ? '启用' : '禁用' }}
            </span>
          </div>
          <button
            class="btn-sm btn-danger"
            @click="handleDeleteCategory(cat)"
            :title="cat.name === 'work' || cat.name === 'home' || cat.name === 'communication' || cat.name === 'health' || cat.name === 'admin' ? '系统内置分类，删除需先清空任务' : ''"
          >
            删除
          </button>
        </div>

        <!-- 新增分类表单 -->
        <div class="cat-add-form">
          <h3>新增分类</h3>
          <div class="form-row">
            <div class="field">
              <label>分类名（英文标识，如 focus）</label>
              <input v-model="newCat.name" placeholder="focus" />
            </div>
            <div class="field">
              <label>i18n Key（如 categories.focus）</label>
              <input v-model="newCat.label_key" placeholder="categories.focus" />
            </div>
          </div>
          <p v-if="catFormError" class="error-msg">{{ catFormError }}</p>
          <button class="btn-primary" @click="handleCreateCategory" :disabled="catFormLoading">
            {{ catFormLoading ? '创建中...' : '新增分类' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminAuthStore } from '@/stores/adminAuth'
import { adminApi, type AdminTask, type AdminCategory } from '@/api/admin'

const router = useRouter()
const adminAuth = useAdminAuthStore()

const tabs = [
  { key: 'tasks', label: '预设任务' },
  { key: 'categories', label: '任务分类' },
]
const activeTab = ref('tasks')

// 任务相关
const tasks = ref<AdminTask[]>([])
const tasksLoading = ref(false)
const showCreateTask = ref(false)
const taskFormLoading = ref(false)
const taskFormError = ref('')
const editingTask = ref<AdminTask | null>(null)
const newTask = ref({ text_key: '', duration_seconds: 60, category: 'work' })

// 分类相关
const categories = ref<AdminCategory[]>([])
const catsLoading = ref(false)
const catFormLoading = ref(false)
const catFormError = ref('')
const newCat = ref({ name: '', label_key: '' })

const activeCategories = computed(() => categories.value.filter((c) => c.is_active))

async function loadTasks() {
  tasksLoading.value = true
  try {
    tasks.value = await adminApi.listTasks()
  } finally {
    tasksLoading.value = false
  }
}

async function loadCategories() {
  catsLoading.value = true
  try {
    categories.value = await adminApi.listCategories()
  } finally {
    catsLoading.value = false
  }
}

async function handleCreateTask() {
  taskFormError.value = ''
  if (!newTask.value.text_key.trim()) {
    taskFormError.value = '请填写 text_key'
    return
  }
  if (newTask.value.duration_seconds < 10 || newTask.value.duration_seconds > 300) {
    taskFormError.value = '时长需在 10-300 秒之间'
    return
  }
  taskFormLoading.value = true
  try {
    await adminApi.createTask(newTask.value)
    newTask.value = { text_key: '', duration_seconds: 60, category: activeCategories.value[0]?.name || 'work' }
    showCreateTask.value = false
    await loadTasks()
  } catch (e: any) {
    taskFormError.value = e.response?.data?.error || '创建失败'
  } finally {
    taskFormLoading.value = false
  }
}

function cancelCreateTask() {
  showCreateTask.value = false
  taskFormError.value = ''
  newTask.value = { text_key: '', duration_seconds: 60, category: activeCategories.value[0]?.name || 'work' }
}

function startEditTask(task: AdminTask) {
  editingTask.value = { ...task }
}

async function handleSaveTask() {
  if (!editingTask.value) return
  const { id, text_key, duration_seconds, category } = editingTask.value
  await adminApi.updateTask(id, { text_key, duration_seconds, category })
  editingTask.value = null
  await loadTasks()
}

async function handleToggleTask(task: AdminTask) {
  await adminApi.updateTask(task.id, { is_active: !task.is_active })
  await loadTasks()
}

async function handleDeleteTask(id: string) {
  if (!confirm('确认删除该预设任务？此操作不可恢复。')) return
  await adminApi.deleteTask(id)
  await loadTasks()
}

async function handleCreateCategory() {
  catFormError.value = ''
  if (!newCat.value.name.trim() || !newCat.value.label_key.trim()) {
    catFormError.value = '名称和 Key 不能为空'
    return
  }
  catFormLoading.value = true
  try {
    await adminApi.createCategory(newCat.value)
    newCat.value = { name: '', label_key: '' }
    await loadCategories()
  } catch (e: any) {
    catFormError.value = e.response?.data?.error || '创建失败，可能已存在同名分类'
  } finally {
    catFormLoading.value = false
  }
}

async function handleDeleteCategory(cat: AdminCategory) {
  if (!confirm(`确认删除分类「${cat.name}」？若有任务引用该分类则无法删除。`)) return
  try {
    await adminApi.deleteCategory(cat.id)
    await loadCategories()
  } catch (e: any) {
    const msg = e.response?.data?.error || '删除失败'
    const count = e.response?.data?.task_count
    alert(count ? `${msg}（共 ${count} 个任务使用该分类）` : msg)
  }
}

function handleLogout() {
  adminAuth.logout()
  router.push('/admin/login')
}

onMounted(() => {
  loadTasks()
  loadCategories()
})
</script>

<style scoped>
.admin-page {
  min-height: 100vh;
  background: #0a0a0a;
  color: #e0e0e0;
  font-size: 0.9rem;
}

/* 顶部导航 */
.admin-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
  height: 52px;
  background: #141414;
  border-bottom: 1px solid #222;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.header-icon {
  opacity: 0.6;
}

.header-title {
  font-weight: 600;
  font-size: 1rem;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-user {
  color: #888;
  font-size: 0.85rem;
}

.logout-btn {
  background: transparent;
  border: 1px solid #333;
  border-radius: 6px;
  color: #aaa;
  padding: 0.3rem 0.75rem;
  font-size: 0.8rem;
  cursor: pointer;
  transition: border-color 0.15s;
}

.logout-btn:hover {
  border-color: #555;
  color: #e0e0e0;
}

/* Tabs */
.tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid #222;
  padding: 0 1.5rem;
}

.tab-btn {
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  color: #666;
  padding: 0.8rem 1.25rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s;
}

.tab-btn.active {
  color: #e0e0e0;
  border-bottom-color: #e0e0e0;
}

/* 内容区 */
.tab-content {
  padding: 1.5rem;
  max-width: 1100px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.section-header h2 {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: #e0e0e0;
}

/* 表单卡片 */
.form-card {
  background: #141414;
  border: 1px solid #222;
  border-radius: 10px;
  padding: 1.25rem;
  margin-bottom: 1.25rem;
}

.form-card h3 {
  font-size: 0.9rem;
  font-weight: 600;
  margin: 0 0 1rem;
  color: #ccc;
}

.form-row {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
  min-width: 180px;
}

.field label {
  font-size: 0.75rem;
  color: #777;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.field input,
.field select {
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  border-radius: 6px;
  padding: 0.5rem 0.7rem;
  color: #e0e0e0;
  font-size: 0.9rem;
  outline: none;
}

.field input:focus,
.field select:focus {
  border-color: #555;
}

.form-actions {
  display: flex;
  gap: 0.6rem;
  margin-top: 1rem;
}

.error-msg {
  color: #f87171;
  font-size: 0.82rem;
  margin: 0.3rem 0 0;
}

/* 按钮 */
.btn-primary {
  background: #e0e0e0;
  color: #0a0a0a;
  border: none;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.15s;
}

.btn-primary:hover:not(:disabled) { opacity: 0.85; }
.btn-primary:disabled { opacity: 0.45; cursor: not-allowed; }

.btn-ghost {
  background: transparent;
  color: #aaa;
  border: 1px solid #333;
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: border-color 0.15s;
}

.btn-ghost:hover { border-color: #555; color: #e0e0e0; }

.btn-warn {
  background: #422;
  color: #f87171;
  border: none;
  border-radius: 6px;
  padding: 0.3rem 0.6rem;
  font-size: 0.78rem;
  cursor: pointer;
}

.btn-danger {
  background: #3a1010;
  color: #ef4444;
  border: none;
  border-radius: 6px;
  padding: 0.3rem 0.6rem;
  font-size: 0.78rem;
  cursor: pointer;
  transition: background 0.15s;
}

.btn-danger:hover { background: #4a1515; }

.btn-sm {
  padding: 0.3rem 0.6rem;
  font-size: 0.78rem;
  border-radius: 5px;
}

/* 任务表格 */
.task-table-wrap {
  overflow-x: auto;
  border: 1px solid #1e1e1e;
  border-radius: 10px;
}

.task-table {
  width: 100%;
  border-collapse: collapse;
}

.task-table th {
  background: #111;
  color: #777;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.65rem 0.9rem;
  text-align: left;
  font-weight: 500;
}

.task-table td {
  padding: 0.6rem 0.9rem;
  border-top: 1px solid #1a1a1a;
  vertical-align: middle;
}

.task-table tr.inactive td {
  opacity: 0.45;
}

.task-table tr:hover td {
  background: #0f0f0f;
}

.id-cell {
  color: #555;
  font-size: 0.78rem;
  font-family: monospace;
}

.cat-badge {
  background: #1e1e1e;
  color: #aaa;
  border-radius: 4px;
  padding: 0.15rem 0.5rem;
  font-size: 0.78rem;
}

.status-dot {
  display: inline-block;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.78rem;
}

.status-dot.active {
  background: #0f2e1a;
  color: #4ade80;
}

.status-dot.inactive {
  background: #1e1a0f;
  color: #888;
}

.action-cell {
  display: flex;
  gap: 0.4rem;
  align-items: center;
}

.inline-input {
  background: #0f0f0f;
  border: 1px solid #333;
  border-radius: 4px;
  color: #e0e0e0;
  padding: 0.25rem 0.5rem;
  font-size: 0.85rem;
  width: 100%;
  outline: none;
}

.inline-input.small {
  width: 70px;
}

.inline-select {
  background: #0f0f0f;
  border: 1px solid #333;
  border-radius: 4px;
  color: #e0e0e0;
  padding: 0.25rem 0.5rem;
  font-size: 0.85rem;
  outline: none;
}

/* 分类列表 */
.cat-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.cat-item {
  background: #141414;
  border: 1px solid #1e1e1e;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.cat-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.cat-name {
  font-weight: 600;
  color: #e0e0e0;
  min-width: 100px;
}

.cat-key {
  color: #666;
  font-size: 0.82rem;
  font-family: monospace;
}

.cat-add-form {
  background: #141414;
  border: 1px solid #222;
  border-radius: 10px;
  padding: 1.25rem;
  margin-top: 1rem;
}

.cat-add-form h3 {
  font-size: 0.9rem;
  font-weight: 600;
  margin: 0 0 1rem;
  color: #ccc;
}

.loading,
.empty {
  color: #555;
  padding: 2rem 0;
  text-align: center;
}
</style>
