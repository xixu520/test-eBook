<template>
  <el-drawer
    v-model="visible"
    title="表单管理"
    size="1000px"
    destroy-on-close
    @open="loadForms"
  >
    <div class="form-manager">
      <div class="left-aside">
        <div class="aside-header">
          <span>模板列表</span>
          <el-button type="primary" link :icon="Plus" @click="addForm">新增表单</el-button>
        </div>
        <el-scrollbar>
          <div class="form-list">
            <div
              v-for="item in forms"
              :key="item.ID"
              class="form-item"
              :class="{ active: currentForm?.ID === item.ID }"
              @click="selectForm(item)"
            >
              <div class="info">
                <div class="name">{{ item.name }}</div>
                <div class="desc">{{ item.description || '无描述' }}</div>
              </div>
              <div class="actions">
                <el-button link type="primary" :icon="Edit" @click.stop="editFormName(item)" />
                <el-button link type="danger" :icon="Delete" @click.stop="handleDeleteForm(item)" />
              </div>
            </div>
          </div>
        </el-scrollbar>
      </div>

      <div class="right-main">
        <template v-if="currentForm">
          <div class="main-header">
            <div class="header-title">
              <span class="label">当前配置：</span>
              <span class="name">{{ currentForm.name }}</span>
            </div>
            <el-button type="primary" :icon="Plus" @click="addField">新增属性字段</el-button>
          </div>

          <div class="table-container">
            <el-table :data="fieldList" border stripe style="width: 100%" class="field-table" empty-text="该模板暂未配置任何属性字段">
              <el-table-column prop="label" label="属性名称 (Label)" width="180">
                <template #default="{ row }">
                  <el-input v-model="row.label" placeholder="如：协议号" size="default" />
                </template>
              </el-table-column>
              <el-table-column prop="field_key" label="字段标示 (Key)" width="180">
                <template #default="{ row }">
                  <el-input 
                    v-model="row.field_key" 
                    placeholder="字母开头, 限字母数字下划线" 
                    size="default"
                    @input="(val: string) => handleKeyInput(row, val)"
                  />
                </template>
              </el-table-column>
              <el-table-column prop="field_type" label="类型" width="140">
                <template #default="{ row }">
                  <el-select v-model="row.field_type" size="default">
                    <el-option label="单行文本" value="input" />
                    <el-option label="数字输入" value="number" />
                    <el-option label="日期选择" value="date" />
                    <el-option label="下拉选择" value="select" />
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column prop="is_required" label="必填项" width="80" align="center">
                <template #default="{ row }">
                  <el-checkbox v-model="row.is_required" />
                </template>
              </el-table-column>
              <el-table-column label="扩展配置 (Options)" min-width="200">
                <template #default="{ row }">
                  <div v-if="row.field_type === 'select'" class="option-config">
                     <el-input v-model="row.options" placeholder="多选项用,隔开 (选项A,选项B)" size="small" />
                  </div>
                  <span v-else class="text-disabled">无需配置</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80" align="center" fixed="right">
                <template #default="{ $index }">
                  <el-button link type="danger" :icon="Delete" @click="removeField($index)">移除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="main-footer">
            <el-button plain @click="loadForms">重置修改</el-button>
            <el-button type="primary" :icon="Check" :loading="saving" @click="handleSaveFields">同步并保存配置</el-button>
          </div>
        </template>
        <el-empty v-else description="请在左侧点击表单模板开始配置字段" :image-size="200" />
      </div>
    </div>

    <!-- 表单名称编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogType === 'add' ? '新增表单模板' : '编辑表单信息'" width="400px" append-to-body>
      <el-form label-width="80px" label-position="top">
        <el-form-item label="名称">
          <el-input v-model="formModel.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formModel.description" type="textarea" placeholder="描述信息" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Plus, Edit, Delete, Check } from '@element-plus/icons-vue'
import { getForms, createForm, updateForm, deleteForm, saveFormFields, type Form, type FormField } from '@/api/form'
import { ElMessage, ElMessageBox } from 'element-plus'

const visible = defineModel<boolean>()
const forms = ref<Form[]>([])
const currentForm = ref<Form | null>(null)
const fieldList = ref<FormField[]>([])
const saving = ref(false)

const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formModel = reactive({
  name: '',
  description: ''
})

const loadForms = async () => {
  try {
    const res = await getForms()
    forms.value = res as any
    // 如果当前有选中的表单，尝试重新绑定引用以触发 UI 更新
    if (currentForm.value) {
      const updated = forms.value.find(f => f.ID === currentForm.value!.ID)
      if (updated) {
        currentForm.value = updated
      }
    }
  } catch (error) {
    console.error('加载列表失败', error)
  }
}

const selectForm = (form: Form) => {
  currentForm.value = form
  // 深拷贝字段列表
  fieldList.value = JSON.parse(JSON.stringify(form.fields || []))
}

const addForm = () => {
  dialogType.value = 'add'
  formModel.name = ''
  formModel.description = ''
  dialogVisible.value = true
}

const editFormName = (form: Form) => {
  dialogType.value = 'edit'
  currentForm.value = form
  formModel.name = form.name
  formModel.description = form.description
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formModel.name) return ElMessage.warning('名称不能为空')
  try {
    if (dialogType.value === 'add') {
      await createForm(formModel)
      ElMessage.success('表单添加成功')
    } else if (currentForm.value) {
      await updateForm(currentForm.value.ID, formModel)
      // 立即更新本地状态，增强 UI 即使感
      currentForm.value.name = formModel.name
      currentForm.value.description = formModel.description
      ElMessage.success('表单信息已更新')
    }
    dialogVisible.value = false
    await loadForms()
  } catch (err) {
    console.error('保存失败', err)
  }
}

const handleDeleteForm = (form: Form) => {
  ElMessageBox.confirm(`确定删除表单「${form.name}」吗？其下所有字段也将被强制删除。`, '警告', {
    type: 'error'
  }).then(async () => {
    await deleteForm(form.ID)
    ElMessage.success('解构化表单已移除')
    if (currentForm.value?.ID === form.ID) currentForm.value = null
    loadForms()
  }).catch(() => {})
}

// 字段逻辑
const addField = () => {
  fieldList.value.push({
    label: '',
    field_key: '',
    field_type: 'input',
    is_required: false,
    order: fieldList.value.length * 10
  })
}

const removeField = (index: number) => {
  fieldList.value.splice(index, 1)
}

const handleKeyInput = (row: any, val: string) => {
  // 1. 移除非法字符（仅保留字母、数字、下划线）
  let filtered = val.replace(/[^a-zA-Z0-9_]/g, '')
  // 2. 强制首位必须是字母（如果首位是数字或下划线，则删除）
  if (filtered.length > 0 && !/^[a-zA-Z]/.test(filtered)) {
    filtered = filtered.replace(/^[^a-zA-Z]+/, '')
  }
  row.field_key = filtered
}

const handleSaveFields = async () => {
  if (!currentForm.value) return
  // 校验
  const keyRegex = /^[a-zA-Z][a-zA-Z0-9_]*$/
  for (const f of fieldList.value) {
    if (!f.label) return ElMessage.warning('请补全字段名称')
    if (!f.field_key) return ElMessage.warning('请补全字段唯一标示')
    if (!keyRegex.test(f.field_key)) {
      return ElMessage.warning(`字段标示 [${f.field_key}] 格式错误：必须以字母开头，且仅包含字母、数字及下划线`)
    }
  }

  saving.value = true
  try {
    await saveFormFields(currentForm.value.ID, fieldList.value)
    ElMessage.success('表单集字段配置已同步')
    loadForms() // 刷新 forms 列表，更新 currentForm 中的 fields
  } catch (err) {} finally {
    saving.value = false
  }
}
</script>

<style scoped lang="scss">
.form-manager {
  display: flex;
  height: calc(100vh - 120px); // 减去 Drawer Header 后的高度
  margin: -20px; // 抵消 el-drawer 的内边距以实现两栏全高
  background-color: var(--el-bg-color-page);
  overflow: hidden;

  .left-aside {
    width: 280px;
    flex-shrink: 0;
    border-right: 1px solid var(--el-border-color-lighter);
    display: flex;
    flex-direction: column;
    background-color: var(--el-bg-color);

    .aside-header {
      padding: 20px 16px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      border-bottom: 1px solid var(--el-border-color-extra-light);
      background-color: var(--el-fill-color-blank);
      flex-shrink: 0;
      
      span {
        font-weight: 600;
        font-size: 15px;
        color: var(--el-text-color-primary);
        white-space: nowrap; // 关键：防止标题被挤压导致垂直排列
      }
    }

    .form-item {
      padding: 16px;
      margin: 8px 12px;
      border-radius: 10px;
      cursor: pointer;
      display: flex;
      justify-content: space-between;
      align-items: center;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      border: 1px solid transparent;

      &:hover {
        background-color: var(--el-fill-color-light);
        transform: translateX(4px);
      }

      &.active {
        background-color: white;
        border-color: var(--el-color-primary-light-3);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        
        .name { color: var(--el-color-primary); font-weight: 600; }
        .arrow-icon { opacity: 1; transform: translateX(0); }
      }

      .info {
        flex: 1;
        min-width: 0;
        .name { font-size: 14px; margin-bottom: 6px; }
        .desc { font-size: 12px; color: var(--el-text-color-secondary); opacity: 0.8; }
      }

      .actions {
        display: flex;
        gap: 4px;
        opacity: 0;
        transition: opacity 0.2s;
      }

      &:hover .actions {
        opacity: 1;
      }
    }
  }

  .right-main {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    background-color: var(--el-bg-color);
    position: relative;

    .main-header {
      padding: 24px 32px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      background-color: var(--el-bg-color);
      border-bottom: 1px solid var(--el-border-color-extra-light);
      flex-shrink: 0;

      .header-title {
        .label { font-size: 13px; color: var(--el-text-color-secondary); }
        .name { font-size: 18px; font-weight: 600; color: var(--el-text-color-primary); }
      }
    }

    .table-container {
      flex: 1;
      padding: 24px 32px;
      overflow: auto;
    }

    .field-table {
      border-radius: 8px;
      overflow: hidden;
      box-shadow: var(--el-box-shadow-lighter);

      :deep(.el-table__header-wrapper) th {
        background-color: var(--el-fill-color-light);
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      :deep(.el-input__inner) {
        border-color: transparent !important;
        background: transparent !important;
        text-align: left;
        &:focus { border-color: var(--el-color-primary) !important; background: white !important; }
      }
    }

    .main-footer {
      padding: 16px 32px;
      background-color: var(--el-fill-color-blank);
      border-top: 1px solid var(--el-border-color-extra-light);
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      flex-shrink: 0;
      box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.03);
    }
  }
}

.text-disabled {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.option-config {
  :deep(.el-input__inner) {
    font-size: 12px;
    font-family: var(--el-font-family-monospace);
  }
}
</style>
