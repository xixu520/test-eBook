<template>
  <div class="audit-log-page">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>审计日志</span>
          <el-form :inline="true" :model="filters" class="filter-form">
            <el-form-item label="操作类型">
              <el-select v-model="filters.action" placeholder="全部类型" clearable style="width: 150px" @change="loadData">
                <el-option label="登录" value="LOGIN" />
                <el-option label="上传" value="UPLOAD" />
                <el-option label="删除" value="DELETE" />
                <el-option label="编辑" value="EDIT" />
                <el-option label="核验" value="VERIFY" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="loadData">刷新</el-button>
            </el-form-item>
          </el-form>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="logList"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="created_at" label="操作时间" width="180" sortable align="center">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="username" label="操作用户" width="120" align="center" />
        <el-table-column prop="action" label="操作类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="150" align="center" />
        <el-table-column label="操作详情" min-width="300">
          <template #default="{ row }">
            <span class="detail-text">{{ truncateText(row.details, 80) }}</span>
            <el-button 
              v-if="row.details && row.details.length > 80" 
              link 
              type="primary" 
              size="small" 
              @click="showDetails(row.details)" 
              style="margin-left: 8px;"
            >
              查看完整
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 详情查看弹窗 -->
      <el-dialog v-model="detailsVisible" title="操作详情" width="600px" destroy-on-close>
        <el-scrollbar max-height="400px">
          <pre class="details-content">{{ currentDetails }}</pre>
        </el-scrollbar>
      </el-dialog>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          layout="total, prev, pager, next"
          @current-change="loadData"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getAuditLogs } from '@/api/audit'
import { ElMessage } from 'element-plus'

interface AuditLog {
  id: number
  created_at: string
  username: string
  action: string
  ip: string
  details: string
}

const loading = ref(false)
const logList = ref<AuditLog[]>([])
const filters = reactive({
  action: ''
})

const pagination = reactive({
  page: 1,
  page_size: 15,
  total: 0
})

const detailsVisible = ref(false)
const currentDetails = ref('')

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getAuditLogs({
      page: pagination.page,
      page_size: pagination.page_size,
      action: filters.action
    })
    logList.value = res.list
    pagination.total = res.total
  } catch (error) {
    ElMessage.error('获取日志失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}

const truncateText = (text: string, length: number) => {
  if (!text) return '-'
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}

const showDetails = (text: string) => {
  currentDetails.value = text
  detailsVisible.value = true
}

const getActionType = (action: string) => {
  const map: any = {
    UPLOAD: 'primary',
    DELETE: 'danger',
    VERIFY: 'success',
    EDIT: 'warning',
    LOGIN: 'info'
  }
  return map[action] || 'info'
}
</script>

<style scoped lang="scss">
.audit-log-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .filter-form {
      margin-top: 18px;
    }
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
  
  .details-content {
    white-space: pre-wrap;
    word-break: break-all;
    background-color: var(--el-fill-color-light);
    padding: 12px;
    border-radius: 4px;
    font-family: var(--el-font-family-monospace, monospace);
    font-size: 13px;
    line-height: 1.5;
    margin: 0;
  }
  
  .detail-text {
    color: var(--el-text-color-regular);
    font-size: 13px;
    font-family: var(--el-font-family-monospace, monospace);
  }
}
</style>
