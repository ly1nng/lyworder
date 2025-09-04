import { useUserStore } from '@/stores/user'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('@/views/HomeView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/tickets',
        name: 'MyTickets',
        component: () => import('@/views/MyTicketsView.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/tickets/detail',
        name: 'TicketDetail',
        component: () => import('@/views/TicketDetailView.vue'),
        meta: { requiresAuth: true },
        props: (route) => ({ id: route.query.id })
    },
    {
        path: '/operator',
        name: 'MyOperator',
        component: () => import('@/views/MyOperatorView.vue'),
        meta: { requiresAuth: true, roles: ['admin', 'operator'] }
    },
    {
        path: '/admin',
        name: 'Admin',
        component: () => import('@/views/AdminView.vue'),
        meta: { requiresAuth: true, roles: ['admin'] }
    },
    {
        path: '/admin/users',
        name: 'UserManagement',
        component: () => import('@/views/UserManagementView.vue'),
        meta: { requiresAuth: true, roles: ['admin'] }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/LoginView.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/NotFoundView.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { top: 0 }
        }
    }
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()

    // 如果用户信息还没有加载，先尝试获取
    if (!userStore.user && userStore.isLoggedIn === null) {
        try {
            await userStore.fetchUserInfo()
        } catch (error) {
            // 获取用户信息失败，可能未登录
            console.log('用户未登录或登录已过期')
        }
    }

    // 检查是否需要认证
    if (to.meta.requiresAuth) {
        if (!userStore.isLoggedIn) {
            // 未登录，重定向到登录页面并跳转到后端登录
            window.location.href = '/login'
            return
        }

        // 检查角色权限
        if (to.meta.roles && !to.meta.roles.includes(userStore.user?.role)) {
            // 权限不足，重定向到首页
            next({ name: 'Home' })
            return
        }
    }

    // 如果已登录访问登录页，重定向到首页
    if (to.name === 'Login' && userStore.isLoggedIn) {
        next({ name: 'Home' })
        return
    }

    next()
})

export default router