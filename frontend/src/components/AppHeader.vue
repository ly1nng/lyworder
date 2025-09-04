<template>
  <header class="app-header glass-header">
    <div class="header-container">
      <div class="header-left">
        <router-link to="/" class="logo-link">
          <div class="logo">
            <el-icon :size="32" class="logo-icon">
              <Ticket />
            </el-icon>
            <span class="logo-text text-gradient">运维工单系统</span>
          </div>
        </router-link>
      </div>
      
      <nav class="header-nav" v-if="!isMobile">
        <router-link 
          v-for="item in navItems" 
          :key="item.name"
          :to="item.path" 
          class="nav-item"
          :class="{ active: $route.path === item.path }"
        >
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <span>{{ item.title }}</span>
        </router-link>
      </nav>
      
      <div class="header-right">
        <!-- 桌面版用户菜单 -->
        <div v-if="!isMobile" class="user-menu">
          <el-dropdown @command="handleCommand" placement="bottom-end">
            <div class="user-info glass-card">
              <el-avatar :size="32" class="user-avatar">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="user-name">{{ userStore.userName }}</span>
              <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="user-dropdown">
                <el-dropdown-item command="tickets">
                  <el-icon><Tickets /></el-icon>
                  我的工单
                </el-dropdown-item>
                <el-dropdown-item v-if="userStore.isOperator" command="operator">
                  <el-icon><Tools /></el-icon>
                  我的运维
                </el-dropdown-item>
                <el-dropdown-item v-if="userStore.isAdmin" command="admin">
                  <el-icon><Setting /></el-icon>
                  管理后台
                </el-dropdown-item>
                <el-dropdown-item v-if="userStore.isAdmin" command="users">
                  <el-icon><UserFilled /></el-icon>
                  用户管理
                </el-dropdown-item>
                <el-dropdown-item divided command="logout" class="logout-item">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        
        <!-- 移动版菜单按钮 -->
        <el-button 
          v-else 
          type="primary" 
          :icon="Menu" 
          circle 
          @click="mobileMenuVisible = true"
        />
      </div>
    </div>
  </header>
  
  <!-- 移动版侧边菜单 -->
  <el-drawer
    v-model="mobileMenuVisible"
    direction="rtl"
    :with-header="false"
    size="280px"
    class="mobile-menu-drawer"
  >
    <div class="mobile-menu">
      <div class="mobile-user-info">
        <el-avatar :size="48" class="user-avatar">
          <el-icon><User /></el-icon>
        </el-avatar>
        <div class="user-details">
          <div class="user-name">{{ userStore.userName }}</div>
          <div class="user-role">{{ getRoleText(userStore.user?.role) }}</div>
        </div>
      </div>
      
      <div class="mobile-nav">
        <router-link 
          v-for="item in allNavItems" 
          :key="item.name"
          :to="item.path" 
          class="mobile-nav-item"
          @click="mobileMenuVisible = false"
        >
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <span>{{ item.title }}</span>
        </router-link>
        
        <div class="mobile-nav-item logout" @click="handleCommand('logout')">
          <el-icon><SwitchButton /></el-icon>
          <span>退出登录</span>
        </div>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { getRoleText } from '@/utils/helpers'
import {
  Ticket,
  User,
  ArrowDown,
  Menu,
  Tickets,
  Tools,
  Setting,
  UserFilled,
  SwitchButton,
  House,
  DataBoard
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const mobileMenuVisible = ref(false)

const isMobile = computed(() => appStore.isMobile)

// 基础导航项
const baseNavItems = [
  { name: 'home', path: '/', title: '提交工单', icon: House },
  { name: 'tickets', path: '/tickets', title: '我的工单', icon: Tickets }
]

// 管理员和运维人员可见的导航项
const adminNavItems = [
  { name: 'operator', path: '/operator', title: '我的运维', icon: Tools, roles: ['admin', 'operator'] },
  { name: 'admin', path: '/admin', title: '管理后台', icon: DataBoard, roles: ['admin'] }
]

// 计算当前用户可见的导航项
const navItems = computed(() => {
  let items = [...baseNavItems]
  
  adminNavItems.forEach(item => {
    if (!item.roles || item.roles.includes(userStore.user?.role)) {
      items.push(item)
    }
  })
  
  return items
})

// 移动端显示所有可访问的导航项
const allNavItems = computed(() => {
  let items = [...baseNavItems]
  
  if (userStore.isOperator) {
    items.push({ name: 'operator', path: '/operator', title: '我的运维', icon: Tools })
  }
  
  if (userStore.isAdmin) {
    items.push({ name: 'admin', path: '/admin', title: '管理后台', icon: DataBoard })
    items.push({ name: 'users', path: '/admin/users', title: '用户管理', icon: UserFilled })
  }
  
  return items
})

const handleCommand = (command) => {
  switch (command) {
    case 'tickets':
      router.push('/tickets')
      break
    case 'operator':
      router.push('/operator')
      break
    case 'admin':
      router.push('/admin')
      break
    case 'users':
      router.push('/admin/users')
      break
    case 'logout':
      userStore.logout()
      break
  }
}
</script>

<style scoped>
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  padding: 0 20px;
  height: 70px;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  flex-shrink: 0;
}

.logo-link {
  text-decoration: none;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-icon {
  background: var(--primary-gradient);
  color: white;
  padding: 6px;
  border-radius: 12px;
  box-shadow: var(--shadow);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.header-nav {
  display: flex;
  gap: 8px;
  flex: 1;
  justify-content: center;
  max-width: 600px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: var(--border-radius-xl);
  text-decoration: none;
  color: var(--text-secondary);
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.nav-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: var(--primary-gradient);
  opacity: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: -1;
}

.nav-item:hover::before,
.nav-item.active::before {
  left: 0;
  opacity: 0.1;
}

.nav-item:hover,
.nav-item.active {
  color: var(--primary-color);
  transform: translateY(-2px);
}

.header-right {
  flex-shrink: 0;
}

.user-menu .user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  border-radius: var(--border-radius-xl);
  background: rgba(255, 255, 255, 0.8);
  /* 移动端优化：简化过渡效果 */
  transition: all 0.2s ease;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: var(--shadow);
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.9);
  transform: translateY(-1px);
  box-shadow: var(--shadow-lg);
}

.user-avatar {
  background: var(--primary-gradient);
}

.user-name {
  font-weight: 500;
  color: var(--text-primary);
}

.dropdown-icon {
  transition: transform 0.3s ease;
}

.user-info:hover .dropdown-icon {
  transform: rotate(180deg);
}

/* 下拉菜单样式 */
.user-dropdown {
  border-radius: var(--border-radius);
  overflow: hidden;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
  padding: 8px 0;
  /* 移动端优化：简化过渡效果 */
  transition: all 0.2s ease;
}

.user-dropdown .el-dropdown-menu__item {
  padding: 12px 20px;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  /* 移动端优化：简化过渡效果 */
  transition: all 0.2s ease;
}

.user-dropdown .el-dropdown-menu__item:hover {
  background: rgba(102, 126, 234, 0.1);
}

.user-dropdown .el-dropdown-menu__item:focus {
  background: rgba(102, 126, 234, 0.1);
}

.user-dropdown .logout-item {
  color: var(--danger-color);
  border-top: 1px solid var(--border-color);
}

.user-dropdown .logout-item:hover {
  background: rgba(231, 57, 70, 0.1);
}

/* 移动端样式 */
.mobile-menu {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--primary-gradient);
  border-radius: var(--border-radius-lg);
  color: white;
  margin-bottom: 30px;
}

.user-details {
  flex: 1;
}

.user-details .user-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 4px;
}

.user-details .user-role {
  font-size: 14px;
  opacity: 0.9;
}

.mobile-nav {
  flex: 1;
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-radius: var(--border-radius);
  text-decoration: none;
  color: var(--text-primary);
  font-weight: 500;
  /* 移动端优化：简化过渡效果 */
  transition: all 0.2s ease;
  margin-bottom: 8px;
}

.mobile-nav-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: var(--primary-color);
}

.mobile-nav-item.logout {
  color: var(--danger-color);
  margin-top: auto;
  cursor: pointer;
}

.mobile-nav-item.logout:hover {
  background: rgba(231, 57, 70, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-header {
    padding: 0 16px;
    height: 60px;
  }
  
  .logo-text {
    font-size: 18px;
  }
  
  .logo-icon {
    padding: 4px;
  }
  
  /* 移动端性能优化：禁用复杂动画 */
  .user-info,
  .user-dropdown,
  .user-dropdown .el-dropdown-menu__item,
  .mobile-nav-item {
    transition: none !important;
    transform: none !important;
  }
  
  .user-info:hover {
    background: rgba(255, 255, 255, 0.9) !important;
    transform: none !important;
    box-shadow: var(--shadow-lg) !important;
  }
  
  .user-dropdown .el-dropdown-menu__item:hover {
    background: rgba(102, 126, 234, 0.1) !important;
  }
  
  .mobile-nav-item:hover {
    background: rgba(102, 126, 234, 0.1) !important;
    color: var(--primary-color) !important;
  }
}

@media (max-width: 480px) {
  .app-header {
    padding: 0 12px;
  }
  
  .logo-text {
    display: none;
  }
}
</style>