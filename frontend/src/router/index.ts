import {createRouter, createWebHistory} from 'vue-router'
import {useUserStore} from '@/stores/user'
import type {RouteRecordRaw} from 'vue-router'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Index',
        redirect: '/profile'
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/auth/Login.vue'),
        meta: {guest: true}
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue'),
        meta: {guest: true}
    },
    {
        path: '/admin',
        component: () => import('@/layouts/AdminLayout.vue'),
        meta: {requiresAuth: true, requiresAdmin: true},
        children: [
            {
                path: 'users',
                name: 'UserManagement',
                component: () => import('@/views/admin/UserManagement.vue'),
                meta: {title: '用户管理'}
            },
            {
                path: 'profile',
                name: 'AdminProfile',
                component: () => import('@/views/profile/UserProfile.vue'),
                meta: {title: '个人中心'}
            },
            {
                path: 'logs',
                name: 'SystemLog',
                component: () => import('@/views/admin/SystemLog.vue'),
                meta: {title: '系统日志'}
            }
        ]
    },
    {
        path: '/profile',
        name: 'UserProfile',
        component: () => import('@/views/profile/UserProfile.vue'),
        meta: {requiresAuth: true, title: '个人中心'}
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFound.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach(async (to, _from, next) => {
    const userStore = useUserStore()

    // 如果是游客页面且用户已登录则跳转到首页
    if (to.meta.guest && userStore.isAuthenticated) {
        return next('/')
    }

    if (to.meta.requiresAuth || to.meta.requiresAdmin) {
        if (!userStore.userInfo && userStore.token) {
            try {
                const res = await userStore.fetchUserInfo()
                if (!res.success || !userStore.userInfo) {
                    alert('认证信息失效，请重新登录！')
                    userStore.logout()
                    return next('/login')
                }
            } catch (error) {
                console.error('路由守卫获取用户信息失败:', error)
                alert('网络错误或认证失败，请重新登录！')
                userStore.logout()
                return next('/login')
            }
        } else if (!userStore.token) {
            alert('您尚未登录，请先登录。')
            return next('/login')
        }
    }

    if (to.meta.requiresAdmin && !userStore.isAdmin) {
        alert('您没有权限访问此页面！')
        return next('/')
    }

    next()
})

router.onError((error, to, from) => {
    console.error('Vue Router 导航错误:', error)
    console.error('跳转目标:', to)
    console.error('跳转来源:', from)
})

export default router
