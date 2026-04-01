<template>
  <div class="admin-documents-page">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="header-icon"><Grid /></el-icon>
            <span>文档管理</span>
            <el-tag type="info" size="small" round class="count-badge">
              共 {{ pagination.total }} 份文档
            </el-tag>
          </div>
          
          <div class="header-actions">
            <!-- 分类选择器 -->
            <el-tree-select
              v-model="searchCategoryID"
              :data="categories"
              placeholder="选择分类查看动态列"
              clearable
              check-strictly
              node-key="ID"
              :props="{ label: 'name', value: 'ID' }"
              class="filter-category"
              @change="handleCategoryChange"
            />
            <el-input
              v-model="searchKeyword"
              placeholder="搜索标题"
              :prefix-icon="Search"
              clearable
              class="search-input"
              @keyup.enter="handleSearch"
              @clear="handleSearch"
            />
            <el-button type="primary" :icon="Plus" @click="handleUpload">
              上传文档
            </el-button>
          </div>
        </div>
      </template>

      <!-- 动态筛选栏 -->
      <div v-if="filterFields.length > 0" class="dynamic-filter-bar">
        <div v-for="field in filterFields" :key="field.ID" class="filter-item">
          <span class="filter-label">{{ field.label }}：</span>
          <el-input 
            v-if="field.field_type === 'input' || field.field_type === 'number'"
            v-model="dynamicFilters[field.ID!]"
            size="small"
            placeholder="全文过滤"
            clearable
            @change="handleSearch"
          />
          <el-select 
            v-else-if="field.field_type === 'select'"
            v-model="dynamicFilters[field.ID!]"
            size="small"
            placeholder="全选"
            clearable
            @change="handleSearch"
          >
            <el-option v-for="opt in (field.options || '').split(',')" :key="opt" :label="opt" :value="opt" />
          </el-select>
          <el-date-picker
            v-else-if="field.field_type === 'date'"
            v-model="dynamicFilters[field.ID!]"
            type="date"
            size="small"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
          />
        </div>
        <el-button link :icon="Refresh" @click="resetFilters">重置筛选</el-button>
      </div>

      <!-- 文档列表表格 -->
      <el-table
        v-loading="loading"
        :data="documentList"
        style="width: 100%"
        class="document-table"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        
        <!-- 固有必选列 -->
        <el-table-column prop="title" label="文档名称" min-width="280" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="title-cell">
              <el-icon class="doc-icon"><Document /></el-icon>
              <span>{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>

        <!-- 动态扩展列 -->
        <template v-for="col in displayColumns" :key="col.ID">
          <el-table-column :prop="col.field_key" :label="col.label" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              {{ getDynamicFieldValue(row, col.ID!) }}
            </template>
          </el-table-column>
        </template>

        <!-- 固有元数据列 -->
        <el-table-column label="所属分类" width="140" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tag size="small" type="info" effect="plain" round>
              {{ getCategoryName(row.category_id) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="file_size" label="文件大小" width="120" align="center">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="created_at" label="上传时间" width="180" align="center" sortable>
          <template #default="{ row }">
            <span class="time-text">{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="处理状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getStatusType(row.status)" 
              size="small" 
              effect="light"
              round
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="sync_status" label="同步状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getSyncStatusType(row.sync_status)" 
              size="small" 
              effect="light"
              round
            >
              {{ getSyncStatusText(row.sync_status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="verify_status" label="核验状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag 
              :type="getVerifyTagType(row.verify_status)" 
              size="small" 
              effect="plain"
              round
            >
              {{ getVerifyStatusText(row.verify_status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="240" align="center" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button link type="primary" size="small" :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">
                预览
              </el-button>
              <el-button link type="primary" size="small" :icon="Edit" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button link type="primary" size="small" :icon="Timer" :disabled="row.status !== 1" @click="handleHistory(row)">
                历史
              </el-button>
              
              <el-button v-if="row.status === 2" link type="warning" size="small" :icon="Refresh" @click="handleRetry(row)">
                重试 OCR
              </el-button>

              <el-button v-if="row.sync_status === 'sync_failed'" link type="warning" size="small" :icon="Refresh" @click="handleRetrySync(row)">
                重试同步
              </el-button>
              
              <el-button link type="danger" size="small" :icon="Delete" @click="handleDelete(row)">
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>

        <!-- 空状态 -->
        <template #empty>
          <el-empty description="暂无文档数据" :image-size="120">
            <el-button type="primary" :icon="Plus" @click="handleUpload">上传第一份文档</el-button>
          </el-empty>
        </template>
      </el-table>

      <!-- 编辑属性弹窗 -->
      <el-dialog
        v-model="editFormVisible"
        title="编辑文档属性"
        width="650px"
        destroy-on-close
      >
        <el-form :model="editForm" ref="editFormRef" label-width="100px" label-position="top">
          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item label="文档名称" required>
                <el-input v-model="editForm.title" placeholder="主要标题" />
              </el-form-item>
            </el-col>
            
            <el-col :span="12">
              <el-form-item label="所属分类" required>
                <el-tree-select
                  v-model="editForm.category_id"
                  :data="categories"
                  placeholder="所属分类"
                  clearable
                  check-strictly
                  node-key="ID"
                  :props="{ label: 'name', value: 'ID' }"
                  style="width: 100%"
                  @change="handleEditCategoryChange"
                />
              </el-form-item>
            </el-col>

            <!-- 动态编辑字段 -->
            <template v-for="f in editFields" :key="f.ID">
              <el-col :span="12">
                <el-form-item :label="f.label" :required="f.is_required">
                   <el-date-picker
                    v-if="f.field_type === 'date'"
                    v-model="editForm.dynamic_fields[f.ID!]"
                    type="date"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                  <el-select 
                    v-else-if="f.field_type === 'select'"
                    v-model="editForm.dynamic_fields[f.ID!]"
                    style="width: 100%"
                  >
                    <el-option v-for="opt in (f.options || '').split(',')" :key="opt" :label="opt" :value="opt" />
                  </el-select>
                  <el-input-number 
                    v-else-if="f.field_type === 'number'"
                    v-model="editForm.dynamic_fields[f.ID!]"
                    style="width: 100%"
                  />
                  <el-input v-else v-model="editForm.dynamic_fields[f.ID!]" />
                </el-form-item>
              </el-col>
            </template>
            
            <el-col :span="12">
              <el-form-item label="核验状态">
                <el-select v-model="editForm.verify_status" placeholder="核验状态" clearable style="width: 100%">
                  <el-option label="待核验" value="pending" />
                  <el-option label="核对通过" value="pass" />
                  <el-option label="需复核" value="retry" />
                </el-select>
              </el-form-item>
            </el-col>

            <el-col :span="12">
              <el-form-item label="解析结果(OCR)">
                <el-radio-group v-model="editForm.status" size="small">
                  <el-radio-button :label="0">解析中</el-radio-button>
                  <el-radio-button :label="1">解析完成</el-radio-button>
                  <el-radio-button :label="2">解析失败</el-radio-button>
                </el-radio-group>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="editFormVisible = false">取消</el-button>
            <el-button type="primary" :loading="editFormLoading" @click="submitEdit">
              保存修改
            </el-button>
          </div>
        </template>
      </el-dialog>

      <!-- 分页区域 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          background
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 文档上传弹窗 (统一化) -->
    <UploadDialog
      v-model="uploadVisible"
      :category-tree="categories"
      @success="loadData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, Delete, Timer, View, Refresh, Grid, Document, Edit } from '@element-plus/icons-vue'
import { getDocuments, deleteDocument, updateDocument, type Document as IDocument, getTaskStatus, retryOCR, retrySync } from '@/api/document'
import { getCategories } from '@/api/category'
import { getForms, type Form as IForm, type FormField } from '@/api/form'
import type { Category } from '@/api/category'
import { ElMessage, ElMessageBox } from 'element-plus'
import UploadDialog from '@/components/document/UploadDialog.vue'

const loading = ref(false)
const documentList = ref<IDocument[]>([])
const searchKeyword = ref('')
const searchCategoryID = ref<number | undefined>()
const dynamicFilters = reactive<Record<number, string>>({})
const forms = ref<IForm[]>([])
const displayColumns = ref<FormField[]>([])
const filterFields = ref<FormField[]>([])

const selectedIds = ref<number[]>([])
const categories = ref<Category[]>([])
const categoryMap = ref<Record<number, string>>({})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const uploadVisible = ref(false)

// 编辑相关
const editFormVisible = ref(false)
const editFormLoading = ref(false)
const editFields = ref<FormField[]>([])
const editForm = reactive({
  id: 0,
  title: '',
  category_id: undefined as number | undefined,
  dynamic_fields: {} as Record<number, string>,
  status: 0,
  verify_status: 'pending'
})
const editFormRef = ref<any>(null)

const handleEdit = async (row: IDocument) => {
  editForm.id = row.id
  editForm.title = row.title
  editForm.category_id = row.category_id
  editForm.status = row.status
  editForm.verify_status = row.verify_status || 'pending'
  
  // 处理动态字段
  editForm.dynamic_fields = {}
  if (row.field_values) {
    row.field_values.forEach(fv => {
      editForm.dynamic_fields[fv.field_id] = fv.value
    })
  }
  
  await updateEditFields(row.category_id)
  editFormVisible.value = true
}

const updateEditFields = async (catID: number) => {
  if (forms.value.length === 0) {
    const res = await getForms()
    forms.value = res as any
  }
  const cat = findCategory(categories.value, catID)
  if (cat && cat.form_id) {
    const f = forms.value.find(form => form.ID === cat.form_id)
    editFields.value = f?.fields || []
  } else {
    editFields.value = []
  }
}

const handleEditCategoryChange = (val: number) => {
  updateEditFields(val)
}

const submitEdit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    editFormLoading.value = true
    try {
      // 构造提交数据
      const payload = {
        ...editForm,
        dynamic_fields: Object.keys(editForm.dynamic_fields).map(key => ({
          field_id: Number(key),
          value: editForm.dynamic_fields[Number(key)]
        }))
      }
      await updateDocument(editForm.id, payload)
      ElMessage.success('文档已更新')
      editFormVisible.value = false
      loadData()
    } catch (e) {
    } finally {
      editFormLoading.value = false
    }
  })
}

onMounted(async () => {
  await loadCategories()
  loadData()
})

// --- 数据处理 ---

const loadCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res as Category[]
    
    // 递归拍平树结构，构建字典
    const flattenCategories = (nodes: Category[]) => {
      nodes.forEach(node => {
        categoryMap.value[node.ID] = node.name
        if (node.children && node.children.length > 0) {
          flattenCategories(node.children)
        }
      })
    }
    flattenCategories(categories.value)
  } catch (error) {
    console.error('获取分类失败', error)
  }
}

const handleCategoryChange = async (catID?: number) => {
  if (forms.value.length === 0) {
    const res = await getForms()
    forms.value = res as any
  }

  if (!catID) {
    displayColumns.value = []
    filterFields.value = []
    handleSearch()
    return
  }

  const cat = findCategory(categories.value, catID)
  if (cat && cat.form_id) {
    const f = forms.value.find(form => form.ID === cat.form_id)
    if (f) {
      displayColumns.value = f.fields.filter(field => field.show_in_list)
      filterFields.value = f.fields.filter(field => field.show_in_filter)
    }
  } else {
    displayColumns.value = []
    filterFields.value = []
  }
  handleSearch()
}

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

const getDynamicFieldValue = (row: IDocument, fieldID: number) => {
  if (!row.field_values) return '-'
  const fv = row.field_values.find(v => v.field_id === fieldID)
  return fv ? fv.value : '-'
}

const resetFilters = () => {
  Object.keys(dynamicFilters).forEach(k => delete dynamicFilters[Number(k)])
  handleSearch()
}

const getCategoryName = (categoryId: number): string => {
  if (categoryId === 0) return '未分类'
  return categoryMap.value[categoryId] || `末知分类 (${categoryId})`
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string): string => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', { 
      year: 'numeric', month: '2-digit', day: '2-digit', 
      hour: '2-digit', minute: '2-digit', second: '2-digit',
      hour12: false
    })
  } catch {
    return dateString
  }
}

const getStatusType = (status: number) => {
  // 0: processing, 1: success, 2: failed
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    case 0: return '处理中'
    case 1: return '已完成'
    case 2: return '失败'
    default: return '未知'
  }
}

const getSyncStatusType = (status: string) => {
  switch (status) {
    case 'synced': return 'success'
    case 'pending_sync': return 'primary'
    case 'syncing': return 'warning'
    case 'sync_failed': return 'danger'
    default: return 'info'
  }
}

const getSyncStatusText = (status: string) => {
  switch (status) {
    case 'synced': return '已同步'
    case 'pending_sync': return '同步中'
    case 'syncing': return '正在同步'
    case 'sync_failed': return '同步失败'
    default: return '已同步'
  }
}

const getVerifyTagType = (status: string) => {
  const types: Record<string, string> = { pass: 'success', pending: 'warning', retry: 'danger' }
  return status ? types[status] || 'info' : 'warning'
}

const getVerifyStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待核验', pass: '核对通过', retry: '需复核' }
  return status ? map[status] || '待核验' : '待核验'
}

// --- 操作方法 ---

const loadData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.size,
      keyword: searchKeyword.value,
      category_id: searchCategoryID.value
    }
    // 动态属性过滤
    Object.keys(dynamicFilters).forEach(fieldID => {
      if (dynamicFilters[Number(fieldID)]) {
        params[`filter[${fieldID}]`] = dynamicFilters[Number(fieldID)]
      }
    })

    const res: any = await getDocuments(params)
    documentList.value = res.list || []
    pagination.total = res.total || 0
  } catch (error) {
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadData()
}

const handleSelectionChange = (selection: IDocument[]) => {
  selectedIds.value = selection.map(i => i.id)
}

const handleUpload = () => {
  uploadVisible.value = true
}

const pollTaskStatus = (taskId: string) => {
  const timer = setInterval(async () => {
    try {
      const res: any = await getTaskStatus(taskId)
      if (res.status === 'completed') {
        clearInterval(timer)
        ElMessage.success('文档解析任务已完成！')
        loadData()
      } else if (res.status === 'failed') {
        clearInterval(timer)
        ElMessage.error(`文档解析任务失败: ${res.error || '未知错误'}`)
        loadData()
      }
    } catch (error) {
      clearInterval(timer) // 出错停止轮询，避免死循环请求
    }
  }, 5000)
}

const handleRetry = async (row: IDocument) => {
  try {
    const res: any = await retryOCR(row.id)
    ElMessage.success('已将文档重新提交至解析队列')
    loadData()
    if (res.task_id) {
      pollTaskStatus(res.task_id)
    }
  } catch (error) {
    // request.ts 已处理错误提示
  }
}

const handleRetrySync = async (row: IDocument) => {
  try {
    await retrySync(row.id)
    ElMessage.success('已重新提交同步任务')
    loadData()
  } catch (error) {
    // request.ts 已处理错误提示
  }
}

const handleHistory = (row: IDocument) => {
  ElMessage.info(`查看标准 ${row.number} 的历史版本记录（敬请期待）`)
}

const handlePreview = (row: IDocument) => {
  ElMessage.info(`正准备预览标准 ${row.number} 的核心内容（敬请期待）`)
}

const handleDelete = (row: IDocument) => {
  ElMessageBox.confirm(`此操作将永久删除文档「${row.title}」，是否继续？`, '危险操作确认', {
    type: 'warning',
    confirmButtonText: '确定删除',
    cancelButtonText: '取消'
  }).then(async () => {
    try {
      await deleteDocument(row.id)
      ElMessage.success('文档删除成功')
      
      // 处理当前页数据删空的情况
      const isCurrentPageEmpty = documentList.value.length === 1 && pagination.page > 1
      if (isCurrentPageEmpty) {
        pagination.page -= 1
      }
      
      loadData()
    } catch (err) {
      // request.ts handles errors uniquely for us 
    }
  }).catch(() => {})
}
</script>

<style scoped lang="scss">
.admin-documents-page {
  .table-card {
    border-radius: 12px;

    :deep(.el-card__header) {
      border-bottom: 1px solid var(--el-border-color-lighter);
      padding: 16px 20px;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      font-size: 16px;

      .header-icon {
        font-size: 20px;
        color: var(--el-color-primary);
      }

      .count-badge {
        font-weight: 400;
        font-size: 12px;
      }
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 16px;

      .search-input {
        width: 280px;
      }
    }
  }

  // ====== 表格核心样式（Premium） ======
  .document-table {
    --el-table-header-bg-color: var(--el-fill-color-light);
    
    // 行 Hover 效果
    :deep(.el-table__row) {
      transition: background-color 0.25s ease;
      
      &:hover > td {
        background-color: var(--el-color-primary-light-9) !important;
      }
    }

    // 标题列带图标样式
    .title-cell {
      display: flex;
      align-items: center;
      gap: 8px;

      .doc-icon {
        color: var(--el-color-primary-light-3);
        font-size: 16px;
        flex-shrink: 0;
      }

      span {
        font-weight: 500;
        color: var(--el-text-color-primary);
      }
    }
    
    .number-text {
      font-family: var(--el-font-family-monospace, monospace);
      font-size: 13px;
    }

    .time-text {
      color: var(--el-text-color-regular);
      font-size: 13px;
    }

    // 操作列按钮悬停过渡
    .action-buttons {
      display: flex;
      justify-content: center;
      gap: 6px;
      
      .el-button {
        padding: 4px 8px;
        transition: all 0.2s;
        
        &:hover:not(.is-disabled) {
          background-color: var(--el-color-primary-light-9);
          border-radius: 4px;
        }
        
        &.el-button--danger:hover:not(.is-disabled) {
          background-color: var(--el-color-danger-light-9);
        }
      }
    }
  }

  .pagination-container {
    margin-top: 24px;
    padding-bottom: 8px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
