<template>
  <div class="dispatch-screen">
    <div v-if="counting" class="countdown-overlay">{{ countdown }}</div>
    <div class="task-card">
      <div class="badge">{{ categoryLabel }}</div>
      <div class="task-text">{{ taskText }}</div>
      <div class="duration">{{ durationLabel }}</div>
    </div>
    <div class="action-area">
      <button class="btn-start" @click="handleStart">{{ t('dispatch.start') }}</button>
      <button class="btn-custom" @click="router.push('/custom-task')">{{ t('dispatch.customTask') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useMissionStore } from '@/stores/mission'

const { t } = useI18n()
const router = useRouter()
const mission = useMissionStore()

const counting = ref(false)
const countdown = ref(3)

const taskText = computed(() => {
  if (!mission.currentTask) return ''
  return t(mission.currentTask.text_key)
})

const categoryLabel = computed(() => {
  if (!mission.currentTask) return ''
  return t(`categories.${mission.currentTask.category}`)
})

const durationLabel = computed(() => {
  if (!mission.currentTask) return ''
  const s = mission.currentTask.duration_seconds
  if (s < 60) return `${s}${t('common.second')}`
  return `${Math.floor(s / 60)}${t('common.minute')}`
})

onMounted(async () => {
  if (!mission.currentTask) {
    await mission.fetchTask()
  }
})

function handleStart() {
  if (counting.value) return
  counting.value = true
  countdown.value = 3
  const iv = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(iv)
      counting.value = false
      router.push('/mission')
    }
  }, 1000)
}
</script>

<style scoped>
.dispatch-screen {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #121212;
  position: relative;
  gap: 2rem;
}
.countdown-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.88);
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 120px;
  font-weight: 700;
  color: #f0f0f0;
}
.task-card {
  width: 85%;
  background: #1a1a1a;
  border-radius: 20px;
  padding: 2.5rem 1.5rem;
  text-align: center;
  border: 1px solid #2a2a2a;
}
.badge {
  display: inline-block;
  padding: 5px 14px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  background: rgba(45, 212, 191, 0.12);
  color: #2dd4bf;
  margin-bottom: 1.5rem;
}
.task-text {
  font-size: 24px;
  font-weight: 600;
  line-height: 1.45;
  color: #f0f0f0;
}
.duration {
  font-size: 14px;
  color: #666;
  margin-top: 1.25rem;
}
.action-area {
  width: 85%;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.btn-start {
  width: 100%;
  padding: 1rem;
  background: #2dd4bf;
  color: #0a0a0a;
  border-radius: 14px;
  font-size: 18px;
  font-weight: 600;
}
.btn-custom {
  width: 100%;
  padding: 0.85rem;
  background: transparent;
  color: #888;
  border-radius: 14px;
  font-size: 15px;
  border: 1px solid #333;
}
</style>
