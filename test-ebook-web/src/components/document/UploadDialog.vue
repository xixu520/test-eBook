<template>
  <el-dialog
    v-model="visible"
    title="上传标准文件"
    width="600px"
    :before-close="handleClose"
    class="upload-dialog"
  >
    <el-form :model="form" label-width="100px" class="upload-form">
      <el-form-item label="文档标题" required>
        <el-input v-model="form.title" placeholder="请输入文档标题" clearable />
      </el-form-item>

      <!-- 动态属性字段渲染 -->
      <template v-for="field in currentFormFields" :key="field.ID">
        <el-form-item :label="field.label" :required="field.is_required">
          <!-- 文本输入 -->
          <el-input 
            v-if="field.field_type === 'input'" 
            v-model="form.dynamic_fields[field.ID!]" 
            :placeholder="'请输入' + field.label" 
            clearable 
          />
          <!-- 数字输入 -->
          <el-input-number 
            v-else-if="field.field_type === 'number'" 
            v-model="form.dynamic_fields[field.ID!]" 
            style="width: 100%"
          />
          <!-- 日期选择 -->
          <el-date-picker
            v-else-if="field.field_type === 'date'"
            v-model="form.dynamic_fields[field.ID!]"
            type="date"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
          <!-- 下拉选择 -->
          <el-select 
            v-else-if="field.field_type === 'select'" 
            v-model="form.dynamic_fields[field.ID!]" 
            style="width: 100%"
          >
            <el-option 
              v-for="opt in (field.options || '').split(',')" 
              :key="opt" 
              :label="opt" 
              :value="opt" 
            />
          </el-select>
          <!-- 复选框组 -->
          <el-checkbox-group
            v-else-if="field.field_type === 'checkbox'"
            v-model="form.dynamic_fields[field.ID!]"
          >
            <el-checkbox 
              v-for="opt in (field.options || '').split(',')" 
              :key="opt" 
              :label="opt" 
              :value="opt" 
            />
          </el-checkbox-group>
        </el-form-item>
      </template>

      <el-form-item label="所属分类" required>
        <el-tree-select
          v-model="form.category_id"
          :data="categoryTree"
          placeholder="请选择分类"
          check-strictly
          node-key="ID"
          :props="{ label: 'name', value: 'ID' }"
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
import { ref, reactive, computed, watch } from 'vue'
import { UploadFilled } from '@element-plus/icons-vue'
import { uploadFile } from '@/api/upload'
import { getForms, type Form, type FormField } from '@/api/form'
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
  title: '',
  category_id: '' as string | number,
  dynamic_fields: {} as Record<number, string>
})

const currentFormFields = ref<FormField[]>([])
const allForms = ref<Form[]>([])

const fileList = ref<any[]>([])
const uploadingFiles = ref<any[]>([])
const isUploading = ref(false)

const canUpload = computed(() => {
  if (!form.title || !form.category_id || fileList.value.length === 0) return false
  // 检查必填动态字段
  for (const field of currentFormFields.value) {
    if (field.is_required && !form.dynamic_fields[field.ID!]) return false
  }
  return true
})

// 分类变化监听
watch(() => form.category_id, async (newVal) => {
  if (!newVal) {
    currentFormFields.value = []
    return
  }
  
  if (allForms.value.length === 0) {
    const res = await getForms()
    allForms.value = res as any
  }
  
  const category = findCategory(props.categoryTree, Number(newVal))
  if (category && category.form_id) {
    const targetForm = allForms.value.find(f => f.ID === category.form_id)
    if (targetForm) {
      currentFormFields.value = targetForm.fields || []
      // 初始化默认值
      const nextFields: Record<number, any> = {}
      currentFormFields.value.forEach(f => {
        if (f.field_type === 'checkbox') {
          nextFields[f.ID!] = f.default_value ? f.default_value.split(',') : []
        } else {
          nextFields[f.ID!] = f.default_value || ''
        }
      })
      form.dynamic_fields = nextFields
    }
  } else {
    currentFormFields.value = []
  }
})

const findCategory = (tree: any[], id: number): any => {
  for (const node of tree) {
    if (node.ID === id) return node
    if (node.children?.length) {
      const found = findCategory(node.children, id)
      if (found) return found
    }
  }
  return null
}

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
  form.title = ''
  form.category_id = ''
  form.dynamic_fields = {}
  currentFormFields.value = []
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
    formData.append('title', form.title)
    formData.append('category_id', String(form.category_id))
    
    // 动态字段处理 (多选需 join)
    const finalizedFields: Record<number, string> = {}
    Object.keys(form.dynamic_fields).forEach(key => {
      const val = form.dynamic_fields[Number(key)]
      finalizedFields[Number(key)] = Array.isArray(val) ? val.join(',') : String(val)
    })
    formData.append('dynamic_fields', JSON.stringify(finalizedFields))

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
