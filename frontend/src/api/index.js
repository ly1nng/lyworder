import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建axios实例
const api = axios.create({
    baseURL: '/api/v1',
    timeout: 10000,
    withCredentials: true, // 发送cookies
    headers: {
        'Content-Type': 'application/json'
    }
})

// 请求拦截器
api.interceptors.request.use(
    (config) => {
        // 如果是FormData，移除Content-Type让浏览器自动设置
        if (config.data instanceof FormData) {
            delete config.headers['Content-Type']
        }

        return config
    },
    (error) => {
        console.error('Request error:', error)
        return Promise.reject(error)
    }
)

// 响应拦截器
api.interceptors.response.use(
    (response) => {
        // 对于201状态码的特殊处理（工单创建成功）
        if (response.status === 201) {
            // 直接返回响应数据，不进行code检查
            return response.data;
        }

        // 检查业务状态码
        if (response.data && response.data.code !== undefined) {
            if (response.data.code !== 200) {
                const message = response.data.message || response.data.error || '请求失败'
                ElMessage.error(message)
                return Promise.reject(new Error(message))
            }
        }

        return response.data
    },
    (error) => {
        console.error('Response error:', error)

        // 处理HTTP状态码错误
        if (error.response) {
            const { status, data } = error.response

            switch (status) {
                case 401:
                    // 未授权，重定向到登录页
                    ElMessage.error('请先登录')
                    window.location.href = '/login'
                    break
                case 403:
                    ElMessage.error('权限不足')
                    break
                case 404:
                    ElMessage.error('请求的资源不存在')
                    break
                case 500:
                    ElMessage.error('服务器内部错误')
                    break
                default:
                    const message = data?.message || data?.error || `请求失败 (${status})`
                    ElMessage.error(message)
            }
        } else if (error.request) {
            // 网络错误
            ElMessage.error('网络连接失败，请检查网络设置')
        } else {
            // 其他错误
            ElMessage.error(error.message || '请求失败')
        }

        return Promise.reject(error)
    }
)

// API方法封装
const apiService = {
    // 通用请求方法
    get: (url, config = {}) => api.get(url, config),
    post: (url, data = {}, config = {}) => api.post(url, data, config),
    put: (url, data = {}, config = {}) => api.put(url, data, config),
    delete: (url, config = {}) => api.delete(url, config),
    patch: (url, data = {}, config = {}) => api.patch(url, data, config),

    // 用户相关接口
    user: {
        // 获取当前用户信息
        getCurrentUser: () => api.get('/users'),

        // 获取所有用户（管理员）
        getAllUsers: (params = {}) => api.get('/operator/users', { params }),

        // 创建用户
        createUser: (userData) => api.post('/operator/users', userData),

        // 更新用户信息
        updateUserInfo: (userName, userData) => api.put(`/operator/users/info?user_name=${encodeURIComponent(userName)}`, userData),

        // 更新用户状态
        updateUserStatus: (userId, status) => api.put('/operator/users/status', { user_id: userId, status }),

        // 更新用户类型
        updateUserType: (userId, opsType) => api.put('/operator/users/type', { user_id: userId, ops_type: opsType }),

        // 删除用户
        deleteUser: (userId) => api.delete('/operator/users', { data: { user_id: userId } }),

        // 根据状态获取用户
        getUsersByStatus: (status) => api.get('/operator/users/status', { params: { status } }),

        // 获取运维人员用户名列表
        getOperatorUsersName: () => api.get('/operator/usersname')
    },

    // 工单相关接口
    ticket: {
        // 获取工单列表
        getTickets: (params = {}) => api.get('/tickets', { params }),

        // 获取工单详情
        getTicketDetail: (ticketId) => api.get(`/tickets/${ticketId}/detail`),

        // 创建工单
        createTicket: (ticketData) => {
            // 始终使用FormData以确保Content-Type为multipart/form-data
            const formData = new FormData()
            formData.append('title', ticketData.title)
            formData.append('content', ticketData.content)

            // 添加截图文件（如果有）
            if (ticketData.screenshots && ticketData.screenshots.length > 0) {
                for (let i = 0; i < ticketData.screenshots.length; i++) {
                    formData.append('screenshots', ticketData.screenshots[i])
                }
            }

            return api.post('/tickets', formData)
        },

        // 更新工单状态
        updateTicketStatus: (ticketId, status) => api.put(`/tickets/${ticketId}/status`, { status }),

        // 更新工单处理人
        updateTicketOperator: (ticketId, operatorName) => {
            // 如果operatorName是数组，转换为逗号分隔的字符串
            const operatorNames = Array.isArray(operatorName) ? operatorName.join(',') : operatorName;
            return api.put(`/tickets/${ticketId}/operatorname`, { operator_name: operatorNames })
        },

        // 更新工单备注
        updateTicketRemark: (ticketId, remark) => api.put(`/tickets/${ticketId}/remark`, { remark }),

        // 更新工单解决方案
        updateTicketSolution: (ticketId, solution) => api.put(`/tickets/${ticketId}/solution`, { solution }),

        // 更新工单类型
        updateTicketType: (ticketId, ticketType) => api.put(`/tickets/${ticketId}/type`, { ticket_type: ticketType })
    },

    // 文件上传相关
    upload: {
        // 上传图片文件
        uploadImages: (files) => {
            const formData = new FormData()
            files.forEach((file, index) => {
                formData.append('images', file)
            })

            return api.post('/upload/images', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
        }
    }
}

export default apiService

// 导出分类API以便向后兼容
export const ticketApi = apiService.ticket
export const userApi = apiService.user
export const uploadApi = apiService.upload