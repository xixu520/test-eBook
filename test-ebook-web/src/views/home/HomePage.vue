<template>
  <div class="home-page">
    <!-- 搜索与筛选区域 -->
    <el-card class="filter-card" shadow="never">
      <!-- 移动端折叠版 -->
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
                  @change="handleCategoryFilter"
                />
              </el-form-item>
              <!-- 动态筛选（全局属性 show_in_home=true 的字段） -->
              <template v-for="f in filterFields" :key="f.ID">
                <el-form-item :label="f.label">
                  <el-select v-if="f.field_type === 'select'" v-model="dynamicFilters[f.ID!]" clearable style="width: 100%">
                    <el-option v-for="opt in (f.options || '').split(',').filter(v => !!v.trim())" :key="opt" :label="opt.trim()" :value="opt.trim()" />
                  </el-select>
                  <el-select
                    v-else-if="f.field_type === 'checkbox'"
                    v-model="dynamicFilters[f.ID!]"
                    multiple collapse-tags collapse-tags-tooltip clearable style="width: 100%"
                  >
                    <el-option v-for="opt in (f.options || '').split(',').filter(v => !!v.trim())" :key="opt" :label="opt.trim()" :value="opt.trim()" />
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

      <!-- 桌面端行内版 -->
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
            @change="handleCategoryFilter"
          />
        </el-form-item>

        <!-- 全局属性动态筛选器（show_in_home = true 的字段） -->
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
              <el-option v-for="opt in (f.options || '').split(',').filter(v => !!v.trim())" :key="opt" :label="opt.trim()" :value="opt.trim()" />
            </el-select>
            <el-select
              v-else-if="f.field_type === 'checkbox'"
              v-model="dynamicFilters[f.ID!]"
              multiple collapse-tags collapse-tags-tooltip clearable placeholder="多选"
              style="width: 140px"
              @change="loadData"
            >
              <el-option v-for="opt in (f.options || '').split(',').filter(v => !!v.trim())" :key="opt" :label="opt.trim()" :value="opt.trim()" />
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
        <el-button type="primary" :icon="Upload" v-if="canUpload" @click="uploadVisible = true">
          {{ isMobile ? '' : '上传文件' }}
        </el-button>
        <template v-if="!isMobile">
          <el-button :icon="Download" :disabled="!selectedIds.length">批量下载</el-button>
        </template>
        <template v-else>
          <el-button-group>
            <el-button :icon="Download" :disabled="!selectedIds.length" />
          </el-button-group>
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

      <!-- 文档名称列：可点击打开详情 -->
      <el-table-column prop="title" label="文档名称" min-width="250" show-overflow-tooltip>
        <template #default="{ row }">
          <el-button link type="primary" style="font-weight: 600; text-align: left; white-space: normal; height: auto;" @click="handleShowDetail(row)">
            {{ row.title }}
          </el-button>
        </template>
      </el-table-column>

      <!-- 全局属性动态展示列（show_in_home = true） -->
      <template v-for="col in displayColumns" :key="col.ID">
        <el-table-column :label="col.label" min-width="150" show-overflow-tooltip v-if="!isMobile">
          <template #default="{ row }">
            <template v-if="col.field_type === 'checkbox'">
              <div class="tag-group">
                <el-tag
                  v-for="tag in (getDynamicFieldValue(row, col.ID!) || '').split(',').filter((v: string) => !!v)"
                  :key="tag"
                  size="small"
                  class="field-tag"
                >{{ tag }}</el-tag>
              </div>
            </template>
            <template v-else>{{ getDynamicFieldValue(row, col.ID!) }}</template>
          </template>
        </el-table-column>
      </template>

      <!-- 分类 -->
      <el-table-column label="分类" width="130" v-if="!isMobile">
        <template #default="{ row }">
          <el-tag size="small" type="info" round effect="plain">{{ row.category?.name || '未分类' }}</el-tag>
        </template>
      </el-table-column>

      <!-- 实施状态 -->
      <el-table-column prop="implementation_status" label="实施状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.implementation_status)" size="small">
            {{ getStatusText(row.implementation_status) }}
          </el-tag>
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column label="操作" :width="isMobile ? 80 : 160" :fixed="isMobile ? false : 'right'" align="center">
        <template #default="{ row }">
          <div v-if="!isMobile" class="action-btns">
            <el-button link type="primary" :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">预览</el-button>
            <el-divider direction="vertical" />
            <el-button link type="primary" :icon="Download" @click="handleDownload(row)">下载</el-button>
          </div>
          <el-dropdown v-else trigger="click">
            <el-button link type="primary" :icon="MoreFilled" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">预览</el-dropdown-item>
                <el-dropdown-item :icon="Download" @click="handleDownload(row)">下载</el-dropdown-item>
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
      <PdfPreview v-if="previewVisible" :url="currentDoc?.url" />
    </el-dialog>

    <!-- 文档上传弹窗 -->
    <UploadDialog
      v-model="uploadVisible"
      :category-tree="categoryTree"
      @success="loadData"
    />

    <!-- 文档详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="文档详情"
      width="700px"
      destroy-on-close
    >
      <div v-if="currentDoc" class="detail-container">
        <!-- 系统固有属性 -->
        <div class="detail-section">
          <h4 class="section-title">系统属性</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="文档名称" :span="2">{{ currentDoc.title }}</el-descriptions-item>
            <el-descriptions-item label="所属分类">{{ currentDoc.category?.name || '未分类' }}</el-descriptions-item>
            <el-descriptions-item label="文件大小">{{ formatFileSize(currentDoc.file_size) }}</el-descriptions-item>
            <el-descriptions-item label="上传时间">{{ formatDate(currentDoc.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="实施状态">
              <el-tag :type="getStatusTagType(currentDoc.implementation_status)" size="small">
                {{ getStatusText(currentDoc.implementation_status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="核验状态">
              <el-tag :type="getVerifyTagType(currentDoc.verify_status)" size="small">
                {{ getVerifyStatusText(currentDoc.verify_status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="处理状态">
              <el-tag :type="getOcrTagType(currentDoc.status)" size="small">
                {{ getOcrStatusText(currentDoc.status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 自定义属性值 -->
        <div v-if="currentDoc.field_values?.length" class="detail-section" style="margin-top: 20px">
          <h4 class="section-title">扩展属性</h4>
          <el-descriptions :column="2" border>
            <el-descriptions-item
              v-for="fv in currentDoc.field_values"
              :key="fv.field_id"
              :label="fv.field?.label || '属性'"
            >
              <template v-if="fv.value && fv.value.includes(',')">
                <el-tag v-for="tag in fv.value.split(',')" :key="tag" size="small" style="margin-right: 4px">{{ tag }}</el-tag>
              </template>
              <template v-else>{{ fv.value || '—' }}</template>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="primary" :disabled="currentDoc?.status !== 1" @click="handlePreview(currentDoc)">
          预览文档
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import {
  Search, RefreshLeft, Upload, Download,
  View, MoreFilled
} from '@element-plus/icons-vue'
import { getDocuments, getDocumentDetail } from '@/api/document'
import { getGlobalForm, type FormField } from '@/api/form'
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
const currentDoc = ref<any>(null)

const filters = reactive({
  keyword: (route.query.keyword as string) || '',
  category_id: (route.query.category_id as string) || ''
})

const dynamicFilters = reactive<Record<number, any>>({})
// 全局属性中 show_in_home=true 的字段，同时作为展示列和筛选项（决策5-C）
const displayColumns = ref<FormField[]>([])
const filterFields = computed(() => displayColumns.value) // 联动：展示列即筛选列

const activeFilters = ref([])
const hasActiveFilters = computed(() => {
  return filters.keyword || filters.category_id || Object.values(dynamicFilters).some(v => !!v)
})

const pagination = reactive({
  page: 1,
  page_size: 10,
  total: 0
})

// 加载全局属性配置
const loadGlobalFields = async () => {
  try {
    const res: any = await getGlobalForm()
    displayColumns.value = (res.fields || []).filter((f: FormField) => f.show_in_home)
  } catch (error) {
    console.error('加载全局属性失败', error)
  }
}

onMounted(() => {
  categoryStore.fetchCategories()
  loadGlobalFields()
  loadData()
})

// 监听路由参数变化（如点击侧边栏分类）
watch(() => route.query, (newQuery) => {
  filters.keyword = (newQuery.keyword as string) || ''
  filters.category_id = (newQuery.category_id as string) || ''
  loadData()
})

const categoryTree = computed(() => categoryStore.categories)

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

// 分类筛选：仅过滤文档，不再影响展示列（全局属性已固定）
const handleCategoryFilter = () => {
  loadData()
}

const getDynamicFieldValue = (row: any, fieldID: number) => {
  if (!row.field_values) return '—'
  const fv = row.field_values.find((v: any) => v.field_id === fieldID)
  return fv ? fv.value : '—'
}

const resetFilters = () => {
  filters.keyword = ''
  filters.category_id = ''
  Object.keys(dynamicFilters).forEach(k => delete dynamicFilters[Number(k)])
  pagination.page = 1
  loadData()
}

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(item => item.id)
}

const handleShowDetail = async (row: any) => {
  try {
    const res: any = await getDocumentDetail(row.id)
    const doc = res.document
    const fv = res.field_values || doc.field_values || []
    doc.field_values = fv.filter((item: any) => item.field?.label)
    currentDoc.value = doc
    detailVisible.value = true
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

const handlePreview = (row: any) => {
  currentDoc.value = {
    ...row,
    url: `/api/v1/documents/${row.id || row.ID}/preview`
  }
  previewVisible.value = true
}

const handleDownload = (row: any) => {
  window.open(`/api/v1/documents/${row.id || row.ID}/download`, '_blank')
}

// ---- 格式化工具函数 ----

const formatFileSize = (bytes: number): string => {
  if (!bytes) return '—'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(2) + ' MB'
}

const formatDate = (dateString: string): string => {
  if (!dateString) return '—'
  try {
    return new Date(dateString).toLocaleString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', hour12: false
    })
  } catch {
    return dateString
  }
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
    margin-bottom: 16px;

    :deep(.el-card__body) {
      padding-bottom: 4px;
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

    .filter-form {
      :deep(.el-form-item) {
        margin-bottom: 12px;
      }
    }
  }

  .action-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 14px;

    &.is-mobile {
      padding: 0 5px;
    }

    .left {
      display: flex;
      gap: 10px;
    }
  }

  .action-btns {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0;
    white-space: nowrap;

    :deep(.el-divider--vertical) {
      margin: 0 4px;
    }
  }

  .tag-group {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;

    .field-tag {
      max-width: 100px;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  // 表格行距优化
  :deep(.el-table) {
    .el-table__cell {
      padding: 10px 0;
    }
    .el-table__header th {
      padding: 12px 0;
      background-color: var(--el-fill-color-light);
      font-weight: 600;
    }
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;

    &.is-mobile {
      justify-content: center;
    }
  }

  // 详情弹窗样式
  .detail-container {
    .detail-section {
      .section-title {
        margin: 0 0 12px 0;
        font-size: 14px;
        font-weight: 600;
        color: var(--el-text-color-regular);
        padding-left: 10px;
        border-left: 3px solid var(--el-color-primary);
      }
    }
  }
}
</style>
