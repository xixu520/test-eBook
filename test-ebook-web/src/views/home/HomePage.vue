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
                <el-input v-model="filters.keyword" placeholder="搜索标准号或名称" clearable @keyup.enter="loadData" />
              </el-form-item>
              <el-form-item label="发布机构">
                <el-select v-model="filters.publisher" placeholder="全部机构" clearable style="width: 100%">
                  <el-option v-for="p in publishers" :key="p" :label="p" :value="p" />
                </el-select>
              </el-form-item>
              <el-form-item label="实施状态">
                <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 100%">
                  <el-option label="现行" value="current" />
                  <el-option label="废止" value="obsolete" />
                  <el-option label="即将实施" value="upcoming" />
                </el-select>
              </el-form-item>
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
          <el-input v-model="filters.keyword" placeholder="搜索标准号或名称" clearable @keyup.enter="loadData" />
        </el-form-item>
        
        <el-form-item label="分类" v-if="!route.query.category_id">
          <el-tree-select
            v-model="filters.category_id"
            :data="categoryTree"
            placeholder="请选择分类"
            clearable
            check-strictly
            node-key="id"
            :props="{ label: 'name', value: 'id' }"
          />
        </el-form-item>
        
        <el-form-item label="发布机构">
          <el-select v-model="filters.publisher" placeholder="全部机构" clearable style="width: 160px">
            <el-option v-for="p in publishers" :key="p" :label="p" :value="p" />
          </el-select>
        </el-form-item>

        <el-form-item label="实施状态">
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 120px">
            <el-option label="现行" value="current" />
            <el-option label="废止" value="obsolete" />
            <el-option label="即将实施" value="upcoming" />
          </el-select>
        </el-form-item>
        
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
      <div class="right" v-if="!isMobile">
        <el-button-group>
          <el-button :icon="Download">导出 Excel</el-button>
          <el-button :icon="Setting">列配置</el-button>
        </el-button-group>
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
      <el-table-column prop="standard_no" label="标准号" width="180" sortable class-name="mono-font" />
      <el-table-column prop="version" label="版本" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_latest ? 'success' : 'info'" size="small">
            {{ row.version || '未知' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="250" show-overflow-tooltip />
      <el-table-column prop="category_name" label="所属分类" width="150" v-if="!isMobile" />
      <el-table-column prop="issue_date" label="发布日期" width="120" sortable class-name="mono-font" v-if="!isMobile" />
      <el-table-column prop="status" label="实施状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)" size="small">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="ocr_status" label="OCR 状态" width="100" align="center" v-if="!isMobile">
        <template #default="{ row }">
          <el-tag :type="getOcrTagType(row.ocr_status)" size="small">
            {{ getOcrStatusText(row.ocr_status) }}
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
            <el-button link type="primary" :icon="View" :disabled="row.ocr_status !== 'completed'" @click="handlePreview(row)">预览</el-button>
            <el-button link type="primary" :icon="Timer" :disabled="row.ocr_status !== 'completed'" @click="handleShowHistory(row)">历史</el-button>
            <el-button link type="primary" :icon="Download" @click="handleDownload(row)">下载</el-button>
          </template>
          <el-dropdown v-else trigger="click">
            <el-button link type="primary" :icon="MoreFilled" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :icon="View" :disabled="row.ocr_status !== 'completed'" @click="handlePreview(row)">预览</el-dropdown-item>
                <el-dropdown-item :icon="Download" @click="handleDownload(row)">下载</el-dropdown-item>
                <el-dropdown-item :icon="Timer" :disabled="row.ocr_status !== 'completed'" @click="handleShowHistory(row)">历史</el-dropdown-item>
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
        v-model:page-size="pagination.size"
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
      :title="`预览 - ${currentDoc?.name}`"
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
      :title="`历史版本 - ${currentHistoryBase?.name}`"
      width="600px"
    >
      <el-table :data="historyList" v-loading="historyLoading" border stripe>
        <el-table-column prop="version" label="版本" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_latest ? 'success' : 'info'">{{ row.version }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="standard_no" label="标准号" width="180" />
        <el-table-column prop="issue_date" label="发布日期" width="120" />
        <el-table-column label="操作" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handlePreview(row)">预览</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { 
  Search, RefreshLeft, Upload, Download, 
  View, Timer, Setting, MoreFilled
} from '@element-plus/icons-vue'
import { getDocuments, getDocumentHistory } from '@/api/document'
import { getCategories } from '@/api/category'
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
const categories = ref<any[]>([])
const selectedIds = ref<number[]>([])

const previewVisible = ref(false)
const uploadVisible = ref(false)
const historyVisible = ref(false)
const currentDoc = ref<any>(null)
const currentHistoryBase = ref<any>(null)
const historyList = ref<any[]>([])
const historyLoading = ref(false)

const filters = reactive({
  keyword: (route.query.keyword as string) || '',
  category_id: (route.query.category_id as string) || '',
  publisher: '',
  status: '',
  dateRange: []
})

const activeFilters = ref([])
const hasActiveFilters = computed(() => {
  return !!(filters.keyword || filters.publisher || filters.status)
})

const publishers = ['住房和城乡建设部', '国家市场监督管理总局', '中国建筑工业出版社']

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

onMounted(() => {
  loadCategories()
  loadData()
})

// 监听路由参数变化（如点击侧边栏分类）
watch(() => route.query, (newQuery) => {
  filters.keyword = (newQuery.keyword as string) || ''
  filters.category_id = (newQuery.category_id as string) || ''
  loadData()
})

const loadCategories = async () => {
  try {
    const res: any = await getCategories()
    categories.value = res
  } catch (error) {
    console.error(error)
  }
}

const categoryTree = computed(() => {
  const map: any = {}
  const roots: any[] = []
  categories.value.forEach(cat => map[cat.id] = { ...cat, children: [] })
  categories.value.forEach(cat => {
    if (cat.parent_id !== 0 && map[cat.parent_id]) {
      map[cat.parent_id].children.push(map[cat.id])
    } else if (cat.parent_id === 0) {
      roots.push(map[cat.id])
    }
  })
  return roots
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getDocuments({
      ...pagination,
      ...filters,
      category_id: filters.category_id ? Number(filters.category_id) : undefined,
      start_date: filters.dateRange?.[0],
      end_date: filters.dateRange?.[1]
    })
    documentList.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  filters.keyword = ''
  filters.category_id = ''
  filters.publisher = ''
  filters.status = ''
  filters.dateRange = []
  loadData()
}

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(item => item.id)
}

const handlePreview = (row: any) => {
  currentDoc.value = {
    ...row,
    url: row.url || 'https://raw.githubusercontent.com/mozilla/pdf.js/ba2edeae/web/compressed.tracemonkey-pldi-09.pdf'
  }
  previewVisible.value = true
}

const handleShowHistory = async (row: any) => {
  currentHistoryBase.value = row
  historyVisible.value = true
  historyLoading.value = true
  try {
    const res: any = await getDocumentHistory(row.standard_no)
    historyList.value = res
  } catch (error) {
    ElMessage.error('获取历史版本失败')
  } finally {
    historyLoading.value = false
  }
}

const handleDownload = (row: any) => {
  // Mock 下载逻辑
  window.open(`https://example.com/mock-download/${row.id}`, '_blank')
}

const getOcrTagType = (status: string) => {
  const types: any = { completed: 'success', pending: 'warning', failed: 'danger' }
  return types[status] || 'info'
}

const getOcrStatusText = (status: string) => {
  const texts: any = { completed: '已完成', pending: '识别中', failed: '失败' }
  return texts[status] || status
}

const getVerifyTagType = (status: string) => {
  const types: any = { pass: 'success', pending: 'warning', retry: 'danger' }
  return types[status] || 'info'
}

const getVerifyStatusText = (status: string) => {
  const map: any = { pending: '待核验', pass: '核验通过', retry: '需重核' }
  return map[status] || status
}

const getStatusTagType = (status: string) => {
  switch (status) {
    case 'current': return 'success'
    case 'obsolete': return 'danger'
    case 'upcoming': return 'warning'
    default: return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'current': return '现行'
    case 'obsolete': return '废止'
    case 'upcoming': return '即将实施'
    default: return '未知'
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
