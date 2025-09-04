import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
    baseURL: '/api/v1',
    timeout: 30000,
    withCredentials: true, // 支持cookie认证
    headers: {
        'Content-Type': 'application/json'
    }
})

// 请求拦截器
api.interceptors.request.use(
    (config) => {
        // 可以在这里添加loading状态
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 响应拦截器
api.interceptors.response.use(
    (response) => {
        // 统一处理响应数据
        return response.data
    },
    (error) => {
        // 统一处理错误
        const message = error.response?.data?.error || error.message || '请求失败'

        if (error.response?.status === 401) {
            // 未认证，重定向到登录页
            ElMessage.error('登录已过期，请重新登录')
            window.location.href = '/login'
        } else if (error.response?.status === 403) {
            // 权限不足
            ElMessage.error('权限不足')
        } else {
            ElMessage.error(message)
        }

        return Promise.reject(error)
    }
)

export { api }
