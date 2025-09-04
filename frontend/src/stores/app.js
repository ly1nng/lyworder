import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
    const loading = ref(false)
    const isMobile = ref(window.innerWidth < 768)

    // 设置加载状态
    const setLoading = (value) => {
        loading.value = value
    }

    // 更新移动端状态
    const updateMobile = () => {
        isMobile.value = window.innerWidth < 768
    }

    // 监听窗口大小变化
    const initResize = () => {
        window.addEventListener('resize', updateMobile)
    }

    // 清理事件监听
    const cleanup = () => {
        window.removeEventListener('resize', updateMobile)
    }

    return {
        loading,
        isMobile,
        setLoading,
        updateMobile,
        initResize,
        cleanup
    }
})