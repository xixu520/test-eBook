<template>
  <div class="category-page">
    <el-card class="category-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-icon class="header-icon"><Grid /></el-icon>
            <span>分类管理</span>
            <el-tag type="info" size="small" round class="count-badge">
              共 {{ totalCount }} 个分类
            </el-tag>
          </div>
          <el-button type="primary" :icon="Plus" @click="handleAdd">添加分类</el-button>
        </div>
      </template>

      <!-- 树状结构展示 -->
      <div v-loading="loading" class="tree-container">
        <el-tree
          v-if="categoryTree.length > 0"
          ref="treeRef"
          :data="categoryTree"
          :props="treeProps"
          node-key="ID"
          default-expand-all
          :expand-on-click-node="false"
          highlight-current
          :indent="28"
          class="category-tree"
        >
          <template #default="{ node, data }">
            <div class="tree-node">
              <div class="node-left">
                <el-icon class="folder-icon" :size="18">
                  <FolderOpened v-if="node.expanded && data.children?.length" />
                  <Folder v-else />
                </el-icon>
                <span class="node-label">{{ data.name }}</span>
                <el-tag
                  :type="data.doc_count > 0 ? 'primary' : 'info'"
                  size="small"
                  effect="plain"
                  round
                  class="doc-tag"
                >
                  {{ data.doc_count ?? 0 }} 篇文档
                </el-tag>
                <el-tag
                  v-if="node.level > 1"
                  type="warning"
                  size="small"
                  effect="light"
                  round
                  class="level-tag"
                >
                  L{{ node.level }}
                </el-tag>
              </div>
              <div class="node-actions">
                <el-button link type="primary" size="small" :icon="Edit" @click.stop="handleEdit(data)">
                  编辑
                </el-button>
                <el-button link type="danger" size="small" :icon="Delete" @click.stop="handleDelete(data)">
                  删除
                </el-button>
              </div>
            </div>
          </template>
        </el-tree>

        <!-- 空状态 -->
        <el-empty v-else-if="!loading" description="暂无分类数据" :image-size="120">
          <el-button type="primary" :icon="Plus" @click="handleAdd">创建第一个分类</el-button>
        </el-empty>
      </div>
    </el-card>

    <!-- 新增/编辑分类对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增分类' : '编辑分类'"
      width="520px"
      destroy-on-close
      @closed="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="90px"
        label-position="right"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="父级分类" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ label: 'name', value: 'ID', children: 'children' }"
            placeholder="无（作为根分类）"
            clearable
            check-strictly
            :render-after-expand="false"
            default-expand-all
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="排序值" prop="order">
          <el-input-number v-model="form.order" :min="0" :max="9999" controls-position="right" />
          <span class="form-tip">数值越小越靠前</span>
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
import { Plus, Edit, Delete, Folder, FolderOpened, Grid } from '@element-plus/icons-vue'
import { getCategories, addCategory, updateCategory, deleteCategory } from '@/api/category'
import type { Category, CategoryForm } from '@/api/category'
import { ElMessageBox, ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

// --- 状态 ---
const categoryTree = ref<Category[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const treeRef = ref()
const editingId = ref<number | undefined>(undefined)

const treeProps = {
  children: 'children',
  label: 'name'
}

const form = ref<CategoryForm>({
  name: '',
  parent_id: 0,
  order: 0
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 1, max: 50, message: '长度在 1 到 50 个字符', trigger: 'blur' }
  ]
}

// --- 统计 ---
const totalCount = computed(() => {
  const count = (nodes: Category[]): number => {
    return nodes.reduce((sum, n) => sum + 1 + count(n.children || []), 0)
  }
  return count(categoryTree.value)
})

// --- 父级分类选择器数据 ---
const parentOptions = computed(() => {
  const rootOption = { ID: 0, name: '根分类（顶级）', children: [] as any[] }

  if (dialogType.value === 'add') {
    return [{ ...rootOption, children: categoryTree.value }]
  }

  const filterSelfAndDescendants = (nodes: Category[]): any[] => {
    return nodes
      .filter(node => node.ID !== editingId.value)
      .map(node => ({
        ...node,
        children: node.children?.length
          ? filterSelfAndDescendants(node.children)
          : []
      }))
  }

  return [{ ...rootOption, children: filterSelfAndDescendants(categoryTree.value) }]
})

// --- 数据加载 ---
onMounted(() => {
  loadData()
})

const loadData = async () => {
  loading.value = true
  try {
    const res = await getCategories()
    categoryTree.value = res as Category[]
  } catch (error) {
    console.error('加载分类数据失败:', error)
  } finally {
    loading.value = false
  }
}

// --- 操作 ---
const handleAdd = () => {
  dialogType.value = 'add'
  editingId.value = undefined
  form.value = {
    name: '',
    parent_id: 0,
    order: 0
  }
  dialogVisible.value = true
}

const handleEdit = (row: Category) => {
  dialogType.value = 'edit'
  editingId.value = row.ID
  form.value = {
    name: row.name,
    parent_id: row.parent_id,
    order: row.order
  }
  dialogVisible.value = true
}

const handleDelete = (row: Category) => {
  ElMessageBox.confirm(
    `确定要删除分类「${row.name}」吗？`,
    '删除确认',
    {
      type: 'warning',
      confirmButtonText: '确定删除',
      cancelButtonText: '取消'
    }
  ).then(async () => {
    try {
      await deleteCategory(row.ID)
      ElMessage.success('分类删除成功')
      loadData()
    } catch (e: any) {
      // request.ts 拦截器已自动弹出后端错误信息
    }
  }).catch(() => {})
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate()

  submitting.value = true
  try {
    if (dialogType.value === 'add') {
      await addCategory(form.value)
      ElMessage.success('分类添加成功')
    } else {
      await updateCategory(editingId.value!, form.value)
      ElMessage.success('分类更新成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (err: any) {
    // request.ts 拦截器已处理错误提示
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  formRef.value?.resetFields()
  editingId.value = undefined
}
</script>

<style scoped lang="scss">
.category-page {
  .category-card {
    border-radius: 12px;

    :deep(.el-card__header) {
      border-bottom: 1px solid var(--el-border-color-lighter);
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
  }

  .tree-container {
    min-height: 200px;
    padding: 8px 0;
  }

  // ====== 树状结构核心样式 ======
  .category-tree {
    --tree-line-color: #dcdfe6;
    --tree-line-color-hover: var(--el-color-primary-light-5);

    background: transparent;

    // 节点行样式
    :deep(.el-tree-node__content) {
      height: 48px;
      border-radius: 8px;
      margin-bottom: 2px;
      padding-right: 16px;
      transition: all 0.25s ease;

      &:hover {
        background-color: var(--el-color-primary-light-9);
      }
    }

    :deep(.el-tree-node.is-current > .el-tree-node__content) {
      background-color: var(--el-color-primary-light-9);
    }

    // ====== 树形连接线 ======
    :deep(.el-tree-node) {
      position: relative;

      // 子节点容器
      .el-tree-node__children {
        position: relative;
        padding-left: 4px;

        // 竖向连接线（从父到最后一个子节点）
        &::before {
          content: '';
          position: absolute;
          left: 13px;
          top: 0;
          bottom: 24px;
          width: 1.5px;
          background: linear-gradient(
            to bottom,
            var(--tree-line-color) 0%,
            var(--tree-line-color) 70%,
            transparent 100%
          );
          border-radius: 1px;
        }
      }
    }

    // 子节点的横向连接线
    :deep(.el-tree-node__children > .el-tree-node) {
      position: relative;

      &::before {
        content: '';
        position: absolute;
        left: 13px;
        top: 24px;
        width: 16px;
        height: 1.5px;
        background-color: var(--tree-line-color);
        border-radius: 1px;
      }
    }

    // 悬停高亮连接线
    :deep(.el-tree-node__children > .el-tree-node:hover) {
      &::before {
        background-color: var(--tree-line-color-hover);
      }
    }

    // 展开/收起箭头样式
    :deep(.el-tree-node__expand-icon) {
      font-size: 14px;
      color: var(--el-text-color-secondary);
      transition: transform 0.25s ease, color 0.2s;
      padding: 4px;

      &:hover {
        color: var(--el-color-primary);
      }

      &.is-leaf {
        color: transparent;
      }
    }

    // 不同层级的节点左边距视觉增强
    :deep(.el-tree-node[aria-level="1"] > .el-tree-node__content) {
      font-weight: 600;
      font-size: 15px;
    }

    :deep(.el-tree-node[aria-level="2"] > .el-tree-node__content) {
      font-weight: 500;
      font-size: 14px;
    }

    :deep(.el-tree-node[aria-level="3"] > .el-tree-node__content) {
      font-weight: 400;
      font-size: 13px;
    }
  }

  // ====== 节点内容布局 ======
  .tree-node {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    padding-right: 8px;

    .node-left {
      display: flex;
      align-items: center;
      gap: 8px;
      flex: 1;
      min-width: 0;

      .folder-icon {
        color: var(--el-color-warning);
        flex-shrink: 0;
        filter: drop-shadow(0 1px 2px rgba(230, 162, 60, 0.3));
      }

      .node-label {
        color: var(--el-text-color-primary);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      .doc-tag {
        flex-shrink: 0;
      }

      .level-tag {
        flex-shrink: 0;
        font-size: 11px;
      }
    }

    .node-actions {
      display: flex;
      align-items: center;
      gap: 4px;
      flex-shrink: 0;
      opacity: 0;
      transition: opacity 0.25s ease;
    }
  }

  // 悬停时显示操作按钮
  :deep(.el-tree-node__content:hover) {
    .node-actions {
      opacity: 1;
    }
  }

  .form-tip {
    margin-left: 12px;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
  }
}
</style>
