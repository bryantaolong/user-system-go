import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { Message } from '@arco-design/web-vue'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/profile',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { guest: true },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/register/index.vue'),
    meta: { guest: true },
  },
  {
    path: '/admin',
    component: () => import('@/components/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/users',
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/views/admin/users/UserManagement.vue'),
      },
      {
        path: 'profile',
        name: 'AdminProfile',
        component: () => import('@/views/profile/index.vue'),
      },
      {
        path: 'logs',
        name: 'SystemLog',
        component: () => import('@/views/admin/logs/SystemLog.vue'),
      },
    ],
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/profile/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach(async (to, from, next) => {
  const store = useUserStore()

  const guestPaths = ['/login', '/register']

  // Authenticated user visiting guest page -> redirect to home
  if (guestPaths.includes(to.path) && store.isAuthenticated) {
    return next('/')
  }

  // Auth required routes
  if (to.meta.requiresAuth && !store.isAuthenticated) {
    Message.error('您尚未登录，请先登录。')
    return next('/login')
  }

  // Fetch user info if authenticated but no user loaded
  if (store.isAuthenticated && !store.user) {
    const res = await store.fetchUserInfo()
    if (!res.success) {
      Message.error('认证信息失效，请重新登录！')
      store.clearToken()
      return next('/login')
    }
  }

  // Admin required
  if (to.meta.requiresAdmin && !store.isAdmin) {
    Message.error('您没有权限访问此页面！')
    return next('/')
  }

  next()
})

export default router
