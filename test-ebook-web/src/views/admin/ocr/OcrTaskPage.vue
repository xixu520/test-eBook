<template>
  <div class="ocr-task-page">
    <el-card shadow="never">
      <template #header>
        <div class="header">
          <span>OCR 任务管理</span>
          <el-button :icon="Refresh" @click="loadData" circle />
        </div>
      </template>

      <el-table :data="taskList" v-loading="loading" border stripe>
        <el-table-column prop="name" label="文件名" min-width="200" show-overflow-tooltip />
        <el-table-column prop="time" label="提交时间" width="180" />
        <el-table-column prop="status" label="状态" width="150" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="200" align="center">
          <template #default="{ row }">
            <el-progress 
              :percentage="row.progress" 
              :status="row.status === 'failed' ? 'exception' : (row.status === 'completed' ? 'success' : '')"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template #default="{ row }">
            <el-button link type="primary" v-if="row.status === 'failed'">重试</el-button>
            <el-button link type="danger" v-if="row.status === 'processing'">取消</el-button>
            <el-button link type="info" @click="viewLogs(row)">日志</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { getOcrTasks } from '@/api/upload'
import { ElMessage } from 'element-plus'

const taskList = ref([])
const loading = ref(false)

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getOcrTasks()
    taskList.value = res
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

const getStatusType = (status: string) => {
  const types: any = { completed: 'success', processing: 'primary', pending: 'info', failed: 'danger' }
  return types[status] || ''
}

const getStatusText = (status: string) => {
  const texts: any = { completed: '已完成', processing: '识别中', pending: '排队中', failed: '失败' }
  return texts[status] || status
}

const viewLogs = (row: any) => {
  ElMessage.info(`查看任务 ${row.id} 的日志`)
}
</script>

<style scoped lang="scss">
.ocr-task-page {
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
