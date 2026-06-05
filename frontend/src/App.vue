<template>
  <SplashScreen v-if="showSplash" @done="showSplash = false" />
  <router-view />
  <button
    v-if="showMyTasksBtn"
    class="my-tasks-btn"
    @click="$router.push('/my-tasks')"
    aria-label="个人任务列表"
  >
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <line x1="8" y1="6" x2="21" y2="6" />
      <line x1="8" y1="12" x2="21" y2="12" />
      <line x1="8" y1="18" x2="21" y2="18" />
      <line x1="3" y1="6" x2="3.01" y2="6" />
      <line x1="3" y1="12" x2="3.01" y2="12" />
      <line x1="3" y1="18" x2="3.01" y2="18" />
    </svg>
  </button>
  <button
    v-if="showCreateTaskBtn"
    class="create-task-btn"
    @click="$router.push('/custom-task')"
    aria-label="创建个人任务"
  >
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 5v14M5 12h14" />
    </svg>
  </button>
  <button
    v-if="showProfileBtn"
    class="profile-btn"
    @click="$router.push('/profile')"
    aria-label="个人页"
  >
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
      <circle cx="12" cy="7" r="4" />
    </svg>
  </button>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import SplashScreen from '@/components/SplashScreen.vue'

const route = useRoute()
const auth = useAuthStore()
const showSplash = ref(true)

const showProfileBtn = computed(() => {
  return auth.isLoggedIn && route.path !== '/profile' && route.path !== '/login' && route.path !== '/my-tasks'
})

const showMyTasksBtn = computed(() => {
  return auth.isLoggedIn && route.path === '/'
})

const showCreateTaskBtn = computed(() => {
  return auth.isLoggedIn && route.path === '/my-tasks'
})
</script>

<style>
.profile-btn {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 50;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}
.profile-btn:hover {
  color: #2dd4bf;
  border-color: #2dd4bf;
}
.profile-btn svg {
  width: 20px;
  height: 20px;
}

.my-tasks-btn {
  position: fixed;
  top: 16px;
  right: 64px;
  z-index: 50;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}
.my-tasks-btn:hover {
  color: #2dd4bf;
  border-color: #2dd4bf;
}
.my-tasks-btn svg {
  width: 20px;
  height: 20px;
}

.create-task-btn {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 50;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  color: #888;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}
.create-task-btn:hover {
  color: #2dd4bf;
  border-color: #2dd4bf;
}
.create-task-btn svg {
  width: 20px;
  height: 20px;
}
</style>
