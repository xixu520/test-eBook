# PLAN-56：后台管理重构 — 分类管理优化、文档属性管理独立化、文件管理增强

> **目标**：重构后台管理模块。(1) 优化分类管理页面显示，保留分类树管理，去除表单模板绑定入口；(2) 建立独立的「文档属性管理」页面，属性可增删改，可配置在主页/管理页面的显示与否；(3) 优化文件管理页面，点击详情显示全部属性；(4) 删除历史控件及相关服务代码。
>
> **完成情况**：✅ 已完成 (2026-04-07)
> - 后端模型与 API 重构已完成
> - 前端属性管理页面已上线
> - 现有页面(分类、文档、首页)适配与清理已完成
> - 历史版本功能全链路移除已完成

---

## 一、当前架构分析

### 1.1 现存问题

| 问题 | 详情 |
|------|------|
| 分类管理页面过载 | `CategoryPage.vue` 同时承载分类树管理 + 表单模板管理（通过 `FormManagerDrawer`），职责不清 |
| 属性管理无独立入口 | 文档属性（FormField）的增删改隐藏在分类管理的抽屉组件里，管理员难以发现和操作 |
| 属性显示控制粒度不足 | 现有 `show_in_list` / `show_in_filter` 无法区分主页和管理页面的显示需求 |
| 文件管理详情缺失 | `AdminDocumentsPage.vue` 无详情弹窗，无法查看文件全部属性 |
| 历史版本功能废弃 | 前后端均保留了历史版本相关代码（`handleHistory`、`GetFileHistory`、`getDocumentHistory`），但功能未完成且将移除 |

### 1.2 目标架构

```
后台管理侧边栏:
├── 仪表盘         (保留)
├── 分类管理       (优化: 纯分类树、去除"表单模板管理"按钮)
├── 文档属性管理   (新增: 独立页面，管理全局属性定义)
├── 文档管理       (增强: 新增详情弹窗、移除历史按钮)
├── OCR 任务       (保留)
├── 回收站         (保留)
├── 用户管理       (保留)
├── 审计日志       (保留)
└── 系统设置       (保留)
```

---

## 二、核心设计决策

### 决策 1：属性管理与分类的关系

**当前方案**：属性通过 `Form（表单模板）→ Category.form_id` 间接绑定到分类。

**新方案**：保留此关联机制不变，但将管理入口从分类页面移至独立的「文档属性管理」页面。分类仍可绑定表单模板，但绑定操作移到新的属性管理页面完成。

### 决策 2：属性显示控制扩展

**当前字段**：`show_in_list`（列表展示）、`show_in_filter`（筛选展示）

**扩展字段**：
- `show_in_home` (bool) — 是否在主页列表展示
- `show_in_admin` (bool) — 是否在管理页面列表展示
- 保留 `show_in_filter` — 筛选展示（主页+管理通用）

> 原有的 `show_in_list` 字段将被废弃，由 `show_in_home` + `show_in_admin` 替代。

### 决策 3：不可变属性（系统固有属性，不参与增删改）

| 字段名 | 说明 | 可编辑 |
|--------|------|--------|
| 所属分类 | category_id | 编辑时可变更 |
| 文档大小 | file_size | ❌ 系统自动 |
| 上传时间 | created_at | ❌ 系统自动 |
| 处理状态 | status (OCR) | ❌ 系统自动 |
| 核验状态 | verify_status | ❌ 系统自动 |
| 同步状态 | sync_status | ❌ 系统自动 |

### 决策 4：历史版本功能 — 全面删除

涉及删除的代码清单：
- 前端：`handleHistory`、`handleShowHistory`、历史版本弹窗、`Timer` 图标引用
- 前端 API：`getDocumentHistory` 函数
- 后端 Handler：`GetFileHistory`
- 后端 Service：`GetFileHistory`
- 后端 Repository：`GetFileHistory`
- 后端 Router：`/documents/history` 路由

---

## 三、后端改动清单

### 3.1 模型层 (`internal/model/`)

#### [MODIFY] `form.go` — FormField 结构体

```
现有字段:
  ShowInList    bool  → 废弃（保留字段，但不再使用）
  ShowInFilter  bool  → 保留

新增字段:
+ ShowInHome   bool   `json:"show_in_home" gorm:"default:true"`   // 主页显示
+ ShowInAdmin  bool   `json:"show_in_admin" gorm:"default:true"`  // 管理页面显示
```

**函数影响**: 无新增函数，仅结构体扩展。GORM AutoMigrate 自动追加列。
**数据迁移**: 需编写一次性迁移逻辑，将现有 `show_in_list=true` 的记录同时设置 `show_in_home=true, show_in_admin=true`。

---

### 3.2 Handler 层 (`internal/handler/`)

#### [MODIFY] `standard.go`

| 函数 | 操作 | 说明 |
|------|------|------|
| `GetFileHistory(c)` | **删除** | 移除历史版本 Handler |
| `GetFileDetail(c)` | **保留** | 已返回 field_values，满足详情弹窗需求 |

#### [MODIFY] `form_handler.go`

| 函数 | 操作 | 说明 |
|------|------|------|
| `BindCategoriesToForm(c)` | **新增** | 处理 `PUT /admin/forms/:id/categories`，接收 `{ category_ids: []uint }` 并调用 Service 层进行批量绑定。 |

---

### 3.3 Service 层 (`internal/service/`)

#### [MODIFY] `standard_service.go`

| 函数 | 操作 | 说明 |
|------|------|------|
| `GetFileHistory(number)` | **删除** | 移除历史版本查询 |

#### [MODIFY] `form_service.go`

| 函数 | 操作 | 说明 |
|------|------|------|
| `BindCategoriesToForm(formID, categoryIDs)` | **新增** | 批量更新分类表中对应的记录的 `form_id` |

---

### 3.4 Repository 层 (`internal/repository/`)

#### [MODIFY] `standard_repo.go`

| 函数 | 操作 | 说明 |
|------|------|------|
| `GetFileHistory(number)` | **删除** | 移除历史版本查询 |

---

### 3.5 Router 层 (`internal/router/router.go`)

#### [MODIFY] `router.go`

```diff
 documents := protected.Group("/documents")
 {
     documents.GET("", standardHandler.ListFiles)
     documents.POST("/upload", standardHandler.UploadFile)
-    documents.GET("/history", standardHandler.GetFileHistory)
     documents.GET("/:id", standardHandler.GetFileDetail)
     ...
 }

 forms := admin.Group("/forms")
 {
     ...
+    forms.PUT("/:id/categories", formHandler.BindCategoriesToForm)
 }
```

---

## 四、前端改动清单

### 4.1 新增页面：文档属性管理 (`views/admin/field-config/`)

#### [NEW] `FieldConfigPage.vue`

独立的文档属性管理页面，将原 `FormManagerDrawer.vue` 的功能重建为全页面版本。

**页面布局**:
```
┌─────────────────────────────────────────────────────┐
│ 页面头部: "文档属性管理" + [新增表单模板] 按钮         │
├──────────────┬──────────────────────────────────────┤
│              │                                      │
│  左栏 280px  │  右栏: 属性字段表格 (可编辑表格)       │
│  表单模板列表 │                                      │
│  - 模板 A    │  属性名 | Key | 类型 | 必填 |         │
│  - 模板 B *  │  主页显示 | 管理显示 | 筛选 |         │
│  - 模板 C    │  默认值 | 扩展配置 | 操作             │
│              │                                      │
│              │  底部: [重置] [保存配置]               │
├──────────────┴──────────────────────────────────────┤
│ 下方区域: 分类绑定关系管理                            │
│ 展示哪些分类绑定了当前选中的模板                      │
│ 支持快速绑定/解绑分类                                │
└─────────────────────────────────────────────────────┘
```

**核心状态与函数**:

| 函数名 | 职责 |
|--------|------|
| `loadForms()` | 加载所有表单模板列表（调用 `getForms()`） |
| `selectForm(form)` | 选中模板，加载其字段列表 |
| `addForm()` | 打开新增模板弹窗 |
| `editFormName(form)` | 编辑模板名称/描述 |
| `handleDeleteForm(form)` | 删除模板（含确认） |
| `addField()` | 新增一行属性字段，`show_in_home: true, show_in_admin: true` |
| `removeField(index)` | 移除一行属性字段 |
| `handleKeyInput(row, val)` | 校验 field_key 输入（字母开头，仅字母数字下划线） |
| `handleSaveFields()` | 校验并保存字段配置（调用 `saveFormFields()`） |
| `loadCategoryBindings()` | 加载使用当前模板的分类列表 |
| `bindCategory(catId)` | 将分类绑定到当前模板 |
| `unbindCategory(catId)` | 解绑分类与当前模板 |

**字段表格列定义**:

| 列名 | 字段 | 宽度 | 控件 |
|------|------|------|------|
| 属性名称 | `label` | 160 | `el-input` |
| 字段标识 | `field_key` | 160 | `el-input` (带校验) |
| 类型 | `field_type` | 130 | `el-select` (input/number/date/select/checkbox) |
| 必填 | `is_required` | 60 | `el-checkbox` |
| 主页显示 | `show_in_home` | 85 | `el-checkbox` |
| 管理显示 | `show_in_admin` | 85 | `el-checkbox` |
| 筛选展示 | `show_in_filter` | 85 | `el-checkbox` |
| 默认值 | `default_value` | 140 | `el-input` / `el-date-picker`(按类型) |
| 扩展配置 | `options` | 180 | `el-input` (select/checkbox 类型才显示) |
| 操作 | — | 70 | `el-button` 移除 |

---

### 4.2 修改页面：分类管理 (`views/admin/category/`)

#### [MODIFY] `CategoryPage.vue`

**删除内容**:
- 移除 `"表单模板管理"` 按钮（第14行 `<el-button :icon="Tickets">` ）
- 移除 `<FormManagerDrawer>` 组件引用和使用（第139行）
- 移除 `import FormManagerDrawer` 语句
- 移除 `drawerVisible` 状态
- 移除 `Tickets` 图标导入
- 移除 `import { getForms, type Form }` 及 `formOptions` 等表单相关状态
- 移除分类弹窗中的"绑定表单"选项（第125-129行）

**保留内容**:
- 分类树展示（el-tree）
- 分类增删改弹窗（name, parent_id, order）
- 树节点操作按钮（编辑、删除）
- 全部 SCSS 样式

**修改的函数**:

| 函数 | 变更 |
|------|------|
| `loadData()` | 移除 `getForms()` 调用，仅加载分类树 |
| `handleAdd()` | 移除 `form_id` 字段初始化 |
| `handleEdit(row)` | 移除 `form_id` 字段赋值 |
| `form` reactive | 移除 `form_id` 字段 |

#### [DELETE] `FormManagerDrawer.vue`

完全删除此文件，其功能迁移至新的 `FieldConfigPage.vue`。

---

### 4.3 修改页面：文档管理 (`views/admin/document/`)

#### [MODIFY] `AdminDocumentsPage.vue`

**新增: 文档详情弹窗**

```vue
<!-- 文档详情弹窗 -->
<el-dialog v-model="detailVisible" title="文档详情" width="650px">
  <div v-if="currentDetailDoc">
    <el-descriptions :column="2" border>
      <!-- 系统固有属性 -->
      <el-descriptions-item label="文档名称" :span="2">
        {{ currentDetailDoc.title }}
      </el-descriptions-item>
      <el-descriptions-item label="所属分类">
        {{ getCategoryName(currentDetailDoc.category_id) }}
      </el-descriptions-item>
      <el-descriptions-item label="文件大小">
        {{ formatFileSize(currentDetailDoc.file_size) }}
      </el-descriptions-item>
      <el-descriptions-item label="上传时间">
        {{ formatDate(currentDetailDoc.created_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="处理状态">
        <el-tag :type="getStatusType(currentDetailDoc.status)">
          {{ getStatusText(currentDetailDoc.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="核验状态">
        <el-tag :type="getVerifyTagType(currentDetailDoc.verify_status)">
          {{ getVerifyStatusText(currentDetailDoc.verify_status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="同步状态">
        <el-tag :type="getSyncStatusType(currentDetailDoc.sync_status)">
          {{ getSyncStatusText(currentDetailDoc.sync_status) }}
        </el-tag>
      </el-descriptions-item>
      <!-- 动态属性值（遍历 field_values） -->
      <el-descriptions-item
        v-for="fv in detailFieldValues"
        :key="fv.field_id"
        :label="fv.field?.label || '属性'"
      >
        <!-- checkbox 类型显示为 tag 组 -->
        <template v-if="isCheckboxField(fv.field_id)">
          <el-tag v-for="tag in fv.value.split(',')" :key="tag" size="small">
            {{ tag }}
          </el-tag>
        </template>
        <template v-else>{{ fv.value || '-' }}</template>
      </el-descriptions-item>
    </el-descriptions>
  </div>
  <template #footer>
    <el-button @click="detailVisible = false">关闭</el-button>
    <el-button type="primary" @click="handlePreview(currentDetailDoc)">
      预览文档
    </el-button>
  </template>
</el-dialog>
```

**新增状态与函数**:

| 名称 | 类型 | 职责 |
|------|------|------|
| `detailVisible` | `ref<boolean>` | 控制详情弹窗显隐 |
| `currentDetailDoc` | `ref<IDocument>` | 当前查看详情的文档 |
| `detailFieldValues` | `ref<FieldValue[]>` | 详情文档的动态属性值列表 |
| `handleShowDetail(row)` | function | 打开详情弹窗，调用 `getDocumentDetail(row.id)` 获取完整数据（含 `field_values`） |
| `isCheckboxField(fieldId)` | function | 判断字段类型是否为 checkbox |

**删除内容 — 历史控件**:

| 删除项 | 位置 |
|--------|------|
| `<el-button :icon="Timer">历史</el-button>` | 第203-205行，操作列 |
| `handleHistory(row)` 函数 | 第707-709行 |
| `Timer` 图标导入 | 第358行 |
| `import { getDocumentHistory }` | 如有引用则移除 |

**修改内容 — 表格列使用新的显示控制字段**:

```typescript
// 替换 displayColumns 的过滤逻辑
// 旧: displayColumns = f.fields.filter(field => field.show_in_list)
// 新: displayColumns = f.fields.filter(field => field.show_in_admin)
```

**新增操作列按钮**:

```vue
<el-button link type="primary" size="small" :icon="InfoFilled"
  @click="handleShowDetail(row)">
  详情
</el-button>
```

---

### 4.4 修改页面：首页 (`views/home/`)

#### [MODIFY] `HomePage.vue`

**删除历史控件**:

| 删除项 | 说明 |
|--------|------|
| 历史版本弹窗 | `<el-dialog>` 历史版本弹窗（第270-289行） |
| `handleShowHistory(row)` | 历史版本查看函数（第490-502行） |
| `historyVisible` / `historyList` / `historyLoading` / `currentHistoryBase` | 历史相关状态变量 |
| 操作列 `历史` 按钮 | `<el-button :icon="Timer">历史</el-button>` |
| `Timer` 图标导入 | |
| `import { getDocumentHistory }` | API 导入 |

**修改表格列显示逻辑**:

```typescript
// 替换 displayColumns 的过滤逻辑
// 旧: displayColumns = f.fields.filter(field => field.show_in_list)
// 新: displayColumns = f.fields.filter(field => field.show_in_home)
```

**优化详情弹窗（已有）**:

现有详情弹窗（第291-323行）结构简单，需增强:
- 展示系统固有属性（文件大小、上传时间、处理状态、核验状态、同步状态）
- 优化布局使用 `el-descriptions` `:column="2"`

---

### 4.5 API 层修改 (`src/api/`)

#### [MODIFY] `form.ts`

```typescript
// FormField 接口扩展
export interface FormField {
  // ... 现有字段
  show_in_list: boolean    // 废弃，保留兼容
  show_in_home: boolean    // 新增：主页显示
  show_in_admin: boolean   // 新增：管理页显示
  show_in_filter: boolean  // 保留
}
```

#### [MODIFY] `document.ts`

```diff
- export function getDocumentHistory(standardNo: string) { ... }
  // 删除历史版本 API 函数
```

#### [MODIFY] `category.ts`

```diff
  // CategoryForm 接口
  export interface CategoryForm {
    name: string
    parent_id: number
    order: number
-   form_id?: number    // 从分类表单中移除（改由属性管理页面处理）
  }
```

> **注意**: `Category` 接口保留 `form_id` 字段（只读展示用），但创建/编辑分类时不再传递。

---

### 4.6 路由配置 (`src/router/index.ts`)

#### [MODIFY] 新增路由

```typescript
{
  path: 'field-config',
  name: 'AdminFieldConfig',
  component: () => import('@/views/admin/field-config/FieldConfigPage.vue'),
  meta: { title: '文档属性管理' }
},
```

---

### 4.7 侧边栏菜单 (`src/components/layout/AdminLayout.vue`)

#### [MODIFY] 菜单项新增

在 "分类管理" 之后新增一项：

```vue
<el-menu-item index="/admin/field-config">
  <el-icon><Setting /></el-icon>
  <span>属性管理</span>
</el-menu-item>
```

---

## 五、删除清单汇总

### 5.1 前端删除

| 文件/代码 | 操作 |
|-----------|------|
| `FormManagerDrawer.vue` | **整文件删除** |
| `CategoryPage.vue` 中表单模板相关代码 | 删除（详见4.2） |
| `AdminDocumentsPage.vue` 历史按钮+函数 | 删除 |
| `HomePage.vue` 历史弹窗+函数+状态 | 删除 |
| `document.ts` → `getDocumentHistory()` | 删除 |

### 5.2 后端删除

| 文件/代码 | 操作 |
|-----------|------|
| `handler/standard.go` → `GetFileHistory()` | 删除 |
| `service/standard_service.go` → `GetFileHistory()` | 删除 |
| `repository/standard_repo.go` → `GetFileHistory()` | 删除 |
| `router/router.go` → `GET /documents/history` 路由 | 删除 |

---

## 六、实施顺序

### 阶段 1：后端模型扩展与历史功能清理
1. 修改 `model/form.go` — FormField 新增 `ShowInHome`, `ShowInAdmin`
2. 删除 `handler/standard.go` → `GetFileHistory`
3. 删除 `service/standard_service.go` → `GetFileHistory`
4. 删除 `repository/standard_repo.go` → `GetFileHistory`
5. 删除 `router/router.go` → `/documents/history` 路由
6. 新增 `service/form_service.go` 和 `handler/form_handler.go` 中关于 `BindCategoriesToForm` 的代码逻辑和 API

### 阶段 2：前端 API 层与路由
7. 修改 `api/form.ts` — FormField 接口扩展
8. 删除 `api/document.ts` → `getDocumentHistory()`
9. 修改 `router/index.ts` — 新增 `field-config` 路由

### 阶段 3：新建文档属性管理页面
10. 新建 `views/admin/field-config/FieldConfigPage.vue`
11. 修改 `AdminLayout.vue` 侧边栏 — 新增"属性管理"菜单项

### 阶段 4：重构现有页面
12. 修改 `CategoryPage.vue` — 去除表单模板管理，精简
13. 删除 `FormManagerDrawer.vue`
14. 修改 `AdminDocumentsPage.vue` — 新增详情弹窗、删除历史控件、使用 `show_in_admin`
15. 修改 `HomePage.vue` — 删除历史功能、使用 `show_in_home`、增强详情弹窗

### 阶段 5：验证
16. `npm run build` 验证前端编译
17. `go build ./...` 验证后端编译
18. 全链路功能测试

---

## 七、文件级变更索引

### 后端 (test-ebook-api)

| 文件 | 操作 | 改动函数 |
|------|------|----------|
| `internal/model/form.go` | MODIFY | FormField 新增 ShowInHome, ShowInAdmin |
| `internal/handler/standard.go` | MODIFY | 删除 GetFileHistory |
| `internal/handler/form_handler.go` | MODIFY | 新增 BindCategoriesToForm |
| `internal/service/standard_service.go` | MODIFY | 删除 GetFileHistory |
| `internal/service/form_service.go` | MODIFY | 新增 BindCategoriesToForm |
| `internal/repository/standard_repo.go` | MODIFY | 删除 GetFileHistory |
| `internal/router/router.go` | MODIFY | 删除 history 路由, 新增 forms/:id/categories 路由 |

### 前端 (test-ebook-web)

| 文件 | 操作 | 关键改动 |
|------|------|----------|
| `src/api/form.ts` | MODIFY | FormField 接口扩展 |
| `src/api/document.ts` | MODIFY | 删除 getDocumentHistory |
| `src/api/category.ts` | MODIFY | CategoryForm 移除 form_id |
| `src/router/index.ts` | MODIFY | 新增 field-config 路由 |
| `src/views/admin/field-config/FieldConfigPage.vue` | **NEW** | 独立属性管理页面 |
| `src/views/admin/category/CategoryPage.vue` | MODIFY | 去除表单相关，精简 |
| `src/views/admin/category/FormManagerDrawer.vue` | **DELETE** | 功能迁移至新页面 |
| `src/views/admin/document/AdminDocumentsPage.vue` | MODIFY | +详情弹窗, -历史控件 |
| `src/views/home/HomePage.vue` | MODIFY | -历史功能, 增强详情弹窗 |
| `src/components/layout/AdminLayout.vue` | MODIFY | 侧边栏新增菜单项 |

---

## 八、风险与注意事项

| 风险 | 应对 |
|------|------|
| `show_in_list` 废弃后旧数据兼容 | 后端迁移时将 `show_in_list=true` 同步设置到 `show_in_home=true, show_in_admin=true` |
| 分类移除 form_id 绑定入口 | 在新的属性管理页面提供分类绑定功能 |
| FormManagerDrawer 删除后 CategoryPage 报错 | 确保先删除所有引用再删文件 |
| 历史版本路由删除后前端 404 | 确保前端所有 `getDocumentHistory` 引用全部清除 |

---

**本 PLAN 生成完毕，等待 xixu520 确认后由 Gemini Flash 执行。**
