<template>
  <div class="category-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>分类管理</span>
          <el-button type="primary" :icon="Plus" @click="handleAdd(0)">新增根分类</el-button>
        </div>
      </template>

      <el-table
        :data="categoryTree"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="order" label="排序值" width="100" align="center" />
        <el-table-column prop="doc_count" label="文档数" width="120" align="center" />
        <el-table-column label="操作" width="250" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleAdd(row.id)">添加子类</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" @click="handleDelete(row)" :disabled="row.children?.length > 0">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 分类表窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增分类' : '编辑分类'"
      width="500px"
    >
      <el-form :model="form" label-width="80px" :rules="rules" ref="formRef">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="排序值" prop="order">
          <el-input-number v-model="form.order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { getCategories, deleteCategory, addCategory, updateCategory } from '@/api/category'
import { ElMessageBox, ElMessage } from 'element-plus'

const categories = ref<any[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogType = ref('add')
const submitting = ref(false)
const formRef = ref()

const form = ref({
  id: undefined,
  name: '',
  parent_id: 0,
  order: 0 // sort_order -> order in backend
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res: any = await getCategories()
    categories.value = res // Backend returns []Category
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
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

const handleAdd = (parentId: number) => {
  dialogType.value = 'add'
  form.value = {
    id: undefined,
    name: '',
    parent_id: parentId,
    order: 0
  }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogType.value = 'edit'
  form.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除分类“${row.name}”吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteCategory(row.id)
      ElMessage.success('删除成功')
      loadData()
    } catch (e) {}
  }).catch(() => {})
}

const submitForm = async () => {
  await formRef.value.validate()
  submitting.value = true
  try {
    if (dialogType.value === 'add') {
      await addCategory(form.value as any)
      ElMessage.success('添加成功')
    } else {
      await updateCategory(form.value.id!, form.value as any)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (err) {
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.category-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
