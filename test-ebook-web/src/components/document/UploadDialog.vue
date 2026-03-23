<template>
  <el-dialog
    v-model="visible"
    title="上传标准文件"
    width="600px"
    :before-close="handleClose"
    class="upload-dialog"
  >
    <el-form :model="form" label-width="100px" class="upload-form">
      <el-form-item label="标准号" required>
        <el-input 
          v-model="form.standard_no" 
          placeholder="请输入标准号" 
          clearable 
          class="mono-font"
        />
      </el-form-item>

      <el-form-item label="版本/年份" required>
        <el-input 
          v-model="form.version" 
          placeholder="请输入年份" 
          clearable 
          class="mono-font"
        />
      </el-form-item>

      <el-form-item label="所属分类" required>
        <el-tree-select
          v-model="form.category_id"
          :data="categoryTree"
          placeholder="请选择分类"
          check-strictly
          node-key="id"
          :props="{ label: 'name', value: 'id' }"
          style="width: 100%"
        />
      </el-form-item>
      
      <el-form-item label="选择文件">
        <el-upload
          class="upload-area"
          drag
          action="#"
          multiple
          :auto-upload="false"
          :on-change="handleFileChange"
          :file-list="fileList"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或 <em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              仅支持 PDF 格式文件，单文件不超过 50MB
            </div>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>

    <!-- 上传进度列表 -->
    <div v-if="uploadingFiles.length" class="progress-list">
      <div v-for="file in uploadingFiles" :key="file.name" class="progress-item">
        <div class="file-info">
          <span class="file-name text-ellipsis">{{ file.name }}</span>
          <span class="file-status">{{ file.status === 'success' ? '已完成' : file.progress + '%' }}</span>
        </div>
        <el-progress 
          :percentage="file.progress" 
          :status="file.status === 'success' ? 'success' : ''"
          :stroke-width="12"
        />
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="isUploading" @click="startUpload" :disabled="!canUpload">
        开始上传
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadFile } from '@/api/upload'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: Boolean,
  categoryTree: Array as any
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const form = reactive({
  standard_no: '',
  version: '',
  category_id: ''
})

const fileList = ref<any[]>([])
const uploadingFiles = ref<any[]>([])
const isUploading = ref(false)

const canUpload = computed(() => form.category_id && form.standard_no && form.version && fileList.value.length > 0)

const handleFileChange = (_file: any, files: any[]) => {
  fileList.value = files
}

const handleClose = () => {
  if (isUploading.value) {
    ElMessage.warning('正在上传中，请稍后...')
    return
  }
  visible.value = false
  resetForm()
}

const resetForm = () => {
  form.standard_no = ''
  form.version = ''
  form.category_id = ''
  fileList.value = []
  uploadingFiles.value = []
}

const startUpload = async () => {
  isUploading.value = true
  uploadingFiles.value = fileList.value.map(f => ({
    name: f.name,
    progress: 0,
    status: 'uploading'
  }))

  const tasks = fileList.value.map((file, index) => {
    const formData = new FormData()
    formData.append('file', file.raw)
    formData.append('standard_no', form.standard_no)
    formData.append('version', form.version)
    formData.append('category_id', form.category_id)

    return uploadFile(formData, (progressEvent) => {
      const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      uploadingFiles.value[index].progress = percent
    }).then(() => {
      uploadingFiles.value[index].status = 'success'
      uploadingFiles.value[index].progress = 100
    }).catch(() => {
      uploadingFiles.value[index].status = 'exception'
    })
  })

  await Promise.all(tasks)
  
  ElMessage.success('文档已提交，后台正在进行 OCR 处理')
  isUploading.value = false
  setTimeout(() => {
    emit('success')
    handleClose()
  }, 1000)
}
</script>

<style scoped lang="scss">
.upload-dialog {
  .upload-area {
    width: 100%;
  }

  .upload-form {
    margin-right: 20px;
  }

  .mono-font {
    :deep(.el-input__inner) {
      font-family: 'JetBrains Mono', 'Courier New', Courier, monospace;
    }
  }
  
  .progress-list {
    margin-top: 20px;
    padding: 10px;
    background: #f8f9fb;
    border-radius: 4px;
    
    .progress-item {
      margin-bottom: 15px;
      .file-info {
        display: flex;
        justify-content: space-between;
        margin-bottom: 5px;
        font-size: 13px;
        color: #606266;
        .file-name { max-width: 70%; }
      }
    }
  }
}

.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
