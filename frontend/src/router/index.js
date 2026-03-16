import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../store/user'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/checkin',
    name: 'Checkin',
    component: () => import('../views/Checkin.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/rewards',
    name: 'Rewards',
    component: () => import('../views/Rewards.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/exchange',
    name: 'Exchange',
    component: () => import('../views/Exchange.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/child-rewards',
    name: 'ChildRewards',
    component: () => import('../views/ChildRewards.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/vocabulary',
    name: 'Vocabulary',
    component: () => import('../views/Vocabulary.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/Admin.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const isAuthenticated = userStore.isAuthenticated
  const isAdmin = userStore.user?.user_type === 1

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (to.meta.requiresAdmin && !isAdmin) {
    // 如果不是管理员，重定向到 home
    next('/home')
  } else {
    next()
  }
})

export default router