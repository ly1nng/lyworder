<template>
  <div class="user-management">
    <div class="page-header">
      <div class="header-content">
        <h1><i class="el-icon-user"></i>用户管理</h1>
        <p>管理系统用户账户和权限</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddUser" class="add-user-btn">
          <el-icon><DocumentAdd /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <div class="stat-card glass-card">
        <div class="stat-icon">
          <el-icon><User /></el-icon>
        </div>
        <div class="stat-info">
          <h3>{{ total }}</h3>
          <p>总用户数</p>
        </div>
      </div>
      
      <div class="stat-card glass-card">
        <div class="stat-icon active">
          <el-icon><UserFilled /></el-icon>
        </div>
        <div class="stat-info">
          <h3>{{ activeUsersCount }}</h3>
          <p>工作中</p>
        </div>
      </div>
      
      <div class="stat-card glass-card">
        <div class="stat-icon inactive">
          <el-icon><User /></el-icon>
        </div>
        <div class="stat-info">
          <h3>{{ inactiveUsersCount }}</h3>
          <p>非工作</p>
        </div>
      </div>
    </div>

    <!-- 用户列表 -->
    <el-card class="users-card glass-card">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索用户名..."
              clearable
              style="width: 200px; margin-right: 12px;"
              @input="handleSearch"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <el-tag>共 {{ total }} 名用户</el-tag>
          </div>
        </div>
      </template>

      <!-- 用户卡片列表 -->
      <div v-if="!loading && users.length > 0" class="users-grid">
        <div 
          v-for="user in users" 
          :key="user.id" 
          class="user-card glass-card"
          :class="{ 'inactive-user': user.status !== 1 }"
        >
          <div class="user-header">
            <div class="user-avatar" :class="getRoleType(user.role)">
              <el-icon><User /></el-icon>
            </div>
            <div class="user-info">
              <h3 class="user-name">{{ user.user_name }}</h3>
              <div class="user-role">
                <el-tag :type="getRoleType(user.role)" size="small">
                  {{ getRoleText(user.role) }}
                </el-tag>
              </div>
            </div>
          </div>
          
          <div class="user-details">
            <div class="detail-item">
              <span class="label">运维类型:</span>
              <span class="value">
                <el-tag v-if="user.ops_type" type="info" size="small">
                  {{ getOpsTypeText(user.ops_type) }}
                </el-tag>
                <span v-else class="text-muted">-</span>
              </span>
            </div>
            
            <div class="detail-item">
              <span class="label">工作状态:</span>
              <span class="value">
                <el-tag 
                  :type="user.status === 1 ? 'success' : 'info'" 
                  size="small"
                >
                  {{ user.status === 1 ? '工作中' : '非工作' }}
                </el-tag>
              </span>
            </div>
            
            <div class="detail-item">
              <span class="label">创建时间:</span>
              <span class="value">{{ formatDate(user.created_at) }}</span>
            </div>
          </div>
          
          <div class="user-actions">
            <el-button 
              type="primary" 
              size="small" 
              @click="editUser(user)"
              class="action-btn"
            >
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            
            <el-button
              :type="user.status === 1 ? 'warning' : 'success'"
              size="small"
              @click="toggleUserStatus(user)"
              class="action-btn"
            >
              <el-icon><Refresh /></el-icon>
              {{ user.status === 1 ? '设为非工作' : '设为工作中' }}
            </el-button>
            
            <el-button 
              type="danger" 
              size="small" 
              @click="deleteUser(user)"
              class="action-btn"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </div>
        </div>
      </div>
      
      <!-- 空状态 -->
      <div v-else-if="!loading && users.length === 0" class="empty-state">
        <el-empty description="暂无用户数据">
          <el-button type="primary" @click="showAddUser">添加用户</el-button>
        </el-empty>
      </div>
      
      <!-- 加载状态 -->
      <div v-else class="loading-state">
        <el-skeleton :rows="4" animated />
      </div>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next, jumper"
          @current-change="loadUsers"
          :page-sizes="[10, 20, 50]"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="showAddUserDialog"
      :title="editingUser ? '编辑用户' : '添加用户'"
      :width="isMobile ? '95%' : '500px'"
      :class="{ 'mobile-dialog': isMobile }"
      :top="isMobile ? '5vh' : '15vh'"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userFormRules"
        label-width="80px"
      >
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="userForm.user_name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select 
            v-model="userForm.role" 
            placeholder="请选择角色" 
            style="width: 100%"
            :disabled="editingUser !== null"
          >
            <el-option label="运维人员" value="operator" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="运维类型" prop="ops_type">
          <el-select v-model="userForm.ops_type" placeholder="请选择运维类型" style="width: 100%" clearable>
            <el-option label="系统运维" value="system" />
            <el-option label="IT运维" value="it" />
          </el-select>
        </el-form-item>
        <el-form-item label="工作状态" prop="status">
          <el-radio-group v-model="userForm.status">
            <el-radio :label="1">工作中</el-radio>
            <el-radio :label="0">非工作</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showAddUserDialog = false">取消</el-button>
          <el-button type="primary" @click="saveUser">
            {{ editingUser ? '更新' : '添加' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  User, 
  UserFilled, 
  Search, 
  Edit, 
  Delete, 
  Refresh,
  DocumentAdd
} from '@element-plus/icons-vue'
import api from '@/api'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 响应式数据
const loading = ref(false)
const users = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const searchDebounceTimer = ref(null)

const showAddUserDialog = ref(false)
const editingUser = ref(null)
const userFormRef = ref()

// 是否为移动端
const isMobile = computed(() => appStore.isMobile)

// 计算属性 - 使用缓存优化
// 缓存的计算属性
const activeUsersCount = computed(() => {
  return users.value.filter(user => user.status === 1).length
})

const inactiveUsersCount = computed(() => {
  return users.value.filter(user => user.status !== 1).length
})

// 分页缓存
const pageCache = new Map()
// 请求防抖定时器
let requestDebounceTimer = null

// 用户表单
const userForm = reactive({
  user_name: '',
  role: 'operator',
  ops_type: '',
  status: 1
})

// 表单验证规则
const userFormRules = {
  user_name: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 防抖搜索功能
const handleSearch = () => {
  // 清除之前的定时器
  if (searchDebounceTimer.value) {
    clearTimeout(searchDebounceTimer.value)
  }
  
  // 设置新的定时器
  searchDebounceTimer.value = setTimeout(() => {
    currentPage.value = 1
    loadUsers(1)
  }, 300) // 300ms防抖延迟
}

// 方法
const loadUsers = async (page = 1, retryCount = 0) => {
  // 请求防抖
  if (requestDebounceTimer) {
    clearTimeout(requestDebounceTimer)
  }
  
  // 避免重复请求
  if (loading.value) return
  
  // 检查缓存
  const cacheKey = `${page}-${pageSize.value}-${searchKeyword.value}`
  if (pageCache.has(cacheKey)) {
    const cachedData = pageCache.get(cacheKey)
    users.value = cachedData.users
    total.value = cachedData.total
    currentPage.value = page
    
    // 预加载下一页
    if (page * pageSize.value < total.value) {
      preloadPage(page + 1)
    }
    
    return
  }
  
  // 设置加载延迟，避免闪烁
  const loadingTimer = setTimeout(() => {
    loading.value = true
  }, 100)
  
  try {
    const params = {
      page,
      limit: pageSize.value
    }
    
    // 如果有搜索关键词，添加到参数中
    if (searchKeyword.value) {
      params.user_name = searchKeyword.value
    }
    
    const response = await api.get('/operator/users', { params })
    // 正确处理分页数据
    if (response && Array.isArray(response.users)) {
      users.value = response.users
      // 如果后端返回了total字段，则使用它；否则使用用户数组长度
      total.value = response.total || response.users.length
      
      // 缓存数据
      pageCache.set(cacheKey, {
        users: response.users,
        total: total.value,
        timestamp: Date.now()
      })
      
      // 清理旧缓存，只保留最近的10页
      if (pageCache.size > 10) {
        const firstKey = pageCache.keys().next().value
        pageCache.delete(firstKey)
      }
      
      // 预加载下一页
      if (page * pageSize.value < total.value) {
        preloadPage(page + 1)
      }
    } else {
      // 如果响应格式不正确，使用空数组
      users.value = []
      total.value = 0
    }
    currentPage.value = page
  } catch (error) {
    // 实现重试机制
    if (retryCount < 2) {
      console.warn(`加载用户列表失败，正在重试 (${retryCount + 1}/2):`, error)
      setTimeout(() => {
        loadUsers(page, retryCount + 1)
      }, 1000 * (retryCount + 1)) // 递增延迟重试
      return
    }
    
    ElMessage.error('加载用户列表失败')
    console.error('Load users error:', error)
    // 出错时清空数据
    users.value = []
    total.value = 0
  } finally {
    clearTimeout(loadingTimer)
    loading.value = false
  }
}

// 预加载页面
const preloadPage = async (page) => {
  // 限制预加载数量，避免过多请求
  if (pageCache.size > 15) return
  
  try {
    const cacheKey = `${page}-${pageSize.value}-${searchKeyword.value}`
    // 如果已缓存则不预加载
    if (pageCache.has(cacheKey)) return
    
    const params = {
      page,
      limit: pageSize.value
    }
    
    // 如果有搜索关键词，添加到参数中
    if (searchKeyword.value) {
      params.user_name = searchKeyword.value
    }
    
    const response = await api.get('/operator/users', { params })
    if (response && Array.isArray(response.users)) {
      // 缓存预加载的数据
      pageCache.set(cacheKey, {
        users: response.users,
        total: response.total || response.users.length,
        timestamp: Date.now()
      })
    }
  } catch (error) {
    console.warn('预加载页面失败:', error)
  }
}

// 处理页面大小变化
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  // 清除缓存并重新加载
  pageCache.clear()
  loadUsers(1)
}

const editUser = (user) => {
  editingUser.value = user
  Object.assign(userForm, {
    user_name: user.user_name,
    role: user.role,
    ops_type: user.ops_type || '',
    status: user.status
  })
  showAddUserDialog.value = true
}

const showAddUser = () => {
  // 重置编辑状态
  editingUser.value = null
  // 重置表单
  resetUserForm()
  // 显示对话框
  showAddUserDialog.value = true
}

const saveUser = async () => {
  const valid = await userFormRef.value.validate()
  if (!valid) return

  try {
    if (editingUser.value) {
      // 编辑用户 - 调用更新用户信息的API
      await api.user.updateUserInfo(editingUser.value.user_name, {
        user_name: userForm.user_name,
        role: userForm.role,
        ops_type: userForm.ops_type,
        status: userForm.status
      })
      ElMessage.success('用户信息更新成功')
      
      // 局部更新用户数据而不是重新加载整个列表
      const index = users.value.findIndex(u => u.id === editingUser.value.id)
      if (index > -1) {
        users.value[index] = {
          ...users.value[index],
          user_name: userForm.user_name,
          role: userForm.role,
          ops_type: userForm.ops_type,
          status: userForm.status
        }
      }
    } else {
      // 创建用户
      const response = await api.post('/operator/users', {
        user_name: userForm.user_name,
        role: userForm.role,
        ops_type: userForm.ops_type,
        status: userForm.status
      })
      ElMessage.success('用户创建成功')
      
      // 如果当前在最后一页且未满，直接添加到列表中
      if (users.value.length < pageSize.value) {
        users.value.push(response.data || {
          id: Date.now(),
          user_name: userForm.user_name,
          role: userForm.role,
          ops_type: userForm.ops_type,
          status: userForm.status,
          created_at: new Date().toISOString()
        })
        total.value++
      }
    }
    
    showAddUserDialog.value = false
    resetUserForm()
    
    // 只有在需要时才重新加载（如分页或搜索结果）
    if (!editingUser.value && users.value.length >= pageSize.value) {
      // 清除缓存并重新加载
      pageCache.clear()
      loadUsers(currentPage.value)
    }
  } catch (error) {
    ElMessage.error(editingUser.value ? '用户更新失败' : '用户创建失败')
    console.error('Save user error:', error)
  }
}

// 状态切换队列，避免并发请求
const statusToggleQueue = new Map()

const toggleUserStatus = async (user) => {
  const newStatus = user.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '设为工作中' : '设为非工作'
  
  // 检查是否正在处理此用户的请求
  if (statusToggleQueue.has(user.user_name)) {
    ElMessage.info('正在处理中，请稍候...')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要${action}用户 ${user.user_name} 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 标记为正在处理
    statusToggleQueue.set(user.user_name, true)
    
    // 使用query参数方式调用后端API
    await api.put(`/operator/users/status?user_name=${encodeURIComponent(user.user_name)}&status=${newStatus}`)
    
    ElMessage.success(`用户${action}成功`)
    // 优化：局部更新而不是重新加载整个列表
    user.status = newStatus
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`用户${action}失败`)
      console.error('Toggle user status error:', error)
    }
  } finally {
    // 清除处理标记
    statusToggleQueue.delete(user.user_name)
  }
}

const deleteUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 ${user.user_name} 吗？此操作不可恢复！`,
      '确认删除',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 使用query参数方式调用后端API
    await api.delete(`/operator/users?user_name=${encodeURIComponent(user.user_name)}`)
    
    ElMessage.success('用户删除成功')
    // 优化：局部更新而不是重新加载整个列表
    const index = users.value.findIndex(u => u.id === user.id)
    if (index > -1) {
      users.value.splice(index, 1)
      total.value--
      
      // 更新缓存
      pageCache.clear()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('用户删除失败')
      console.error('Delete user error:', error)
    }
  }
}

const resetUserForm = () => {
  Object.assign(userForm, {
    user_name: '',
    role: 'operator', // 默认为运维人员
    ops_type: '',
    status: 1
  })
  editingUser.value = null
}

// 角色类型映射 - 避免重复计算
const roleTypeMap = {
  admin: 'danger',
  operator: 'warning'
}

const roleTextMap = {
  admin: '管理员',
  operator: '运维人员'
}

const opsTypeTextMap = {
  system: '系统运维',
  it: 'IT运维'
}

const getRoleType = (role) => {
  return roleTypeMap[role] || 'info'
}

const getRoleText = (role) => {
  return roleTextMap[role] || role
}

// 添加获取运维类型文本的函数
const getOpsTypeText = (opsType) => {
  return opsTypeTextMap[opsType] || opsType
}

// 日期格式化缓存
const formatDateCache = new Map()

const formatDate = (dateString) => {
  if (!dateString) return '-'
  
  // 使用缓存避免重复格式化
  if (formatDateCache.has(dateString)) {
    return formatDateCache.get(dateString)
  }
  
  // 限制缓存大小，避免内存泄漏
  if (formatDateCache.size > 100) {
    const firstKey = formatDateCache.keys().next().value
    formatDateCache.delete(firstKey)
  }
  
  const formatted = new Date(dateString).toLocaleString('zh-CN')
  formatDateCache.set(dateString, formatted)
  return formatted
}

// 清理防抖定时器
const clearSearchDebounce = () => {
  if (searchDebounceTimer.value) {
    clearTimeout(searchDebounceTimer.value)
    searchDebounceTimer.value = null
  }
}

// 监听搜索关键词变化，实现即时搜索
watch(searchKeyword, (newVal) => {
  // 只有当搜索关键词变化较大时才清除缓存
  if (!newVal || newVal.length < 2) {
    pageCache.clear()
  }
  
  // 使用防抖搜索
  handleSearch()
})

// 生命周期
onMounted(() => {
  // 检查是否有缓存的第一页数据，如果有则快速渲染
  const firstPageCacheKey = `1-${pageSize.value}-${searchKeyword.value}`
  if (pageCache.has(firstPageCacheKey)) {
    const cachedData = pageCache.get(firstPageCacheKey)
    users.value = cachedData.users
    total.value = cachedData.total
    currentPage.value = 1
    
    // 异步加载最新数据
    setTimeout(() => {
      loadUsers(1)
    }, 0)
  } else {
    loadUsers()
  }
})

// 组件卸载时清理
// onUnmounted(() => {
//   clearSearchDebounce()
// })

</script>

<style scoped>
.user-management {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px 0;
}

.header-content h1 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 2rem;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  font-size: 1rem;
}

.add-user-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.add-user-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

/* 统计卡片 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 20px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.2);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  font-size: 24px;
  margin-right: 16px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.stat-icon.active {
  background: linear-gradient(135deg, #10b981, #059669);
}

.stat-icon.inactive {
  background: linear-gradient(135deg, #6b7280, #4b5563);
}

.stat-info h3 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.stat-info p {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
}

.users-card {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
}

.card-header .header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 用户卡片网格 */
.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
  padding: 20px;
}

.user-card {
  border-radius: 12px;
  padding: 20px;
  background: white;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  border: 1px solid #e5e7eb;
}

.user-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
}

.user-card.inactive-user {
  opacity: 0.7;
  background: #f9fafb;
}

.user-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  margin-right: 16px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.user-avatar.warning {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.user-avatar.danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.user-info {
  flex: 1;
}

.user-name {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.user-role {
  display: flex;
  align-items: center;
}

.user-details {
  margin-bottom: 20px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
}

.detail-item:last-child {
  margin-bottom: 0;
}

.detail-item .label {
  font-weight: 500;
  color: #6b7280;
  font-size: 14px;
}

.detail-item .value {
  font-weight: 500;
  color: #1f2937;
  font-size: 14px;
}

.user-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.action-btn {
  flex: 1;
  min-width: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

/* 空状态 */
.empty-state {
  padding: 60px 20px;
  text-align: center;
}

/* 加载状态 */
.loading-state {
  padding: 40px 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding: 20px;
  border-top: 1px solid #e5e7eb;
}

.text-muted {
  color: #9ca3af;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
}

/* Element Plus 样式自定义 */
:deep(.el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  padding: 20px;
}

:deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

:deep(.el-dialog__close) {
  color: white;
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-select .el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-radio__inner) {
  border-radius: 50%;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .user-management {
    padding: 15px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .stats-cards {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-card {
    padding: 16px;
  }
  
  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
    margin-right: 12px;
  }
  
  .stat-info h3 {
    font-size: 24px;
  }
  
  .users-grid {
    grid-template-columns: 1fr;
    gap: 16px;
    padding: 16px;
  }
  
  .user-card {
    padding: 16px;
  }
  
  .user-header {
    margin-bottom: 16px;
    padding-bottom: 12px;
  }
  
  .user-avatar {
    width: 45px;
    height: 45px;
    font-size: 18px;
    margin-right: 12px;
  }
  
  .user-name {
    font-size: 16px;
  }
  
  .detail-item {
    margin-bottom: 10px;
    padding: 6px 0;
  }
  
  .user-actions {
    gap: 8px;
  }
  
  .action-btn {
    min-width: 70px;
    font-size: 13px;
    padding: 8px 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 16px;
  }
  
  .card-header .header-actions {
    width: 100%;
    flex-direction: column;
    gap: 12px;
  }
  
  :deep(.el-input) {
    width: 100% !important;
  }
  
  .pagination-container {
    margin-top: 20px;
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .user-management {
    padding: 12px;
  }
  
  .header-content h1 {
    font-size: 1.5rem;
  }
  
  .stats-cards {
    gap: 12px;
  }
  
  .stat-card {
    padding: 12px;
  }
  
  .stat-icon {
    width: 45px;
    height: 45px;
    font-size: 18px;
    margin-right: 10px;
  }
  
  .stat-info h3 {
    font-size: 20px;
  }
  
  .stat-info p {
    font-size: 12px;
  }
  
  .users-grid {
    gap: 12px;
    padding: 12px;
  }
  
  .user-card {
    padding: 12px;
  }
  
  .user-header {
    margin-bottom: 12px;
    padding-bottom: 10px;
  }
  
  .user-avatar {
    width: 40px;
    height: 40px;
    font-size: 16px;
    margin-right: 10px;
  }
  
  .user-name {
    font-size: 15px;
  }
  
  .detail-item {
    margin-bottom: 8px;
    padding: 4px 0;
  }
  
  .detail-item .label,
  .detail-item .value {
    font-size: 13px;
  }
  
  .user-actions {
    gap: 6px;
  }
  
  .action-btn {
    min-width: 60px;
    font-size: 12px;
    padding: 6px 8px;
    gap: 2px;
  }
  
  .dialog-footer {
    flex-direction: column;
    gap: 8px;
  }
  
  .dialog-footer .el-button {
    width: 100%;
  }
}
</style>