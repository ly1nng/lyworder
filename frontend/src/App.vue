<template>
  <div id="app">
    <div v-if="appStore.loading" class="app-loading">
      <div class="loading-content">
        <div class="loading-spinner"></div>
        <p>加载中...</p>
      </div>
    </div>
    
    <template v-else>
      <!-- 顶部导航 -->
      <AppHeader v-if="showHeader" />
      
      <!-- 主要内容区域 -->
      <main class="app-main" :class="{ 'with-header': showHeader }">
        <router-view v-slot="{ Component, route }">
          <transition 
            :name="route.meta.transition || 'fade'"
            mode="out-in"
          >
            <component :is="Component" :key="route.path" />
          </transition>
        </router-view>
      </main>
    </template>
    
    <!-- 全局错误提示 -->
    <el-backtop 
      :right="20" 
      :bottom="80"
      :visibility-height="300"
    >
      <div class="backtop-btn">
        <el-icon><CaretTop /></el-icon>
      </div>
    </el-backtop>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import AppHeader from '@/components/AppHeader.vue'
import { CaretTop } from '@element-plus/icons-vue'

const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

// 计算是否显示头部导航
const showHeader = computed(() => {
  const hideHeaderRoutes = ['/login']
  return !hideHeaderRoutes.includes(route.path) && userStore.isLoggedIn
})

// 初始化应用
onMounted(async () => {
  appStore.setLoading(true)
  appStore.initResize()
  
  try {
    // 尝试获取用户信息
    await userStore.fetchUserInfo()
  } catch (error) {
    // 用户未登录或登录过期，这是正常情况
    console.log('用户未登录')
  } finally {
    appStore.setLoading(false)
  }
})

// 清理事件监听
onUnmounted(() => {
  appStore.cleanup()
})
</script>

<style>
.app-loading {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.loading-content {
  text-align: center;
  color: var(--text-primary);
}

.loading-content p {
  margin-top: 16px;
  font-size: 16px;
  font-weight: 500;
}

.app-main {
  min-height: 100vh;
  transition: all 0.3s ease;
}

.app-main.with-header {
  padding-top: 70px;
  min-height: calc(100vh - 70px);
}

/* 回到顶部按钮样式 */
.backtop-btn {
  width: 40px;
  height: 40px;
  background: var(--primary-gradient);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: var(--shadow-lg);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.backtop-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-xl);
}

.backtop-btn .el-icon {
  font-size: 20px;
  font-weight: bold;
}

/* 页面过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-main.with-header {
    padding-top: 60px;
    min-height: calc(100vh - 60px);
  }
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: var(--primary-gradient);
  border-radius: 4px;
  transition: all 0.3s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, var(--primary-dark), var(--secondary-color));
}

/* 修复Element Plus组件的边距 */
.el-main {
  padding: 0;
}

/* 优化移动端体验 */
@media (max-width: 480px) {
  .loading-content p {
    font-size: 14px;
  }
  
  .backtop-btn {
    width: 36px;
    height: 36px;
  }
  
  .backtop-btn .el-icon {
    font-size: 18px;
  }
}
</style>