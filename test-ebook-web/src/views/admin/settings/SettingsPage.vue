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
              <el-input v-model="settings.ocr.paddle_config.token" placeholder="请输入 AI Studio 访问令牌" show-password />
            </el-form-item>
            <el-form-item label="云端解析模型">
              <el-select v-model="settings.ocr.paddle_config.model" placeholder="请选择模型" style="width: 100%">
                <el-option label="PaddleOCR-VL-1.5 (推荐)" value="PaddleOCR-VL-1.5" />
                <el-option label="PaddleOCR-VL" value="PaddleOCR-VL" />
                <el-option label="PP-OCRv5" value="PP-OCRv5" />
                <el-option label="PP-StructureV3" value="PP-StructureV3" />
              </el-select>
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
              <el-radio label="aliyun_oss">阿里云 OSS</el-radio>
              <el-radio label="tencent_cos">腾讯云 COS</el-radio>
              <el-radio label="cstcloud">中科院数据胶囊 (S3)</el-radio>
            </el-radio-group>
          </el-form-item>

          <!-- 本地存储配置 -->
          <template v-if="settings.storage.type === 'local'">
            <el-form-item label="本地数据目录">
              <el-input v-model="settings.storage.local_path" placeholder="默认为 uploads/" />
            </el-form-item>
          </template>

          <!-- 阿里云 OSS 配置 -->
          <template v-if="settings.storage.type === 'aliyun_oss'">
            <el-form-item label="Endpoint" required>
              <el-input v-model="settings.storage.aliyun_endpoint" placeholder="oss-cn-beijing.aliyuncs.com" />
            </el-form-item>
            <el-form-item label="AccessKey ID" required>
              <el-input v-model="settings.storage.aliyun_access_key_id" placeholder="AccessKey ID" />
            </el-form-item>
            <el-form-item label="AccessKey Secret" required>
              <el-input v-model="settings.storage.aliyun_access_key_secret" placeholder="AccessKey Secret" show-password />
            </el-form-item>
            <el-form-item label="Bucket 名称" required>
              <el-input v-model="settings.storage.aliyun_bucket" placeholder="Bucket Name" />
            </el-form-item>
          </template>

          <!-- 腾讯云 COS 配置 -->
          <template v-if="settings.storage.type === 'tencent_cos'">
            <el-form-item label="Bucket URL" required>
              <el-input v-model="settings.storage.tencent_bucket_url" placeholder="example-1250000000.cos.ap-beijing.myqcloud.com" />
            </el-form-item>
            <el-form-item label="Secret ID" required>
              <el-input v-model="settings.storage.tencent_secret_id" placeholder="Secret ID" />
            </el-form-item>
            <el-form-item label="Secret Key" required>
              <el-input v-model="settings.storage.tencent_secret_key" placeholder="Secret Key" show-password />
            </el-form-item>
          </template>

          <!-- 中科院数据胶囊 (S3) 配置 -->
          <template v-if="settings.storage.type === 'cstcloud'">
            <el-form-item label="Endpoint" required>
              <el-input v-model="settings.storage.cstcloud_endpoint" placeholder="s3.cstcloud.cn" />
            </el-form-item>
            <el-form-item label="Access Key" required>
              <el-input v-model="settings.storage.cstcloud_access_key" placeholder="Access Key" />
            </el-form-item>
            <el-form-item label="Secret Key" required>
              <el-input v-model="settings.storage.cstcloud_secret_key" placeholder="Secret Key" show-password />
            </el-form-item>
            <el-form-item label="Bucket 名称" required>
              <el-input v-model="settings.storage.cstcloud_bucket" placeholder="Bucket Name" />
            </el-form-item>
          </template>
          <el-form-item label="最大文件限制">
            <el-input v-model="settings.storage.max_size_mb" placeholder="50">
              <template #append>MB</template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="success" plain @click="handleTestStorageConnection" :loading="testingStorage">测试存储连接</el-button>
            <el-button type="primary" @click="handleSave" :loading="saving">确认修改</el-button>
          </el-form-item>

          <el-divider />
          <el-form-item label="孤儿文件清理">
            <div style="display: flex; flex-direction: column; gap: 8px;">
              <div style="color: var(--el-text-color-secondary); font-size: 13px;">
                扫描并清理无主的暂存文件、过期软删除记录，以及云端存储中未被数据库引用的文件。
              </div>
              <el-button type="warning" plain @click="handleOrphanScan" :loading="scanningOrphans" :icon="Search">
                扫描孤儿文件
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-tab-pane>

    </el-tabs>

    <!-- 孤儿扫描结果弹窗 -->
    <el-dialog v-model="orphanResultVisible" title="孤儿文件扫描结果" width="600px">
      <div v-if="orphanResult" class="orphan-result">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="清理暂存文件">{{ orphanResult.staging_cleaned }} 个</el-descriptions-item>
          <el-descriptions-item label="清理过期记录">{{ orphanResult.expired_cleaned }} 个</el-descriptions-item>
        </el-descriptions>
        <div v-if="orphanResult.cloud_orphans && orphanResult.cloud_orphans.length > 0" style="margin-top: 16px;">
          <h4>云端孤儿文件（{{ orphanResult.cloud_orphans.length }} 个）</h4>
          <el-table :data="orphanResult.cloud_orphans.map((f: string) => ({ path: f }))" max-height="300" size="small">
            <el-table-column prop="path" label="文件路径" show-overflow-tooltip />
          </el-table>
          <el-alert type="info" :closable="false" style="margin-top: 8px;">以上文件在云端存储中存在，但数据库中无对应记录。系统已记录审计日志，不会自动删除。</el-alert>
        </div>
        <div v-if="orphanResult.errors && orphanResult.errors.length > 0" style="margin-top: 16px;">
          <el-alert v-for="(err, idx) in orphanResult.errors" :key="idx" :title="err" type="error" :closable="false" style="margin-bottom: 4px;" />
        </div>
      </div>
      <template #footer>
        <el-button @click="orphanResultVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { getSettings, saveSettings, testOcrConnection, testStorageConnection, orphanScan } from '@/api/settings'
import { ElMessage } from 'element-plus'

const activeTab = ref('ocr')
const saving = ref(false)
const testing = ref(false)
const testingStorage = ref(false)
const scanningOrphans = ref(false)
const orphanResultVisible = ref(false)
const orphanResult = ref<any>(null)

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
    local_path: 'uploads',
    max_size_mb: 50,
    aliyun_endpoint: '',
    aliyun_access_key_id: '',
    aliyun_access_key_secret: '',
    aliyun_bucket: '',
    tencent_secret_id: '',
    tencent_secret_key: '',
    tencent_bucket_url: '',
    cstcloud_endpoint: '',
    cstcloud_access_key: '',
    cstcloud_secret_key: '',
    cstcloud_bucket: ''
  }
})

// 深度合并函数（保留 target 中的默认值，仅覆盖 source 中存在的值）
const deepMerge = (target: any, source: any) => {
  if (!source || typeof source !== 'object') return
  for (const key in source) {
    if (typeof source[key] === 'object' && source[key] !== null && !Array.isArray(source[key])) {
      if (!target[key] || typeof target[key] !== 'object') {
        target[key] = {}
      }
      deepMerge(target[key], source[key])
    } else if (source[key] !== undefined && source[key] !== '') {
      target[key] = source[key]
    }
  }
}

const loadData = async () => {
  try {
    const res: any = await getSettings()
    deepMerge(settings, res)
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

const handleTestStorageConnection = async () => {
  testingStorage.value = true
  try {
    await testStorageConnection(settings.storage)
    ElMessage.success('存储连接测试成功')
  } catch (error) {
    // 错误处理在拦截器中
  } finally {
    testingStorage.value = false
  }
}

const handleOrphanScan = async () => {
  scanningOrphans.value = true
  try {
    const res: any = await orphanScan()
    orphanResult.value = res
    orphanResultVisible.value = true
    ElMessage.success('孤儿文件扫描完成')
  } catch (error) {
    // 错误处理在拦截器中
  } finally {
    scanningOrphans.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.settings-page {

  .advanced-options {
    .options-wrapper {
      display: flex;
      flex-direction: column;
    }
  }
}
</style>
