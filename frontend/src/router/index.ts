import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import LoginView from '@/views/LoginView.vue'
import DispatchView from '@/views/DispatchView.vue'
import MissionView from '@/views/MissionView.vue'
import VictoryView from '@/views/VictoryView.vue'
import CustomTaskView from '@/views/CustomTaskView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView, meta: { public: true } },
    { path: '/', name: 'dispatch', component: DispatchView },
    { path: '/mission', name: 'mission', component: MissionView },
    { path: '/victory', name: 'victory', component: VictoryView },
    { path: '/custom-task', name: 'custom-task', component: CustomTaskView },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.token) {
    return '/login'
  }
  if (to.path === '/login' && auth.token) {
    return '/'
  }
})

export default router
