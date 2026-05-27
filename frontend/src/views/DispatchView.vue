<template>
  <div class="dispatch-screen">
    <div v-if="counting" class="countdown-overlay">{{ countdown }}</div>

    <div class="gesture-hint-bar">
      <span class="hint-item hint-pull">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M12 5v14M5 12l7 7 7-7"/>
        </svg>
        {{ t('dispatch.pullToStart') }}
      </span>
      <span class="hint-item hint-swipe">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <path d="M5 12h14M12 5l7 7-7 7"/>
        </svg>
        {{ t('dispatch.swipeToChange') }}
      </span>
    </div>

    <div class="card-wrapper">
      <div
        class="task-card"
        :class="{
          'is-dragging': cardState === 'dragging',
          'is-entering': cardState === 'entering'
        }"
        :style="cardStyle"
        @touchstart="onTouchStart"
        @touchmove.prevent="onTouchMove"
        @touchend="onTouchEnd"
        @touchcancel="onTouchEnd"
        @mousedown="onMouseDown"
      >
        <div class="pull-hint" :class="{ 'visible': pullHintVisible }">
          <div class="arrow-down"></div>
          <span>{{ t('dispatch.releaseToStart') }}</span>
        </div>

        <div class="swipe-hint" :class="{ 'visible': swipeHintVisible }">
          <span>{{ t('dispatch.swipeReleaseToChange') }}</span>
          <div class="arrow-right"></div>
        </div>

        <div class="badge">{{ categoryLabel }}</div>
        <div class="task-text">{{ taskText }}</div>
        <div class="duration">{{ durationLabel }}</div>

        <div class="gesture-icons">
          <div class="gesture-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12l7 7 7-7"/>
            </svg>
          </div>
          <div class="gesture-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div class="action-area">
      <button class="btn-custom" @click="router.push('/custom-task')">
        {{ t('dispatch.customTask') }}
      </button>
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

type CardState = 'idle' | 'dragging' | 'exiting' | 'entering'
const cardState = ref<CardState>('idle')
const cardOffset = ref({ x: 0, y: 0 })

const PULL_THRESHOLD = 80
const SWIPE_THRESHOLD = 100
const MAX_PULL = 150
const MAX_SWIPE = 220

let startX = 0
let startY = 0
let currentDirection: 'pull' | 'swipe' | null = null
let isPointerDown = false

const pullHintVisible = computed(() => {
  return currentDirection === 'pull' && cardOffset.value.y > PULL_THRESHOLD * 0.5
})

const swipeHintVisible = computed(() => {
  return currentDirection === 'swipe' && cardOffset.value.x > SWIPE_THRESHOLD * 0.5
})

const cardStyle = computed(() => {
  if (cardState.value === 'entering') {
    return {}
  }
  if (cardState.value === 'exiting') {
    return {
      transform: 'translateX(120%) rotate(4deg)',
      opacity: '0',
      transition: 'transform 0.35s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s'
    }
  }
  const isDragging = cardState.value === 'dragging'
  return {
    transform: `translate(${cardOffset.value.x}px, ${cardOffset.value.y}px) scale(${isDragging ? 0.97 : 1})`,
    opacity: isDragging ? '0.85' : '1',
    transition: isDragging ? 'none' : 'transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s'
  }
})

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

// ===== Pointer Events =====

function onTouchStart(e: TouchEvent) {
  const touch = e.touches[0]
  beginGesture(touch.clientX, touch.clientY)
}

function onMouseDown(e: MouseEvent) {
  if (e.button !== 0) return
  beginGesture(e.clientX, e.clientY)
  document.addEventListener('mousemove', onMouseMove)
  document.addEventListener('mouseup', onMouseUp)
}

function beginGesture(x: number, y: number) {
  isPointerDown = true
  startX = x
  startY = y
  currentDirection = null
  cardState.value = 'idle'
  cardOffset.value = { x: 0, y: 0 }
}

function onTouchMove(e: TouchEvent) {
  const touch = e.touches[0]
  updateGesture(touch.clientX, touch.clientY)
}

function onMouseMove(e: MouseEvent) {
  if (!isPointerDown) return
  updateGesture(e.clientX, e.clientY)
}

function updateGesture(x: number, y: number) {
  if (!isPointerDown) return

  const dx = x - startX
  const dy = y - startY

  if (currentDirection === null) {
    if (dy > 8 && dy > Math.abs(dx) * 0.8) {
      currentDirection = 'pull'
    } else if (dx > 8 && dx > Math.abs(dy) * 0.8) {
      currentDirection = 'swipe'
    }
  }

  if (currentDirection === 'pull') {
    if (dy > 0) {
      cardState.value = 'dragging'
      cardOffset.value = {
        x: 0,
        y: Math.min(dy * 0.55, MAX_PULL)
      }
    }
  } else if (currentDirection === 'swipe') {
    if (dx > 0) {
      cardState.value = 'dragging'
      cardOffset.value = {
        x: Math.min(dx * 0.7, MAX_SWIPE),
        y: 0
      }
    }
  }
}

function onTouchEnd() {
  endGesture()
}

function onMouseUp() {
  endGesture()
  document.removeEventListener('mousemove', onMouseMove)
  document.removeEventListener('mouseup', onMouseUp)
}

function endGesture() {
  if (!isPointerDown) return
  isPointerDown = false

  const wasPull = currentDirection === 'pull'
  const wasSwipe = currentDirection === 'swipe'

  if (wasPull && cardOffset.value.y >= PULL_THRESHOLD) {
    cardState.value = 'idle'
    cardOffset.value = { x: 0, y: 0 }
    currentDirection = null
    handleStart()
    return
  }

  if (wasSwipe && cardOffset.value.x >= SWIPE_THRESHOLD) {
    handleSwap()
    currentDirection = null
    return
  }

  // Bounce back
  cardState.value = 'idle'
  cardOffset.value = { x: 0, y: 0 }
  currentDirection = null
}

// ===== Actions =====

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

async function handleSwap() {
  cardState.value = 'exiting'
  cardOffset.value = { x: 0, y: 0 }

  await new Promise(resolve => setTimeout(resolve, 350))

  await mission.fetchTask()

  cardState.value = 'entering'
  cardOffset.value = { x: 0, y: 0 }

  setTimeout(() => {
    cardState.value = 'idle'
  }, 350)
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
  gap: 1.5rem;
}

/* Countdown overlay */
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

/* Gesture hint bar */
.gesture-hint-bar {
  display: flex;
  justify-content: center;
  gap: 24px;
  font-size: 11px;
  color: #555;
  margin-bottom: -4px;
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.hint-pull {
  color: #2dd4bf;
}

.hint-swipe {
  color: #f59e0b;
}

/* Card wrapper */
.card-wrapper {
  width: 85%;
  overflow: hidden;
  display: flex;
  justify-content: center;
}

/* Task card */
.task-card {
  width: 100%;
  max-width: 400px;
  background: #1a1a1a;
  border-radius: 20px;
  padding: 2.5rem 1.5rem;
  text-align: center;
  border: 1px solid #2a2a2a;
  position: relative;
  user-select: none;
  touch-action: none;
  cursor: grab;
}

.task-card:active {
  cursor: grabbing;
}

.task-card.is-entering {
  animation: slideInFromLeft 0.35s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

@keyframes slideInFromLeft {
  from {
    transform: translateX(-110%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* Pull hint */
.pull-hint {
  position: absolute;
  top: -28px;
  left: 50%;
  transform: translateX(-50%);
  color: #2dd4bf;
  font-size: 12px;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.2s, top 0.2s;
  white-space: nowrap;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.pull-hint.visible {
  opacity: 1;
  top: -36px;
}

.arrow-down {
  width: 0;
  height: 0;
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 6px solid #2dd4bf;
}

/* Swipe hint */
.swipe-hint {
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%) translateX(80%);
  color: #f59e0b;
  font-size: 12px;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.2s, transform 0.2s;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 6px;
}

.swipe-hint.visible {
  opacity: 1;
  transform: translateY(-50%) translateX(100%);
}

.arrow-right {
  width: 0;
  height: 0;
  border-top: 5px solid transparent;
  border-bottom: 5px solid transparent;
  border-left: 6px solid #f59e0b;
}

/* Badge */
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

/* Task text */
.task-text {
  font-size: 24px;
  font-weight: 600;
  line-height: 1.45;
  color: #f0f0f0;
}

/* Duration */
.duration {
  font-size: 14px;
  color: #666;
  margin-top: 1.25rem;
}

/* Gesture icons */
.gesture-icons {
  margin-top: 1.5rem;
  display: flex;
  justify-content: center;
  gap: 36px;
}

.gesture-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #3a3a3a;
}

.gesture-icon svg {
  width: 20px;
  height: 20px;
  stroke: #3a3a3a;
}

/* Action area */
.action-area {
  width: 85%;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

/* Custom task button */
.btn-custom {
  width: 100%;
  padding: 0.85rem;
  background: transparent;
  color: #888;
  border-radius: 14px;
  font-size: 15px;
  border: 1px solid #333;
  cursor: pointer;
  transition: border-color 0.2s, color 0.2s;
}

.btn-custom:hover {
  border-color: #555;
  color: #aaa;
}
</style>
