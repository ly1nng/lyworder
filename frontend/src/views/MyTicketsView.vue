<template>
  <div class="my-tickets-view">
    <div class="container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">
          <el-icon class="title-icon"><Tickets /></el-icon>
          我的工单
        </h1>
        <p class="page-subtitle">查看和管理您提交的所有工单</p>
      </div>
      
      <!-- 搜索表单 -->
      <SearchForm
        ref="searchFormRef"
        :loading="loading"
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
            description="暂无工单记录"
            :image-size="120"
          >
            <template #image>
              <el-icon class="empty-icon"><DocumentRemove /></el-icon>
            </template>
            <el-button 
              type="primary" 
              @click="$router.push('/')"
              :icon="DocumentAdd"
            >
              提交新工单
            </el-button>
          </el-empty>
        </div>
        
        <div v-else class="tickets-list">
          <!-- 移动端优化：使用v-for.keyed以提高列表渲染性能 -->
          <OperatorTicketCard
            v-for="ticket in tickets"
            :key="ticket.id"
            :ticket="ticket"
            :editable-status="false"
            @view-screenshots="handleViewScreenshots"
            @remark="viewRemark"
            @solution="viewSolution"
          >
            <template #actions="{ ticket }">
              <el-button 
                type="success" 
                size="small" 
                @click="viewRemark(ticket)"
                :icon="EditPen"
                class="gradient-button"
              >
                备注
              </el-button>
              <el-button 
                type="warning" 
                size="small" 
                @click="viewSolution(ticket)"
                :icon="Opportunity"
                class="gradient-button"
              >
                解决方案
              </el-button>
              <el-button 
                type="info" 
                size="small" 
                @click="viewDetail(ticket)"
                :icon="InfoFilled"
                class="gradient-button"
              >
                详情
              </el-button>
            </template>
          </OperatorTicketCard>
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
    
    <!-- 备注查看对话框 -->
    <el-dialog
      v-model="remarkDialogVisible"
      title="工单备注"
      :width="isMobile ? '95%' : '700px'"
      class="content-dialog"
      :class="{ 'mobile-dialog': isMobile }"
      destroy-on-close
      :top="isMobile ? '5vh' : '15vh'"
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <div class="section-card">
          <h4 class="section-title">
            <el-icon><EditPen /></el-icon>
            备注信息
          </h4>
          <div v-if="currentTicket?.remark" class="content-box">
            {{ currentTicket.remark }}
          </div>
          <div v-else class="empty-state">
            <el-icon><EditPen /></el-icon>
            <p>暂无备注信息</p>
          </div>
        </div>
        
        <div class="update-info" v-if="currentTicket?.updated_at">
          <span>最后更新: {{ formatDate(currentTicket.updated_at) }}</span>
        </div>
      </div>
    </el-dialog>
    
    <!-- 解决方案查看对话框 -->
    <el-dialog
      v-model="solutionDialogVisible"
      title="解决方案"
      :width="isMobile ? '95%' : '700px'"
      class="content-dialog"
      :class="{ 'mobile-dialog': isMobile }"
      destroy-on-close
      :top="isMobile ? '5vh' : '15vh'"
    >
      <div class="dialog-content">
        <div class="ticket-info">
          <h4>{{ currentTicket?.title }}</h4>
          <p class="ticket-id">工单ID: {{ currentTicket?.id }}</p>
        </div>
        
        <div class="section-card">
          <h4 class="section-title">
            <el-icon><Opportunity /></el-icon>
            解决方案
          </h4>
          <div v-if="currentTicket?.solution" class="content-box">
            {{ currentTicket.solution }}
          </div>
          <div v-else class="empty-state">
            <el-icon><Opportunity /></el-icon>
            <p>暂无解决方案</p>
          </div>
        </div>
        
        <div class="solution-meta" v-if="currentTicket?.solution">
          <div class="meta-row">
            <span>处理人员: {{ currentTicket.operator_name || '未知' }}</span>
          </div>
          <div class="meta-row">
            <span>处理时间: {{ formatDate(currentTicket.updated_at || currentTicket.created_at) }}</span>
          </div>
        </div>
      </div>
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
            <!-- 移动端优化：使用lazy加载和简化图片处理 -->
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
import { ref, onMounted, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { ticketApi } from '@/api'
import { formatDate } from '@/utils/helpers'
import SearchForm from '@/components/SearchForm.vue'
import OperatorTicketCard from '@/components/OperatorTicketCard.vue'
import AppPagination from '@/components/AppPagination.vue'
import EnhancedScreenshotsViewer from '@/components/EnhancedScreenshotsViewer.vue'
import {
  Tickets,
  DocumentRemove,
  DocumentAdd,
  EditPen,
  Opportunity,
  InfoFilled,
  Picture,
  Loading,
  ZoomIn
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

// 组件引用
const searchFormRef = ref(null)

// 数据状态
const loading = ref(false)
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

// 是否为移动端
const isMobile = computed(() => appStore.isMobile)

// 获取工单列表 - 添加性能优化
const fetchTickets = async () => {
  try {
    loading.value = true
    
    // 移动端显示全部工单，不分页
    const params = {
      user_name: userStore.userName,
      ...(isMobile.value ? {} : { page: currentPage.value, limit: pageSize.value }),
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
  fetchTickets()
}

// 查看备注
const viewRemark = (ticket) => {
  console.log('点击查看备注', ticket)
  currentTicket.value = ticket
  remarkDialogVisible.value = true
  console.log('对话框状态', remarkDialogVisible.value)
}

// 查看解决方案
const viewSolution = (ticket) => {
  console.log('点击查看解决方案', ticket)
  currentTicket.value = ticket
  solutionDialogVisible.value = true
  console.log('对话框状态', solutionDialogVisible.value)
}

// 查看详情
const viewDetail = (ticket) => {
  router.push(`/tickets/detail/${ticket.id}`)
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

// 获取图片URL - 添加移动端优化
const getImageUrl = (screenshot, useThumbnail = false) => {
  if (!screenshot) return ''
  
  const url = typeof screenshot === 'string' ? screenshot : screenshot.url || screenshot
  
  // 移动端优化：总是使用缩略图以减少数据传输
  if (isMobile.value || useThumbnail) {
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

// 页面可见性变化监听器
let visibilityChangeListener = null

// 页面初始化
onMounted(() => {
  fetchTickets()
  
  // 添加页面可见性监听
  visibilityChangeListener = () => {
    if (document.hidden) {
      // 页面隐藏时暂停不必要的操作
      console.log('页面隐藏，暂停不必要的操作')
    } else {
      // 页面显示时恢复操作
      console.log('页面显示，恢复操作')
    }
  }
  
  document.addEventListener('visibilitychange', visibilityChangeListener)
})

// 清理事件监听
onUnmounted(() => {
  if (visibilityChangeListener) {
    document.removeEventListener('visibilitychange', visibilityChangeListener)
  }
})
</script>

<style scoped>
.my-tickets-view {
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
.content-dialog {
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

.section-card {
  background: var(--bg-secondary);
  border-radius: var(--border-radius);
  padding: 20px;
  margin-bottom: 20px;
}

.section-title {
  margin: 0 0 16px;
  font-size: 16px;
  font-weight: 600;
  color: var(--primary-color);
  display: flex;
  align-items: center;
  gap: 8px;
}

.content-box {
  background: white;
  border-radius: var(--border-radius);
  padding: 16px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  border-left: 4px solid var(--primary-color);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
}

.empty-state .el-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  margin: 0;
  font-size: 16px;
}

.screenshots-action {
  text-align: center;
  padding: 20px 0;
}

.screenshots-action .el-button {
  font-weight: 500;
}

.update-info {
  text-align: right;
  font-size: 14px;
  color: var(--text-muted);
  margin-top: 16px;
}

.solution-meta {
  background: var(--bg-secondary);
  border-radius: var(--border-radius);
  padding: 16px;
  font-size: 14px;
  color: var(--text-secondary);
}

.meta-row {
  margin-bottom: 8px;
}

.meta-row:last-child {
  margin-bottom: 0;
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

.mobile-dialog .section-card {
  padding: 16px;
  margin-bottom: 16px;
}

.mobile-dialog .section-title {
  font-size: 14px;
  margin-bottom: 12px;
}

.mobile-dialog .content-box {
  padding: 12px;
  font-size: 14px;
}

.mobile-dialog .empty-state {
  padding: 30px 15px;
}

.mobile-dialog .empty-state .el-icon {
  font-size: 36px;
  margin-bottom: 12px;
}

.mobile-dialog .empty-state p {
  font-size: 14px;
}

.mobile-dialog .update-info {
  text-align: center;
  font-size: 12px;
  margin-top: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .my-tickets-view {
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
  
  .content-dialog {
    width: 95% !important;
    margin: 5vh auto;
  }
  
  .dialog-content {
    padding: 4px 0;
  }
  
  .section-card {
    padding: 16px;
    margin-bottom: 16px;
  }
  
  .tickets-list {
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .my-tickets-view {
    padding: 16px 12px;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .title-icon {
    font-size: 28px;
    padding: 4px;
  }
  
  .empty-icon {
    font-size: 60px;
  }
  
  .section-card {
    padding: 12px;
  }
  
  .content-box {
    padding: 12px;
  }
  
  .empty-state {
    padding: 30px 15px;
  }
  
  .empty-state .el-icon {
    font-size: 36px;
  }
  
  .empty-state p {
    font-size: 14px;
  }
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
