import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '@/api/client'

export interface Task {
  id: string
  text_key: string
  duration_seconds: number
  category: string
}

export const useMissionStore = defineStore('mission', () => {
  const currentTask = ref<Task | null>(null)
  const timeLeft = ref(0)
  const isRunning = ref(false)

  async function fetchTask() {
    const { data } = await client.get('/dispatch')
    currentTask.value = data
    timeLeft.value = data.duration_seconds
    isRunning.value = false
    return data
  }

  function startTimer() {
    isRunning.value = true
  }

  function stopTimer() {
    isRunning.value = false
  }

  function tick() {
    if (timeLeft.value > 0) {
      timeLeft.value--
    }
    if (timeLeft.value <= 0) {
      isRunning.value = false
    }
  }

  function setTask(task: Task) {
    currentTask.value = task
    timeLeft.value = task.duration_seconds
    isRunning.value = false
  }

  async function complete(outcome: 'completed' | 'skipped') {
    if (!currentTask.value) return
    await client.post('/history', { task_id: currentTask.value.id, outcome })
    currentTask.value = null
    timeLeft.value = 0
    isRunning.value = false
  }

  return { currentTask, timeLeft, isRunning, fetchTask, setTask, startTimer, stopTimer, tick, complete }
})
