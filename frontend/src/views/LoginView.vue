<template>
  <div class="login-page">
    <div class="login-background">
      <div class="bg-pattern"></div>
      <div class="floating-elements">
        <div class="element element-1"></div>
        <div class="element element-2"></div>
        <div class="element element-3"></div>
      </div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="logo-section">
            <div class="logo">
              <el-icon size="40"><Ticket /></el-icon>
            </div>
            <h1>运维工单系统</h1>
            <p>Operations Ticket Management System</p>
          </div>
        </div>

        <div class="login-body">
          <div class="welcome-text">
            <h2>欢迎登录</h2>
            <p>请使用您的账户登录系统</p>
          </div>

          <div class="login-methods">
            <el-button 
              type="primary" 
              size="large" 
              class="oauth-button"
              @click="handleOAuthLogin"
              :loading="loginLoading"
            >
              <el-icon><User /></el-icon>
              OAuth 登录
            </el-button>

            <div class="divider">
              <span>或</span>
            </div>

            <el-form
              ref="loginFormRef"
              :model="loginForm"
              :rules="loginRules"
              @submit.prevent="handleManualLogin"
            >
              <el-form-item prop="username">
                <el-input
                  v-model="loginForm.username"
                  placeholder="用户名"
                  size="large"
                  prefix-icon="User"
                  clearable
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  v-model="loginForm.password"
                  type="password"
                  placeholder="密码"
                  size="large"
                  prefix-icon="Lock"
                  show-password
                  @keyup.enter="handleManualLogin"
                />
              </el-form-item>
              <el-button
                type="primary"
                size="large"
                class="login-button"
                :loading="loginLoading"
                @click="handleManualLogin"
              >
                登录
              </el-button>
            </el-form>
          </div>
        </div>

        <div class="login-footer">
          <div class="tips">
            <el-icon><InfoFilled /></el-icon>
            <span>首次登录请联系管理员分配账户</span>
          </div>
          <div class="version">
            Version 2.0.0
          </div>
        </div>
      </div>

      <div class="features-section">
        <h3>系统特性</h3>
        <div class="features-grid">
          <div class="feature-item">
            <el-icon size="24"><Document /></el-icon>
            <h4>工单管理</h4>
            <p>高效的工单提交、跟踪和处理流程</p>
          </div>
          <div class="feature-item">
            <el-icon size="24"><UserFilled /></el-icon>
            <h4>角色权限</h4>
            <p>细粒度的用户角色和权限控制</p>
          </div>
          <div class="feature-item">
            <el-icon size="24"><ChatLineRound /></el-icon>
            <h4>实时协作</h4>
            <p>团队成员间的实时沟通和协作</p>
          </div>
          <div class="feature-item">
            <el-icon size="24"><DataAnalysis /></el-icon>
            <h4>数据分析</h4>
            <p>详细的工单统计和分析报告</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  User, 
  Lock, 
  Ticket, 
  InfoFilled,
  Document,
  UserFilled,
  ChatLineRound,
  DataAnalysis
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

// 响应式数据
const loginLoading = ref(false)
const loginFormRef = ref()

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 50, message: '用户名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' }
  ]
}

// 方法
const handleOAuthLogin = () => {
  loginLoading.value = true
  // 重定向到OAuth登录端点
  window.location.href = '/login'
}

const handleManualLogin = async () => {
  const valid = await loginFormRef.value.validate()
  if (!valid) return

  loginLoading.value = true
  try {
    // 这里可以实现手动登录逻辑
    // 目前系统主要使用OAuth登录，这里作为备用方案
    ElMessage.info('请使用 OAuth 登录')
  } catch (error) {
    ElMessage.error('登录失败，请重试')
    console.error('Manual login error:', error)
  }
  loginLoading.value = false
}

const checkLoginStatus = async () => {
  try {
    const user = await userStore.fetchUserInfo()
    if (user) {
      // 已登录，重定向到首页
      router.push('/')
    }
  } catch (error) {
    // 未登录，停留在登录页面
    console.log('Not logged in')
  }
}

// 生命周期
onMounted(() => {
  checkLoginStatus()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  z-index: -1;
}

.bg-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.1) 2px, transparent 2px),
    radial-gradient(circle at 75% 75%, rgba(255, 255, 255, 0.1) 2px, transparent 2px);
  background-size: 60px 60px;
  animation: patternMove 20s linear infinite;
}

@keyframes patternMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(60px, 60px); }
}

.floating-elements {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.element {
  position: absolute;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  animation: float 6s ease-in-out infinite;
}

.element-1 {
  width: 80px;
  height: 80px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.element-2 {
  width: 60px;
  height: 60px;
  top: 60%;
  right: 15%;
  animation-delay: 2s;
}

.element-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
    opacity: 0.5;
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
    opacity: 0.8;
  }
}

.login-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  max-width: 1200px;
  width: 100%;
  padding: 40px;
  align-items: center;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.logo {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 8px 30px rgba(102, 126, 234, 0.3);
}

.logo-section h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #2d3748;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo-section p {
  margin: 0;
  color: #718096;
  font-size: 0.9rem;
  font-weight: 500;
}

.welcome-text {
  text-align: center;
  margin-bottom: 32px;
}

.welcome-text h2 {
  margin: 0 0 8px 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: #2d3748;
}

.welcome-text p {
  margin: 0;
  color: #718096;
  font-size: 0.95rem;
}

.login-methods {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.oauth-button {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.oauth-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.divider {
  position: relative;
  text-align: center;
  color: #a0aec0;
  font-size: 0.9rem;
}

.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background: #e2e8f0;
  z-index: 0;
}

.divider span {
  background: rgba(255, 255, 255, 0.95);
  padding: 0 16px;
  position: relative;
  z-index: 1;
}

.el-form-item {
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.login-footer {
  margin-top: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 24px;
  border-top: 1px solid #e2e8f0;
}

.tips {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #718096;
  font-size: 0.85rem;
}

.version {
  color: #a0aec0;
  font-size: 0.8rem;
  font-weight: 500;
}

.features-section {
  color: white;
}

.features-section h3 {
  margin: 0 0 24px 0;
  font-size: 1.5rem;
  font-weight: 600;
  text-align: center;
}

.features-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.feature-item {
  text-align: center;
  padding: 24px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.feature-item:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.15);
}

.feature-item h4 {
  margin: 12px 0 8px 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.feature-item p {
  margin: 0;
  font-size: 0.9rem;
  opacity: 0.9;
  line-height: 1.5;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .login-container {
    grid-template-columns: 1fr;
    gap: 40px;
    max-width: 500px;
  }

  .features-section {
    order: -1;
  }

  .features-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 768px) {
  .login-container {
    padding: 20px;
    gap: 30px;
  }

  .login-card {
    padding: 30px 24px;
    border-radius: 16px;
  }

  .logo {
    width: 60px;
    height: 60px;
  }

  .logo-section h1 {
    font-size: 1.5rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .feature-item {
    padding: 20px;
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 15px;
  }

  .login-card {
    padding: 24px 20px;
  }

  .login-footer {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }
}
</style>