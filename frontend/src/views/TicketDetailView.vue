<template>
  <div class="ticket-detail-view">
    <div class="container">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-section">
        <div class="loading-content">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
      
      <!-- 错误状态 -->
      <div v-else-if="error" class="error-section">
        <el-result
          icon="error"
          title="加载失败"
          :sub-title="error"
        >
          <template #extra>
            <el-button type="primary" @click="fetchTicketDetail">
              重新加载
            </el-button>
            <el-button @click="$router.back()">
              返回
            </el-button>
          </template>
        </el-result>
      </div>
      
      <!-- 工单详情内容 -->
      <div v-else-if="ticket" class="ticket-detail-content">
        <!-- 页面头部 -->
        <div class="detail-header glass-card">
          <div class="header-content gradient-bg">
            <div class="header-left">
              <h1 class="ticket-title">
                <el-icon><DocumentCopy /></el-icon>
                工单详情
              </h1>
              <div class="ticket-meta">
                <span class="ticket-id">ID: {{ ticket.id }}</span>
                <el-tag 
                  :type="getStatusType(ticket.status)"
                  class="status-tag"
                  effect="dark"
                >
                  {{ getStatusText(ticket.status) }}
                </el-tag>
              </div>
            </div>
            <div class="header-actions">
              <el-button 
                @click="$router.back()"
                :icon="ArrowLeft"
                class="back-btn"
              >
                返回
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 详情内容 -->
        <div class="detail-body">
          <!-- 基本信息 -->
          <div class="section-card">
            <h2 class="section-title">{{ ticket.title }}</h2>
            <div class="meta-grid">
              <div class="meta-item">
                <span class="meta-label">工单ID</span>
                <span class="meta-value">{{ ticket.id }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">创建人</span>
                <span class="meta-value">{{ ticket.user_name || '未知' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">处理人</span>
                <span class="meta-value">{{ ticket.operator_name || '未分配' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">工单类型</span>
                <span class="meta-value">
                  <el-tag size="small" class="type-tag">
                    {{ ticket.ticket_type || '未分类' }}
                  </el-tag>
                </span>
              </div>
              <div class="meta-item">
                <span class="meta-label">创建时间</span>
                <span class="meta-value">{{ formatDateTimeToSecond(ticket.created_at) }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">更新时间</span>
                <span class="meta-value">{{ formatDateTimeToSecond(ticket.updated_at || ticket.created_at) }}</span>
              </div>
            </div>
          </div>
          
          <!-- 问题描述 -->
          <div class="section-card">
            <h3 class="section-title">
              <el-icon><Document /></el-icon>
              问题描述
            </h3>
            <div class="content-box">
              {{ ticket.content }}
            </div>
          </div>
          
          <!-- 备注信息 -->
          <div class="section-card">
            <h3 class="section-title">
              <el-icon><EditPen /></el-icon>
              备注信息
            </h3>
            <div v-if="ticket.remark" class="content-box">
              {{ ticket.remark }}
            </div>
            <div v-else class="empty-state">
              <el-icon><EditPen /></el-icon>
              <p>暂无备注信息</p>
            </div>
          </div>
          
          <!-- 解决方案 -->
          <div class="section-card">
            <h3 class="section-title">
              <el-icon><Opportunity /></el-icon>
              解决方案
            </h3>
            <div v-if="ticket.solution" class="content-box">
              {{ ticket.solution }}
            </div>
            <div v-else class="empty-state">
              <el-icon><Opportunity /></el-icon>
              <p>暂无解决方案</p>
            </div>
          </div>
          
          <!-- 问题截图 -->
          <div class="section-card">
            <h3 class="section-title">
              <el-icon><Picture /></el-icon>
              问题截图
            </h3>
            <div v-if="ticket.screenshots && ticket.screenshots.length > 0" class="screenshots-list">
              <div 
                v-for="(screenshot, index) in ticket.screenshots"
                :key="index"
                class="screenshot-item"
                @click="viewImageInFullscreen(index)"
              >
                <el-image
                  :src="getImageUrl(screenshot, true)"
                  :alt="`截图 ${index + 1}`"
                  fit="cover"
                  lazy
                  class="screenshot-image"
                >
                  <template #placeholder>
                    <div class="image-placeholder">
                      <el-icon><Loading /></el-icon>
                      <span>加载中...</span>
                    </div>
                  </template>
                  <template #error>
                    <div class="image-error">
                      <el-icon><Picture /></el-icon>
                      <span>图片{{ index + 1 }}</span>
                    </div>
                  </template>
                </el-image>
                <div class="screenshot-overlay">
                  <span class="screenshot-index">{{ index + 1 }}</span>
                  <div class="screenshot-actions">
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
            <div v-else class="empty-state">
              <el-icon><Picture /></el-icon>
              <p>暂无问题截图</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
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
      :thumbnail-suffix="'.thumb'"
      @close="closeImageViewer"
      @show-grid="showGridFromViewer"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ticketApi } from '@/api'
import { formatDate, formatDateTimeToSecond, getStatusText, getStatusType } from '@/utils/helpers'
import ScreenshotsDisplay from '@/components/ScreenshotsDisplay.vue'
import EnhancedScreenshotsViewer from '@/components/EnhancedScreenshotsViewer.vue'
import {  
  DocumentCopy,
  ArrowLeft,
  Document,
  EditPen,
  Opportunity,
  Picture,
  Camera,
  ZoomIn,
  Loading
} from '@element-plus/icons-vue'

const props = defineProps({
  id: {
    type: String,
    required: true
  }
})

const route = useRoute()
const router = useRouter()

// 数据状态
const loading = ref(false)
const error = ref('')
const ticket = ref(null)

// 截图对话框状态
const screenshotsDialogVisible = ref(false)
const showImageViewer = ref(false)
const currentTicket = ref(null)
const currentScreenshots = ref([])
const currentImageIndex = ref(0)

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

// 在全屏对话框中查看图片
const viewImageInFullscreen = (index) => {
  // 先确保当前截图已经加载到currentScreenshots中
  if (!currentScreenshots.value.length && ticket.value?.screenshots) {
    currentScreenshots.value = ticket.value.screenshots.map(screenshot => {
      // 确保使用原始图片URL，而不是缩略图URL
      return typeof screenshot === 'string' ? screenshot : screenshot.url || screenshot
    })
    currentTicket.value = ticket.value
  }
  
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

// 处理截图点击事件
const handleImageClick = (index, screenshot) => {
  // 不需要处理，因为现在使用统一的查看逻辑
}

// 处理查看所有截图事件
const handleViewAll = (screenshots) => {
  // 不需要处理，因为现在使用统一的查看逻辑
}

// 获取工单详情
const fetchTicketDetail = async () => {
  try {
    loading.value = true
    error.value = ''
    
    // 从查询参数中获取工单ID
    const ticketId = route.query.id
    if (!ticketId) {
      throw new Error('工单ID不能为空')
    }
    
    const response = await ticketApi.getTicketDetail(ticketId)
    ticket.value = response
    
  } catch (err) {
    console.error('获取工单详情失败:', err)
    error.value = err.response?.data?.error || err.message || '获取工单详情失败'
    ElMessage.error(error.value)
  } finally {
    loading.value = false
  }
}

// 页面初始化
onMounted(() => {
  fetchTicketDetail()
})
</script>

<style scoped>
.ticket-detail-view {
  min-height: 100vh;
  padding: 40px 20px;
}

.container {
  max-width: 1000px;
  margin: 0 auto;
}

.loading-section,
.error-section {
  min-height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-content {
  text-align: center;
  color: var(--text-secondary);
}

.loading-content p {
  margin-top: 16px;
  font-size: 16px;
}

.ticket-detail-content {
  margin: 0 auto;
}

.detail-header {
  margin-bottom: 32px;
  overflow: hidden;
}

.header-content {
  padding: 32px;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 20px;
}

.header-left {
  flex: 1;
  min-width: 0;
}

.ticket-title {
  margin: 0 0 16px;
  font-size: 32px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 16px;
  color: white;
}

.ticket-title .el-icon {
  font-size: 36px;
  background: rgba(255, 255, 255, 0.2);
  padding: 8px;
  border-radius: var(--border-radius);
}

.ticket-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.ticket-id {
  font-size: 16px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.2);
  padding: 6px 12px;
  border-radius: var(--border-radius);
}

.status-tag {
  font-weight: 500;
  letter-spacing: 0.5px;
}

.header-actions {
  flex-shrink: 0;
}

.back-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  color: white;
  transform: translateY(-2px);
}

.detail-body {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: var(--border-radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-lg);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.section-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.section-title {
  margin: 0 0 20px;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 12px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--border-color);
}

.section-title:first-child {
  font-size: 24px;
  color: var(--primary-color);
  border-bottom: 2px solid var(--primary-color);
}

.section-title .el-icon {
  color: var(--primary-color);
  font-size: 22px;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meta-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-muted);
}

.meta-value {
  font-size: 16px;
  color: var(--text-primary);
  font-weight: 500;
}

.type-tag {
  background: var(--primary-gradient);
  color: white;
  border: none;
  font-weight: 500;
}

.content-box {
  background: var(--bg-secondary);
  border-radius: var(--border-radius);
  padding: 20px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  border-left: 4px solid var(--primary-color);
  font-size: 16px;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-muted);
  background: var(--bg-secondary);
  border-radius: var(--border-radius);
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

/* 截图列表样式 */
.screenshots-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.screenshot-item {
  aspect-ratio: 1;
  border-radius: var(--border-radius);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--border-color);
  background: var(--bg-secondary);
}

.screenshot-item:hover {
  transform: scale(1.03);
  border-color: var(--primary-color);
  box-shadow: var(--shadow);
}

.screenshot-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.screenshot-overlay {
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

.screenshot-item:hover .screenshot-overlay {
  opacity: 1;
}

.screenshot-index {
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  align-self: flex-start;
}

.screenshot-actions {
  display: flex;
  justify-content: center;
  align-items: center;
}

.image-placeholder,
.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
  color: var(--text-muted);
  background: var(--bg-secondary);
  font-size: 12px;
}

.screenshots-action {
  text-align: center;
  padding: 20px 0;
}

.screenshots-action .el-button {
  font-weight: 500;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .ticket-detail-view {
    padding: 20px 16px;
  }
  
  .header-content {
    padding: 24px 20px;
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
  
  .ticket-title {
    font-size: 24px;
    flex-direction: column;
    text-align: center;
    gap: 12px;
  }
  
  .ticket-title .el-icon {
    font-size: 28px;
    align-self: center;
  }
  
  .ticket-meta {
    justify-content: center;
  }
  
  .header-actions {
    align-self: center;
  }
  
  .section-card {
    padding: 20px 16px;
  }
  
  .section-title {
    font-size: 18px;
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }
  
  .section-title:first-child {
    font-size: 20px;
  }
  
  .meta-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .meta-item {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
  
  .meta-label {
    flex-shrink: 0;
  }
  
  .screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 16px;
  }
  
  .screenshots-list {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
}

@media (max-width: 480px) {
  .ticket-detail-view {
    padding: 16px 12px;
  }
  
  .header-content {
    padding: 20px 16px;
  }
  
  .ticket-title {
    font-size: 20px;
  }
  
  .ticket-title .el-icon {
    font-size: 24px;
  }
  
  .section-card {
    padding: 16px 12px;
  }
  
  .section-title {
    font-size: 16px;
  }
  
  .section-title:first-child {
    font-size: 18px;
  }
  
  .content-box {
    padding: 16px;
    font-size: 14px;
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
  
  .screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 12px;
  }
  
  .screenshots-list {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 8px;
  }
  
  .screenshot-index {
    font-size: 10px;
    padding: 3px 6px;
  }
}
</style>