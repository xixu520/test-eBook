<template>
  <div class="home-page">
    <!-- 搜索与筛选区域 -->
    <el-card class="filter-card" shadow="never">
      <el-collapse v-if="isMobile" v-model="activeFilters" class="mobile-filter-collapse">
        <el-collapse-item name="filters">
          <template #title>
            <div class="mobile-filter-header">
              <el-icon><Search /></el-icon>
              <span>筛选条件</span>
              <el-tag v-if="hasActiveFilters" size="small" type="primary" round style="margin-left: 10px">已启用</el-tag>
            </div>
          </template>
          <div class="mobile-filter-content">
            <el-form :model="filters" label-position="top">
              <el-form-item label="关键字">
                <el-input v-model="filters.keyword" placeholder="搜索关键词" clearable @keyup.enter="loadData" />
              </el-form-item>
              <el-form-item label="分类">
                <el-tree-select
                  v-model="filters.category_id"
                  :data="categoryTree"
                  placeholder="选择分类"
                  clearable
                  check-strictly
                  node-key="ID"
                  :props="{ label: 'name', value: 'ID' }"
                  @change="handleCategoryChange"
                />
              </el-form-item>
              <!-- 动态字段 (Mobile) -->
              <template v-for="f in filterFields" :key="f.ID">
                <el-form-item :label="f.label">
                  <el-select v-if="f.field_type === 'select'" v-model="dynamicFilters[f.ID!]" clearable style="width: 100%">
                    <el-option v-for="opt in (f.options || '').split(',')" :key="opt" :label="opt" :value="opt" />
                  </el-select>
                  <el-date-picker v-else-if="f.field_type === 'date'" v-model="dynamicFilters[f.ID!]" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
                  <el-input v-else v-model="dynamicFilters[f.ID!]" placeholder="输入过滤值" clearable />
                </el-form-item>
              </template>
              <div class="filter-actions">
                <el-button type="primary" @click="loadData">查询</el-button>
                <el-button @click="resetFilters">重置</el-button>
              </div>
            </el-form>
          </div>
        </el-collapse-item>
      </el-collapse>

      <el-form v-else :inline="true" :model="filters" class="filter-form">
        <el-form-item label="关键字">
          <el-input v-model="filters.keyword" placeholder="搜索关键词" clearable @keyup.enter="loadData" />
        </el-form-item>
        
        <el-form-item label="分类">
          <el-tree-select
            v-model="filters.category_id"
            :data="categoryTree"
            placeholder="请选择分类"
            clearable
            check-strictly
            node-key="ID"
            :props="{ label: 'name', value: 'ID' }"
            @change="handleCategoryChange"
          />
        </el-form-item>
        
        <!-- 动态筛选器 -->
        <template v-for="f in filterFields" :key="f.ID">
          <el-form-item :label="f.label">
            <el-select 
              v-if="f.field_type === 'select'"
              v-model="dynamicFilters[f.ID!]"
              placeholder="请选择"
              clearable
              style="width: 140px"
              @change="loadData"
            >
              <el-option v-for="opt in (f.options || '').split(',')" :key="opt" :label="opt" :value="opt" />
            </el-select>
            <el-date-picker
              v-else-if="f.field_type === 'date'"
              v-model="dynamicFilters[f.ID!]"
              type="date"
              placeholder="选择日期"
              value-format="YYYY-MM-DD"
              style="width: 140px"
              @change="loadData"
            />
            <el-input 
              v-else 
              v-model="dynamicFilters[f.ID!]" 
              placeholder="搜索值" 
              clearable 
              style="width: 140px"
              @keyup.enter="loadData"
            />
          </el-form-item>
        </template>
        
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="loadData">查询</el-button>
          <el-button :icon="RefreshLeft" @click="resetFilters">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作工具栏 -->
    <div class="action-bar" :class="{ 'is-mobile': isMobile }">
      <div class="left">
        <el-button type="primary" :icon="Upload" v-if="canUpload" @click="uploadVisible = true">{{ isMobile ? '' : '上传文件' }}</el-button>
        <el-button-group v-if="isMobile">
          <el-button :icon="Download" :disabled="!selectedIds.length" />
        </el-button-group>
        <template v-else>
          <el-button :icon="Download" :disabled="!selectedIds.length">批量下载</el-button>
        </template>
      </div>
    </div>

    <!-- 数据表格 -->
    <el-table
      v-loading="loading"
      :data="documentList"
      border
      stripe
      style="width: 100%"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" align="center" />
      
      <el-table-column prop="title" label="文档名称" min-width="250" show-overflow-tooltip>
        <template #default="{ row }">
           <el-button link type="primary" style="font-weight:600" @click="handleShowDetail(row)">{{ row.title }}</el-button>
        </template>
      </el-table-column>

      <!-- 动态展示列 -->
      <template v-for="col in displayColumns" :key="col.ID">
        <el-table-column 
          :prop="col.field_key" 
          :label="col.label" 
          min-width="150" 
          v-if="!isMobile"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            {{ getDynamicFieldValue(row, col.ID!) }}
          </template>
        </el-table-column>
      </template>

      <el-table-column label="分类" width="130" v-if="!isMobile">
        <template #default="{ row }">
          <el-tag size="small" type="info" round effect="plain">{{ row.category?.name || '未分类' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="implementation_status" label="实施状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.implementation_status)" size="small">
            {{ getStatusText(row.implementation_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="OCR 状态" width="100" align="center" v-if="!isMobile">
        <template #default="{ row }">
          <el-tag :type="getOcrTagType(row.status)" size="small">
            {{ getOcrStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="verify_status" label="核验状态" width="100" align="center" v-if="!isMobile">
        <template #default="{ row }">
          <el-tag :type="getVerifyTagType(row.verify_status)" effect="plain" size="small">
            {{ getVerifyStatusText(row.verify_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" :width="isMobile ? 80 : 200" :fixed="isMobile ? false : 'right'" align="center">
        <template #default="{ row }">
          <template v-if="!isMobile">
            <el-button link type="primary" :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">预览</el-button>
            <el-button link type="primary" :icon="Timer" :disabled="row.status !== 1" @click="handleShowHistory(row)">历史</el-button>
            <el-button link type="primary" :icon="Download" @click="handleDownload(row)">下载</el-button>
          </template>
          <el-dropdown v-else trigger="click">
            <el-button link type="primary" :icon="MoreFilled" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">预览</el-dropdown-item>
                <el-dropdown-item :icon="Download" @click="handleDownload(row)">下载</el-dropdown-item>
                <el-dropdown-item :icon="Timer" :disabled="row.status !== 1" @click="handleShowHistory(row)">历史</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页器 -->
    <div class="pagination-container" :class="{ 'is-mobile': isMobile }">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.page_size"
        :page-sizes="[10, 20, 50, 100]"
        :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
        :total="pagination.total"
        :small="isMobile"
        @size-change="loadData"
        @current-change="loadData"
      />
    </div>
    <!-- PDF 预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      :title="`预览 - ${currentDoc?.title}`"
      width="80%"
      top="5vh"
      destroy-on-close
    >
      <PdfPreview 
        v-if="previewVisible" 
        :url="currentDoc?.url" 
        :standard-no="currentDoc?.standard_no"
        :current-version="currentDoc?.version"
        @version-change="handlePreview"
      />
    </el-dialog>

    <!-- 文档上传弹窗 -->
    <UploadDialog
      v-model="uploadVisible"
      :category-tree="categoryTree"
      @success="loadData"
    />

    <!-- 历史版本弹窗 -->
    <el-dialog
      v-model="historyVisible"
      :title="`历史版本 - ${currentHistoryBase?.title}`"
      width="600px"
    >
      <el-table :data="historyList" v-loading="historyLoading" border stripe>
        <el-table-column prop="version" label="版本" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_latest ? 'success' : 'info'">{{ row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="number" label="标准号" width="180" />
        <el-table-column prop="implementation_date" label="实施日期" width="120" />
        <el-table-column label="操作" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handlePreview(row)">预览</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 文档详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="文档详情"
      width="500px"
      custom-class="detail-dialog"
    >
      <div v-if="currentDoc" class="detail-container">
        <div class="detail-header">
           <el-icon class="doc-icon"><Document /></el-icon>
           <div class="title">{{ currentDoc.title }}</div>
        </div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="所属分类">
            {{ currentDoc.category?.name || '未分类' }}
          </el-descriptions-item>
          <el-descriptions-item label="文件大小">
            {{ (currentDoc.file_size / 1024 / 1024).toFixed(2) }} MB
          </el-descriptions-item>
          <el-descriptions-item label="上传时间">
            {{ currentDoc.created_at }}
          </el-descriptions-item>
          <!-- 动态分配解析展示 -->
          <el-descriptions-item v-for="fv in (currentDoc.field_values || [])" :key="fv.field_id" :label="fv.field?.label || '属性'">
             {{ fv.value }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button type="primary" @click="handlePreview(currentDoc)">预览文档内容</el-button>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { 
  Search, RefreshLeft, Upload, Download, 
  View, Timer, MoreFilled, Document
} from '@element-plus/icons-vue'
import { getDocuments, getDocumentHistory } from '@/api/document'
import { getForms, type Form as IForm, type FormField } from '@/api/form'
import { useCategoryStore } from '@/stores/category'
import { useAuthStore } from '@/stores/auth'
import { useRoute } from 'vue-router'
import PdfPreview from '@/components/document/PdfPreview.vue'
import UploadDialog from '@/components/document/UploadDialog.vue'
import { ElMessage } from 'element-plus'
import { useResponsive } from '@/composables/useResponsive'

const auth = useAuthStore()
const route = useRoute()
const { isMobile } = useResponsive()

const canUpload = computed(() => auth.user && ['admin', 'editor'].includes(auth.user.role))

const loading = ref(false)
const documentList = ref<any[]>([])
const categoryStore = useCategoryStore()
const selectedIds = ref<number[]>([])

const previewVisible = ref(false)
const detailVisible = ref(false)
const uploadVisible = ref(false)
const historyVisible = ref(false)
const currentDoc = ref<any>(null)
const currentHistoryBase = ref<any>(null)
const historyList = ref<any[]>([])
const historyLoading = ref(false)

const filters = reactive({
  keyword: (route.query.keyword as string) || '',
  category_id: (route.query.category_id as string) || ''
})

const dynamicFilters = reactive<Record<number, string>>({})
const displayColumns = ref<FormField[]>([])
const filterFields = ref<FormField[]>([])
const allForms = ref<IForm[]>([])

const activeFilters = ref([])
const hasActiveFilters = computed(() => {
  return filters.keyword || Object.values(dynamicFilters).some(v => !!v)
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

onMounted(() => {
  categoryStore.fetchCategories()
  loadData()
})

// 监听路由参数变化（如点击侧边栏分类）
watch(() => route.query, (newQuery) => {
  filters.keyword = (newQuery.keyword as string) || ''
  filters.category_id = (newQuery.category_id as string) || ''
  loadData()
})

const categoryTree = computed(() => {
  return categoryStore.categories
})

const loadData = async () => {
  loading.value = true
  try {
    const query: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      keyword: filters.keyword,
      category_id: filters.category_id ? Number(filters.category_id) : undefined
    }
    
    // 动态属性过滤
    Object.keys(dynamicFilters).forEach(fid => {
      if (dynamicFilters[Number(fid)]) {
        query[`filter[${fid}]`] = dynamicFilters[Number(fid)]
      }
    })

    const res: any = await getDocuments(query)
    documentList.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCategoryChange = async (catID: number) => {
  if (allForms.value.length === 0) {
    const res = await getForms()
    allForms.value = res as any
  }
  
  const cat = findCategory(categoryTree.value, catID)
  if (cat && cat.form_id) {
    const f = allForms.value.find(form => form.ID === cat.form_id)
    if (f) {
      displayColumns.value = f.fields.filter(field => field.show_in_list)
      filterFields.value = f.fields.filter(field => field.show_in_filter)
    }
  } else {
    displayColumns.value = []
    filterFields.value = []
  }
  loadData()
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

const getDynamicFieldValue = (row: any, fieldID: number) => {
  if (!row.field_values) return '-'
  const fv = row.field_values.find((v: any) => v.field_id === fieldID)
  return fv ? fv.value : '-'
}

const resetFilters = () => {
  filters.keyword = ''
  filters.category_id = ''
  Object.keys(dynamicFilters).forEach(k => delete dynamicFilters[Number(k)])
  displayColumns.value = []
  filterFields.value = []
  loadData()
}

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(item => item.id)
}

const handleShowDetail = (row: any) => {
  currentDoc.value = row
  detailVisible.value = true
}

const handlePreview = (row: any) => {
  currentDoc.value = {
    ...row,
    url: `/api/v1/documents/${row.ID}/preview`
  }
  previewVisible.value = true
}

const handleShowHistory = async (row: any) => {
  currentHistoryBase.value = row
  historyVisible.value = true
  historyLoading.value = true
  try {
    const res: any = await getDocumentHistory(row.number)
    historyList.value = res
  } catch (error) {
    ElMessage.error('获取历史版本失败')
  } finally {
    historyLoading.value = false
  }
}

const handleDownload = (row: any) => {
  window.open(`/api/v1/documents/${row.ID}/download`, '_blank')
}

const getOcrTagType = (status: number) => {
  const types: Record<number, string> = { 1: 'success', 0: 'warning', 2: 'danger' }
  return types[status] || 'info'
}

const getOcrStatusText = (status: number) => {
  const texts: Record<number, string> = { 1: '已完成', 0: '处理中', 2: '失败' }
  return texts[status] || '未知'
}

const getVerifyTagType = (status: string) => {
  const types: Record<string, string> = { pass: 'success', pending: 'warning', retry: 'danger' }
  return status ? types[status] || 'info' : 'warning'
}

const getVerifyStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待核验', pass: '核对通过', retry: '需复核' }
  return status ? map[status] || '待核验' : '待核验'
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'current': return 'success'
    case 'obsolete': return 'danger'
    case 'upcoming': return 'warning'
    default: return 'success'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'current': return '现行'
    case 'obsolete': return '废止'
    case 'upcoming': return '即将实施'
    default: return '现行'
  }
}
</script>

<style scoped lang="scss">
.home-page {
  padding: 20px;
  
  .filter-card {
    margin-bottom: 20px;
    :deep(.el-card__body) {
      padding-bottom: 2px;
    }
    
    .mobile-filter-header {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
    }
    
    .filter-actions {
      display: flex;
      gap: 10px;
      margin-top: 10px;
      .el-button {
        flex: 1;
      }
    }
  }
  
  .action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    &.is-mobile {
      padding: 0 5px;
    }
    
    .left {
      display: flex;
      gap: 10px;
    }
  }
  
  .mono-font {
    font-family: 'JetBrains Mono', 'Courier New', Courier, monospace;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    
    &.is-mobile {
      justify-content: center;
    }
  }
}
</style>
