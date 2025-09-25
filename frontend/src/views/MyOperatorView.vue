<template>
  <div class="my-operator-view">
    <div class="container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon class="title-icon"><Tools /></el-icon>
          我的运维
        </h1>
        <p class="page-subtitle">处理和管理分配给您的运维工单</p>
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
              <p class="stat-label">进行中</p>
            </div>
          </div>
          
          <div class="stat-card glass-card">
            <div class="stat-icon closed">
              <el-icon><CircleCheck /></el-icon>
            </div>
            <div class="stat-content">
              <h3 class="stat-number">{{ stats.closed }}</h3>
              <p class="stat-label">已解决</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 搜索表单 -->
      <SearchForm
        ref="searchFormRef"
        :loading="loading"
        show-creator
        @search="handleSearch"
        @reset="handleReset"
      />
      
      <!-- 工单列表 -->
      <div class="tickets-section">
        <div v-if="loading" class="loading-section">
          <div class="loading-content">
            <div class="loading-spinner"></div>
            <p>加载中...</p>
          </div>
        </div>
        
        <div v-else-if="tickets.length === 0" class="empty-section">
          <el-empty 
            description="暂无分配的工单"
            :image-size="120"
          >
            <template #image>
              <el-icon class="empty-icon"><FolderRemove /></el-icon>
            </template>
          </el-empty>
        </div>
        
        <div v-else class="tickets-list">
          <OperatorTicketCard
            v-for="ticket in tickets"
            :key="ticket.id"
            :ticket="ticket"
            @status-change="handleStatusChange"
            @remark="handleRemark"
            @solution="handleSolution"
            @view-screenshots="handleViewScreenshots"
          />
        </div>
      </div>
      
      <!-- 分页 -->
      <AppPagination
        v-if="!loading && tickets.length > 0 && !isMobile"
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
      :width="isMobile ? '95%' : '600px'"
      class="edit-dialog"
      :class="{ 'mobile-dialog': isMobile }"
      :top="isMobile ? '5vh' : '15vh'"
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
      :width="isMobile ? '95%' : '600px'"
      class="edit-dialog"
      :class="{ 'mobile-dialog': isMobile }"
      :top="isMobile ? '5vh' : '15vh'"
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
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ticketApi } from '@/api'
import SearchForm from '@/components/SearchForm.vue'
import OperatorTicketCard from '@/components/OperatorTicketCard.vue'
import AppPagination from '@/components/AppPagination.vue'
import EnhancedScreenshotsViewer from '@/components/EnhancedScreenshotsViewer.vue'
import {
  Tools,
  Tickets,
  FolderOpened,
  CircleCheck,
  FolderRemove,
  Loading,
  Picture,
  ZoomIn
} from '@element-plus/icons-vue'

const userStore = useUserStore()
const appStore = useAppStore()

// 组件引用
const searchFormRef = ref(null)

// 数据状态
const loading = ref(false)
const saving = ref(false)
const tickets = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchParams = ref({})

// 对话框状态
const remarkDialogVisible = ref(false)
const solutionDialogVisible = ref(false)
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

// 统计数据
const stats = computed(() => {
  const openCount = tickets.value.filter(t => t.status === 'open').length
  const closedCount = tickets.value.filter(t => t.status === 'closed').length
  
  return {
    total: total.value,
    open: openCount,
    closed: closedCount
  }
})

// 是否为移动端
const isMobile = computed(() => appStore.isMobile)

// 获取工单列表
const fetchTickets = async () => {
  try {
    loading.value = true
    
    // 移动端显示全部工单，不分页
    const isMobile = appStore.isMobile
    const params = {
      operator_name: userStore.userName,
      ...(isMobile ? {} : { page: currentPage.value, limit: pageSize.value }),
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
  fetchTickets().then(() => {
    // 分页切换后自动滚动到筛选栏
    scrollToSearchForm()
  })
}

// 滚动到筛选栏
const scrollToSearchForm = () => {
  // 查找搜索表单元素
  const searchForm = document.querySelector('.search-form')
  if (searchForm) {
    // 获取搜索表单的位置
    const rect = searchForm.getBoundingClientRect()
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop
    // 滚动到搜索表单位置
    const offsetTop = rect.top + scrollTop
    window.scrollTo({ top: offsetTop, behavior: 'smooth' })
  } else {
    // 如果找不到搜索表单，滚动到页面顶部
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
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
  fetchTickets()
})
</script>

<style scoped>
.my-operator-view {
  min-height: 100vh;
  padding: 40px 20px;
}

.container {
  max-width: 1200px;
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
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
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

/* 移动端对话框优化 */
.mobile-dialog {
  margin: 0 auto !important;
}

.mobile-dialog :deep(.el-dialog__body) {
  padding: 16px;
}

.mobile-dialog .ticket-info h4 {
  font-size: 16px;
}

.mobile-dialog .ticket-id {
  font-size: 12px;
}

.mobile-dialog :deep(.el-form-item__label) {
  font-size: 14px;
}

.mobile-dialog :deep(.el-textarea__inner) {
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .my-operator-view {
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
  
  .edit-dialog {
    width: 95% !important;
  }
  
  .tickets-list {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .my-operator-view {
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

/* 移动端对话框优化 */
.mobile-dialog {
  margin: 0 auto !important;
}

.mobile-dialog :deep(.el-dialog__body) {
  padding: 16px;
}

.mobile-dialog .ticket-info h4 {
  font-size: 16px;
}

.mobile-dialog .ticket-id {
  font-size: 12px;
}

.mobile-dialog :deep(.el-form-item__label) {
  font-size: 14px;
}

.mobile-dialog :deep(.el-textarea__inner) {
  font-size: 14px;
}
</style>