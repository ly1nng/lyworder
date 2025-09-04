<template>
  <div class="user-management">
    <div class="page-header">
      <div class="header-content">
        <h1><i class="el-icon-user"></i>用户管理</h1>
        <p>管理系统用户账户和权限</p>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="showAddUserDialog = true">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <el-card class="search-card">
      <div class="search-section">
        <div class="search-inputs">
          <el-input
            v-model="searchForm.username"
            placeholder="搜索用户名"
            @keyup.enter="searchUsers"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select v-model="searchForm.role" placeholder="角色筛选" clearable>
            <el-option label="运维人员" value="operator" />
            <el-option label="管理员" value="admin" />
          </el-select>
          <el-select v-model="searchForm.status" placeholder="状态筛选" clearable>
            <el-option label="激活" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </div>
        <div class="search-actions">
          <el-button @click="searchUsers">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </div>
      </div>
    </el-card>

    <!-- 用户列表 -->
    <el-card class="users-card">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-tag>共 {{ total }} 名用户</el-tag>
        </div>
      </template>

      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_name" label="用户名" min-width="120" />
        <el-table-column prop="ops_type" label="运维类型" min-width="100">
          <template #default="{ row }">
            <el-tag v-if="row.ops_type" type="info" size="small">
              {{ getOpsTypeText(row.ops_type) }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getRoleType(row.role)"
              size="small"
            >
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 1 ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status === 1 ? '激活' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button
                type="primary"
                size="small"
                @click="editUser(row)"
              >
                编辑
              </el-button>
              <el-button
                :type="row.status === 1 ? 'warning' : 'success'"
                size="small"
                @click="toggleUserStatus(row)"
              >
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="deleteUser(row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </el-card>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="showAddUserDialog"
      :title="editingUser ? '编辑用户' : '添加用户'"
      width="500px"
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
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="userForm.status">
            <el-radio :label="1">激活</el-radio>
            <el-radio :label="0">禁用</el-radio>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import api from '@/api'

// 响应式数据
const loading = ref(false)
const users = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const showAddUserDialog = ref(false)
const editingUser = ref(null)
const userFormRef = ref()

// 搜索表单
const searchForm = reactive({
  username: '',
  role: '',
  status: ''
})

// 用户表单
const userForm = reactive({
  user_name: '',
  role: 'user',
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

// 方法
const loadUsers = async (page = 1) => {
  loading.value = true
  try {
    const params = {
      page,
      limit: pageSize.value,
      ...searchForm
    }
    const response = await api.get('/operator/users', { params })
    // 后端返回的是 { users: [...] }，没有total字段
    // 先处理用户列表
    users.value = response.users || []
    // 暂时使用用户数组长度作为总数（后续需要后端支持分页）
    total.value = users.value.length
    currentPage.value = page
  } catch (error) {
    ElMessage.error('加载用户列表失败')
    console.error('Load users error:', error)
  }
  loading.value = false
}

const searchUsers = () => {
  currentPage.value = 1
  loadUsers(1)
}

const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    role: '',
    status: ''
  })
  searchUsers()
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
    } else {
      // 创建用户
      await api.post('/operator/users', {
        user_name: userForm.user_name,
        role: userForm.role,
        ops_type: userForm.ops_type,
        status: userForm.status
      })
      ElMessage.success('用户创建成功')
    }
    
    showAddUserDialog.value = false
    resetUserForm()
    loadUsers(currentPage.value)
  } catch (error) {
    ElMessage.error(editingUser.value ? '用户更新失败' : '用户创建失败')
    console.error('Save user error:', error)
  }
}

const toggleUserStatus = async (user) => {
  const newStatus = user.status === 1 ? 0 : 1
  const action = newStatus === 1 ? '启用' : '禁用'
  
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
    
    // 使用query参数方式调用后端API
    await api.put(`/operator/users/status?user_name=${encodeURIComponent(user.user_name)}&status=${newStatus}`)
    
    ElMessage.success(`用户${action}成功`)
    loadUsers(currentPage.value)
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`用户${action}失败`)
      console.error('Toggle user status error:', error)
    }
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
    loadUsers(currentPage.value)
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

const getRoleType = (role) => {
  const types = {
    admin: 'danger',
    operator: 'warning'
  }
  return types[role] || 'info'
}

const getRoleText = (role) => {
  const texts = {
    admin: '管理员',
    operator: '运维人员'
  }
  return texts[role] || role
}

// 添加获取运维类型文本的函数
const getOpsTypeText = (opsType) => {
  const texts = {
    system: '系统运维',
    it: 'IT运维'
  }
  return texts[opsType] || opsType
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadUsers()
})
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
  align-items: flex-end;
  margin-bottom: 24px;
  padding: 20px 0;
}

.header-content h1 {
  color: white;
  margin: 0 0 8px 0;
  font-size: 2rem;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  font-size: 1rem;
}

.search-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.search-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.search-inputs {
  display: flex;
  gap: 12px;
  flex: 1;
}

.search-inputs .el-input,
.search-inputs .el-select {
  min-width: 150px;
}

.search-actions {
  display: flex;
  gap: 8px;
}

.users-card {
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.95);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.text-muted {
  color: #999;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* 操作按钮水平排列样式 */
.action-buttons {
  display: flex;
  gap: 8px;
  padding: 8px 0;
}

.action-buttons .el-button {
  margin: 0;
  min-width: 60px;
  padding: 8px 12px;
}

.action-buttons .el-button[size="small"] {
  font-size: 12px;
  height: 32px;
  line-height: 1;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-management {
    padding: 15px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .search-section {
    flex-direction: column;
    align-items: stretch;
  }

  .search-inputs {
    flex-direction: column;
  }

  .search-actions {
    justify-content: flex-end;
  }

  .el-table {
    font-size: 14px;
  }

  .el-dialog {
    width: 90%;
    margin: 0 5%;
  }
  
  /* 移动端操作按钮优化 */
  .el-table__body .el-button {
    margin-right: 6px;
    min-width: 50px;
    padding: 6px 10px;
  }
  
  .el-table__body .el-button[size="small"] {
    font-size: 12px;
    height: 28px;
  }
  
  /* 移动端操作按钮水平排列优化 */
  .action-buttons {
    gap: 6px;
    padding: 6px 0;
  }
}
</style>