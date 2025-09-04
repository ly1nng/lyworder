<template>
  <div class="admin-view">
    <div class="container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon class="title-icon"><DataBoard /></el-icon>
          管理后台
        </h1>
        <p class="page-subtitle">系统管理和工单统计分析</p>
      </div>
      
      <!-- 统计卡片 -->
      <div class="stats-section">
        <div class="stats-grid">
          <div class="stat-card glass-card">
            <div class="stat-icon total">
              <el-icon><Tickets /></el-icon>
            </div>
            <div class="stat-content">
              <h3 class="stat-number">{{ stats.total }}</h3>
              <p class="stat-label">总工单数</p>
            </div>
          </div>
          
          <div class="stat-card glass-card">
            <div class="stat-icon open">
              <el-icon><FolderOpened /></el-icon>
            </div>
            <div class="stat-content">
              <h3 class="stat-number">{{ stats.open }}</h3>
              <p class="stat-label">进行中的工单</p>
            </div>
          </div>
          
          <div class="stat-card glass-card">
            <div class="stat-icon closed">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <h3 class="stat-number">{{ stats.closed }}</h3>
              <p class="stat-label">已解决的工单</p>
            </div>
          </div>
          
          <div class="stat-card glass-card">
            <div class="stat-icon operators">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stat-content">
              <h3 class="stat-number">{{ stats.operators }}</h3>
              <p class="stat-label">处理人员</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 快速操作 -->
      <div class="quick-actions-section">
        <div class="section-title">
          <h2>快速操作</h2>
          <p>常用的管理功能</p>
        </div>
        <div class="actions-grid">
          <div class="action-card glass-card" @click="$router.push('/admin/users')">
            <div class="action-icon users">
              <el-icon><User /></el-icon>
            </div>
            <div class="action-content">
              <h4>用户管理</h4>
              <p>管理系统用户和权限</p>
            </div>
            <div class="action-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
          
          <div class="action-card glass-card" @click="refreshData">
            <div class="action-icon refresh">
              <el-icon><Refresh /></el-icon>
            </div>
            <div class="action-content">
              <h4>刷新数据</h4>
              <p>重新加载统计数据</p>
            </div>
            <div class="action-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 搜索表单 -->
      <SearchForm
        ref="searchFormRef"
        :loading="loading"
        show-creator
        show-operator
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- 工单列表 -->
      <div class="tickets-section">
        <div class="section-title">
          <h2>工单管理</h2>
          <p>查看和管理所有工单</p>
        </div>
        
        <div v-if="loading" class="loading-section">
          <div class="loading-content">
            <div class="loading-spinner"></div>
            <p>加载中...</p>
          </div>
        </div>
        
        <div v-else-if="tickets.length === 0" class="empty-section">
          <el-empty 
            description="暂无工单记录"
            :image-size="120"
          >
            <template #image>
              <el-icon class="empty-icon"><DocumentRemove /></el-icon>
            </template>
          </el-empty>
        </div>
        
        <div v-else class="tickets-list">
          <OperatorTicketCard
            v-for="ticket in tickets"
            :key="ticket.id"
            :ticket="ticket"
            :editable-type="true"
            :editable-operator="true"
            @view-screenshots="handleViewScreenshots"
            @remark="handleRemark"
            @solution="handleSolution"
            @status-change="handleStatusChange"
            @update-ticket-type="handleTypeChange"
            @update-ticket-operator="handleOperatorChange"
          />
        </div>
      </div>
      
      <!-- 分页 -->
      <AppPagination
        v-if="!loading && tickets.length > 0"
        :total="total"
        :page="currentPage"
        :limit="pageSize"
        @change="handlePageChange"
      />
    </div>
    
    <!-- 备注编辑对话框 -->
    <el-dialog
      v-model="remarkDialogVisible"
      title="编辑备注"
      width="600px"
      class="edit-dialog"
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <el-form :model="remarkForm" label-position="top">
          <el-form-item label="备注内容">
            <el-input
              v-model="remarkForm.remark"
              type="textarea"
              :rows="6"
              placeholder="请输入备注信息..."
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="remarkDialogVisible = false">
            取消
          </el-button>
          <el-button 
            type="primary" 
            @click="saveRemark"
            :loading="saving"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 解决方案编辑对话框 -->
    <el-dialog
      v-model="solutionDialogVisible"
      title="编辑解决方案"
      width="600px"
      class="edit-dialog"
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <el-form :model="solutionForm" label-position="top">
          <el-form-item label="解决方案">
            <el-input
              v-model="solutionForm.solution"
              type="textarea"
              :rows="6"
              placeholder="请输入解决方案..."
              maxlength="2000"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="solutionDialogVisible = false">
            取消
          </el-button>
          <el-button 
            type="primary" 
            @click="saveSolution"
            :loading="saving"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 工单类型修改对话框 -->
    <el-dialog
      v-model="typeDialogVisible"
      title="修改工单类型"
      width="500px"
      class="edit-dialog"
      destroy-on-close
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <el-form 
          :model="typeForm" 
          label-width="80px"
          @submit.prevent="saveType"
        >
          <el-form-item label="工单类型">
            <el-select 
              v-model="typeForm.ticket_type" 
              placeholder="请选择工单类型"
              style="width: 100%"
            >
              <el-option label="IT" value="IT" />
              <el-option label="系统" value="系统" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="typeDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="saveType"
            :loading="saving"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 处理人修改对话框 -->
    <el-dialog
      v-model="operatorDialogVisible"
      title="分配处理人"
      width="500px"
      class="edit-dialog"
      destroy-on-close
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <el-form 
          :model="operatorForm" 
          label-width="80px"
          @submit.prevent="saveOperator"
        >
          <el-form-item label="处理人">
            <el-select 
              v-model="operatorForm.operator_name" 
              placeholder="请选择处理人"
              style="width: 100%"
            >
              <el-option 
                v-for="operator in operators" 
                :key="operator" 
                :label="operator" 
                :value="operator" 
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="operatorDialogVisible = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="saveOperator"
            :loading="saving"
          >
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 截图全屏展示对话框 -->
    <el-dialog
      v-model="screenshotsDialogVisible"
      :title="`所有截图 (${currentScreenshots.length}张) - ${currentTicket?.title || '工单'}`"
      width="95%"
      class="screenshots-fullscreen-dialog"
      destroy-on-close
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      top="2vh"
    >
      <div class="fullscreen-screenshots-content">
        <!-- 对话框头部信息 -->
        <div class="dialog-header-info">
          <p class="dialog-description">点击任意图片可全屏查看，支持左右切换</p>
        </div>
        
        <!-- 截图网格 - 使用更大的空间 -->
        <div class="fullscreen-screenshots-grid">
          <div 
            v-for="(screenshot, index) in currentScreenshots"
            :key="index"
            class="fullscreen-screenshot-thumb"
            @click="viewImageInFullscreen(index)"
          >
            <el-image
              :src="getImageUrl(screenshot, true)"
              :alt="`截图 ${index + 1}`"
              fit="cover"
              lazy
              class="fullscreen-thumb-image"
            >
              <template #placeholder>
                <div class="thumb-placeholder">
                  <el-icon><Loading /></el-icon>
                  <span>加载中...</span>
                </div>
              </template>
              <template #error>
                <div class="thumb-error">
                  <el-icon><Picture /></el-icon>
                  <span>图片{{ index + 1 }}</span>
                </div>
              </template>
            </el-image>
            
            <!-- 图片序号和操作按钮 -->
            <div class="fullscreen-thumb-overlay">
              <span class="thumb-index">{{ index + 1 }}</span>
              <div class="thumb-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  circle 
                  :icon="ZoomIn"
                  @click.stop="viewImageInFullscreen(index)"
                  title="查看大图"
                />
              </div>
            </div>
          </div>
        </div>
        
        <!-- 对话框底部操作 -->
        <div class="dialog-footer-actions">
          <el-button 
            size="large" 
            @click="screenshotsDialogVisible = false"
          >
            关闭
          </el-button>
          <el-button 
            type="primary" 
            size="large" 
            @click="viewImageInFullscreen(0)"
            :icon="ZoomIn"
            v-if="currentScreenshots.length > 0"
          >
            查看第一张
          </el-button>
        </div>
      </div>
    </el-dialog>
    
    <!-- 增强版图片查看器 -->
    <EnhancedScreenshotsViewer
      v-if="showImageViewer"
      :images="currentScreenshots"
      :initial-index="currentImageIndex"
      @close="closeImageViewer"
      @show-grid="showGridFromViewer"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { ticketApi, userApi } from '@/api'
import SearchForm from '@/components/SearchForm.vue'
import OperatorTicketCard from '@/components/OperatorTicketCard.vue'
import AppPagination from '@/components/AppPagination.vue'
import EnhancedScreenshotsViewer from '@/components/EnhancedScreenshotsViewer.vue'
import {
  DataBoard,
  Tickets,
  FolderOpened,
  CircleCheck,
  UserFilled,
  User,
  ArrowRight,
  Refresh,
  DocumentRemove,
  Loading,
  Picture,
  ZoomIn
} from '@element-plus/icons-vue'

// 组件引用
const searchFormRef = ref(null)

// 数据状态
const loading = ref(false)
const saving = ref(false)
const tickets = ref([])
const operators = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchParams = ref({})

// 对话框状态
const remarkDialogVisible = ref(false)
const solutionDialogVisible = ref(false)
const typeDialogVisible = ref(false) // 新增工单类型对话框状态
const operatorDialogVisible = ref(false) // 新增处理人对话框状态
const screenshotsDialogVisible = ref(false)
const showImageViewer = ref(false)
const currentTicket = ref(null)
const currentScreenshots = ref([])
const currentImageIndex = ref(0)

// 表单数据
const remarkForm = ref({
  remark: ''
})

const solutionForm = ref({
  solution: ''
})

const typeForm = ref({
  ticket_type: ''
}) // 新增工单类型表单数据

const operatorForm = ref({
  operator_name: ''
}) // 新增处理人表单数据

// 统计数据
const stats = computed(() => {
  const openCount = tickets.value.filter(t => t.status === 'open').length
  const closedCount = tickets.value.filter(t => t.status === 'closed').length
  
  return {
    total: total.value,
    open: openCount,
    closed: closedCount,
    operators: operators.value.length
  }
})

// 获取工单列表
const fetchTickets = async () => {
  try {
    loading.value = true
    
    const params = {
      page: currentPage.value,
      limit: pageSize.value,
      ...searchParams.value
    }
    
    const response = await ticketApi.getTickets(params)
    tickets.value = response.tickets || []
    total.value = response.total || 0
    
  } catch (error) {
    console.error('获取工单列表失败:', error)
    ElMessage.error('获取工单列表失败')
  } finally {
    loading.value = false
  }
}

// 获取运维人员列表
const fetchOperators = async () => {
  try {
    const response = await userApi.getOperatorUsersName()
    operators.value = response.users || []
  } catch (error) {
    console.error('获取运维人员列表失败:', error)
  }
}

// 刷新数据
const refreshData = async () => {
  ElMessage.info('正在刷新数据...')
  await Promise.all([
    fetchTickets(),
    fetchOperators()
  ])
  ElMessage.success('数据刷新完成')
}

// 处理搜索
const handleSearch = (params) => {
  searchParams.value = params
  currentPage.value = 1
  fetchTickets()
}

// 处理重置
const handleReset = () => {
  searchParams.value = {}
  currentPage.value = 1
  fetchTickets()
}

// 处理分页变化
const handlePageChange = ({ page, limit }) => {
  currentPage.value = page
  pageSize.value = limit
  fetchTickets()
}

// 处理状态变更
const handleStatusChange = async (ticketId, newStatus) => {
  try {
    await ticketApi.updateTicketStatus(ticketId, newStatus)
    ElMessage.success('状态更新成功')
    
    // 更新本地数据
    const ticket = tickets.value.find(t => t.id === ticketId)
    if (ticket) {
      ticket.status = newStatus
      ticket.updated_at = new Date().toISOString()
    }
    
  } catch (error) {
    console.error('更新状态失败:', error)
    ElMessage.error('更新状态失败')
  }
}

// 处理备注
const handleRemark = (ticket) => {
  currentTicket.value = ticket
  remarkForm.value.remark = ticket.remark || ''
  remarkDialogVisible.value = true
}

// 保存备注
const saveRemark = async () => {
  try {
    saving.value = true
    
    await ticketApi.updateTicketRemark(currentTicket.value.id, remarkForm.value.remark)
    ElMessage.success('备注保存成功')
    
    // 更新本地数据
    const ticket = tickets.value.find(t => t.id === currentTicket.value.id)
    if (ticket) {
      ticket.remark = remarkForm.value.remark
      ticket.updated_at = new Date().toISOString()
    }
    
    remarkDialogVisible.value = false
    
  } catch (error) {
    console.error('保存备注失败:', error)
    ElMessage.error('保存备注失败')
  } finally {
    saving.value = false
  }
}

// 处理解决方案
const handleSolution = (ticket) => {
  currentTicket.value = ticket
  solutionForm.value.solution = ticket.solution || ''
  solutionDialogVisible.value = true
}

// 处理截图查看
const handleViewScreenshots = (ticket) => {
  if (!ticket.screenshots || ticket.screenshots.length === 0) {
    ElMessage.info('该工单暂无截图')
    return
  }
  
  currentTicket.value = ticket
  currentScreenshots.value = ticket.screenshots
  screenshotsDialogVisible.value = true
}

// 保存解决方案
const saveSolution = async () => {
  try {
    saving.value = true
    
    await ticketApi.updateTicketSolution(currentTicket.value.id, solutionForm.value.solution)
    ElMessage.success('解决方案保存成功')
    
    // 更新本地数据
    const ticket = tickets.value.find(t => t.id === currentTicket.value.id)
    if (ticket) {
      ticket.solution = solutionForm.value.solution
      ticket.updated_at = new Date().toISOString()
    }
    
    solutionDialogVisible.value = false
    
  } catch (error) {
    console.error('保存解决方案失败:', error)
    ElMessage.error('保存解决方案失败')
  } finally {
    saving.value = false
  }
}

// 处理工单类型修改
const handleTypeChange = (ticket) => {
  currentTicket.value = ticket
  typeForm.value.ticket_type = ticket.ticket_type || ''
  typeDialogVisible.value = true
}

// 保存工单类型
const saveType = async () => {
  try {
    saving.value = true
    
    await ticketApi.updateTicketType(currentTicket.value.id, typeForm.value.ticket_type)
    ElMessage.success('工单类型保存成功')
    
    // 更新本地数据
    const ticket = tickets.value.find(t => t.id === currentTicket.value.id)
    if (ticket) {
      ticket.ticket_type = typeForm.value.ticket_type
      ticket.updated_at = new Date().toISOString()
    }
    
    typeDialogVisible.value = false
    
  } catch (error) {
    console.error('保存工单类型失败:', error)
    ElMessage.error('保存工单类型失败')
  } finally {
    saving.value = false
  }
}

// 处理工单处理人修改
const handleOperatorChange = async (ticket) => {
  try {
    // 确保在打开对话框前刷新处理人列表
    await fetchOperators()
    
    currentTicket.value = ticket
    operatorForm.value.operator_name = ticket.operator_name || ''
    operatorDialogVisible.value = true
  } catch (error) {
    console.error('获取处理人列表失败:', error)
    ElMessage.error('获取处理人列表失败')
  }
}

// 保存工单处理人
const saveOperator = async () => {
  try {
    saving.value = true
    
    await ticketApi.updateTicketOperator(currentTicket.value.id, operatorForm.value.operator_name)
    ElMessage.success('处理人分配成功')
    
    // 更新本地数据
    const ticket = tickets.value.find(t => t.id === currentTicket.value.id)
    if (ticket) {
      ticket.operator_name = operatorForm.value.operator_name
      ticket.updated_at = new Date().toISOString()
    }
    
    operatorDialogVisible.value = false
    
  } catch (error) {
    console.error('分配处理人失败:', error)
    ElMessage.error('分配处理人失败')
  } finally {
    saving.value = false
  }
}

// 获取图片URL
const getImageUrl = (screenshot, useThumbnail = false) => {
  if (!screenshot) return ''
  
  const url = typeof screenshot === 'string' ? screenshot : screenshot.url || screenshot
  
  if (useThumbnail) {
    // 生成缩略图URL
    const lastDotIndex = url.lastIndexOf('.')
    if (lastDotIndex > -1) {
      return url.substring(0, lastDotIndex) + '.thumb' + url.substring(lastDotIndex)
    }
  }
  
  return url
}

// 在全屏对话框中查看图片
const viewImageInFullscreen = (index) => {
  screenshotsDialogVisible.value = false
  // 等待对话框关闭动画完成后再显示查看器
  setTimeout(() => {
    currentImageIndex.value = index
    showImageViewer.value = true
  }, 300)
}

// 关闭图片查看器
const closeImageViewer = () => {
  showImageViewer.value = false
}

// 从查看器显示网格
const showGridFromViewer = () => {
  showImageViewer.value = false
  screenshotsDialogVisible.value = true
}

// 页面初始化
onMounted(() => {
  Promise.all([
    fetchTickets(),
    fetchOperators()
  ])
})
</script>

<style scoped>
.admin-view {
  min-height: 100vh;
  padding: 40px 20px;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 16px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.title-icon {
  font-size: 40px;
  background: var(--primary-gradient);
  color: white;
  padding: 8px;
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-lg);
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0;
}

.stats-section {
  margin-bottom: 40px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}

.stat-card {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: var(--border-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 28px;
  box-shadow: var(--shadow);
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea, #764ba2);
}

.stat-icon.open {
  background: linear-gradient(135deg, #10b981, #059669);
}

.stat-icon.closed {
  background: linear-gradient(135deg, #6b7280, #4b5563);
}

.stat-icon.operators {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px;
  color: var(--text-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.quick-actions-section {
  margin-bottom: 40px;
}

.section-title {
  margin-bottom: 24px;
}

.section-title h2 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px;
  color: var(--text-primary);
}

.section-title p {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.action-card {
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.action-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.action-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  box-shadow: var(--shadow);
}

.action-icon.users {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
}

.action-icon.refresh {
  background: linear-gradient(135deg, #06b6d4, #0891b2);
}

.action-content {
  flex: 1;
}

.action-content h4 {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 4px;
  color: var(--text-primary);
}

.action-content p {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
}

.action-arrow {
  color: var(--text-muted);
  font-size: 18px;
  transition: all 0.3s ease;
}

.action-card:hover .action-arrow {
  color: var(--primary-color);
  transform: translateX(4px);
}

.tickets-section {
  margin-bottom: 40px;
}

.loading-section {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.loading-content {
  text-align: center;
  color: var(--text-secondary);
}

.loading-content p {
  margin-top: 16px;
  font-size: 16px;
}

.empty-section {
  min-height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon {
  font-size: 80px;
  color: var(--text-muted);
  opacity: 0.5;
}

.tickets-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 对话框样式 */
.edit-dialog {
  border-radius: var(--border-radius-lg);
}

.dialog-content {
  padding: 8px 0;
}

.ticket-info {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--border-color);
}

.ticket-info h4 {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.ticket-id {
  margin: 0;
  font-size: 14px;
  color: var(--text-muted);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* Element Plus 自定义样式 */
:deep(.el-dialog) {
  border-radius: var(--border-radius-lg);
}

:deep(.el-dialog__header) {
  background: var(--primary-gradient);
  color: white;
  padding: 20px 24px;
  border-radius: var(--border-radius-lg) var(--border-radius-lg) 0 0;
}

:deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

:deep(.el-dialog__close) {
  color: white;
}

:deep(.el-dialog__close:hover) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-dialog__body) {
  padding: 24px;
}

:deep(.el-dialog__footer) {
  padding: 20px 24px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

/* 截图全屏对话框样式 */
:deep(.screenshots-fullscreen-dialog) {
  .el-dialog {
    margin: 2vh auto;
    max-height: 96vh;
    display: flex;
    flex-direction: column;
  }
  
  .el-dialog__header {
    background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
    color: white;
    padding: 12px 20px;
    margin: 0;
    border-radius: var(--el-border-radius-base) var(--el-border-radius-base) 0 0;
    flex-shrink: 0;
  }
  
  .el-dialog__title {
    color: white;
    font-weight: 600;
    font-size: 16px;
  }
  
  .el-dialog__close {
    color: white;
    font-size: 18px;
  }
  
  .el-dialog__close:hover {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .el-dialog__body {
    padding: 16px 20px 20px;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }
}

.fullscreen-screenshots-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.dialog-header-info {
  padding: 0 0 12px 0;
  margin-bottom: 16px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.dialog-description {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 13px;
  text-align: center;
  background: var(--el-color-info-light-9);
  padding: 8px 12px;
  border-radius: var(--el-border-radius-base);
  border-left: 3px solid var(--el-color-primary);
}

.fullscreen-screenshots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  flex: 1;
  overflow-y: auto;
  padding: 0 4px;
  margin-bottom: 16px;
  min-height: 0;
  max-height: 70vh;
}

.fullscreen-screenshot-thumb {
  aspect-ratio: 1;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--el-border-color-light);
  background: var(--el-bg-color-page);
}

.fullscreen-screenshot-thumb:hover {
  transform: scale(1.02);
  border-color: var(--el-color-primary);
  box-shadow: var(--el-box-shadow);
}

.fullscreen-thumb-image {
  width: 100%;
  height: 100%;
}

.fullscreen-thumb-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.fullscreen-screenshot-thumb:hover .fullscreen-thumb-overlay {
  opacity: 1;
}

.thumb-index {
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  align-self: flex-start;
}

.thumb-actions {
  display: flex;
  justify-content: center;
  align-items: center;
}

.dialog-footer-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.thumb-placeholder,
.thumb-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
  color: var(--el-text-color-placeholder);
  background: var(--el-bg-color-page);
  font-size: 12px;
}

/* 截图全屏对话框样式 */
:deep(.screenshots-fullscreen-dialog) {
  .el-dialog {
    margin: 2vh auto;
    max-height: 96vh;
    display: flex;
    flex-direction: column;
  }
  
  .el-dialog__header {
    background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
    color: white;
    padding: 12px 20px;
    margin: 0;
    border-radius: var(--el-border-radius-base) var(--el-border-radius-base) 0 0;
    flex-shrink: 0;
  }
  
  .el-dialog__title {
    color: white;
    font-weight: 600;
    font-size: 16px;
  }
  
  .el-dialog__close {
    color: white;
    font-size: 18px;
  }
  
  .el-dialog__close:hover {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .el-dialog__body {
    padding: 16px 20px 20px;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }
}

.fullscreen-screenshots-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.dialog-header-info {
  padding: 0 0 12px 0;
  margin-bottom: 16px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.dialog-description {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 13px;
  text-align: center;
  background: var(--el-color-info-light-9);
  padding: 8px 12px;
  border-radius: var(--el-border-radius-base);
  border-left: 3px solid var(--el-color-primary);
}

.fullscreen-screenshots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  flex: 1;
  overflow-y: auto;
  padding: 0 4px;
  margin-bottom: 16px;
  min-height: 0;
  max-height: 70vh;
}

.fullscreen-screenshot-thumb {
  aspect-ratio: 1;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--el-border-color-light);
  background: var(--el-bg-color-page);
}

.fullscreen-screenshot-thumb:hover {
  transform: scale(1.02);
  border-color: var(--el-color-primary);
  box-shadow: var(--el-box-shadow);
}

.fullscreen-thumb-image {
  width: 100%;
  height: 100%;
}

.fullscreen-thumb-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.fullscreen-screenshot-thumb:hover .fullscreen-thumb-overlay {
  opacity: 1;
}

.thumb-index {
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  align-self: flex-start;
}

.thumb-actions {
  display: flex;
  justify-content: center;
  align-items: center;
}

.dialog-footer-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.thumb-placeholder,
.thumb-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
  color: var(--el-text-color-placeholder);
  background: var(--el-bg-color-page);
  font-size: 12px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-primary);
}

:deep(.el-textarea__inner) {
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
}

:deep(.el-textarea__inner:focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .admin-view {
    padding: 20px 16px;
  }
  
  .page-title {
    font-size: 28px;
    flex-direction: column;
    gap: 12px;
  }
  
  .title-icon {
    font-size: 32px;
    padding: 6px;
  }
  
  .page-subtitle {
    font-size: 14px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-card {
    padding: 20px;
    gap: 16px;
  }
  
  .stat-icon {
    width: 56px;
    height: 56px;
    font-size: 24px;
  }
  
  .stat-number {
    font-size: 28px;
  }
  
  .actions-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .action-card {
    padding: 20px;
    gap: 16px;
  }
  
  .action-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .edit-dialog {
    width: 95% !important;
  }
  
  .tickets-list {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .admin-view {
    padding: 16px 12px;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .title-icon {
    font-size: 28px;
    padding: 4px;
  }
  
  .stat-card {
    padding: 16px;
    gap: 12px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
    font-size: 20px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .action-card {
    padding: 16px;
    gap: 12px;
  }
  
  .action-icon {
    width: 40px;
    height: 40px;
    font-size: 18px;
  }
  
  .action-content h4 {
    font-size: 14px;
  }
  
  .action-content p {
    font-size: 12px;
  }
  
  .empty-icon {
    font-size: 60px;
  }
  
  .dialog-content {
    padding: 4px 0;
  }
  
  .dialog-footer {
    flex-direction: column;
  }
  
  .dialog-footer .el-button {
    width: 100%;
    margin: 0;
  }
}

/* 移动端截图对话框优化 */
@media (max-width: 768px) {
  :deep(.screenshots-fullscreen-dialog) {
    .el-dialog {
      margin: 3vh auto;
      max-height: 94vh;
      width: 95% !important;
    }
    
    .el-dialog__header {
      padding: 10px 16px;
    }
    
    .el-dialog__title {
      font-size: 14px;
    }
    
    .el-dialog__body {
      padding: 12px 16px 16px;
    }
  }
  
  .fullscreen-screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
    max-height: 60vh;
  }
  
  .dialog-description {
    font-size: 12px;
    padding: 6px 10px;
  }
  
  .dialog-footer-actions {
    flex-direction: column;
    gap: 8px;
    padding-top: 12px;
  }
  
  .dialog-footer-actions .el-button {
    width: 100%;
  }
}

@media (max-width: 480px) {
  :deep(.screenshots-fullscreen-dialog) {
    .el-dialog {
      margin: 2vh auto;
      max-height: 96vh;
      width: 98% !important;
    }
  }
  
  .fullscreen-screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 10px;
    max-height: 65vh;
  }
  
  .thumb-index {
    font-size: 11px;
    padding: 3px 6px;
  }
}
</style>