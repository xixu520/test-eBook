<template>
  <div class="field-config-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="title-area">
        <h2 class="page-title">文档属性管理</h2>
        <p class="page-subtitle">管理全局文档属性定义，控制其在首页的展示方式。</p>
      </div>
      <el-button type="primary" :icon="Plus" @click="addField">新增属性</el-button>
    </div>

    <!-- 系统固有属性（只读） -->
    <el-card class="fixed-attrs-card" shadow="never">
      <template #header>
        <div class="card-header-inner">
          <el-icon><Lock /></el-icon>
          <span>系统固有属性（只读，不可增删改）</span>
        </div>
      </template>
      <div class="fixed-attrs-grid">
        <div v-for="attr in fixedAttributes" :key="attr.label" class="fixed-attr-item">
          <el-tag type="info" effect="plain" size="large" class="attr-tag">
            <el-icon class="tag-icon"><component :is="attr.icon" /></el-icon>
            {{ attr.label }}
          </el-tag>
          <span class="attr-desc">{{ attr.desc }}</span>
        </div>
      </div>
    </el-card>

    <!-- 自定义属性表格 -->
    <el-card class="attrs-card" shadow="never">
      <template #header>
        <div class="card-header-inner">
          <el-icon><Setting /></el-icon>
          <span>自定义属性</span>
          <el-tag type="primary" size="small" round style="margin-left: 8px">
            {{ fieldList.length }} 个
          </el-tag>
        </div>
      </template>

      <el-table
        :data="fieldList"
        border
        stripe
        style="width: 100%"
        class="field-table"
        v-loading="loading"
      >
        <!-- 属性名称 -->
        <el-table-column label="属性名称" min-width="160">
          <template #default="{ row }">
            <el-input v-model="row.label" placeholder="如：发布机构" size="small" />
          </template>
        </el-table-column>

        <!-- 类型 -->
        <el-table-column label="类型" width="130">
          <template #default="{ row }">
            <el-select v-model="row.field_type" size="small" style="width: 100%">
              <el-option label="单行文本" value="input" />
              <el-option label="数字" value="number" />
              <el-option label="日期" value="date" />
              <el-option label="下拉选择" value="select" />
              <el-option label="多选框" value="checkbox" />
            </el-select>
          </template>
        </el-table-column>

        <!-- 选项（仅 select/checkbox 显示） -->
        <el-table-column label="可选项（逗号隔开）" min-width="200">
          <template #default="{ row }">
            <el-input
              v-if="row.field_type === 'select' || row.field_type === 'checkbox'"
              v-model="row.options"
              placeholder="如：北京,上海,广州"
              size="small"
            />
            <span v-else class="text-disabled">—</span>
          </template>
        </el-table-column>

        <!-- 默认值 -->
        <el-table-column label="默认值" width="160">
          <template #default="{ row }">
            <el-date-picker
              v-if="row.field_type === 'date'"
              v-model="row.default_value"
              type="date"
              value-format="YYYY-MM-DD"
              size="small"
              style="width: 100%"
            />
            <el-input
              v-else
              v-model="row.default_value"
              placeholder="可选"
              size="small"
            />
          </template>
        </el-table-column>

        <!-- 首页展示 -->
        <el-table-column label="首页展示" width="90" align="center">
          <template #header>
            <span>首页展示</span>
            <el-tooltip content="开启后，该属性将在首页文档列表和搜索筛选中显示" placement="top">
              <el-icon style="margin-left: 4px; cursor: help; color: #909399"><QuestionFilled /></el-icon>
            </el-tooltip>
          </template>
          <template #default="{ row }">
            <el-switch v-model="row.show_in_home" size="small" />
          </template>
        </el-table-column>

        <!-- 必填 -->
        <el-table-column label="必填" width="70" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.is_required" size="small" />
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="80" align="center" fixed="right">
          <template #default="{ $index }">
            <el-button
              link
              type="danger"
              :icon="Delete"
              @click="removeField($index)"
            />
          </template>
        </el-table-column>

        <!-- 空状态 -->
        <template #empty>
          <el-empty description="暂无自定义属性，点击右上角「新增属性」添加" :image-size="100" />
        </template>
      </el-table>

      <!-- 底部操作栏 -->
      <div class="config-footer">
        <span class="save-tip" v-if="isDirty">您有未保存的更改</span>
        <el-button @click="resetFields" :disabled="!isDirty">恢复</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveFields">
          保存配置
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Plus, Delete, Lock, Setting, QuestionFilled } from '@element-plus/icons-vue'
import { getGlobalForm, saveFormFields, type FormField } from '@/api/form'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const globalFormId = ref<number>(0)
const fieldList = ref<FormField[]>([])
const originalFields = ref<FormField[]>([])

// 系统固有属性定义（只读展示）
const fixedAttributes = [
  { label: '所属分类', desc: '文档所属的分类节点', icon: 'Folder' },
  { label: '文件大小', desc: '上传文件的字节大小', icon: 'Files' },
  { label: '上传时间', desc: '文档创建时间', icon: 'Calendar' },
  { label: '处理状态', desc: 'OCR 解析处理状态', icon: 'Loading' },
  { label: '核验状态', desc: '文档人工核验状态', icon: 'CircleCheck' },
  { label: '同步状态', desc: '远程同步状态', icon: 'Refresh' },
]

const isDirty = computed(() => {
  return JSON.stringify(fieldList.value) !== JSON.stringify(originalFields.value)
})

const loadGlobalForm = async () => {
  loading.value = true
  try {
    const res: any = await getGlobalForm()
    globalFormId.value = res.ID
    fieldList.value = JSON.parse(JSON.stringify(res.fields || []))
    originalFields.value = JSON.parse(JSON.stringify(res.fields || []))
  } catch (error) {
    ElMessage.error('加载全局属性失败')
  } finally {
    loading.value = false
  }
}

const addField = () => {
  fieldList.value.push({
    label: '',
    field_key: `field_${Date.now()}`,
    field_type: 'input',
    is_required: false,
    order: fieldList.value.length * 10,
    show_in_home: true,
    show_in_filter: true,
    show_in_list: true,
    show_in_admin: true,
    default_value: '',
    options: ''
  })
}

const removeField = (index: number) => {
  ElMessageBox.confirm('确定删除该属性？删除后，已有文档的该属性值将无法展示。', '提示', {
    type: 'warning',
    confirmButtonText: '确定删除',
    cancelButtonText: '取消'
  }).then(() => {
    fieldList.value.splice(index, 1)
  }).catch(() => {})
}

const resetFields = () => {
  fieldList.value = JSON.parse(JSON.stringify(originalFields.value))
}

const handleSaveFields = async () => {
  if (!globalFormId.value) return ElMessage.error('未找到全局属性配置，请刷新页面')

  // 校验属性名称不能为空
  for (const f of fieldList.value) {
    if (!f.label.trim()) {
      return ElMessage.warning('属性名称不能为空，请检查后再保存')
    }
  }

  // show_in_filter 与 show_in_home 联动（决策5-C）
  const fieldsToSave = fieldList.value.map((f, index) => ({
    ...f,
    label: f.label.trim(),
    order: index * 10,
    show_in_filter: f.show_in_home,  // 联动
    show_in_admin: f.show_in_home,   // 保持兼容性
    show_in_list: f.show_in_home,    // 保持兼容性
    // field_key 保持不变，由前端 addField 时自动生成
  }))

  saving.value = true
  try {
    await saveFormFields(globalFormId.value, fieldsToSave)
    ElMessage.success('属性配置已保存')
    await loadGlobalForm()
  } catch (err) {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadGlobalForm()
})
</script>

<style scoped lang="scss">
.field-config-page {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;

    .page-title {
      margin: 0 0 6px 0;
      font-size: 22px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .page-subtitle {
      margin: 0;
      font-size: 13px;
      color: var(--el-text-color-secondary);
    }
  }

  .card-header-inner {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    font-size: 14px;
  }

  .fixed-attrs-card {
    :deep(.el-card__header) {
      padding: 14px 20px;
      background-color: var(--el-fill-color-lighter);
    }

    .fixed-attrs-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 12px;
      padding: 4px 0;

      .fixed-attr-item {
        display: flex;
        align-items: center;
        gap: 10px;

        .attr-tag {
          display: flex;
          align-items: center;
          gap: 4px;
          flex-shrink: 0;

          .tag-icon {
            font-size: 13px;
          }
        }

        .attr-desc {
          font-size: 12px;
          color: var(--el-text-color-placeholder);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
  }

  .attrs-card {
    .field-table {
      :deep(.el-table__row) {
        transition: background-color 0.2s;
        &:hover > td {
          background-color: var(--el-color-primary-light-9) !important;
        }
      }
    }

    .config-footer {
      display: flex;
      justify-content: flex-end;
      align-items: center;
      gap: 12px;
      margin-top: 20px;
      padding-top: 16px;
      border-top: 1px solid var(--el-border-color-lighter);

      .save-tip {
        font-size: 13px;
        color: var(--el-color-warning);
        margin-right: auto;

        &::before {
          content: '● ';
        }
      }
    }
  }
}

.text-disabled {
  color: var(--el-text-color-placeholder);
  font-size: 13px;
}
</style>
