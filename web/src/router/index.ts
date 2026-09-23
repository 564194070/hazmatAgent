import { createRouter, createWebHistory } from 'vue-router'
import { useAuth } from '../composables/useAuth.ts'
import ChatView from '../views/ChatView.vue'
import LoginView from '../views/LoginView.vue'

// 给路由 meta 补上 requiresAuth，模板和守卫里才能类型安全地读它
declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
  }
}

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/',
      name: 'chat',
      component: ChatView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const { isLoggedIn } = useAuth()

  // 聊天页必须已登录。redirect 记下来，登录成功后回到用户原本要去的地址
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  // 已登录就不要再停在登录页
  if (to.path === '/login' && isLoggedIn.value) {
    return { path: '/' }
  }
})

export default router