<template>
  <div class="screenshots-display">
    <!-- 无截图状态 -->
    <div v-if="!screenshots || screenshots.length === 0" class="empty-screenshots">
      <div class="empty-state">
        <el-icon class="empty-icon"><Camera /></el-icon>
        <p class="empty-text">{{ emptyText }}</p>
      </div>
    </div>

    <!-- 有截图显示 -->
    <div v-else class="screenshots-container">
      <!-- 单张截图显示 -->
      <template v-if="screenshots.length === 1">
        <div 
          class="screenshot-item single-screenshot"
          @click="handleImageClick(0)"
        >
          <el-image
            :src="getImageUrl(screenshots[0], showThumbnail)"
            :alt="getAltText(0)"
            :fit="imageFit"
            :lazy="lazy"
            class="screenshot-image"
          >
            <template #placeholder>
              <div class="image-placeholder">
                <el-icon class="loading-icon"><Loading /></el-icon>
                <span>加载中...</span>
              </div>
            </template>
            <template #error>
              <div class="image-error">
                <el-icon class="error-icon"><Picture /></el-icon>
                <span>加载失败</span>
              </div>
            </template>
          </el-image>
          <div class="screenshot-overlay">
            <el-icon class="overlay-icon"><ZoomIn /></el-icon>
            <span class="overlay-text">点击查看</span>
          </div>
        </div>
      </template>

      <!-- 多张截图显示 -->
      <template v-else-if="displayMode === 'grid'">
        <div class="screenshots-grid" :class="gridClass">
          <div 
            v-for="(screenshot, index) in visibleScreenshots"
            :key="index"
            class="screenshot-item grid-screenshot"
            @click="handleImageClick(index)"
          >
            <el-image
              :src="getImageUrl(screenshot, showThumbnail)"
              :alt="getAltText(index)"
              :fit="imageFit"
              :lazy="lazy"
              class="screenshot-image"
            >
              <template #placeholder>
                <div class="image-placeholder">
                  <el-icon class="loading-icon"><Loading /></el-icon>
                </div>
              </template>
              <template #error>
                <div class="image-error">
                  <el-icon class="error-icon"><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="screenshot-overlay">
              <el-icon class="overlay-icon"><ZoomIn /></el-icon>
            </div>
          </div>
          
          <!-- 更多截图按钮 -->
          <div 
            v-if="hasMoreScreenshots" 
            class="screenshot-item more-screenshots"
            @click="handleViewAll"
          >
            <div class="more-content">
              <el-icon class="more-icon"><MoreFilled /></el-icon>
              <span class="more-text">+{{ remainingCount }}</span>
              <span class="more-description">张截图</span>
            </div>
          </div>
        </div>
      </template>

      <!-- 计数器模式（紧凑显示） -->
      <template v-else-if="displayMode === 'counter'">
        <div 
          class="screenshots-counter"
          @click="handleViewAll"
        >
          <div class="counter-content">
            <el-icon class="counter-icon"><Picture /></el-icon>
            <span class="counter-text">{{ screenshots.length }} 张问题截图</span>
            <el-icon class="arrow-icon"><ArrowRight /></el-icon>
          </div>
        </div>
      </template>
    </div>

    <!-- 全屏图片查看器 -->
    <el-image-viewer
      v-if="showViewer && !useEnhancedViewer"
      :url-list="fullImageUrls"
      :initial-index="currentIndex"
      @close="closeViewer"
      :hide-on-click-modal="hideOnClickModal"
      teleported
    />
    
    <!-- 增强版图片查看器 -->
    <EnhancedScreenshotsViewer
      v-if="showViewer && useEnhancedViewer"
      :images="screenshots"
      :initial-index="currentIndex"
      :thumbnail-suffix="thumbnailSuffix"
      @close="closeViewer"
      @show-grid="showGridFromViewer"
    />

    <!-- 所有截图对话框 -->
    <el-dialog
      v-model="showAllDialog"
      :title="`所有截图 (${screenshots.length}张)`"
      width="95%"
      :before-close="closeAllDialog"
      class="screenshots-dialog"
      destroy-on-close
      :close-on-click-modal="true"
      :close-on-press-escape="true"
      top="5vh"
    >
      <div class="dialog-content">
        <!-- 对话框头部信息 -->
        <div class="dialog-header-info">
          <p class="dialog-description">点击任意图片可全屏查看，支持左右切换</p>
        </div>
        
        <!-- 截图网格 -->
        <div class="all-screenshots-grid">
          <div 
            v-for="(screenshot, index) in screenshots"
            :key="index"
            class="screenshot-thumb"
            @click="viewImageInDialog(index)"
          >
            <el-image
              :src="getImageUrl(screenshot, true)"
              :alt="getAltText(index)"
              fit="cover"
              lazy
              class="thumb-image"
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
            <div class="thumb-overlay">
              <span class="thumb-index">{{ index + 1 }}</span>
              <div class="thumb-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  circle 
                  :icon="ZoomIn"
                  @click.stop="viewImageInDialog(index)"
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
            @click="closeAllDialog"
          >
            关闭
          </el-button>
          <el-button 
            type="primary" 
            size="large" 
            @click="viewImageInDialog(0)"
            :icon="ZoomIn"
          >
            查看第一张
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import EnhancedScreenshotsViewer from '@/components/EnhancedScreenshotsViewer.vue'
import { 
  Camera, 
  Picture, 
  ZoomIn, 
  Loading, 
  MoreFilled, 
  ArrowRight 
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  // 截图数组
  screenshots: {
    type: Array,
    default: () => []
  },
  // 显示模式：grid(网格), counter(计数器)
  displayMode: {
    type: String,
    default: 'grid',
    validator: (value) => ['grid', 'counter'].includes(value)
  },
  // 最大显示数量（网格模式）
  maxVisible: {
    type: Number,
    default: 4
  },
  // 网格列数
  gridColumns: {
    type: Number,
    default: 3
  },
  // 是否显示缩略图
  showThumbnail: {
    type: Boolean,
    default: true
  },
  // 图片适配方式
  imageFit: {
    type: String,
    default: 'cover'
  },
  // 是否懒加载
  lazy: {
    type: Boolean,
    default: true
  },
  // 空状态文本
  emptyText: {
    type: String,
    default: '暂无问题截图'
  },
  // 是否点击模态框背景关闭查看器
  hideOnClickModal: {
    type: Boolean,
    default: true
  },
  // 自定义缩略图后缀
  thumbnailSuffix: {
    type: String,
    default: '.thumb'
  },
  // 是否使用增强版查看器
  useEnhancedViewer: {
    type: Boolean,
    default: true
  }
})

// Emits
const emit = defineEmits(['image-click', 'view-all', 'viewer-close'])

// 响应式数据
const showViewer = ref(false)
const showAllDialog = ref(false)
const currentIndex = ref(0)

// 计算属性
const visibleScreenshots = computed(() => {
  if (props.displayMode === 'counter') return []
  return props.screenshots.slice(0, props.maxVisible)
})

const hasMoreScreenshots = computed(() => {
  return props.displayMode === 'grid' && props.screenshots.length > props.maxVisible
})

const remainingCount = computed(() => {
  return Math.max(0, props.screenshots.length - props.maxVisible)
})

const gridClass = computed(() => {
  return `grid-cols-${Math.min(props.gridColumns, visibleScreenshots.value.length)}`
})

const fullImageUrls = computed(() => {
  return props.screenshots.map(screenshot => getImageUrl(screenshot, false))
})

// 方法
const getImageUrl = (screenshot, useThumbnail = false) => {
  if (!screenshot) return ''
  
  const url = typeof screenshot === 'string' ? screenshot : screenshot.url || screenshot
  
  if (useThumbnail && props.showThumbnail) {
    // 生成缩略图URL
    const lastDotIndex = url.lastIndexOf('.')
    if (lastDotIndex > -1) {
      return url.substring(0, lastDotIndex) + props.thumbnailSuffix + url.substring(lastDotIndex)
    }
  }
  
  return url
}

const getAltText = (index) => {
  return `问题截图 ${index + 1}`
}

const handleImageClick = (index) => {
  currentIndex.value = index
  showViewer.value = true
  emit('image-click', index, props.screenshots[index])
}

const handleViewAll = () => {
  showAllDialog.value = true
  emit('view-all', props.screenshots)
}

const viewImageInDialog = (index) => {
  showAllDialog.value = false
  // 等待对话框关闭动画完成后再显示查看器
  setTimeout(() => {
    handleImageClick(index)
  }, 300)
}

const closeViewer = () => {
  showViewer.value = false
  emit('viewer-close')
}

const closeAllDialog = () => {
  showAllDialog.value = false
}

const showGridFromViewer = () => {
  showViewer.value = false
  showAllDialog.value = true
}

// 监听截图变化，重置状态
watch(() => props.screenshots, () => {
  showViewer.value = false
  showAllDialog.value = false
  currentIndex.value = 0
})
</script>

<style scoped>
.screenshots-display {
  width: 100%;
}

/* 空状态样式 */
.empty-screenshots {
  padding: 24px;
  text-align: center;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--el-text-color-placeholder);
}

.empty-icon {
  font-size: 48px;
  opacity: 0.5;
}

.empty-text {
  margin: 0;
  font-size: 14px;
}

/* 截图容器 */
.screenshots-container {
  width: 100%;
}

/* 单张截图样式 */
.single-screenshot {
  max-width: 300px;
  aspect-ratio: 4/3;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--el-border-color-light);
  background: var(--el-bg-color-page);
}

.single-screenshot:hover {
  transform: scale(1.02);
  border-color: var(--el-color-primary);
  box-shadow: var(--el-box-shadow-light);
}

/* 网格布局 */
.screenshots-grid {
  display: grid;
  gap: 12px;
  width: 100%;
}

.grid-cols-1 { grid-template-columns: repeat(1, 1fr); }
.grid-cols-2 { grid-template-columns: repeat(2, 1fr); }
.grid-cols-3 { grid-template-columns: repeat(3, 1fr); }
.grid-cols-4 { grid-template-columns: repeat(4, 1fr); }

.grid-screenshot {
  aspect-ratio: 1;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--el-border-color-light);
  background: var(--el-bg-color-page);
}

.grid-screenshot:hover {
  transform: scale(1.05);
  border-color: var(--el-color-primary);
  box-shadow: var(--el-box-shadow);
}

/* 更多截图按钮 */
.more-screenshots {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-bg-color-page);
  border: 2px dashed var(--el-border-color);
  color: var(--el-text-color-placeholder);
  cursor: pointer;
  transition: all 0.3s ease;
}

.more-screenshots:hover {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.more-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.more-icon {
  font-size: 24px;
}

.more-text {
  font-size: 18px;
  font-weight: 600;
}

.more-description {
  font-size: 12px;
}

/* 计数器模式 */
.screenshots-counter {
  padding: 16px;
  border: 2px solid var(--el-border-color);
  border-radius: var(--el-border-radius-base);
  cursor: pointer;
  transition: all 0.3s ease;
  background: var(--el-bg-color-page);
}

.screenshots-counter:hover {
  border-color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.counter-content {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--el-text-color-regular);
}

.counter-icon {
  font-size: 20px;
  color: var(--el-color-primary);
}

.counter-text {
  flex: 1;
  font-weight: 500;
}

.arrow-icon {
  font-size: 16px;
  color: var(--el-text-color-placeholder);
}

/* 图片样式 */
.screenshot-image {
  width: 100%;
  height: 100%;
}

/* 悬停遮罩 */
.screenshot-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.screenshot-item:hover .screenshot-overlay {
  opacity: 1;
}

.overlay-icon {
  font-size: 24px;
}

.overlay-text {
  font-size: 12px;
  font-weight: 500;
}

/* 加载和错误状态 */
.image-placeholder,
.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 8px;
  color: var(--el-text-color-placeholder);
}

.loading-icon,
.error-icon {
  font-size: 24px;
}

.image-placeholder span,
.image-error span {
  font-size: 12px;
}

/* 对话框中的截图网格 */
.dialog-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
}

.dialog-header-info {
  padding: 0 0 12px 0;
  margin-bottom: 12px;
  border-bottom: 1px solid var(--el-border-color-light);
  flex-shrink: 0;
}

.dialog-description {
  margin: 0;
  color: var(--el-text-color-regular);
  font-size: 12px;
  text-align: center;
  background: var(--el-color-info-light-9);
  padding: 6px 10px;
  border-radius: var(--el-border-radius-base);
  border-left: 3px solid var(--el-color-primary);
}

.all-screenshots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 10px;
  flex: 1;
  overflow-y: auto;
  padding: 0 2px;
  margin-bottom: 12px;
  min-height: 0;
  max-height: 45vh;
}

.screenshot-thumb {
  aspect-ratio: 1;
  border-radius: var(--el-border-radius-base);
  overflow: hidden;
  position: relative;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid var(--el-border-color-light);
  background: var(--el-bg-color-page);
}

.screenshot-thumb:hover {
  transform: scale(1.02);
  border-color: var(--el-color-primary);
  box-shadow: var(--el-box-shadow);
}

.thumb-image {
  width: 100%;
  height: 100%;
}

.thumb-overlay {
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

.screenshot-thumb:hover .thumb-overlay {
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
  padding-top: 12px;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .screenshots-grid {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 8px;
  }
  
  .single-screenshot {
    max-width: 100%;
  }
  
  .dialog-description {
    font-size: 13px;
    padding: 10px 12px;
  }
}

@media (max-width: 480px) {
  .screenshots-grid {
    grid-template-columns: 1fr !important;
  }
  
  .thumb-index {
    font-size: 11px;
    padding: 3px 6px;
  }
}

/* 对话框样式覆盖 */
:deep(.screenshots-dialog) {
  .el-dialog {
    margin: 8vh auto;
    max-height: 75vh;
    max-width: min(90vw, 800px);
    display: flex;
    flex-direction: column;
  }
  
  .el-dialog__header {
    background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
    color: white;
    padding: 10px 16px;
    margin: 0;
    border-radius: var(--el-border-radius-base) var(--el-border-radius-base) 0 0;
    flex-shrink: 0;
  }
  
  .el-dialog__title {
    color: white;
    font-weight: 600;
    font-size: 14px;
  }
  
  .el-dialog__close {
    color: white;
    font-size: 16px;
  }
  
  .el-dialog__close:hover {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .el-dialog__body {
    padding: 12px 16px 16px;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }
}

/* 移动端对话框优化 */
@media (max-width: 768px) {
  :deep(.screenshots-dialog) {
    .el-dialog {
      margin: 5vh auto;
      max-height: 80vh;
      width: 95% !important;
    }
    
    .el-dialog__header {
      padding: 8px 12px;
    }
    
    .el-dialog__title {
      font-size: 13px;
    }
    
    .el-dialog__body {
      padding: 10px 12px 12px;
    }
  }
  
  .all-screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(110px, 1fr));
    gap: 8px;
    max-height: 40vh;
  }
  
  .dialog-header-info {
    padding: 0 0 8px 0;
    margin-bottom: 8px;
  }
  
  .dialog-description {
    font-size: 11px;
    padding: 5px 8px;
  }
  
  .dialog-footer-actions {
    flex-direction: column;
    gap: 6px;
    padding-top: 8px;
  }
  
  .dialog-footer-actions .el-button {
    width: 100%;
    padding: 8px 12px;
  }
}

/* 小屏幕设备进一步优化 */
@media (max-width: 480px) {
  :deep(.screenshots-dialog) {
    .el-dialog {
      margin: 3vh auto;
      max-height: 85vh;
      width: 98% !important;
    }
  }
  
  .all-screenshots-grid {
    grid-template-columns: repeat(auto-fill, minmax(90px, 1fr));
    gap: 6px;
    max-height: 50vh;
  }
}
</style>