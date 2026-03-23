<template>
  <div class="settings-page">
    <el-tabs v-model="activeTab" type="border-card">
      <!-- OCR 设置 -->
      <el-tab-pane label="OCR 引擎配置" name="ocr">
        <el-form :model="settings.ocr" label-width="120px" style="max-width: 600px">
          <el-form-item label="OCR 引擎">
            <el-select v-model="settings.ocr.engine" placeholder="请选择 OCR 引擎" style="width: 100%">
              <el-option label="PaddleOCR (官方 API)" value="paddleocr" />
              <el-option label="Tesseract OCR (开源多语言)" value="tesseract" />
              <el-option label="百度 AI OCR (公有云 API)" value="baidu" />
            </el-select>
          </el-form-item>

          <template v-if="settings.ocr.engine === 'paddleocr'">
            <el-form-item label="Token (Auth)" required>
              <el-input v-model="settings.ocr.paddle_config.token" placeholder="请输入身份 Token (Bearer)" show-password />
            </el-form-item>
            <el-form-item label="云端解析模型">
              <el-input v-model="settings.ocr.paddle_config.model" placeholder="模型名称 (如 PaddleOCR-VL-1.5)" clearable />
            </el-form-item>
            <el-form-item label="增强处理选项" class="advanced-options">
              <div class="options-wrapper">
                <el-checkbox v-model="settings.ocr.paddle_config.use_doc_orientation_classify">应用文档方向分类还原 (useDocOrientationClassify)</el-checkbox>
                <el-checkbox v-model="settings.ocr.paddle_config.use_doc_unwarping">应用文档去畸变修正 (useDocUnwarping)</el-checkbox>
                <el-checkbox v-model="settings.ocr.paddle_config.use_chart_recognition">应用图表解析 (useChartRecognition)</el-checkbox>
              </div>
            </el-form-item>
            <el-form-item>
              <el-button type="success" plain @click="handleTestConnection" :loading="testing">测试连接</el-button>
            </el-form-item>
          </template>

          <template v-if="settings.ocr.engine === 'baidu'">
            <el-form-item label="App ID" required>
              <el-input v-model="settings.ocr.baidu_config.app_id" placeholder="请输入百度 OCR 的 App ID" />
            </el-form-item>
            <el-form-item label="API Key" required>
              <el-input v-model="settings.ocr.baidu_config.api_key" placeholder="请输入 API Key" show-password />
            </el-form-item>
            <el-form-item label="Secret Key" required>
              <el-input v-model="settings.ocr.baidu_config.secret_key" placeholder="请输入 Secret Key" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="success" plain @click="handleTestConnection" :loading="testing">测试连接</el-button>
            </el-form-item>
          </template>

          <el-form-item label="识别语言">
            <el-checkbox-group v-model="settings.ocr.language">
              <el-checkbox label="ch">简体中文</el-checkbox>
              <el-checkbox label="en">英文</el-checkbox>
              <el-checkbox label="jp">日文</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          <el-form-item label="并发线程数">
            <el-input-number v-model="settings.ocr.threads" :min="1" :max="16" />
          </el-form-item>
          <el-form-item label="启用 GPU 加速">
            <el-switch v-model="settings.ocr.use_gpu" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSave" :loading="saving">保存设置</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 存储设置 -->
      <el-tab-pane label="存储与路径" name="storage">
        <el-form :model="settings.storage" label-width="120px" style="max-width: 600px">
          <el-form-item label="存储方式">
            <el-radio-group v-model="settings.storage.type">
              <el-radio label="local">本地存储</el-radio>
              <el-radio label="s3">S3 兼容存储</el-radio>
              <el-radio label="oss">阿里云 OSS</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="本地数据目录" v-if="settings.storage.type === 'local'">
            <el-input v-model="settings.storage.local_path" placeholder="/data/ebook/storage" />
          </el-form-item>
          <el-form-item label="最大文件限制">
            <el-input v-model="settings.storage.max_size" placeholder="50">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSave" :loading="saving">确认修改</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 系统状态 -->
      <el-tab-pane label="系统运行状态" name="status">
        <div class="status-container">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card shadow="hover" class="status-card">
                <template #header>CPU 使用率</template>
                <el-progress type="dashboard" :percentage="systemStatus.cpu" :color="cpuColor" />
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="status-card">
                <template #header>内存占用</template>
                <el-progress type="dashboard" :percentage="systemStatus.memory" :color="memColor" />
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card shadow="hover" class="status-card">
                <template #header>存储空间 (已用)</template>
                <el-progress type="dashboard" :percentage="systemStatus.disk" />
              </el-card>
            </el-col>
          </el-row>
          <div class="uptime-info">
            <el-descriptions title="详细运行指标" :column="2" border>
              <el-descriptions-item label="运行时间">{{ systemStatus.uptime }}</el-descriptions-item>
              <el-descriptions-item label="系统版本">v1.2.4-stable</el-descriptions-item>
              <el-descriptions-item label="API 节点">Node-01 (Beijing)</el-descriptions-item>
              <el-descriptions-item label="数据库状态">
                <el-tag type="success">健康 (0.5ms)</el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { getSettings, saveSettings, getSystemStatus, testOcrConnection } from '@/api/settings'
import { ElMessage } from 'element-plus'

const activeTab = ref('ocr')
const saving = ref(false)
const testing = ref(false)
const timer = ref<any>(null)

const settings = reactive({
  ocr: {
    engine: 'paddleocr',
    language: ['ch', 'en'],
    threads: 4,
    use_gpu: true,
    baidu_config: {
      app_id: '',
      api_key: '',
      secret_key: ''
    },
    paddle_config: {
      token: '',
      model: 'PaddleOCR-VL-1.5',
      use_doc_orientation_classify: false,
      use_doc_unwarping: false,
      use_chart_recognition: false
    }
  },
  storage: {
    type: 'local',
    local_path: '',
    max_size: 50
  }
})

const systemStatus = reactive({
  cpu: 0,
  memory: 0,
  disk: 0,
  uptime: '0d 0h 0m'
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

const loadData = async () => {
  try {
    const res: any = await getSettings()
    Object.assign(settings, res)
    await refreshStatus()
  } catch (error) {
    console.error(error)
  }
}

const refreshStatus = async () => {
  try {
    const res: any = await getSystemStatus()
    Object.assign(systemStatus, res)
  } catch (error) {
    console.error(error)
  }
}

const handleSave = async () => {
  saving.value = true
  try {
    await saveSettings(settings)
    ElMessage.success('设置保存成功')
  } catch (error) {
    console.error(error)
  } finally {
    saving.value = false
  }
}

const handleTestConnection = async () => {
  testing.value = true
  try {
    await testOcrConnection(settings.ocr)
    ElMessage.success('连接测试成功，配置有效')
  } catch (error) {
    // 错误在 request.ts 中自动被拦截弹窗
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  loadData()
  timer.value = setInterval(refreshStatus, 3000)
})

onUnmounted(() => {
  if (timer.value) clearInterval(timer.value)
})
</script>

<style scoped lang="scss">
.settings-page {
  .status-container {
    padding: 10px;
    
    .status-card {
      text-align: center;
      margin-bottom: 20px;
      :deep(.el-card__header) {
        font-weight: bold;
        color: #606266;
      }
    }
    
    .uptime-info {
      margin-top: 10px;
    }
  }

  .advanced-options {
    .options-wrapper {
      display: flex;
      flex-direction: column;
    }
  }
}
</style>
