<template>
  <div class="pdf-preview-container" v-loading="loading">
    <div class="preview-toolbar">
      <div class="page-nav">
        <el-button-group>
          <el-button :icon="ArrowLeft" @click="prevPage" :disabled="currentPage <= 1" size="small" />
          <el-button size="small" disabled>{{ currentPage }} / {{ totalPages }}</el-button>
          <el-button :icon="ArrowRight" @click="nextPage" :disabled="currentPage >= totalPages" size="small" />
        </el-button-group>
      </div>
      <div class="zoom-controls">
        <el-button :icon="Minus" @click="zoomOut" size="small" />
        <span class="zoom-text">{{ Math.round(scale * 100) }}%</span>
        <el-button :icon="Plus" @click="zoomIn" size="small" />
      </div>
      <div class="version-switch" v-if="standardNo">
        <el-dropdown trigger="click" @command="handleVersionChange">
          <el-button size="small" :icon="Timer">
            版本: {{ currentVersion }} <el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item 
                v-for="ver in versions" 
                :key="ver.id" 
                :command="ver"
                :disabled="ver.version === currentVersion"
              >
                {{ ver.version }} ({{ ver.standard_no }})
                <el-tag v-if="ver.is_latest" size="small" type="success" style="margin-left: 8px">最新</el-tag>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <el-button :icon="Download" @click="handleDownload" size="small" circle title="下载文件" />
    </div>

    <div class="canvas-wrapper">
      <canvas ref="canvasRef"></canvas>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as pdfjsLib from 'pdfjs-dist'
import { ArrowLeft, ArrowRight, Plus, Minus, Download, Timer, ArrowDown } from '@element-plus/icons-vue'
import { getDocumentHistory } from '@/api/document'

// 设置 worker 路径
pdfjsLib.GlobalWorkerOptions.workerSrc = `https://cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjsLib.version}/pdf.worker.min.mjs`

const props = defineProps({
  url: { type: String, required: true },
  standardNo: { type: String, default: '' },
  currentVersion: { type: String, default: '' }
})

const emit = defineEmits(['version-change'])

const loading = ref(true)
const pdfDoc = ref<any>(null)
const currentPage = ref(1)
const totalPages = ref(0)
const scale = ref(1.0)
const canvasRef = ref<HTMLCanvasElement | null>(null)

const versions = ref<any[]>([])
const historyLoading = ref(false)

const renderPage = async (num: number) => {
  if (!pdfDoc.value || !canvasRef.value) return
  
  const page = await pdfDoc.value.getPage(num)
  const viewport = page.getViewport({ scale: scale.value })
  const canvas = canvasRef.value
  const context = canvas.getContext('2d')
  
  canvas.height = viewport.height
  canvas.width = viewport.width

  const renderContext = {
    canvasContext: context,
    viewport: viewport,
  }
  
  await page.render(renderContext).promise
}

const loadPdf = async () => {
  loading.value = true
  try {
    const loadingTask = pdfjsLib.getDocument(props.url)
    pdfDoc.value = await loadingTask.promise
    totalPages.value = pdfDoc.value.numPages
    await renderPage(currentPage.value)
  } catch (error) {
    console.error('Failed to load PDF:', error)
  } finally {
    loading.value = false
  }
}

const prevPage = () => {
  if (currentPage.value <= 1) return
  currentPage.value--
  renderPage(currentPage.value)
}

const nextPage = () => {
  if (currentPage.value >= totalPages.value) return
  currentPage.value++
  renderPage(currentPage.value)
}

const zoomIn = () => {
  scale.value += 0.2
  renderPage(currentPage.value)
}

const zoomOut = () => {
  if (scale.value <= 0.4) return
  scale.value -= 0.2
  renderPage(currentPage.value)
}

const handleDownload = () => {
  const link = document.createElement('a')
  link.href = props.url
  link.download = 'document.pdf'
  link.click()
}

const loadHistory = async () => {
  if (!props.standardNo) return
  historyLoading.value = true
  try {
    const res: any = await getDocumentHistory(props.standardNo)
    versions.value = res
  } catch (error) {
    console.error('Failed to load history:', error)
  } finally {
    historyLoading.value = false
  }
}

const handleVersionChange = (ver: any) => {
  emit('version-change', ver)
}

watch(() => props.url, () => {
  currentPage.value = 1
  loadPdf()
})

watch(() => props.standardNo, () => {
  loadHistory()
}, { immediate: true })

onMounted(() => {
  loadPdf()
})
</script>

<style scoped lang="scss">
.pdf-preview-container {
  display: flex;
  flex-direction: column;
  height: 70vh;
  background-color: #f5f7fa;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;

  .preview-toolbar {
    height: 48px;
    background-color: #fff;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    padding: 0 16px;
    z-index: 10;
    
    .zoom-text {
      font-size: 14px;
      color: #606266;
      display: inline-block;
      width: 50px;
      text-align: center;
    }
    
    .page-nav, .zoom-controls {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }

  .canvas-wrapper {
    flex: 1;
    overflow: auto;
    padding: 20px;
    display: flex;
    justify-content: center;
    
    canvas {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      background-color: #fff;
    }
  }
}
</style>
