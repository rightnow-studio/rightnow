<template>
  <div class="hud-screen">
    <div class="hud-task">{{ taskText }}</div>
    <div class="timer-ring">
      <svg viewBox="0 0 100 100" class="ring-svg">
        <circle cx="50" cy="50" r="45" fill="none" stroke="#2a2a2a" stroke-width="4" />
        <circle
          cx="50" cy="50" r="45" fill="none" stroke="#2dd4bf" stroke-width="4"
          stroke-dasharray="283" :stroke-dashoffset="progressOffset"
          stroke-linecap="round" transform="rotate(-90 50 50)"
        />
      </svg>
      <div class="time">{{ formattedTime }}</div>
    </div>
    <div class="hud-buttons">
      <button class="btn-done" @click="handleDone">{{ t('mission.done') }}</button>
      <button class="btn-skip" @click="handleSkip">{{ t('mission.skip') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useMissionStore } from '@/stores/mission'

const { t } = useI18n()
const router = useRouter()
const mission = useMissionStore()

const taskText = computed(() => {
  if (!mission.currentTask) return ''
  return t(mission.currentTask.text_key)
})

const formattedTime = computed(() => {
  const m = Math.floor(mission.timeLeft / 60)
  const s = mission.timeLeft % 60
  return `${m}:${s.toString().padStart(2, '0')}`
})

const progressOffset = computed(() => {
  if (!mission.currentTask) return 283
  const total = mission.currentTask.duration_seconds
  const ratio = mission.timeLeft / total
  return 283 * ratio
})

let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  if (!mission.currentTask) {
    router.replace('/')
    return
  }
  mission.startTimer()
  timer = setInterval(() => {
    mission.tick()
    if (mission.timeLeft <= 0 && timer) {
      clearInterval(timer)
    }
  }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

async function handleDone() {
  if (timer) clearInterval(timer)
  await mission.complete('completed')
  router.push('/victory')
}

async function handleSkip() {
  if (timer) clearInterval(timer)
  await mission.complete('skipped')
  router.replace('/')
}
</script>

<style scoped>
.hud-screen {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #121212;
  position: relative;
}
.hud-task {
  position: absolute;
  top: 90px;
  width: 80%;
  text-align: center;
  font-size: 18px;
  font-weight: 500;
  color: #e0e0e0;
  line-height: 1.5;
}
.timer-ring {
  width: 200px;
  height: 200px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ring-svg {
  position: absolute;
  width: 100%;
  height: 100%;
}
.time {
  font-size: 40px;
  font-weight: 700;
  color: #f0f0f0;
  z-index: 1;
}
.hud-buttons {
  position: absolute;
  bottom: 70px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}
.btn-done {
  width: 72%;
  padding: 1.1rem;
  background: #059669;
  color: #fff;
  border-radius: 14px;
  font-size: 18px;
  font-weight: 600;
}
.btn-skip {
  font-size: 14px;
  color: #555;
  background: none;
}
</style>
