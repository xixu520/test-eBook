# PLAN-54：动态表单属性体系重构（分类管理页面升级）

> **目标**：将文档的「硬编码属性」（版本、发布机构、实施状态、标准号、实施日期等）从 `StandardFile` 模型中剥离，转由动态表单系统（`Form` + `FormField`）驱动。实现属性可增删改、属性值随文档存储、可控制前端展示/隐藏的能力。
>
> **核心理念**：保留 `所属分类 (category_id)` 和 `OCR 解析状态 (status)` 为固有字段，其余所有业务元数据全部动态化。

---

## 一、整体架构设计

### 1.1 当前架构问题

| 维度 | 现状 | 改进成果 |
|------|------|------|
| 数据模型 | `StandardFile` 包含硬编码字段 | **已解耦**。引入 `DocumentFieldValue` EAV 模型存储动态属性 |
| 表单系统 | 分类绑定模板，但未与文档关联 | **已打通**。文档上传/更新时自动同步动态属性值 |
| 前端展示 | 硬编码属性列 | **已重构**。列表列、筛选器、编辑表单均由模板动态驱动 |
| 完成情况 | 进度：0% | **进度：100%** |

### 1.2 目标架构

```
┌───────────────────────────────────────────────────────────┐
│  FormField（属性定义）                                      │
│  + show_in_list (bool)   ← 控制是否在列表页展示              │
│  + show_in_filter (bool) ← 控制是否在筛选区展示              │
│  + default_value (string)← 默认值                           │
├───────────────────────────────────────────────────────────┤
│  DocumentFieldValue（属性值存储）  ← 新增模型                 │
│  document_id + field_id → value                             │
├───────────────────────────────────────────────────────────┤
│  StandardFile（文档）                                        │
│  保留: id, title, category_id, status, file_path,           │
│        file_size, ocr_content, sync_status, tags            │
│  移除: number, year, version, publisher, implementation_*,  │
│        verify_status                                        │
└───────────────────────────────────────────────────────────┘
```

---

## 二、需确认的设计决策

### 决策 1：硬编码属性的迁移策略

> **方案 A（推荐 ✅ 渐进式）**：保留 `StandardFile` 中的现有硬编码字段不删除，但标记为 `legacy`。新文档使用动态属性存储。前端优先读取动态属性值，如不存在则回退读取旧字段。待数据完全迁移后再移除旧字段。
>


### 决策 2：动态属性是否绑定到分类

> **方案 A（推荐 ✅ 分类绑定）**：保留现有 `Category.form_id` 关联。每个分类拥有独立的表单模板，文档上传时根据所选分类自动加载对应的属性表单。
>


### 决策 3：是否创建全局默认表单

> 对于没有绑定表单的分类，是否创建一个「全局默认表单」，包含标准号、版本、发布机构等基础属性？
>
> **方案 A（推荐 ✅）**：创建默认表单，所有未绑定表单的分类自动使用。


---

## 三、后端改动清单

### 3.1 模型层 (`internal/model/`)

#### [MODIFY] `form.go`

```
FormField 结构体新增字段:
+ ShowInList    bool   `json:"show_in_list" gorm:"default:true"`
+ ShowInFilter  bool   `json:"show_in_filter" gorm:"default:false"`
+ DefaultValue  string `json:"default_value" gorm:"type:varchar(255)"`
```

**函数级影响**: 无新增函数，仅结构体扩展。GORM AutoMigrate 将自动追加列。

#### [NEW] `document_field_value.go`

新增模型：文档动态属性值存储表

```go
// DocumentFieldValue 文档动态属性值
type DocumentFieldValue struct {
    ID         uint   `gorm:"primarykey" json:"id"`
    DocumentID uint   `json:"document_id" gorm:"index;not null"`
    FieldID    uint   `json:"field_id" gorm:"index;not null"`
    Value      string `json:"value" gorm:"type:text"`
    // 关联
    Field      FormField `json:"field" gorm:"foreignKey:FieldID"`
}
```

**说明**: `(document_id, field_id)` 加联合唯一索引防重复。

#### [MODIFY] `standard.go` — StandardFile

如果采用方案 A（渐进式），模型暂不删除字段。但需在 JSON tag 中做标记：

```go
// 以下字段为 legacy，未来将移除。新数据使用 DocumentFieldValue 存储。
Number               string `json:"number" gorm:"..."` // legacy
Year                 string `json:"year" gorm:"..."`   // legacy
...
```

无需删减字段，仅注释标记。

---

### 3.2 数据库迁移 (`internal/database/migrate.go`)

#### [MODIFY] `migrate.go`

```diff
 func AutoMigrate() error {
     return WriteDB.AutoMigrate(
         ...
         &model.Form{},
         &model.FormField{},
+        &model.DocumentFieldValue{},
     )
 }
```

---

### 3.3 仓库层 (`internal/repository/`)

#### [NEW] `document_field_value_repo.go`

新增 Repository，包含以下函数：

| 函数签名 | 职责 |
|----------|------|
| `NewDocFieldValueRepo(db *gorm.DB) *DocFieldValueRepo` | 构造函数 |
| `GetByDocumentID(docID uint) ([]DocumentFieldValue, error)` | 按文档 ID 查询所有动态字段值（Preload Field） |
| `BatchSave(docID uint, values []DocumentFieldValue) error` | 事务内批量保存（先删旧值再插入新值） |
| `DeleteByDocumentID(docID uint) error` | 删除文档时级联清理 |
| `GetByFieldIDs(fieldIDs []uint, filterValue string) ([]uint, error)` | 按字段值筛选，返回匹配的 document_id 列表（用于搜索） |

#### [MODIFY] `form_repo.go`

无需修改。现有函数已满足需求。

#### [MODIFY] `standard_repo.go`

| 函数 | 变更说明 |
|------|----------|
| `ListFiles(...)` | 新增参数 `dynamicFilters map[uint]string`（field_id → 搜索值）。当存在动态筛选条件时，通过子查询 JOIN `document_field_values` 表过滤。 |
| `FindFileByID(id)` | 保持不变（详情页的动态属性值由 Service 层另行查询拼装）。 |

---

### 3.4 服务层 (`internal/service/`)

#### [MODIFY] `form_service.go`

| 函数 | 变更说明 |
|------|----------|
| `SaveFormFields(formID, fields)` | 需同步保存新增的 `show_in_list`、`show_in_filter`、`default_value` 字段。当前已支持全量覆盖，无需额外逻辑。 |
| `GetForms()` | 无需修改（已 Preload Fields）。 |
| [NEW] `GetFormByID(id uint) (*Form, error)` | 按 ID 获取单个表单（含 Fields），用于文档编辑时加载表单模板。 |
| [NEW] `GetDefaultFormID() uint` | 返回全局默认表单 ID（若采用决策 3 方案 A）。 |

#### [NEW] `document_field_service.go`

新增服务，职责为文档动态属性值的读写：

| 函数签名 | 职责 |
|----------|------|
| `NewDocFieldService(repo, formRepo, fieldValueRepo)` | 构造函数 |
| `GetFieldValues(docID uint) ([]DocumentFieldValue, error)` | 获取文档的所有动态属性值 |
| `SaveFieldValues(docID uint, values []FieldValueInput) error` | 批量保存文档属性值。`FieldValueInput` = `{field_id, value}` |
| `GetListColumns(formID uint) ([]FormField, error)` | 获取指定表单中 `show_in_list=true` 的字段列表 |
| `GetFilterFields(formID uint) ([]FormField, error)` | 获取指定表单中 `show_in_filter=true` 的字段列表 |

#### [MODIFY] `standard_service.go`

| 函数 | 变更说明 |
|------|----------|
| `UploadFile(...)` | 1. 签名新增 `dynamicFields map[string]string` 参数。<br>2. 创建 `StandardFile` 后，调用 `DocFieldService.SaveFieldValues()` 保存动态属性。 |
| `UpdateFile(...)` | 1. 签名简化：移除 `number, version, publisher, implementationDate, implStatus, verifyStatus` 等硬编码参数。<br>2. 新增 `dynamicFields []FieldValueInput` 参数。<br>3. 更新 `StandardFile` 固有字段后，调用 `DocFieldService.SaveFieldValues()` 更新动态值。 |
| `GetFileDetail(id)` | 返回值增强：除 `*StandardFile` 外，额外返回 `[]DocumentFieldValue`。或新增一个 `GetFileDetailFull()` 返回聚合结构。 |
| `SearchFiles(...)` | 新增 `dynamicFilters map[uint]string` 参数透传到 Repo 层。 |

---

### 3.5 处理层 (`internal/handler/`)

#### [MODIFY] `standard.go`

| 函数 | 变更说明 |
|------|----------|
| `UploadFile(c)` | 1. 从 FormData 中读取 `dynamic_fields` JSON 字符串。<br>2. 解析为 `map[string]string` 后传递给 Service。 |
| `UpdateFile(c)` | 1. Input 结构体改为接受 `title`、`category_id`、`status` + `dynamic_fields []FieldValueInput`。<br>2. 移除 `number`, `version`, `publisher` 等硬编码字段。 |
| `GetFileDetail(c)` | 返回值中追加 `field_values` 数组。 |
| `ListFiles(c)` | 1. 解析 URL query 中的动态筛选参数 `filter[field_id]=value`。<br>2. 透传给 Service。 |

#### [NEW] `document_field_handler.go`

或直接在 `standard.go` 中追加：

| 函数签名 | API 路径 | 职责 |
|----------|----------|------|
| `GetDocumentFields(c)` | `GET /documents/:id/fields` | 获取文档的动态属性值列表 |
| `SaveDocumentFields(c)` | `PUT /documents/:id/fields` | 批量保存文档的动态属性值 |

#### [MODIFY] `form_handler.go`

无需修改。现有 CRUD 接口已满足表单模板管理需求。`FormField` 结构体扩展后 JSON 自动携带新字段。

---

### 3.6 路由层 (`internal/router/router.go`)

#### [MODIFY] `router.go`

```diff
 documents := protected.Group("/documents")
 {
     ...
+    documents.GET("/:id/fields", standardHandler.GetDocumentFields)
+    documents.PUT("/:id/fields", standardHandler.SaveDocumentFields)
 }
```

---

## 四、前端改动清单

### 4.1 API 层 (`src/api/`)

#### [MODIFY] `form.ts`

```diff
 export interface FormField {
   ...
+  show_in_list: boolean
+  show_in_filter: boolean
+  default_value: string
 }
```

#### [NEW] `document-field.ts`

```typescript
// 获取文档动态属性值
export function getDocumentFields(docId: number): Promise<FieldValue[]>

// 保存文档动态属性值
export function saveDocumentFields(docId: number, values: FieldValueInput[]): Promise<void>

// 类型定义
export interface FieldValue {
  id: number
  field_id: number
  value: string
  field: FormField  // 含 label、field_type 等定义
}

export interface FieldValueInput {
  field_id: number
  value: string
}
```

#### [MODIFY] `document.ts`

- `Document` 接口新增可选的 `field_values?: FieldValue[]`
- `updateDocument()` 函数的 data 参数结构调整

---

### 4.2 分类管理页 (`views/admin/category/`)

#### [MODIFY] `CategoryPage.vue`

**保留原有功能**，仅微调：

| 函数/区域 | 变更说明 |
|-----------|----------|
| 对话框表单 `<el-form>` | "绑定表单" 的 `<el-select>` 保持不变 |
| `handleAdd` / `handleEdit` / `submitForm` | 保持不变 |
| 树节点展示 | 可选：节点上展示绑定的表单模板名称（如 `[建筑标准模板]`） |

**整体无破坏性变更。**

#### [MODIFY] `FormManagerDrawer.vue`

表格列扩展：

| 区域 | 变更说明 |
|------|----------|
| `<el-table>` 字段表 | 新增 3 列：<br>1. `列表展示` — `<el-checkbox v-model="row.show_in_list" />`<br>2. `筛选区展示` — `<el-checkbox v-model="row.show_in_filter" />`<br>3. `默认值` — `<el-input v-model="row.default_value" />` |
| `addField()` 函数 | 新增字段对象追加 `show_in_list: true, show_in_filter: false, default_value: ''` |
| `handleSaveFields()` | 无需修改（已全量提交 fieldList 数组） |

**Drawer 宽度**：考虑从 `1000px` 扩展到 `1100px` 以容纳新列。

---

### 4.3 文档管理页 (`views/admin/document/AdminDocumentsPage.vue`)

此页面改动最大，核心变化是 **表格列和编辑弹窗由硬编码转为动态渲染**。

#### 数据加载改造

| 函数 | 变更说明 |
|------|----------|
| `loadData()` | 1. 同时请求 `getForms()` 获取所有表单模板定义。<br>2. 构建 `formFieldMap: Map<number, FormField[]>`（form_id → fields）。<br>3. 构建 `listableFields: FormField[]`（`show_in_list=true` 的字段集合）。 |
| `loadCategories()` | 保持不变。 |
| [NEW] `loadFormDefinitions()` | 调用 `getForms()` 并处理上述映射逻辑。 |

#### 表格列动态渲染

```vue
<!-- 固定列：保留 -->
<el-table-column prop="title" label="文档名称" />
<el-table-column label="所属分类" />
<el-table-column label="文件大小" />
<el-table-column label="上传时间" />
<el-table-column label="OCR 状态" />

<!-- 动态列：根据 listableFields 动态渲染 -->
<el-table-column
  v-for="field in listableFields"
  :key="field.ID"
  :label="field.label"
  :width="getFieldColumnWidth(field)"
>
  <template #default="{ row }">
    {{ getFieldDisplayValue(row, field) }}
  </template>
</el-table-column>

<!-- 操作列：保留 -->
<el-table-column label="操作" />
```

| 新增函数 | 职责 |
|----------|------|
| `getFieldColumnWidth(field)` | 根据 `field_type` 返回合理列宽（date→120, input→180, select→120） |
| `getFieldDisplayValue(row, field)` | 从 `row.field_values` 中按 `field_id` 查找并格式化显示值。对 `select` 类型显示标签名而非原始值。 |

#### 编辑弹窗动态化

**移除**原有的硬编码表单项（`number`、`version`、`publisher`、`implementation_date`、`implementation_status`、`verify_status`）。

**保留**：
- `title`（文档名称）
- `category_id`（所属分类）
- `status`（OCR 状态 — Radio 组）

**新增动态表单区域**：

```vue
<!-- 动态属性编辑区 -->
<el-divider content-position="left">扩展属性</el-divider>
<el-row :gutter="20">
  <el-col :span="12" v-for="field in editableFields" :key="field.ID">
    <el-form-item :label="field.label" :prop="'dynamic_' + field.field_key">
      <!-- 根据 field_type 动态渲染控件 -->
      <el-input v-if="field.field_type === 'input'" v-model="dynamicValues[field.ID]" />
      <el-input-number v-else-if="field.field_type === 'number'" v-model="dynamicValues[field.ID]" />
      <el-date-picker v-else-if="field.field_type === 'date'" v-model="dynamicValues[field.ID]"
        type="date" value-format="YYYY-MM-DD" />
      <el-select v-else-if="field.field_type === 'select'" v-model="dynamicValues[field.ID]">
        <el-option v-for="opt in parseOptions(field.options)" :key="opt" :label="opt" :value="opt" />
      </el-select>
    </el-form-item>
  </el-col>
</el-row>
```

| 新增函数 | 职责 |
|----------|------|
| `loadEditableFields(categoryId)` | 通过 `categoryId` 找到关联的 `form_id`，加载该表单的所有 `FormField` |
| `initDynamicValues(row)` | 编辑时从 `row.field_values` 初始化 `dynamicValues` 对象 |
| `parseOptions(optionsStr)` | 解析 `,` 分隔的选项字符串为数组 |
| `submitEdit()` 改造 | 除提交固有字段外，额外调用 `saveDocumentFields(id, dynamicValues)` |

#### 详情弹窗（需新增）

**点击文件详情时弹出的窗口，显示全部属性**：

| 区域 | 说明 |
|------|------|
| [NEW] `<el-dialog>` 详情弹窗 | 展示文档全部信息（固有字段 + 全部动态属性值） |
| [NEW] `handleDetail(row)` | 调用 `getDocumentFields(row.id)` 获取全量属性值 |
| 渲染方式 | `<el-descriptions>` 组件，每个字段为一行 |

---

### 4.4 首页 (`views/home/HomePage.vue`)

#### 表格列动态化

与 `AdminDocumentsPage.vue` 相似的改造策略：

| 区域 | 变更说明 |
|------|----------|
| 固定列 | 保留：标准号（未来从动态属性读取但兼容 legacy）、名称、所属分类、OCR 状态 |
| 动态列 | 根据 `show_in_list=true` 的 fields 动态渲染 |
| 筛选区 | 移除硬编码的 `发布机构`、`实施状态` 筛选。改为根据 `show_in_filter=true` 的 fields 动态渲染筛选控件 |

| 修改函数 | 变更说明 |
|----------|----------|
| `loadData()` | 额外加载表单定义 |
| `filters` reactive | 从硬编码改为 `dynamicFilters: Record<number, string>` |
| `resetFilters()` | 重置动态筛选器 |
| 硬编码的 `publishers` 数组 | 移除，由动态属性的 `options` 提供 |
| 详情弹窗 | 新增文档详情弹窗，点击行时弹出全属性 |

#### 筛选区动态渲染

```vue
<el-form-item v-for="field in filterableFields" :key="field.ID" :label="field.label">
  <el-input v-if="field.field_type === 'input'" v-model="dynamicFilters[field.ID]" clearable />
  <el-select v-else-if="field.field_type === 'select'" v-model="dynamicFilters[field.ID]" clearable>
    <el-option v-for="opt in parseOptions(field.options)" :key="opt" :label="opt" :value="opt" />
  </el-select>
  <el-date-picker v-else-if="field.field_type === 'date'" v-model="dynamicFilters[field.ID]"
    type="date" value-format="YYYY-MM-DD" clearable />
</el-form-item>
```

---

### 4.5 上传弹窗 (`components/document/UploadDialog.vue`)

#### 动态属性输入

| 区域 | 变更说明 |
|------|----------|
| 表单 `<el-form>` | 保留：文件标题、所属分类、选择文件。<br>移除：标准号、发布年份、版本、发布机构、实施日期、实施状态（硬编码项）。 |
| [NEW] 动态属性区 | 当用户选择 `category_id` 后，自动加载该分类绑定的表单模板 Fields，动态渲染输入控件。 |
| `startUpload()` | FormData 中追加 `dynamic_fields` JSON 字符串。 |
| `canUpload` computed | 校验逻辑改为：`title` + `category_id` + 文件已选 + 所有 `is_required=true` 的动态字段均已填写。 |

**关键交互**：用户切换分类时，表单区域自动刷新为对应的属性模板。

| 新增函数 | 职责 |
|----------|------|
| `watch(() => form.category_id, loadCategoryForm)` | 监听分类切换 |
| `loadCategoryForm(categoryId)` | 1. 从 `categories` 中查找该分类的 `form_id`。<br>2. 从 `formList` 中查找对应 Form 的 fields。<br>3. 更新 `currentFields` 和 `dynamicValues` 响应式对象。 |
| `resetDynamicValues()` | 切换分类时重置动态属性值 |

---

## 五、新增 API 端点汇总

| HTTP Method | Path | Handler 函数 | 说明 |
|-------------|------|-------------|------|
| `GET` | `/api/v1/documents/:id/fields` | `GetDocumentFields` | 获取文档的全部动态属性值 |
| `PUT` | `/api/v1/documents/:id/fields` | `SaveDocumentFields` | 批量保存文档的动态属性值 |

现有 API 变更：

| 端点 | 变更 |
|------|------|
| `POST /documents/upload` | 请求体新增 `dynamic_fields` 表单字段 |
| `PUT /documents/:id` | 请求体结构调整（移除硬编码属性，新增 `dynamic_fields`） |
| `GET /documents/:id` | 响应体追加 `field_values` 数组 |
| `GET /documents` | 支持 `filter[field_id]=value` 查询参数 |

---

## 六、实施顺序建议

按依赖关系从底层到顶层分 5 个阶段：

### 阶段 1：后端模型与仓库（无前端依赖）
1. 修改 `model/form.go` — 扩展 `FormField`
2. 新增 `model/document_field_value.go`
3. 修改 `database/migrate.go` — 添加新模型迁移
4. 新增 `repository/document_field_value_repo.go`
5. 修改 `repository/standard_repo.go` — `ListFiles` 支持动态筛选

### 阶段 2：后端服务与处理层
6. 新增 `service/document_field_service.go`
7. 修改 `service/form_service.go` — 添加 `GetFormByID`
8. 修改 `service/standard_service.go` — `UploadFile`、`UpdateFile`、`GetFileDetail` 改造
9. 修改 `handler/standard.go` — 对应 Handler 函数改造
10. 修改 `router/router.go` — 新增路由
11. 修改 `cmd/server/main.go` — 注入新 repository/service

### 阶段 3：前端 API 与公共逻辑
12. 修改 `api/form.ts` — 扩展 FormField 类型
13. 新增 `api/document-field.ts`
14. 修改 `api/document.ts` — 调整类型定义

### 阶段 4：前端页面改造
15. 修改 `FormManagerDrawer.vue` — 表格新增列
16. 修改 `UploadDialog.vue` — 动态化
17. 修改 `AdminDocumentsPage.vue` — 表格、编辑弹窗、详情弹窗
18. 修改 `HomePage.vue` — 表格列和筛选区动态化

### 阶段 5：集成验证
19. 全量回归测试
20. 创建全局默认表单（Seed 脚本或手动创建）
21. Legacy 数据兼容性验证

---

## 七、风险点

| 风险 | 应对措施 |
|------|----------|
| 动态列过多导致表格宽度溢出 | 表格强制 `min-width`，超出列横向滚动 |
| 筛选性能（JOIN 查询）| `document_field_values` 表加联合索引 `(field_id, value)` |
| Legacy 数据兼容 | 前端 `getFieldDisplayValue` 增加回退逻辑读取旧字段 |
| Category 无绑定 Form 时的行为 | 使用全局默认表单兜底 |

---

## 八、完成状态

- [x] 决策 1 确认
- [x] 决策 2 确认
- [x] 决策 3 确认
- [ ] 阶段 1：后端模型与仓库
- [ ] 阶段 2：后端服务与处理层
- [ ] 阶段 3：前端 API 与公共逻辑
- [ ] 阶段 4：前端页面改造
- [ ] 阶段 5：集成验证
