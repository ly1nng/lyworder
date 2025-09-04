<template>
  <div class="not-found-page">
    <div class="not-found-background">
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>
    </div>

    <div class="not-found-container">
      <div class="error-content">
        <div class="error-animation">
          <div class="error-code">
            <span class="digit">4</span>
            <span class="digit middle">0</span>
            <span class="digit">4</span>
          </div>
          <div class="error-icon">
            <el-icon size="120"><Document /></el-icon>
          </div>
        </div>

        <div class="error-text">
          <h1>页面未找到</h1>
          <p>抱歉，您访问的页面不存在或已被移动</p>
          <div class="error-details">
            <p>可能的原因：</p>
            <ul>
              <li>页面地址输入错误</li>
              <li>页面已被删除或移动</li>
              <li>您没有访问该页面的权限</li>
              <li>服务器暂时无法响应</li>
            </ul>
          </div>
        </div>

        <div class="error-actions">
          <el-button 
            type="primary" 
            size="large" 
            @click="goHome"
            class="action-button"
          >
            <el-icon><HomeFilled /></el-icon>
            返回首页
          </el-button>
          <el-button 
            size="large" 
            @click="goBack"
            class="action-button"
          >
            <el-icon><Back /></el-icon>
            返回上页
          </el-button>
          <el-button 
            size="large" 
            @click="refresh"
            class="action-button"
          >
            <el-icon><Refresh /></el-icon>
            刷新页面
          </el-button>
        </div>

        <div class="help-section">
          <el-card class="help-card">
            <template #header>
              <div class="help-header">
                <el-icon><QuestionFilled /></el-icon>
                <span>需要帮助？</span>
              </div>
            </template>
            <div class="help-content">
              <p>如果您认为这是一个错误，请联系系统管理员</p>
              <div class="contact-info">
                <div class="contact-item">
                  <el-icon><Message /></el-icon>
                  <span>support@company.com</span>
                </div>
                <div class="contact-item">
                  <el-icon><Phone /></el-icon>
                  <span>400-000-0000</span>
                </div>
              </div>
            </div>
          </el-card>
        </div>

        <div class="suggestions">
          <h3>您可能想要访问：</h3>
          <div class="suggestion-links">
            <router-link to="/" class="suggestion-link">
              <el-icon><Ticket /></el-icon>
              <span>提交工单</span>
            </router-link>
            <router-link to="/tickets" class="suggestion-link">
              <el-icon><Document /></el-icon>
              <span>我的工单</span>
            </router-link>
            <router-link to="/admin" class="suggestion-link" v-if="userStore.isAdmin">
              <el-icon><Setting /></el-icon>
              <span>管理后台</span>
            </router-link>
            <router-link to="/operator" class="suggestion-link" v-if="userStore.isOperator">
              <el-icon><Tools /></el-icon>
              <span>运维工作台</span>
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { 
  Document, 
  HomeFilled, 
  Back, 
  Refresh, 
  QuestionFilled,
  Message,
  Phone,
  Ticket,
  Setting,
  Tools
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

// 方法
const goHome = () => {
  router.push('/')
}

const goBack = () => {
  router.go(-1)
}

const refresh = () => {
  window.location.reload()
}
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.not-found-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.shape {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  animation: float 8s ease-in-out infinite;
}

.shape-1 {
  width: 100px;
  height: 100px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 20%;
  right: 20%;
  animation-delay: 2s;
}

.shape-3 {
  width: 80px;
  height: 80px;
  bottom: 30%;
  left: 30%;
  animation-delay: 4s;
}

.shape-4 {
  width: 120px;
  height: 120px;
  bottom: 10%;
  right: 10%;
  animation-delay: 6s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.3;
  }
  50% {
    transform: translateY(-30px) rotate(180deg);
    opacity: 0.6;
  }
}

.not-found-container {
  max-width: 800px;
  width: 100%;
  padding: 40px;
  text-align: center;
}

.error-content {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 60px 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.error-animation {
  position: relative;
  margin-bottom: 40px;
}

.error-code {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.digit {
  font-size: 8rem;
  font-weight: 900;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  animation: bounce 2s ease-in-out infinite;
}

.digit.middle {
  animation-delay: 0.5s;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
  60% {
    transform: translateY(-10px);
  }
}

.error-icon {
  color: #a0aec0;
  opacity: 0.6;
  animation: pulse 3s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.6;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.05);
  }
}

.error-text {
  margin-bottom: 40px;
}

.error-text h1 {
  margin: 0 0 16px 0;
  font-size: 2.5rem;
  font-weight: 700;
  color: #2d3748;
}

.error-text > p {
  margin: 0 0 24px 0;
  font-size: 1.2rem;
  color: #718096;
  line-height: 1.6;
}

.error-details {
  background: #f7fafc;
  border-radius: 12px;
  padding: 24px;
  margin: 24px 0;
  text-align: left;
  border-left: 4px solid #667eea;
}

.error-details p {
  margin: 0 0 12px 0;
  font-weight: 600;
  color: #4a5568;
}

.error-details ul {
  margin: 0;
  padding-left: 20px;
  color: #718096;
}

.error-details li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.error-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.action-button {
  min-width: 140px;
  height: 48px;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.action-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.help-section {
  margin-bottom: 40px;
}

.help-card {
  border-radius: 16px;
  border: none;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.help-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #4a5568;
}

.help-content p {
  margin: 0 0 16px 0;
  color: #718096;
  line-height: 1.6;
}

.contact-info {
  display: flex;
  justify-content: center;
  gap: 32px;
  flex-wrap: wrap;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #667eea;
  font-weight: 500;
}

.suggestions {
  text-align: center;
}

.suggestions h3 {
  margin: 0 0 24px 0;
  font-size: 1.3rem;
  font-weight: 600;
  color: #4a5568;
}

.suggestion-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.suggestion-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px 24px;
  background: #f7fafc;
  border-radius: 12px;
  text-decoration: none;
  color: #4a5568;
  font-weight: 500;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.suggestion-link:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .not-found-container {
    padding: 20px;
  }

  .error-content {
    padding: 40px 24px;
    border-radius: 20px;
  }

  .digit {
    font-size: 5rem;
  }

  .error-text h1 {
    font-size: 2rem;
  }

  .error-text > p {
    font-size: 1.1rem;
  }

  .error-actions {
    flex-direction: column;
    align-items: center;
  }

  .action-button {
    width: 100%;
    max-width: 280px;
  }

  .contact-info {
    flex-direction: column;
    gap: 16px;
  }

  .suggestion-links {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .not-found-container {
    padding: 15px;
  }

  .error-content {
    padding: 30px 20px;
  }

  .digit {
    font-size: 4rem;
  }

  .error-code {
    gap: 10px;
  }

  .error-text h1 {
    font-size: 1.8rem;
  }

  .error-details {
    padding: 20px;
  }

  .action-button {
    min-width: 120px;
    height: 44px;
  }
}
</style>