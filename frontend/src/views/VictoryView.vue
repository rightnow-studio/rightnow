<template>
  <div class="victory-screen">
    <div class="check">&#10003;</div>
    <h2>{{ t('victory.title') }}</h2>
    <div class="sub">{{ taskSummary }}</div>
    <button class="btn-next" @click="handleNext">{{ t('victory.next') }}</button>
    <button class="btn-end" @click="handleEnd">{{ t('victory.end') }}</button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useMissionStore } from '@/stores/mission'

const { t } = useI18n()
const router = useRouter()
const mission = useMissionStore()

const taskSummary = computed(() => {
  if (!mission.currentTask) return ''
  const text = t(`task.${mission.currentTask.text_key}`)
  const s = mission.currentTask.duration_seconds
  const d = s < 60 ? `${s}秒` : `${Math.floor(s / 60)}分钟`
  return `${text} — ${d}`
})

function handleNext() {
  router.replace('/')
}

function handleEnd() {
  mission.currentTask = null
  router.replace('/')
}
</script>

<style scoped>
.victory-screen {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #121212;
  padding: 1.5rem;
}
.check {
  font-size: 80px;
  color: #2dd4bf;
  margin-bottom: 0.5rem;
  font-weight: 700;
}
h2 {
  font-size: 28px;
  font-weight: 700;
  color: #f0f0f0;
  margin-bottom: 6px;
}
.sub {
  font-size: 14px;
  color: #888;
  margin-bottom: 3rem;
  text-align: center;
  line-height: 1.5;
}
.btn-next {
  width: 72%;
  padding: 1rem;
  background: #2dd4bf;
  color: #0a0a0a;
  border-radius: 14px;
  font-size: 17px;
  font-weight: 600;
  margin-bottom: 0.75rem;
}
.btn-end {
  width: 72%;
  padding: 1rem;
  background: transparent;
  color: #888;
  border: 1px solid #333;
  border-radius: 14px;
  font-size: 16px;
}
</style>
