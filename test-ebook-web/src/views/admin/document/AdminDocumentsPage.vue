<template>
  <div class="admin-documents-page">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>文档列表库</span>
          <div class="header-actions">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索标准号或名称"
              :prefix-icon="Search"
              clearable
              style="width: 250px; margin-right: 16px"
              @keyup.enter="loadData"
            />
            <el-button type="primary" :icon="Plus" @click="handleUpload">上传文档</el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="documentList"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column prop="number" label="标准号" width="180" sortable />
        <el-table-column prop="title" label="名称" min-width="250" show-overflow-tooltip />
        <el-table-column prop="category_id" label="分类ID" width="100" />
        <el-table-column prop="file_size" label="大小" width="100" align="center">
          <template #default="{ row }">
            {{ (row.file_size / 1024 / 1024).toFixed(2) }} MB
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="180" sortable align="center">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" :icon="View" :disabled="row.status !== 1" @click="handlePreview(row)">预览</el-button>
          <el-button link type="primary" :icon="Timer" :disabled="row.status !== 1" @click="handleHistory(row)">历史</el-button>
          <el-button link type="primary" :icon="Edit" :disabled="row.status !== 1" @click="handleEdit(row)">编辑</el-button>
          <el-button link type="danger" :icon="Delete" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑文档信息" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="标准号">
          <el-input v-model="editForm.standard_no" />
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="editForm.name" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="editForm.category_id" style="width: 100%">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, Edit, Delete, Timer, View } from '@element-plus/icons-vue'
import { getDocuments, uploadFile, deleteDocument, type Document } from '@/api/document'
import { getCategories } from '@/api/category'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const documentList = ref<Document[]>([])
const searchKeyword = ref('')
const selectedIds = ref<number[]>([])
const categories = ref<any[]>([])

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const editVisible = ref(false)
const editForm = reactive({
  id: 0,
  standard_no: '',
  name: '',
  category_id: 0
})

onMounted(() => {
  loadCategories()
  loadData()
})

const loadCategories = async () => {
  try {
    const res: any = await getCategories()
    categories.value = res
  } catch (error) {}
}

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getDocuments({
      page: pagination.page,
      size: pagination.size,
      keyword: searchKeyword.value
    })
    documentList.value = res.list
    pagination.total = res.total
  } catch (error) {
    // request.ts handles message
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(i => i.id)
}

const handleUpload = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.pdf,.doc,.docx'
  input.onchange = async (e: any) => {
    const file = e.target.files[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)
    formData.append('title', file.name.split('.')[0])
    formData.append('number', 'NEW-' + Date.now()) // Temporary placeholder
    formData.append('category_id', '1') // Default category

    loading.value = true
    try {
      await uploadFile(formData)
      ElMessage.success('上传成功，OCR 正在排队处理...')
      loadData()
    } catch (err) {
      console.error(err)
    } finally {
      loading.value = false
    }
  }
  input.click()
}

const handleEdit = (row: any) => {
  Object.assign(editForm, {
    id: row.id,
    standard_no: row.standard_no,
    name: row.name,
    category_id: row.category_id
  })
  editVisible.value = true
}

const submitEdit = () => {
  ElMessage.success('修改成功（模拟）')
  editVisible.value = false
  loadData()
}

const handleHistory = (row: any) => {
  ElMessage.info(`查看标准 ${row.standard_no} 的历史记录`)
}

const handlePreview = (row: any) => {
  ElMessage.info(`预览标准 ${row.standard_no} 的文件`)
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除文件 ${row.number} 吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteDocument(row.id)
      ElMessage.success('删除成功')
      loadData()
    } catch (err) {}
  })
}

const getStatusType = (status: number) => {
  const map: any = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: any = { 0: '处理中', 1: '已完成', 2: '失败' }
  return map[status] || '等待中'
}
</script>

<style scoped lang="scss">
.admin-documents-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
  }
  
  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
