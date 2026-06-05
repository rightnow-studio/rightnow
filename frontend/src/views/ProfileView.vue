<template>
  <div class="profile-screen">
    <!-- 临时用户绑定邮箱横幅 -->
    <div v-if="user?.is_temp" class="guest-banner">
      <div class="banner-hint">
        <span>{{ t('profile.guestHint') }}</span>
        <button class="btn-bind-toggle" @click="showBind = !showBind">{{ t('profile.bindEmail') }}</button>
      </div>
      <form v-if="showBind" class="bind-form" @submit.prevent="handleBind">
        <input
          v-model="bindEmailVal"
          type="email"
          class="bind-input"
          :placeholder="t('login.email')"
          required
        />
        <input
          v-model="bindPasswordVal"
          type="password"
          class="bind-input"
          :placeholder="t('login.password')"
          required
        />
        <button type="submit" class="btn-bind-confirm" :disabled="binding">
          {{ binding ? t('profile.loading') : t('profile.bindConfirm') }}
        </button>
      </form>
    </div>

    <div class="profile-header">
      <div class="avatar">{{ avatarLetter }}</div>
      <div class="user-info">
        <div class="email">{{ user?.email }}</div>
        <div class="joined">{{ t('profile.joined', { date: joinedDate }) }}</div>
      </div>
    </div>

    <div class="heatmap-section">
      <div class="section-title">{{ t('profile.activity') }}</div>
      <div class="heatmap-card">
        <div class="heatmap-months">
          <div v-for="m in monthLabels" :key="m" class="month-label">{{ m }}</div>
        </div>
        <div class="heatmap-grid">
          <div
            v-for="(cell, idx) in heatmapCells"
            :key="idx"
            class="heatmap-cell"
            :class="`level-${cell.level}`"
            :title="cell.date ? `${cell.date}: ${cell.count} ${t('profile.tasks')}` : ''"
          />
        </div>
        <div class="heatmap-legend">
          <span>{{ t('profile.less') }}</span>
          <div class="legend-cell level-0" />
          <div class="legend-cell level-1" />
          <div class="legend-cell level-2" />
          <div class="legend-cell level-3" />
          <div class="legend-cell level-4" />
          <span>{{ t('profile.more') }}</span>
        </div>
      </div>
    </div>

    <div class="history-section">
      <div class="section-title">{{ t('profile.history') }}</div>
      <div
        ref="historyRef"
        class="history-scroll"
        @scroll="onHistoryScroll"
      >
        <div v-if="visibleGroups.length === 0" class="empty">{{ t('profile.empty') }}</div>
        <div v-for="group in visibleGroups" :key="group.date" class="day-group">
          <div class="day-label">{{ group.date }}</div>
          <div class="day-list">
            <div v-for="item in group.items" :key="item.id" class="history-item">
              <div class="item-dot" :class="item.outcome" />
              <div class="item-text" :class="{ 'is-completed': item.outcome === 'completed' }">{{ t(item.text_key) }}</div>
              <div class="item-meta">
                <span class="category">{{ t(`categories.${item.category}`) }}</span>
                <span class="duration">{{ formatDuration(item.duration_seconds) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-if="loadingMore" class="loading-more">{{ t('profile.loading') }}</div>
      </div>
    </div>

    <button class="btn-back" @click="router.replace('/')">{{ t('common.back') }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { getMe, getHistory, bindEmail } from '@/api/profile'
import type { User, HistoryItem } from '@/api/profile'

const router = useRouter()
const { t } = useI18n()
const auth = useAuthStore()

const user = ref<User | null>(null)
const history = ref<HistoryItem[]>([])
const historyRef = ref<HTMLElement | null>(null)

const displayLimit = ref(5)
const loadingMore = ref(false)

const showBind = ref(false)
const bindEmailVal = ref('')
const bindPasswordVal = ref('')
const binding = ref(false)

const avatarLetter = computed(() => {
  if (user.value?.is_temp) return '?'
  return user.value?.email?.[0]?.toUpperCase() ?? '?'
})

const joinedDate = computed(() => {
  if (!user.value?.created_at) return ''
  const d = new Date(user.value.created_at)
  return d.toLocaleDateString('zh-CN')
})

interface HeatmapCell {
  date: string
  count: number
  level: number
}

const heatmapCells = computed<HeatmapCell[]>(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const start = new Date(today)
  start.setDate(start.getDate() - 90)

  const countMap = new Map<string, number>()
  for (const item of history.value) {
    if (item.outcome !== 'completed') continue
    const d = new Date(item.created_at)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    countMap.set(key, (countMap.get(key) ?? 0) + 1)
  }

  const cells: HeatmapCell[] = []
  const dayOfWeek = start.getDay()
  for (let i = 0; i < dayOfWeek; i++) {
    cells.push({ date: '', count: 0, level: -1 })
  }

  for (let i = 0; i <= 90; i++) {
    const d = new Date(start)
    d.setDate(d.getDate() + i)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
    const count = countMap.get(key) ?? 0
    let level = 0
    if (count === 0) level = 0
    else if (count === 1) level = 1
    else if (count <= 3) level = 2
    else if (count <= 5) level = 3
    else level = 4
    cells.push({ date: key, count, level })
  }

  return cells
})

const monthLabels = computed(() => {
  const today = new Date()
  const labels: string[] = []
  for (let i = 2; i >= 0; i--) {
    const d = new Date(today.getFullYear(), today.getMonth() - i, 1)
    labels.push(`${d.getMonth() + 1}月`)
  }
  return labels
})

interface DayGroup {
  date: string
  items: HistoryItem[]
}

const groupedHistory = computed<DayGroup[]>(() => {
  const map = new Map<string, HistoryItem[]>()
  for (const item of history.value) {
    const d = new Date(item.created_at)
    const key = d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(item)
  }
  const groups: DayGroup[] = []
  for (const [date, items] of map) {
    // 未完成的在前，已完成的在后
    const sorted = [...items].sort((a, b) => {
      if (a.outcome === b.outcome) return 0
      return a.outcome === 'completed' ? 1 : -1
    })
    groups.push({ date, items: sorted })
  }
  return groups
})

const visibleGroups = computed(() => groupedHistory.value.slice(0, displayLimit.value))

function formatDuration(s: number) {
  if (s < 60) return `${s}${t('common.second')}`
  return `${Math.floor(s / 60)}${t('common.minute')}`
}

function onHistoryScroll() {
  const el = historyRef.value
  if (!el || loadingMore.value) return
  const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 20
  if (nearBottom && displayLimit.value < groupedHistory.value.length) {
    loadingMore.value = true
    setTimeout(() => {
      displayLimit.value += 5
      loadingMore.value = false
    }, 300)
  }
}

onMounted(async () => {
  try {
    const [u, h] = await Promise.all([getMe(), getHistory()])
    user.value = u
    history.value = h
  } catch (e) {
    console.error(e)
  }
})

async function handleBind() {
  binding.value = true
  try {
    const result = await bindEmail(bindEmailVal.value, bindPasswordVal.value)
    auth.setTokens(result.access_token, result.refresh_token)
    user.value = await getMe()
    showBind.value = false
    bindEmailVal.value = ''
    bindPasswordVal.value = ''
  } catch (e: unknown) {
    const err = e as { response?: { data?: { error?: string } } }
    alert(t('common.error', { msg: err.response?.data?.error ?? t('common.unknownError') }))
  } finally {
    binding.value = false
  }
}
</script>

<style scoped>
.profile-screen {
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

.profile-header {
  width: 100%;
  max-width: 600px;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.5rem;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: #2dd4bf;
  color: #0a0a0a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.email {
  font-size: 18px;
  font-weight: 600;
  color: #f0f0f0;
}

.joined {
  font-size: 12px;
  color: #666;
}

.heatmap-section {
  width: 100%;
  max-width: 600px;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #e0e0e0;
  margin-bottom: 0.5rem;
}

.heatmap-card {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  padding: 0.75rem;
}

.heatmap-months {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.month-label {
  font-size: 10px;
  color: #666;
  flex: 1;
  text-align: center;
}

.heatmap-grid {
  display: grid;
  grid-template-rows: repeat(7, 16px);
  grid-auto-flow: column;
  gap: 2px;
  overflow-x: auto;
  padding-bottom: 2px;
}

.heatmap-cell {
  width: 16px;
  height: 16px;
  border-radius: 2px;
}

.heatmap-cell.level--1 {
  background: transparent;
}

.heatmap-cell.level-0 {
  background: #2a2a2a;
}

.heatmap-cell.level-1 {
  background: rgba(45, 212, 191, 0.25);
}

.heatmap-cell.level-2 {
  background: rgba(45, 212, 191, 0.5);
}

.heatmap-cell.level-3 {
  background: rgba(45, 212, 191, 0.75);
}

.heatmap-cell.level-4 {
  background: #2dd4bf;
}

.heatmap-legend {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 3px;
  margin-top: 6px;
  font-size: 10px;
  color: #666;
}

.legend-cell {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}

.history-section {
  width: 100%;
  max-width: 600px;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  margin-bottom: 0.75rem;
}

.history-scroll {
  flex: 1;
  overflow-y: auto;
  min-height: 0;
}

.empty {
  font-size: 13px;
  color: #555;
  text-align: center;
  padding: 1.5rem 0;
}

.day-group {
  margin-bottom: 0.75rem;
}

.day-label {
  font-size: 12px;
  font-weight: 600;
  color: #888;
  margin-bottom: 0.35rem;
  padding-left: 4px;
}

.day-list {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  padding: 0.35rem 0.75rem;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid #222;
}

.history-item:last-child {
  border-bottom: none;
}

.item-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;
}

.item-dot.completed {
  background: #2dd4bf;
}

.item-dot.skipped {
  background: #555;
}

.item-text {
  flex: 1;
  font-size: 13px;
  color: #e0e0e0;
  line-height: 1.3;
}

.item-text.is-completed {
  text-decoration: line-through;
  color: #666;
}

.item-meta {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  flex-shrink: 0;
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

.loading-more {
  text-align: center;
  font-size: 12px;
  color: #666;
  padding: 0.75rem 0;
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
  margin-bottom: 0.5rem;
  cursor: pointer;
  flex-shrink: 0;
}

.btn-back:hover {
  border-color: #444;
  color: #aaa;
}

.guest-banner {
  width: 100%;
  max-width: 600px;
  background: rgba(45, 212, 191, 0.08);
  border: 1px solid rgba(45, 212, 191, 0.25);
  border-radius: 12px;
  padding: 0.75rem 1rem;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.banner-hint {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.banner-hint span {
  font-size: 13px;
  color: #aaa;
  line-height: 1.4;
}

.btn-bind-toggle {
  font-size: 13px;
  font-weight: 600;
  color: #2dd4bf;
  background: transparent;
  border: 1px solid rgba(45, 212, 191, 0.4);
  border-radius: 8px;
  padding: 0.3rem 0.75rem;
  cursor: pointer;
  flex-shrink: 0;
}

.btn-bind-toggle:hover {
  background: rgba(45, 212, 191, 0.1);
}

.bind-form {
  margin-top: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.bind-input {
  width: 100%;
  padding: 0.6rem 0.75rem;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 8px;
  color: #f0f0f0;
  font-size: 14px;
  box-sizing: border-box;
}

.bind-input:focus {
  outline: none;
  border-color: rgba(45, 212, 191, 0.5);
}

.btn-bind-confirm {
  width: 100%;
  padding: 0.65rem;
  background: #2dd4bf;
  color: #0a0a0a;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
}

.btn-bind-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
