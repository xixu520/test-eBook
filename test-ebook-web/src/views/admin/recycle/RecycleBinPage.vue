<template>
  <div class="recycle-bin-page">
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>回收站</span>
          <div class="header-actions">
            <el-button type="danger" plain :icon="Delete" @click="handleEmpty">清空回收站</el-button>
            <el-button type="primary" :disabled="!selectedIds.length" @click="handleBatchRestore">批量还原</el-button>
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
        <el-table-column prop="number" label="标准号" width="180" />
        <el-table-column prop="title" label="名称" min-width="250" show-overflow-tooltip />
        <el-table-column prop="created_at" label="上传时间" width="180" align="center">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" :icon="RefreshLeft" @click="handleRestore(row)">还原</el-button>
            <el-button link type="danger" :icon="Delete" @click="handleDelete(row)">彻底删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Delete, RefreshLeft } from '@element-plus/icons-vue'
import { getRecycleBinDocuments, restoreDocuments, batchDeleteDocuments, type Document } from '@/api/document'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const documentList = ref<Document[]>([])
const selectedIds = ref<number[]>([])

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getRecycleBinDocuments()
    documentList.value = res.list
  } catch (error) {
    ElMessage.error('获取回收站列表失败')
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(i => i.id)
}

const handleRestore = (row: any) => {
  restoreDocuments([row.id]).then(() => {
    ElMessage.success('还原成功')
    loadData()
  })
}

const handleBatchRestore = () => {
  restoreDocuments(selectedIds.value).then(() => {
    ElMessage.success('批量还原成功')
    loadData()
  })
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm('彻底删除后将无法还原，确定吗？', '警告', {
    type: 'error',
    confirmButtonText: '确定删除',
    confirmButtonClass: 'el-button--danger'
  }).then(() => {
    batchDeleteDocuments([row.id]).then(() => {
      ElMessage.success('彻底删除成功')
      loadData()
    })
  })
}

const handleEmpty = () => {
  ElMessageBox.confirm('确定要清空回收站吗？此操作不可逆！', '清空确认', {
    type: 'error',
    confirmButtonText: '清空',
    confirmButtonClass: 'el-button--danger'
  }).then(() => {
    batchDeleteDocuments([], true).then(() => {
      ElMessage.success('回收站已清空')
      loadData()
    })
  })
}
</script>

<style scoped lang="scss">
.recycle-bin-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: bold;
  }
}
</style>
