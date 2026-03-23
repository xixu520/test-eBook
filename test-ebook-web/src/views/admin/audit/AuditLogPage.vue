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
        <el-table-column prop="timestamp" label="操作时间" width="180" sortable align="center" />
        <el-table-column prop="username" label="操作用户" width="120" align="center" />
        <el-table-column prop="action" label="操作类型" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getActionType(row.action)">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP 地址" width="150" align="center" />
        <el-table-column prop="details" label="操作详情" min-width="300" />
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
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
  timestamp: string
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
  size: 15,
  total: 0
})

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getAuditLogs({
      page: pagination.page,
      size: pagination.size,
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
}
</style>
