<template>
  <div class="custom-screen">
    <h2>{{ t('customTask.title') }}</h2>
    <label>{{ t('customTask.content') }}</label>
    <input v-model="content" :placeholder="t('customTask.placeholder')" />
    <label>{{ t('customTask.duration') }}</label>
    <select v-model="duration">
      <option :value="30">{{ t('durations.30') }}</option>
      <option :value="60">{{ t('durations.60') }}</option>
      <option :value="120">{{ t('durations.120') }}</option>
      <option :value="180">{{ t('durations.180') }}</option>
      <option :value="300">{{ t('durations.300') }}</option>
    </select>
    <label>{{ t('customTask.category') }}</label>
    <select v-model="category">
      <option value="work">{{ t('categories.work') }}</option>
      <option value="home">{{ t('categories.home') }}</option>
      <option value="communication">{{ t('categories.communication') }}</option>
      <option value="health">{{ t('categories.health') }}</option>
      <option value="admin">{{ t('categories.admin') }}</option>
    </select>
    <button class="btn-save" @click="handleSave">{{ t('customTask.save') }}</button>
    <p class="hint">{{ t('customTask.hint') }}</p>
    <button class="btn-back" @click="router.replace('/')">&larr; {{ t('common.back') }}</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import client from '@/api/client'

const { t } = useI18n()
const router = useRouter()

const content = ref('')
const duration = ref(60)
const category = ref('work')

async function handleSave() {
  if (!content.value.trim()) return
  try {
    await client.post('/tasks', {
      text_key: content.value.trim(),
      duration_seconds: duration.value,
      category: category.value,
    })
    router.replace('/')
  } catch (e) {
    const msg = (e as any).response?.data?.error || t('common.unknownError')
    alert(t('common.error', { msg }))
  }
}
</script>

<style scoped>
.custom-screen {
  width: 100%;
  height: 100%;
  background: #121212;
  color: #f0f0f0;
  padding: 2rem 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  justify-content: flex-start;
  overflow-y: auto;
}
h2 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: #f0f0f0;
}
label {
  font-size: 12px;
  color: #888;
  margin-bottom: 0.4rem;
  display: block;
}
input, select {
  width: 100%;
  padding: 0.75rem 1rem;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 10px;
  color: #f0f0f0;
  font-size: 15px;
  margin-bottom: 1rem;
  outline: none;
}
input:focus, select:focus {
  border-color: #2dd4bf;
}
.btn-save {
  width: 100%;
  padding: 0.9rem;
  background: #2dd4bf;
  color: #0a0a0a;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  margin-top: 0.5rem;
}
.hint {
  font-size: 12px;
  color: #555;
  margin-top: 1rem;
  text-align: center;
}
.btn-back {
  margin-top: 1.5rem;
  font-size: 14px;
  color: #666;
  background: none;
  text-align: center;
}
</style>
