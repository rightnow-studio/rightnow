import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAdminAuthStore } from '@/stores/adminAuth'
import LoginView from '@/views/LoginView.vue'
import DispatchView from '@/views/DispatchView.vue'
import MissionView from '@/views/MissionView.vue'
import VictoryView from '@/views/VictoryView.vue'
import CustomTaskView from '@/views/CustomTaskView.vue'
import ProfileView from '@/views/ProfileView.vue'
import TaskListView from '@/views/TaskListView.vue'
import AdminLoginView from '@/views/AdminLoginView.vue'
import AdminView from '@/views/AdminView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/login', name: 'login', component: LoginView, meta: { public: true } },
    { path: '/', name: 'dispatch', component: DispatchView },
    { path: '/mission', name: 'mission', component: MissionView },
    { path: '/victory', name: 'victory', component: VictoryView },
    { path: '/custom-task', name: 'custom-task', component: CustomTaskView },
    { path: '/my-tasks', name: 'my-tasks', component: TaskListView },
    { path: '/profile', name: 'profile', component: ProfileView },
    { path: '/admin/login', name: 'admin-login', component: AdminLoginView, meta: { adminPublic: true } },
    { path: '/admin', name: 'admin', component: AdminView, meta: { requiresAdmin: true } },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const adminAuth = useAdminAuthStore()

  // admin 路由守卫
  if (to.meta.requiresAdmin) {
    if (!adminAuth.isLoggedIn) return '/admin/login'
    return
  }
  if (to.path === '/admin/login' && adminAuth.isLoggedIn) {
    return '/admin'
  }

  // 普通用户路由守卫
  if (to.path === '/login' && auth.token) {
    return '/'
  }
})

export default router
