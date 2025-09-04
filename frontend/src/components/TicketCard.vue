<template>
  <div class="ticket-card glass-card">
    <div class="card-header">
      <div class="ticket-title">
        <h3>{{ ticket.title }}</h3>
        <el-tag 
          :type="getStatusType(ticket.status)" 
          class="status-tag"
          effect="dark"
        >
          {{ getStatusText(ticket.status) }}
        </el-tag>
      </div>
      <div class="ticket-type">
        <el-tag size="small" class="type-tag">
          {{ ticket.ticket_type || '未分类' }}
        </el-tag>
      </div>
    </div>
    
    <div class="card-body">
      <div class="content-section">
        <div class="content-text" ref="contentRef">
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
          <ScreenshotsDisplay 
            :screenshots="ticket.screenshots"
            display-mode="counter"
            empty-text="暂无问题截图"
            @image-click="handleImageClick"
            @view-all="handleViewAll"
          />
        </div>
      </div>
    </div>
    
    <div class="card-footer">
      <div class="actions">
        <slot name="actions" :ticket="ticket">
          <el-button 
            type="info" 
            size="small" 
            @click="viewDetail"
            :icon="InfoFilled"
          >
            详情
          </el-button>
        </slot>
      </div>
      
      <div class="meta-info">
        <div class="meta-row">
          <div class="meta-item">
            <span class="meta-label">创建人:</span>
            <span class="meta-value">{{ ticket.user_name || '未知' }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">处理人:</span>
            <span class="meta-value">{{ ticket.operator_name || '未分配' }}</span>
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
import { formatDate, getStatusText, getStatusType } from '@/utils/helpers'
import ScreenshotsDisplay from '@/components/ScreenshotsDisplay.vue'
import { 
  Picture, 
  Camera, 
  InfoFilled 
} from '@element-plus/icons-vue'

const props = defineProps({
  ticket: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const isExpanded = ref(false)
const contentRef = ref(null)

// 切换内容展开/收起
const toggleContent = () => {
  isExpanded.value = !isExpanded.value
}

// 处理截图点击事件
const handleImageClick = (index, screenshot) => {
  // ScreenshotsDisplay组件已经处理了图片查看逻辑
}

// 处理查看所有截图事件
const handleViewAll = (screenshots) => {
  // ScreenshotsDisplay组件已经处理了查看所有截图的逻辑
}

// 获取缩略图URL（保留兼容性）
const getThumbUrl = (originalUrl) => {
  if (!originalUrl) return ''
  return originalUrl.replace(/\.([^.]+)$/, '.thumb.$1')
}

// 查看详情
const viewDetail = () => {
  router.push(`/tickets/detail/${props.ticket.id}`)
}
</script>

<style scoped>
.ticket-card {
  margin-bottom: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.ticket-card:hover {
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
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.4;
  word-break: break-word;
}

.status-tag {
  border-radius: var(--border-radius);
  font-weight: 500;
  letter-spacing: 0.5px;
}

.ticket-type {
  flex-shrink: 0;
}

.type-tag {
  background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: white;
  border: none;
  border-radius: var(--border-radius);
  font-weight: 500;
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
  
  .meta-row {
    flex-direction: column;
    gap: 8px;
  }
  
  .meta-item {
    flex-direction: column;
    gap: 4px;
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

}
</style>