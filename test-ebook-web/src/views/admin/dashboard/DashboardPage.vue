<template>
  <div class="dashboard-page">
    <el-row :gutter="20">
      <el-col :span="6" v-for="card in statCards" :key="card.title">
        <el-card shadow="hover" class="stat-card" :body-style="{ padding: '20px' }">
          <div class="card-content">
            <div class="info">
              <div class="title">{{ card.title }}</div>
              <div class="value">{{ stats[card.key] ?? 0 }}</div>
            </div>
            <div :class="['icon-wrapper', card.color]">
              <el-icon class="icon"><component :is="card.icon" /></el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card shadow="never" class="recent-card">
          <template #header>
            <div class="card-header">
              <el-icon><Operation /></el-icon>
              <span>最近动态</span>
            </div>
          </template>
          <el-timeline v-if="stats.recent_activities?.length">
            <el-timeline-item
              v-for="activity in stats.recent_activities"
              :key="activity.id"
              :timestamp="activity.time"
              :type="getActivityType(activity.type)"
              hollow
            >
              <span class="activity-content">{{ activity.content }}</span>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else description="暂无动态" :image-size="80" />
        </el-card>
      </el-col>
      
      <el-col :span="8">
        <el-card shadow="never" class="storage-card">
          <template #header>
            <div class="card-header">
              <el-icon><PieChart /></el-icon>
              <span>存储使用 (10GB 总上限)</span>
            </div>
          </template>
          <div class="storage-info">
            <el-progress 
              type="dashboard" 
              :percentage="storagePercentage" 
              :color="colors" 
              :stroke-width="10"
              :width="160"
            />
            <div class="storage-details">
              <div class="usage-text">已使用 {{ formatFileSize(stats.storage_used || 0) }}</div>
              <div class="limit-text">当前配额充足</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card shadow="never" class="system-status-card">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>系统监控指标</span>
            </div>
          </template>
          <div class="system-status-content">
            <el-row :gutter="20" justify="space-around">
              <el-col :span="5">
                <div class="progress-box">
                  <div class="label">CPU 使用率</div>
                  <el-progress type="dashboard" :percentage="Math.round(systemStatus.cpu)" :color="cpuColor" :width="120" />
                </div>
              </el-col>
              <el-col :span="5">
                <div class="progress-box">
                  <div class="label">内存占用</div>
                  <el-progress type="dashboard" :percentage="Math.round(systemStatus.memory)" :color="memColor" :width="120" />
                </div>
              </el-col>
              <el-col :span="5">
                <div class="progress-box">
                  <div class="label">磁盘空间</div>
                  <el-progress type="dashboard" :percentage="Math.round(systemStatus.disk)" :width="120" />
                </div>
              </el-col>
              <el-col :span="9">
                <div class="detail-box">
                  <el-descriptions :column="1" border size="small">
                    <el-descriptions-item label="系统运行时间">{{ systemStatus.uptime }}</el-descriptions-item>
                    <el-descriptions-item label="系统版本">{{ systemStatus.version }}</el-descriptions-item>
                    <el-descriptions-item label="数据库状态">
                      <el-tag type="success" size="small">{{ systemStatus.db_status }}</el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="应用占用内存">{{ systemStatus.app_mem?.toFixed(2) }} MB</el-descriptions-item>
                  </el-descriptions>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, reactive } from 'vue'
import { Files, Upload, Grid, Warning, Operation, PieChart, Monitor } from '@element-plus/icons-vue'
import { getDashboardStats, getSystemStatus } from '@/api/stats'

const stats = ref<any>({
  total_documents: 0,
  total_categories: 0,
  today_uploaded: 0,
  pending_ocr: 0,
  storage_used: 0,
  recent_activities: []
})

const statCards = [
  { title: '文档总数', key: 'total_documents', icon: Files, color: 'blue' },
  { title: '分类总数', key: 'total_categories', icon: Grid, color: 'orange' },
  { title: '今日上传', key: 'today_uploaded', icon: Upload, color: 'green' },
  { title: '处理中 OCR', key: 'pending_ocr', icon: Warning, color: 'red' },
]

const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]

const systemStatus = reactive({
  cpu: 0,
  memory: 0,
  disk: 0,
  uptime: '计算中...',
  version: '-',
  db_status: '-',
  app_mem: 0
})

const cpuColor = (percentage: number) => {
  if (percentage < 30) return '#67C23A'
  if (percentage < 70) return '#E6A23C'
  return '#F56C6C'
}

const memColor = (percentage: number) => {
  if (percentage < 50) return '#409EFF'
  if (percentage < 85) return '#E6A23C'
  return '#F56C6C'
}

let timer: any = null

const updateStatus = async () => {
  try {
    const res: any = await getSystemStatus()
    Object.assign(systemStatus, res)
  } catch (error) {
    // 失败时不报错，由 api 层的 silent: true 保证
    console.warn('System status update failed silently')
  } finally {
    // 无论成功失败，5秒后再次尝试，除非组件已卸载
    if (timer !== null) {
      timer = setTimeout(updateStatus, 5000)
    }
  }
}

onMounted(async () => {
  try {
    const res: any = await getDashboardStats()
    stats.value = res
  } catch (error) {
    console.error('Failed to load dashboard stats:', error)
  }
  
  // 启动系统状态轮询
  timer = 0 // 给个非 null 的初始值表示轮询已启动
  updateStatus()
})

onUnmounted(() => {
  if (timer !== null) {
    clearTimeout(timer)
    timer = null
  }
})

const storagePercentage = computed(() => {
  const used = stats.value.storage_used || 0
  const total = 10 * 1024 * 1024 * 1024 // 10GB in bytes
  return Math.min(Math.round((used / total) * 100), 100)
})

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getActivityType = (type: string) => {
  const types: any = { upload: 'primary', ocr: 'success', verify: 'warning' }
  return types[type] || 'info'
}
</script>

<style scoped lang="scss">
.dashboard-page {
  padding: 4px;

  .stat-card {
    border: none;
    border-radius: 12px;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 12px 24px rgba(0, 0, 0, 0.08);
    }

    .card-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      .info {
        .title { 
          color: var(--el-text-color-secondary); 
          font-size: 14px; 
          margin-bottom: 8px;
          font-weight: 500;
        }
        .value { 
          color: var(--el-text-color-primary); 
          font-size: 28px; 
          font-weight: 700;
          font-family: 'Dosis', sans-serif;
        }
      }
      
      .icon-wrapper {
        width: 48px;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 12px;
        
        .icon {
          font-size: 24px;
        }

        &.blue { color: #409eff; background: rgba(64, 158, 255, 0.1); }
        &.green { color: #67c23a; background: rgba(103, 194, 58, 0.1); }
        &.orange { color: #e6a23c; background: rgba(230, 162, 60, 0.1); }
        &.red { color: #f56c6c; background: rgba(245, 108, 108, 0.1); }
      }
    }
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    
    .el-icon {
      color: var(--el-color-primary);
      font-size: 18px;
    }
  }

  .recent-card, .storage-card, .system-status-card {
    border-radius: 12px;
    height: 400px;
    
    :deep(.el-card__header) {
      padding: 16px 20px;
      border-bottom: 1px solid var(--el-border-color-lighter);
    }
  }

  .system-status-card {
    height: auto;
    min-height: 250px;
    
    .system-status-content {
      padding: 10px 0;
      
      .progress-box {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
        
        .label {
          font-size: 14px;
          color: var(--el-text-color-secondary);
          font-weight: 500;
        }
      }
      
      .detail-box {
        height: 100%;
        display: flex;
        align-items: center;
        
        :deep(.el-descriptions) {
          width: 100%;
        }
      }
    }
  }

  .activity-content {
    font-size: 14px;
    color: var(--el-text-color-regular);
  }
  
  .storage-info {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 20px 0;

    .storage-details {
      margin-top: -20px;
      text-align: center;
      
      .usage-text {
        color: var(--el-text-color-primary);
        font-size: 16px;
        font-weight: 600;
        margin-bottom: 4px;
      }
      
      .limit-text {
        color: var(--el-text-color-placeholder);
        font-size: 12px;
      }
    }
  }
}
</style>
