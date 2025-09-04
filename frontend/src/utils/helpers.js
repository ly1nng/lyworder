// 验证图片文件
export const validateImageFile = (file) => {
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif']
    const maxSize = 5 * 1024 * 1024 // 5MB

    if (!file) {
        return { valid: false, message: '请选择文件' }
    }

    if (!allowedTypes.includes(file.type)) {
        return { valid: false, message: '只支持 JPG、PNG、WebP 和 GIF 格式的图片' }
    }

    if (file.size > maxSize) {
        return { valid: false, message: '图片大小不能超过 5MB' }
    }

    return { valid: true }
}

// 格式化文件大小
export const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 B'

    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))

    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期时间 - 显示更友好的日期时间格式
export const formatDateTime = (dateString) => {
    if (!dateString) return '-'

    const date = new Date(dateString)
    const now = new Date()

    // 如果是今天
    if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString('zh-CN', {
            hour: '2-digit',
            minute: '2-digit'
        })
    }

    // 如果是今年
    if (date.getFullYear() === now.getFullYear()) {
        return date.toLocaleDateString('zh-CN', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        })
    }

    // 其他情况显示完整日期
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

// 格式化日期（修改为显示更完整的日期时间）
export const formatDate = (dateString) => {
    if (!dateString) return '-'
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    })
}

// 格式化相对时间
export const formatRelativeTime = (dateString) => {
    if (!dateString) return '-'

    const date = new Date(dateString)
    const now = new Date()
    const diff = now.getTime() - date.getTime()

    const minute = 60 * 1000
    const hour = 60 * minute
    const day = 24 * hour
    const week = 7 * day
    const month = 30 * day

    if (diff < minute) {
        return '刚刚'
    } else if (diff < hour) {
        return Math.floor(diff / minute) + ' 分钟前'
    } else if (diff < day) {
        return Math.floor(diff / hour) + ' 小时前'
    } else if (diff < week) {
        return Math.floor(diff / day) + ' 天前'
    } else if (diff < month) {
        return Math.floor(diff / week) + ' 周前'
    } else {
        return Math.floor(diff / month) + ' 个月前'
    }
}

// 获取工单状态文本和类型
export const getTicketStatusInfo = (status) => {
    const statusMap = {
        'open': { text: '待处理', type: 'warning' },
        'in_progress': { text: '处理中', type: 'primary' },
        'resolved': { text: '已解决', type: 'success' },
        'closed': { text: '已关闭', type: 'info' },
        'rejected': { text: '已拒绝', type: 'danger' }
    }

    return statusMap[status] || { text: status, type: 'info' }
}

// 获取状态文本
export const getStatusText = (status) => {
    return getTicketStatusInfo(status).text
}

// 获取状态类型
export const getStatusType = (status) => {
    return getTicketStatusInfo(status).type
}

// 获取工单类型文本和颜色
export const getTicketTypeInfo = (type) => {
    const typeMap = {
        'bug': { text: '故障报修', color: '#ff4757' },
        'feature': { text: '功能需求', color: '#3742fa' },
        'support': { text: '技术支持', color: '#2ed573' },
        'question': { text: '咨询问题', color: '#ffa502' },
        'other': { text: '其他', color: '#747d8c' }
    }

    return typeMap[type] || { text: type, color: '#747d8c' }
}

// 获取用户角色文本和类型
export const getUserRoleInfo = (role) => {
    const roleMap = {
        'admin': { text: '管理员', type: 'danger' },
        'operator': { text: '运维人员', type: 'warning' },
        'user': { text: '用户', type: 'info' }
    }

    return roleMap[role] || { text: role, type: 'info' }
}

// 获取角色文本
export const getRoleText = (role) => {
    return getUserRoleInfo(role).text
}

// 防抖函数
export const debounce = (func, wait) => {
    let timeout
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout)
            func(...args)
        }
        clearTimeout(timeout)
        timeout = setTimeout(later, wait)
    }
}

// 节流函数
export const throttle = (func, limit) => {
    let inThrottle
    return function (...args) {
        if (!inThrottle) {
            func.apply(this, args)
            inThrottle = true
            setTimeout(() => inThrottle = false, limit)
        }
    }
}

// 深拷贝
export const deepClone = (obj) => {
    if (obj === null || typeof obj !== 'object') return obj
    if (obj instanceof Date) return new Date(obj.getTime())
    if (obj instanceof Array) return obj.map(item => deepClone(item))
    if (typeof obj === 'object') {
        const clonedObj = {}
        for (let key in obj) {
            if (obj.hasOwnProperty(key)) {
                clonedObj[key] = deepClone(obj[key])
            }
        }
        return clonedObj
    }
}

// 生成随机ID
export const generateId = (length = 8) => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let result = ''
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return result
}

// 检查是否为移动设备
export const isMobile = () => {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}

// 复制文本到剪贴板
export const copyToClipboard = async (text) => {
    try {
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(text)
            return true
        } else {
            // 降级处理
            const textArea = document.createElement('textarea')
            textArea.value = text
            textArea.style.position = 'fixed'
            textArea.style.left = '-999999px'
            textArea.style.top = '-999999px'
            document.body.appendChild(textArea)
            textArea.focus()
            textArea.select()
            const success = document.execCommand('copy')
            textArea.remove()
            return success
        }
    } catch (error) {
        console.error('复制失败:', error)
        return false
    }
}