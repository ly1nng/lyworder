<template>
  <div class="operator-ticket-card glass-card">
    <div class="card-header">
      <div class="ticket-title">
        <h3>{{ ticket.title }}</h3>
      </div>
      <div class="ticket-tags">
        <div class="status-controls">
          <!-- 可编辑状态 -->
          <el-select 
            v-if="editableStatus"
            :model-value="ticket.status" 
            @change="handleStatusChange"
            size="small"
            class="status-select"
          >
            <el-option label="进行中" value="open" />
            <el-option label="已关闭" value="closed" />
          </el-select>
          
          <!-- 只读状态 -->
          <el-tag
            v-else
            :type="ticket.status === 'closed' ? 'info' : 'success'"
            size="small"
            class="status-tag"
          >
            {{ ticket.status === 'closed' ? '已关闭' : '进行中' }}
          </el-tag>
        </div>
        <div class="ticket-meta">
          <el-tag 
            size="small" 
            class="type-tag"
            :class="{ 'editable-tag': editableType }"
            :type="getStatusType(ticket.status)"
            @click="editableType && handleTypeChange()"
          >
            {{ ticket.ticket_type || '未分类' }}
            <el-icon v-if="editableType" class="edit-icon"><Edit /></el-icon>
          </el-tag>
        </div>
      </div>
    </div>
    
    <div class="card-body">
      <div class="content-section">
        <div class="content-text" ref="contentRef">
          <!-- 移动端优化：简化长文本处理 -->
          <p 
            v-if="!isExpanded && ticket.content.length > 150"
            class="content-preview"
          >
            {{ ticket.content.substring(0, 150) }}...
            <el-button 
              type="primary" 
              link 
              @click="toggleContent"
              class="expand-btn"
            >
              展开
            </el-button>
          </p>
          <p v-else class="content-full">
            {{ ticket.content }}
            <el-button 
              v-if="ticket.content.length > 150"
              type="primary" 
              link 
              @click="toggleContent"
              class="expand-btn"
            >
              收起
            </el-button>
          </p>
        </div>
        
        <!-- 截图区域 -->
        <div class="screenshots-section">
          <!-- 截图展示区域 -->
          <div v-if="ticket.screenshots && ticket.screenshots.length > 0" class="screenshots-display">
            <div class="screenshots-counter" @click="handleFullscreenView">
              <div class="counter-content">
                <el-icon class="counter-icon"><Picture /></el-icon>
                <span class="counter-text">{{ ticket.screenshots.length }} 张问题截图</span>
                <el-icon class="arrow-icon"><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
          
          <!-- 空状态 -->
          <div v-else class="empty-screenshots">
            <div class="empty-state">
              <el-icon class="empty-icon"><Camera /></el-icon>
              <p class="empty-text">暂无问题截图</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="card-footer">
      <div class="actions">
        <el-button 
          type="success" 
          size="small" 
          @click="handleRemark"
          :icon="EditPen"
          class="gradient-button"
        >
          备注
        </el-button>
        <el-button 
          type="warning" 
          size="small" 
          @click="handleSolution"
          :icon="Opportunity"
          class="gradient-button"
        >
          解决方案
        </el-button>
        <el-button 
          v-if="editableOperator"
          type="primary" 
          size="small" 
          @click="handleOperatorChange"
          :icon="User"
          class="gradient-button"
        >
          分配处理人
        </el-button>
        <el-button 
          type="info" 
          size="small" 
          @click="viewDetail"
          :icon="InfoFilled"
          class="gradient-button"
        >
          详情
        </el-button>
      </div>
      
      <div class="meta-info">
        <div class="meta-row">
          <div class="meta-item">
            <span class="meta-label">创建人:</span>
            <span class="meta-value">{{ ticket.user_name || '未知' }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">处理人:</span>
            <span class="meta-value">{{ formatOperatorNames(ticket.operator_name) || '未分配' }}</span>
          </div>
        </div>
        <div class="meta-row">
          <div class="meta-item">
            <span class="meta-label">创建时间:</span>
            <span class="meta-value">{{ formatDate(ticket.created_at) }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">更新时间:</span>
            <span class="meta-value">{{ formatDate(ticket.updated_at || ticket.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate, getStatusType } from '@/utils/helpers'
import { 
  Picture, 
  Camera, 
  InfoFilled, 
  EditPen, 
  Opportunity,
  ArrowRight,
  Edit,
  User
} from '@element-plus/icons-vue'

const props = defineProps({
  ticket: {
    type: Object,
    required: true
  },
  editableType: {
    type: Boolean,
    default: false
  },
  editableStatus: {
    type: Boolean,
    default: true
  },
  editableOperator: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['status-change', 'remark', 'solution', 'view-screenshots', 'update-ticket-type', 'update-ticket-operator'])

const router = useRouter()

const isExpanded = ref(false)

// 切换内容展开/收起
const toggleContent = () => {
  isExpanded.value = !isExpanded.value
}

// 处理全屏查看截图
const handleFullscreenView = () => {
  emit('view-screenshots', props.ticket)
}

// 获取缩略图URL（保留兼容性）
const getThumbUrl = (originalUrl) => {
  if (!originalUrl) return ''
  return originalUrl.replace(/\.([^.]+)$/, '.thumb.$1')
}

// 处理状态变更
const handleStatusChange = async (newStatus) => {
  try {
    if (newStatus === 'closed') {
      // 检查是否有解决方案
      if (!props.ticket.solution) {
        ElMessage.warning('请先填写解决方案再关闭工单')
        return
      }
      
      // 确认关闭
      await ElMessageBox.confirm('确定要关闭这个工单吗？', '确认操作', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
    }
    
    emit('status-change', props.ticket.id, newStatus)
  } catch (error) {
    // 用户取消操作
  }
}

// 处理备注
const handleRemark = () => {
  emit('remark', props.ticket)
}

// 处理解决方案
const handleSolution = () => {
  emit('solution', props.ticket)
}

// 处理工单类型修改
const handleTypeChange = () => {
  if (props.editableType) {
    emit('update-ticket-type', props.ticket)
  }
}

// 处理工单处理人修改
const handleOperatorChange = () => {
  if (props.editableOperator) {
    emit('update-ticket-operator', props.ticket)
  }
}

// 查看详情
const viewDetail = () => {
  router.push(`/tickets/detail?id=${props.ticket.id}`)
}

// 格式化处理人名称显示
const formatOperatorNames = (operatorNames) => {
  if (!operatorNames) return ''
  
  // 如果是字符串且包含逗号，分割成数组
  if (typeof operatorNames === 'string' && operatorNames.includes(',')) {
    return operatorNames.split(',').map(name => name.trim()).join(', ')
  }
  
  return operatorNames
}

</script>

<style scoped>
.operator-ticket-card {
  margin-bottom: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.operator-ticket-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-xl);
}

.card-header {
  padding: 20px 20px 0;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 12px;
}

.ticket-title {
  flex: 1;
  min-width: 0;
}

.ticket-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.4;
  word-break: break-word;
}

.status-controls {
  flex-shrink: 0;
}

.status-select {
  min-width: 100px;
}

.ticket-tags {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.ticket-meta {
  display: flex;
  align-items: center;
}

.type-tag {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  border: none;
  border-radius: 20px;
  font-weight: 500;
  cursor: default;
  position: relative;
  padding: 2px 16px;
  font-size: 13px;
  height: auto;
  line-height: 1.8;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.type-tag:hover {
  transform: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.editable-tag {
  padding-right: 24px;
  cursor: pointer;
}

.editable-tag:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.edit-icon {
  position: absolute;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 12px;
  opacity: 0.8;
}

.card-body {
  padding: 20px;
}

.content-section {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.content-text {
  flex: 1;
  min-width: 0;
}

.content-preview,
.content-full {
  margin: 0;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
}

.expand-btn {
  margin-left: 8px;
  padding: 0;
  font-size: 14px;
  font-weight: 500;
}

.screenshots-section {
  flex-shrink: 0;
  width: 200px;
}

/* 截图展示样式 */
.screenshots-display {
  width: 100%;
}

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
  transform: translateY(-2px);
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
  transition: all 0.3s ease;
}

.screenshots-counter:hover .arrow-icon {
  color: var(--el-color-primary);
  transform: translateX(4px);
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

.card-footer {
  padding: 0 20px 20px;
  border-top: 1px solid var(--border-color);
  padding-top: 16px;
}

.actions {
  margin-bottom: 16px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.actions .el-button {
  border-radius: var(--border-radius-xl);
  font-weight: 500;
  transition: all 0.3s ease;
}

.actions .el-button:hover {
  transform: translateY(-1px);
}

.meta-info {
  font-size: 14px;
}

.meta-row {
  display: flex;
  gap: 24px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.meta-row:last-child {
  margin-bottom: 0;
}

.meta-item {
  display: flex;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.meta-label {
  color: var(--text-muted);
  font-weight: 500;
  flex-shrink: 0;
}

.meta-value {
  color: var(--text-secondary);
  word-break: break-word;
}

/* Element Plus 自定义样式 */
:deep(.status-select .el-input__wrapper) {
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
}

:deep(.status-select .el-input__wrapper:hover) {
  box-shadow: var(--shadow);
}

.status-tag {
  font-size: 13px;
  padding: 2px 12px;
  height: auto;
  line-height: 1.8;
  border-radius: 20px;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-section {
    flex-direction: column;
    gap: 16px;
  }
  
  .screenshots-section {
    width: 100%;
    max-width: none;
  }
  
  .card-header {
    padding: 16px 16px 0;
  }
  
  .card-body {
    padding: 16px;
  }
  
  .card-footer {
    padding: 0 16px 16px;
  }
  
  .ticket-title {
    margin-bottom: 12px;
  }
  
  .ticket-tags {
    align-self: stretch;
    justify-content: flex-end;
  }
  
  .meta-row {
    flex-direction: column;
    gap: 8px;
  }
  
  .meta-item {
    flex-direction: column;
    gap: 4px;
  }
  
  /* 移动端优化：简化卡片悬停效果 */
  .operator-ticket-card {
    transition: none;
  }
  
  .operator-ticket-card:hover {
    transform: none;
    box-shadow: var(--shadow-xl);
  }
  
  .actions .el-button {
    transition: none;
  }
  
  .actions .el-button:hover {
    transform: none;
  }
  
  .type-tag {
    transition: none;
  }
  
  .type-tag:hover {
    transform: none;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .editable-tag:hover {
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
  
  .screenshots-counter {
    transition: none;
  }
  
  .screenshots-counter:hover {
    transform: none;
  }
  
  .arrow-icon {
    transition: none;
  }
  
  .screenshots-counter:hover .arrow-icon {
    transform: none;
  }
}

@media (max-width: 480px) {
  .ticket-title h3 {
    font-size: 16px;
  }
  
  .card-header {
    padding: 12px 12px 0;
  }
  
  .card-body {
    padding: 12px;
  }
  
  .card-footer {
    padding: 0 12px 12px;
  }
  
  .actions {
    justify-content: center;
  }
  
  .actions .el-button {
    flex: 1;
    min-width: 80px;
  }
  
  /* 移动端优化：进一步简化样式 */
  :deep(.status-select .el-input__wrapper) {
    transition: none;
  }
  
  :deep(.status-select .el-input__wrapper:hover) {
    box-shadow: var(--shadow);
  }
}
</style>