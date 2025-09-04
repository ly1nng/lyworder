<template>
  <div class="search-form glass-card">
    <el-form 
      :model="searchParams" 
      @submit.prevent="handleSearch"
      label-position="top"
      class="search-form-content"
    >
      <div class="search-grid">
        <el-form-item label="工单标题" class="search-item">
          <el-input
            v-model="searchParams.title"
            placeholder="请输入工单标题"
            clearable
            :prefix-icon="Search"
          />
        </el-form-item>
        
        <el-form-item label="工单状态" class="search-item">
          <el-select
            v-model="searchParams.status"
            placeholder="请选择状态"
            clearable
            style="width: 100%"
          >
            <el-option label="全部状态" value="" />
            <el-option label="进行中" value="open" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        
        <el-form-item 
          v-if="showCreator" 
          label="创建人" 
          class="search-item"
        >
          <el-input
            v-model="searchParams.user_name"
            placeholder="请输入创建人"
            clearable
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item 
          v-if="showOperator" 
          label="处理人" 
          class="search-item"
        >
          <el-input
            v-model="searchParams.operator_name"
            placeholder="请输入处理人"
            clearable
            :prefix-icon="UserFilled"
          />
        </el-form-item>
        
        <el-form-item label="日期范围" class="search-item date-range">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
      </div>
      
      <div class="search-actions">
        <el-button 
          type="primary" 
          @click="handleSearch"
          :loading="loading"
          :icon="Search"
          class="search-btn"
        >
          搜索
        </el-button>
        <el-button 
          @click="handleReset"
          :icon="RefreshLeft"
          class="reset-btn"
        >
          重置
        </el-button>
        <slot name="extra-actions"></slot>
      </div>
    </el-form>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { Search, RefreshLeft, User, UserFilled } from '@element-plus/icons-vue'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  showCreator: {
    type: Boolean,
    default: false
  },
  showOperator: {
    type: Boolean,
    default: false
  },
  defaultValues: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['search', 'reset'])

const searchParams = ref({
  title: '',
  status: '',
  user_name: '',
  operator_name: '',
  ...props.defaultValues
})

const dateRange = ref(null)

// 监听日期范围变化
watch(dateRange, (newValue) => {
  if (newValue && newValue.length === 2) {
    searchParams.value.start_date = newValue[0]
    searchParams.value.end_date = newValue[1]
  } else {
    searchParams.value.start_date = ''
    searchParams.value.end_date = ''
  }
})

// 搜索
const handleSearch = () => {
  // 过滤空值
  const params = Object.fromEntries(
    Object.entries(searchParams.value).filter(([_, value]) => value !== '' && value != null)
  )
  emit('search', params)
}

// 重置
const handleReset = () => {
  searchParams.value = {
    title: '',
    status: '',
    user_name: '',
    operator_name: '',
    start_date: '',
    end_date: '',
    ...props.defaultValues
  }
  dateRange.value = null
  emit('reset')
}

// 暴露方法给父组件
defineExpose({
  handleSearch,
  handleReset,
  searchParams
})
</script>

<style scoped>
.search-form {
  margin-bottom: 24px;
  padding: 24px;
}

.search-form-content {
  margin: 0;
}

.search-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px 24px;
  margin-bottom: 24px;
}

.search-item {
  margin-bottom: 0;
}

.date-range {
  grid-column: span 1;
}

.search-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.search-btn {
  background: var(--primary-gradient);
  border: none;
  padding: 10px 24px;
  border-radius: var(--border-radius-xl);
  font-weight: 500;
  box-shadow: var(--shadow);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.reset-btn {
  border-radius: var(--border-radius-xl);
  padding: 10px 20px;
  transition: all 0.3s ease;
}

/* Element Plus 自定义样式 */
:deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

:deep(.el-input__wrapper) {
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  box-shadow: var(--shadow-sm);
}

:deep(.el-input__wrapper:hover) {
  box-shadow: var(--shadow);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

:deep(.el-select .el-input__wrapper) {
  border-radius: var(--border-radius);
}

:deep(.el-date-editor .el-input__wrapper) {
  border-radius: var(--border-radius);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .search-form {
    padding: 20px 16px;
    margin-bottom: 20px;
  }
  
  .search-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .date-range {
    grid-column: span 1;
  }
  
  .search-actions {
    justify-content: center;
  }
  
  .search-btn,
  .reset-btn {
    flex: 1;
    min-width: 120px;
  }
}

@media (max-width: 480px) {
  .search-form {
    padding: 16px 12px;
    margin-bottom: 16px;
  }
  
  .search-grid {
    gap: 12px;
  }
  
  .search-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .search-btn,
  .reset-btn {
    width: 100%;
    margin: 0;
  }
}

/* 高级搜索样式优化 */
@media (min-width: 1200px) {
  .search-grid {
    grid-template-columns: repeat(4, 1fr);
  }
  
  .date-range {
    grid-column: span 1;
  }
}

@media (min-width: 992px) and (max-width: 1199px) {
  .search-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (min-width: 768px) and (max-width: 991px) {
  .search-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .date-range {
    grid-column: span 2;
  }
}
</style>