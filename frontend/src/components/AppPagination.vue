<template>
  <div class="app-pagination">
    <div class="pagination-info">
      <span class="total-count">
        共 {{ total }} 条记录
      </span>
    </div>
    
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="pageSizes"
      :total="total"
      :background="true"
      layout="sizes, prev, pager, next, jumper"
      class="pagination-controls"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

const props = defineProps({
  total: {
    type: Number,
    default: 0
  },
  page: {
    type: Number,
    default: 1
  },
  limit: {
    type: Number,
    default: 10
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  }
})

const emit = defineEmits(['update:page', 'update:limit', 'change'])

const currentPage = ref(props.page)
const pageSize = ref(props.limit)

// 监听props变化
watch(() => props.page, (newVal) => {
  currentPage.value = newVal
})

watch(() => props.limit, (newVal) => {
  pageSize.value = newVal
})

// 处理页数变化
const handleCurrentChange = (page) => {
  emit('update:page', page)
  emit('change', { page, limit: pageSize.value })
}

// 处理每页数量变化
const handleSizeChange = (size) => {
  emit('update:limit', size)
  // 当改变每页数量时，重置到第一页
  const newPage = 1
  currentPage.value = newPage
  emit('update:page', newPage)
  emit('change', { page: newPage, limit: size })
}

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(props.total / pageSize.value) || 1
})
</script>

<style scoped>
.app-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 24px 0;
  padding: 20px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(12px);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.pagination-info {
  flex-shrink: 0;
}

.total-count {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 500;
}

.pagination-controls {
  flex: 1;
  justify-content: flex-end;
}

/* Element Plus 分页组件自定义样式 */
:deep(.el-pagination) {
  display: flex;
  align-items: center;
  gap: 8px;
}

:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next) {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  min-width: 36px;
  height: 36px;
}

:deep(.el-pagination .btn-prev:hover),
:deep(.el-pagination .btn-next:hover) {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
  transform: translateY(-1px);
  box-shadow: var(--shadow);
}

:deep(.el-pagination .el-pager li) {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  min-width: 36px;
  height: 36px;
  line-height: 34px;
  margin: 0 2px;
}

:deep(.el-pagination .el-pager li:hover) {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
  transform: translateY(-1px);
  box-shadow: var(--shadow);
}

:deep(.el-pagination .el-pager li.is-active) {
  background: var(--primary-gradient);
  color: white;
  border-color: transparent;
  box-shadow: var(--shadow);
}

:deep(.el-pagination .el-pagination__sizes .el-select .el-input__wrapper) {
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

:deep(.el-pagination .el-pagination__jump input) {
  border-radius: var(--border-radius);
  border: 1px solid var(--border-color);
  transition: all 0.3s ease;
  text-align: center;
}

:deep(.el-pagination .el-pagination__jump input:focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .app-pagination {
    flex-direction: column;
    gap: 16px;
    padding: 16px;
    margin: 20px 0;
  }
  
  .pagination-info {
    order: 2;
  }
  
  .pagination-controls {
    order: 1;
    justify-content: center;
  }
  
  :deep(.el-pagination) {
    flex-wrap: wrap;
    justify-content: center;
  }
  
  :deep(.el-pagination .el-pagination__sizes) {
    order: -1;
    margin-bottom: 8px;
  }
}

@media (max-width: 480px) {
  .app-pagination {
    padding: 12px;
    margin: 16px 0;
  }
  
  :deep(.el-pagination) {
    font-size: 14px;
  }
  
  :deep(.el-pagination .btn-prev),
  :deep(.el-pagination .btn-next),
  :deep(.el-pagination .el-pager li) {
    min-width: 32px;
    height: 32px;
    line-height: 30px;
    font-size: 14px;
  }
  
  /* 隐藏跳转器在小屏幕上 */
  :deep(.el-pagination .el-pagination__jump) {
    display: none;
  }
}

/* 优化小屏幕下的布局 */
@media (max-width: 360px) {
  :deep(.el-pagination .el-pagination__total),
  :deep(.el-pagination .el-pagination__jump) {
    display: none;
  }
}
</style>