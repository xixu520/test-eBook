<template>
  <div class="dashboard-page">
    <el-row :gutter="20">
      <el-col :span="6" v-for="card in statCards" :key="card.title">
        <el-card shadow="hover" class="stat-card">
          <div class="card-content">
            <div class="info">
              <div class="title">{{ card.title }}</div>
              <div class="value">{{ stats[card.key] }}</div>
            </div>
            <el-icon :class="['icon', card.color]"><component :is="card.icon" /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card shadow="never" header="最近动态">
          <el-timeline>
            <el-timeline-item
              v-for="activity in stats.recent_activities"
              :key="activity.id"
              :timestamp="activity.time"
              :type="getActivityType(activity.type)"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card shadow="never" header="存储使用">
          <div class="storage-info">
            <el-progress type="dashboard" :percentage="45" :color="colors" />
            <div class="label">已使用 {{ stats.storage_used }} / 10 GB</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Files, Upload, List, Warning } from '@element-plus/icons-vue'
import { getDashboardStats } from '@/api/stats'

const stats = ref<any>({})
const statCards = [
  { title: '文件总数', key: 'total_documents', icon: Files, color: 'blue' },
  { title: '今日上传', key: 'today_uploaded', icon: Upload, color: 'green' },
  { title: '待核验', key: 'pending_verify', icon: List, color: 'orange' },
  { title: '待处理 OCR', key: 'pending_ocr', icon: Warning, color: 'red' },
]

const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]

onMounted(async () => {
  try {
    const res: any = await getDashboardStats()
    stats.value = res
  } catch (error) {
    console.error(error)
  }
})

const getActivityType = (type: string) => {
  const types: any = { upload: 'primary', ocr: 'success', verify: 'warning' }
  return types[type] || 'info'
}
</script>

<style scoped lang="scss">
.dashboard-page {
  .stat-card {
    .card-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .info {
        .title { color: #909399; font-size: 14px; margin-bottom: 8px; }
        .value { color: #303133; font-size: 24px; font-weight: bold; }
      }
      
      .icon {
        font-size: 32px;
        padding: 10px;
        border-radius: 8px;
        &.blue { color: #409eff; background: #ecf5ff; }
        &.green { color: #67c23a; background: #f0f9eb; }
        &.orange { color: #e6a23c; background: #fdf6ec; }
        &.red { color: #f56c6c; background: #fef0f0; }
      }
    }
  }
  
  .storage-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 0;
    .label { margin-top: 10px; color: #606266; font-size: 14px; }
  }
}
</style>
