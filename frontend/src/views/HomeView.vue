<template>
  <div class="home-view">
    <div class="container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title text-gradient">
          <el-icon class="title-icon"><Ticket /></el-icon>
          运维工单系统
        </h1>
        <p class="page-subtitle">快速提交您的运维需求，我们将第一时间为您处理</p>
      </div>
      
      <!-- 工单提交表单 -->
      <div class="ticket-form glass-card">
        <div class="form-header gradient-bg">
          <h2 class="form-title">
            <el-icon><DocumentAdd /></el-icon>
            提交新工单
          </h2>
        </div>
        
        <div class="form-body">
          <el-form
            ref="formRef"
            :model="ticketForm"
            :rules="formRules"
            label-position="top"
            @submit.prevent="submitTicket"
          >
            <el-form-item label="工单标题" prop="title" required>
              <el-input
                v-model="ticketForm.title"
                placeholder="请简要描述您遇到的问题"
                size="large"
              />
            </el-form-item>
            
            <el-form-item label="问题描述" prop="content" required>
              <el-input
                v-model="ticketForm.content"
                type="textarea"
                placeholder="请详细描述您遇到的问题或需要的帮助，支持粘贴图片"
                :rows="6"
                @paste="handlePaste"
              />
            </el-form-item>
            
            <el-form-item label="上传图片">
              <div class="upload-section">
                <div class="upload-controls">
                  <el-upload
                    ref="uploadRef"
                    :file-list="fileList"
                    :before-upload="beforeUpload"
                    :on-change="handleFileChange"
                    :on-remove="handleFileRemove"
                    :auto-upload="false"
                    multiple
                    accept="image/*"
                    drag
                    class="upload-area"
                  >
                    <div class="upload-content">
                      <el-icon class="upload-icon"><UploadFilled /></el-icon>
                      <div class="upload-text">
                        <p class="upload-title">点击或拖拽图片到此处上传</p>
                        <p class="upload-subtitle">支持 JPG、PNG、WebP 和 GIF 格式，单张图片最大 5MB</p>
                      </div>
                    </div>
                  </el-upload>
                  
                  <div class="paste-hint">
                    <el-icon><Picture /></el-icon>
                    <span>您也可以直接粘贴图片</span>
                  </div>
                </div>
                
                <!-- 图片预览区域 -->
                <div v-if="previewImages.length > 0" class="preview-section">
                  <h4 class="preview-title">图片预览</h4>
                  <div class="preview-grid">
                    <div 
                      v-for="(image, index) in previewImages" 
                      :key="index"
                      class="preview-item"
                    >
                      <el-image
                        :src="image.url"
                        fit="contain"
                        class="preview-image"
                        :preview-src-list="[image.url]"
                        preview-teleported
                      >
                        <template #error>
                          <div class="image-error">
                            <el-icon><Picture /></el-icon>
                            <span>加载失败</span>
                          </div>
                        </template>
                      </el-image>
                      <div class="preview-actions">
                        <el-button 
                          type="danger" 
                          size="small" 
                          :icon="Delete"
                          circle
                          @click="removePreviewImage(index)"
                        />
                      </div>
                      <div class="preview-info">
                        <span class="file-name">{{ image.name }}</span>
                        <span class="file-size">{{ formatFileSize(image.size) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </el-form-item>
          </el-form>
        </div>
        
        <div class="form-footer">
          <el-button 
            type="primary" 
            size="large"
            @click="submitTicket"
            :loading="submitting"
            :disabled="!canSubmit"
            class="submit-btn"
          >
            <el-icon><DocumentAdd /></el-icon>
            {{ submitting ? '提交中...' : '提交工单' }}
          </el-button>
        </div>
      </div>
      
      <!-- 提交成功对话框 -->
      <el-dialog
        v-model="showSuccessDialog"
        title="工单提交成功"
        width="500px"
        :before-close="handleCloseSuccess"
        class="success-dialog"
      >
        <div class="success-content">
          <el-result
            icon="success"
            title="提交成功！"
            :sub-title="`工单 ID: ${newTicketId}`"
          >
            <template #extra>
              <el-button 
                type="primary" 
                @click="viewMyTickets"
              >
                查看我的工单
              </el-button>
              <el-button @click="createAnother">
                继续提交
              </el-button>
            </template>
          </el-result>
        </div>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import api from '@/api'
import { validateImageFile, formatFileSize } from '@/utils/helpers'
import {
  Ticket,
  DocumentAdd,
  UploadFilled,
  Picture,
  Delete
} from '@element-plus/icons-vue'

const router = useRouter()

// 表单引用
const formRef = ref(null)
const uploadRef = ref(null)

// 表单数据
const ticketForm = ref({
  title: '',
  content: ''
})

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: '请输入工单标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请描述您遇到的问题', trigger: 'blur' }
  ]
}

// 文件和图片处理
const fileList = ref([])
const previewImages = ref([])
const submitting = ref(false)
const showSuccessDialog = ref(false)
const newTicketId = ref('')

// 计算是否可以提交
const canSubmit = computed(() => {
  return ticketForm.value.title.trim() && 
         ticketForm.value.content.trim() && 
         !submitting.value
})

// 处理粘贴事件
const handlePaste = (event) => {
  const items = event.clipboardData?.items
  if (!items) return
  
  for (let item of items) {
    if (item.kind === 'file' && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        addImageFile(file, 'paste')
        event.preventDefault()
      }
    }
  }
}

// 文件上传前检查
const beforeUpload = (file) => {
  const validation = validateImageFile(file)
  if (!validation.valid) {
    ElMessage.error(validation.message)
    return false
  }
  return false // 阻止自动上传
}

// 处理文件变化
const handleFileChange = (file) => {
  const validation = validateImageFile(file.raw)
  if (validation.valid) {
    addImageFile(file.raw, 'upload')
  }
}

// 处理文件移除
const handleFileRemove = (file) => {
  // 从预览列表中移除对应的图片
  const index = previewImages.value.findIndex(img => img.uid === file.uid)
  if (index > -1) {
    removePreviewImage(index)
  }
}

// 添加图片文件
const addImageFile = (file, source) => {
  const validation = validateImageFile(file)
  if (!validation.valid) {
    ElMessage.error(validation.message)
    return
  }
  
  // 创建预览对象
  const reader = new FileReader()
  reader.onload = (e) => {
    const imageData = {
      file,
      url: e.target.result,
      name: file.name,
      size: file.size,
      uid: Date.now() + Math.random(),
      source
    }
    previewImages.value.push(imageData)
    
    // 同步到文件列表
    if (source === 'paste') {
      fileList.value.push({
        name: file.name,
        size: file.size,
        uid: imageData.uid,
        raw: file
      })
    }
  }
  reader.readAsDataURL(file)
}

// 移除预览图片
const removePreviewImage = (index) => {
  const image = previewImages.value[index]
  previewImages.value.splice(index, 1)
  
  // 从文件列表中移除
  const fileIndex = fileList.value.findIndex(file => file.uid === image.uid)
  if (fileIndex > -1) {
    fileList.value.splice(fileIndex, 1)
  }
  
  // 如果是通过上传组件添加的，需要清理上传组件
  if (image.source === 'upload') {
    nextTick(() => {
      uploadRef.value?.clearFiles()
      // 重新添加剩余的文件
      fileList.value.forEach(file => {
        if (file.raw && previewImages.value.some(img => img.uid === file.uid)) {
          uploadRef.value?.handleStart(file.raw)
        }
      })
    })
  }
}

// 提交工单
const submitTicket = async () => {
  try {
    // 验证表单
    await formRef.value.validate()
    
    submitting.value = true
    
    // 提交工单
    const response = await api.ticket.createTicket({
      title: ticketForm.value.title,
      content: ticketForm.value.content,
      screenshots: previewImages.value.map(img => img.file)
    })
    
    newTicketId.value = response.data.id
    showSuccessDialog.value = true
    
    ElMessage.success('工单提交成功！')
    
  } catch (error) {
    console.error('提交工单失败:', error)
    if (error.response?.data?.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('提交失败，请重试')
    }
  } finally {
    submitting.value = false
  }
}

// 查看我的工单
const viewMyTickets = () => {
  showSuccessDialog.value = false
  router.push('/tickets')
}

// 继续提交新工单
const createAnother = () => {
  showSuccessDialog.value = false
  resetForm()
}

// 重置表单
const resetForm = () => {
  ticketForm.value = {
    title: '',
    content: ''
  }
  fileList.value = []
  previewImages.value = []
  formRef.value?.resetFields()
  uploadRef.value?.clearFiles()
}

// 关闭成功对话框
const handleCloseSuccess = () => {
  showSuccessDialog.value = false
  resetForm()
}
</script>

<style scoped>
.home-view {
  min-height: 100vh;
  padding: 40px 20px;
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-title {
  font-size: 48px;
  font-weight: 700;
  margin: 0 0 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.title-icon {
  font-size: 52px;
  background: var(--primary-gradient);
  color: white;
  padding: 12px;
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-lg);
}

.page-subtitle {
  font-size: 18px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
}

.ticket-form {
  overflow: hidden;
}

.form-header {
  padding: 24px 32px;
  color: white;
  text-align: center;
}

.form-title {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.form-body {
  padding: 32px;
}

.upload-section {
  width: 100%;
}

.upload-controls {
  margin-bottom: 24px;
}

.upload-area {
  width: 100%;
}

.upload-content {
  padding: 40px 20px;
  text-align: center;
}

.upload-icon {
  font-size: 48px;
  color: var(--primary-color);
  margin-bottom: 16px;
}

.upload-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.upload-subtitle {
  font-size: 14px;
  color: var(--text-muted);
  margin: 0;
}

.paste-hint {
  margin-top: 16px;
  padding: 12px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: var(--primary-color);
  font-size: 14px;
  font-weight: 500;
}

.preview-section {
  border-top: 1px solid var(--border-color);
  padding-top: 24px;
}

.preview-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 16px;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.preview-item {
  position: relative;
  border: 2px solid var(--border-color);
  border-radius: var(--border-radius);
  overflow: hidden;
  transition: all 0.3s ease;
}

.preview-item:hover {
  border-color: var(--primary-color);
  box-shadow: var(--shadow);
}

.preview-image {
  width: 100%;
  height: 150px;
  display: block;
}

.preview-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.preview-item:hover .preview-actions {
  opacity: 1;
}

.preview-info {
  padding: 12px;
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  word-break: break-all;
}

.file-size {
  font-size: 12px;
  color: var(--text-muted);
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
  gap: 8px;
}

.form-footer {
  padding: 24px 32px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  text-align: center;
}

.submit-btn {
  min-width: 200px;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  background: var(--primary-gradient);
  border: none;
  border-radius: var(--border-radius-xl);
  box-shadow: var(--shadow);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

/* Element Plus 样式自定义 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

:deep(.el-input__wrapper) {
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: var(--shadow);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-textarea__inner) {
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

:deep(.el-textarea__inner:hover) {
  box-shadow: var(--shadow);
}

:deep(.el-textarea__inner:focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-upload-dragger) {
  border: 2px dashed var(--border-color);
  border-radius: var(--border-radius-lg);
  background: var(--bg-secondary);
  transition: all 0.3s ease;
}

:deep(.el-upload-dragger:hover) {
  border-color: var(--primary-color);
  background: rgba(102, 126, 234, 0.05);
}

/* 成功对话框样式 */
.success-dialog .success-content {
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-view {
    padding: 20px 16px;
  }
  
  .page-title {
    font-size: 36px;
    flex-direction: column;
    gap: 12px;
  }
  
  .title-icon {
    font-size: 40px;
    padding: 8px;
  }
  
  .page-subtitle {
    font-size: 16px;
  }
  
  .form-header {
    padding: 20px 24px;
  }
  
  .form-title {
    font-size: 20px;
  }
  
  .form-body {
    padding: 24px 20px;
  }
  
  .form-footer {
    padding: 20px 24px;
  }
  
  .submit-btn {
    width: 100%;
    min-width: auto;
  }
  
  .preview-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
  
  .preview-image {
    height: 120px;
  }
  
  .upload-content {
    padding: 30px 15px;
  }
  
  .upload-icon {
    font-size: 36px;
  }
}

@media (max-width: 480px) {
  .home-view {
    padding: 16px 12px;
  }
  
  .page-title {
    font-size: 28px;
  }
  
  .page-subtitle {
    font-size: 14px;
  }
  
  .form-header {
    padding: 16px 20px;
  }
  
  .form-title {
    font-size: 18px;
  }
  
  .form-body {
    padding: 20px 16px;
  }
  
  .form-footer {
    padding: 16px 20px;
  }
  
  .preview-grid {
    grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
    gap: 8px;
  }
  
  .preview-image {
    height: 100px;
  }
  
  .upload-content {
    padding: 20px 10px;
  }
  
  .upload-title {
    font-size: 14px;
  }
  
  .upload-subtitle {
    font-size: 12px;
  }
}
</style>