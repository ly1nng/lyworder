import api from '@/api'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useUserStore = defineStore('user', () => {
    const user = ref(null)
    const isLoggedIn = ref(null) // null表示未知状态，true表示已登录，false表示未登录

    // 计算属性
    const isAdmin = computed(() => user.value?.role === 'admin')
    const isOperator = computed(() => user.value?.role === 'operator' || user.value?.role === 'admin')
    const userName = computed(() => user.value?.user_name || '')

    // 获取用户信息
    const fetchUserInfo = async () => {
        try {
            const response = await api.user.getCurrentUser()
            user.value = response.data
            isLoggedIn.value = true
            return response.data
        } catch (error) {
            user.value = null
            isLoggedIn.value = false
            throw error
        }
    }

    // 登出
    const logout = async () => {
        try {
            user.value = null
            isLoggedIn.value = false
            // 重定向到后端登出接口
            window.location.href = '/logout'
        } catch (error) {
            console.error('登出失败:', error)
        }
    }

    // 登录（重定向到后端登录）
    const login = () => {
        window.location.href = '/login'
    }

    // 重置状态
    const reset = () => {
        user.value = null
        isLoggedIn.value = null
    }

    return {
        user,
        isLoggedIn,
        isAdmin,
        isOperator,
        userName,
        fetchUserInfo,
        logout,
        login,
        reset
    }
})