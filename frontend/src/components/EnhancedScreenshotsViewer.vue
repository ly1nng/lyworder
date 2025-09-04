<template>
  <div class="enhanced-screenshots-viewer">
    <!-- 快速操作栏 -->
    <div class="viewer-toolbar" :class="{ 'mobile-hidden': isMobile && !showControls }">
      <div class="toolbar-left">
        <span class="image-counter">{{ currentIndex + 1 }} / {{ images.length }}</span>
      </div>
      <div class="toolbar-right">
        <el-button
          type="info"
          size="small"
          circle
          :icon="Grid"
          @click="showGridView"
          title="网格视图"
        />
        <el-button
          type="info"
          size="small"
          circle
          :icon="FullScreen"
          @click="toggleFullscreen"
          title="全屏"
        />
        <el-button
          type="info"
          size="small"
          circle
          :icon="Close"
          @click="closeViewer"
          title="关闭"
        />
      </div>
    </div>

    <!-- 主图片显示区域 -->
    <div class="main-image-container" @click="toggleControlsVisibility">
      <el-image
        :src="currentImageUrl"
        :alt="`图片 ${currentIndex + 1}`"
        fit="contain"
        class="main-image"
        :preview-src-list="[]"
      >
        <template #error>
          <div class="image-error">
            <el-icon class="error-icon"><Picture /></el-icon>
            <span>图片加载失败</span>
          </div>
        </template>
      </el-image>

      <!-- 导航箭头 -->
      <div v-if="images.length > 1" class="navigation-arrows" :class="{ hidden: !showControls }">
        <el-button
          v-if="currentIndex > 0"
          type="info"
          size="large"
          circle
          :icon="ArrowLeft"
          class="nav-button nav-prev"
          @click.stop="prevImage"
        />
        <el-button
          v-if="currentIndex < images.length - 1"
          type="info"
          size="large"
          circle
          :icon="ArrowRight"
          class="nav-button nav-next"
          @click.stop="nextImage"
        />
      </div>
    </div>

    <!-- 缩略图导航条 -->
    <div v-if="images.length > 1" class="thumbnail-strip" :class="{ hidden: !showControls }">
      <div class="thumbnail-container">
        <div
          v-for="(image, index) in images"
          :key="index"
          class="thumbnail-item"
          :class="{ active: index === currentIndex }"
          @click="goToImage(index)"
        >
          <el-image
            :src="getImageUrl(image, true)"
            :alt="`缩略图 ${index + 1}`"
            fit="cover"
            class="thumbnail-image"
          >
            <template #error>
              <div class="thumbnail-error">
                <el-icon><Picture /></el-icon>
              </div>
            </template>
          </el-image>
          <span class="thumbnail-index">{{ index + 1 }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAppStore } from '@/stores/app'
import {
  ArrowLeft,
  ArrowRight,
  Close,
  Picture,
  Grid,
  FullScreen
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  images: {
    type: Array,
    required: true,
    default: () => []
  },
  initialIndex: {
    type: Number,
    default: 0
  },
  thumbnailSuffix: {
    type: String,
    default: '.thumb'
  }
})

// Emits
const emit = defineEmits(['close', 'show-grid'])

// 使用appStore检测是否为移动端
const appStore = useAppStore()
const isMobile = computed(() => appStore.isMobile)

// 响应式数据
const currentIndex = ref(props.initialIndex)
const showControls = ref(true)
const controlsTimeout = ref(null)

// 计算属性
const currentImageUrl = computed(() => {
  return getImageUrl(props.images[currentIndex.value], false)
})

// 方法
const getImageUrl = (image, useThumbnail = false) => {
  if (!image) return ''
  
  const url = typeof image === 'string' ? image : image.url || image
  
  if (useThumbnail) {
    const lastDotIndex = url.lastIndexOf('.')
    if (lastDotIndex > -1) {
      return url.substring(0, lastDotIndex) + props.thumbnailSuffix + url.substring(lastDotIndex)
    }
  }
  
  return url
}

const prevImage = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--
  }
}

const nextImage = () => {
  if (currentIndex.value < props.images.length - 1) {
    currentIndex.value++
  }
}

const goToImage = (index) => {
  currentIndex.value = index
}

const closeViewer = () => {
  emit('close')
}

const showGridView = () => {
  emit('show-grid')
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const toggleControlsVisibility = () => {
  showControls.value = !showControls.value
  
  if (showControls.value && isMobile.value) {
    // 移动端3秒后自动隐藏控制元素
    clearTimeout(controlsTimeout.value)
    controlsTimeout.value = setTimeout(() => {
      showControls.value = false
    }, 3000)
  } else {
    clearTimeout(controlsTimeout.value)
  }
}

// 键盘事件处理
const handleKeydown = (event) => {
  switch (event.key) {
    case 'ArrowLeft':
      event.preventDefault()
      prevImage()
      break
    case 'ArrowRight':
      event.preventDefault()
      nextImage()
      break
    case 'Escape':
      event.preventDefault()
      closeViewer()
      break
    case ' ':
      event.preventDefault()
      toggleControlsVisibility()
      break
  }
}

// 生命周期
onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
  
  // 初始化控制元素自动隐藏
  if (isMobile.value) {
    controlsTimeout.value = setTimeout(() => {
      showControls.value = false
    }, 3000)
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  clearTimeout(controlsTimeout.value)
})
</script>

<style scoped>
.enhanced-screenshots-viewer {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.95);
  z-index: 3000;
  display: flex;
  flex-direction: column;
}

.viewer-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  transition: all 0.3s ease;
}

.viewer-toolbar.hidden {
  transform: translateY(-100%);
  opacity: 0;
}

/* 移动端优化：在移动端默认隐藏工具栏 */
@media (max-width: 768px) {
  .viewer-toolbar {
    transform: translateY(0);
    opacity: 1;
  }
  
  .viewer-toolbar.mobile-hidden {
    transform: translateY(-100%);
    opacity: 0;
  }
}

.toolbar-left .image-counter {
  font-size: 16px;
  font-weight: 500;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.main-image-container {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  overflow: hidden;
}

.main-image {
  max-width: 90vw;
  max-height: 80vh;
  object-fit: contain;
  width: auto;
  height: auto;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 18px;
}

.error-icon {
  font-size: 48px;
}

.navigation-arrows {
  position: absolute;
  top: 50%;
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 0 40px;
  pointer-events: none;
  transition: all 0.3s ease;
}

.navigation-arrows.hidden {
  opacity: 0;
}

.nav-button {
  pointer-events: all;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  backdrop-filter: blur(10px);
  transform: translateY(-50%);
}

.nav-button:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-50%) scale(1.1);
}

.thumbnail-strip {
  padding: 16px 20px;
  background: rgba(0, 0, 0, 0.8);
  transition: all 0.3s ease;
}

.thumbnail-strip.hidden {
  transform: translateY(100%);
  opacity: 0;
}

.thumbnail-container {
  display: flex;
  justify-content: center;
  gap: 8px;
  overflow-x: auto;
  padding: 4px 0;
}

.thumbnail-item {
  position: relative;
  width: 60px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.thumbnail-item.active {
  border-color: var(--el-color-primary);
  transform: scale(1.1);
}

.thumbnail-item:hover {
  transform: scale(1.05);
}

.thumbnail-image {
  width: 100%;
  height: 100%;
}

.thumbnail-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
}

.thumbnail-index {
  position: absolute;
  top: 2px;
  right: 2px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  font-size: 10px;
  padding: 2px 4px;
  border-radius: 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .viewer-toolbar {
    padding: 12px 16px;
  }
  
  .toolbar-left .image-counter {
    font-size: 14px;
  }
  
  .main-image {
    max-width: 95vw;
    max-height: 85vh;
    width: auto;
    height: auto;
  }
  
  .navigation-arrows {
    padding: 0 20px;
  }
  
  .nav-button {
    width: 40px;
    height: 40px;
  }
  
  .thumbnail-strip {
    padding: 12px 16px;
  }
  
  .thumbnail-item {
    width: 50px;
    height: 50px;
  }
  
  /* 移动端优化：禁用复杂动画以提高性能 */
  .thumbnail-item {
    transition: none;
  }
  
  .thumbnail-item:hover {
    transform: none;
  }
  
  .thumbnail-item.active {
    transform: none;
  }
  
  .nav-button {
    transition: none;
  }
  
  .nav-button:hover {
    transform: translateY(-50%) scale(1.1);
    background: rgba(255, 255, 255, 0.3);
  }
}

@media (max-width: 480px) {
  .viewer-toolbar {
    padding: 8px 12px;
  }
  
  .toolbar-right {
    gap: 4px;
  }
  
  .main-image {
    max-width: 98vw;
    max-height: 80vh;
    width: auto;
    height: auto;
  }
  
  .navigation-arrows {
    padding: 0 10px;
  }
  
  .nav-button {
    width: 36px;
    height: 36px;
  }
  
  .thumbnail-item {
    width: 45px;
    height: 45px;
  }
  
  .thumbnail-index {
    font-size: 9px;
    padding: 1px 3px;
  }
  
  /* 移动端优化：进一步简化样式 */
  .viewer-toolbar,
  .thumbnail-strip {
    transition: none;
  }
  
  .navigation-arrows {
    transition: none;
  }
}
</style>