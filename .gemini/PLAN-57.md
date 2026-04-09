# PLAN-57：全局属性重构 + 主页/文档管理页面 UI 优化

> **目标**：(1) 将属性体系从"表单模板-分类绑定"架构重构为全局属性架构；(2) 属性管理页面精简重构；(3) 主页 UI 优化与功能梳理；(4) 文档管理页面 UI 优化 + 文档名称点击详情。
>
> **状态**：✅ 已完成 (2026-04-09)

---

## 一、现状分析

### 1.1 当前架构问题（属性体系）

当前的属性体系是 **"表单模板 → 分类绑定 → 文档属性"** 三层结构：

```
Form（表单模板）
  └── FormField[]（属性字段）
       └── 通过 form_id 绑定到 Category → 文档继承
```

**问题清单：**

| 问题 | 说明 |
|------|------|
| 架构复杂 | 用户需先创建模板，再绑定分类，操作链路长 |
| 属性非全局 | 同一属性需在多个模板中重复定义 |
| `field_key` 需用户手动填写 | 用户体验差，容易出错 |
| `show_in_filter`、`show_in_admin` 等字段 | 与新需求冲突，需精简 |

### 1.2 新架构目标

```
GlobalField（全局属性）  ← 属性管理页面统一管理
  └── DocumentFieldValue[]（每个文档的属性值）
```

所有文档共享同一套全局属性定义，废弃"表单模板"和"分类绑定"概念。

---

## 二、需要决策的问题（请 xixu520 选择）

> [!IMPORTANT]
> 以下 5 个问题需要你明确选择后，我才开始修改代码。

### ❓ 决策 1：后端架构重构深度

**选项 A（推荐）：复用现有 Form/FormField 表，引入"全局唯一 Form"约定**

- 后端数据库表不变（Form、FormField、DocumentFieldValue 结构不动）
- 约定只保留一个 Form 作为系统全局属性集
- 前端属性管理页面不再显示模板列表，直接管理该 Form 下的所有字段
- **优点**：改动最小，无需数据库迁移，后端改动量极少
- **缺点**：历史遗留多个 Form 需初始化逻辑处理

**选项 B（彻底重构）：新建 GlobalField 数据库表**

- 新建 `global_fields` 表，废弃 Form 体系
- **优点**：架构语义清晰
- **缺点**：需数据库迁移、历史数据迁移，改动量大，风险高

> **建议选 A**

---

### ❓ 决策 2：`field_key` 的生成方式

当前用户需要手动填写 `field_key`（如 `sign_date`），这对普通用户不友好。

**选项 A（推荐）：前端自动根据时间戳生成，用户不可见**

- 用户只填属性名称
- 前端自动生成 `field_key = "field_" + timestamp`

**选项 B：后端保存时自动生成**

- 前端不展示 `field_key` 输入框
- 后端在创建/更新时若 field_key 为空，自动生成 `field_<id>`

**选项 C：用中文 label 直接作为 field_key**

- 去掉 field_key 概念，后端需改查询逻辑

> **建议选 A**

---

### ❓ 决策 3：属性排序功能

属性管理页面是否需要支持**拖拽排序**？

- **选项 A（推荐）**：不支持拖拽，列表顺序固定
- **选项 B**：支持表格内拖拽排序（需引入 `sortablejs` 库）

> **建议选 A（暂不实现，后续可加）**

---

### ❓ 决策 4：主页布局重构方向

当前问题：动态展示列依赖"分类→form_id→fields"链路，改为全局属性后需修复。

**选项 A：侧边筛选抽屉**

- 精简上方区域，只有搜索框和"筛选"按钮
- 点击筛选弹出抽屉，内有分类树 + 全局属性筛选项

**选项 B（推荐）：保持当前行内筛选布局，修复数据链路**

- 筛选区域保持行内展示
- 改为加载全局属性 `show_in_home = true` 的字段
- 做 UI 微调，不大改布局

> **建议选 B**

---

### ❓ 决策 5："筛选展示"开关是否保留

每个属性目前有 `show_in_home`、`show_in_admin`、`show_in_filter` 三个开关。  
你要求属性管理只需"是否在首页展示"一个开关，那么：

- **选项 A**：只保留 `show_in_home`，废弃筛选功能
- **选项 B**：保留 `show_in_home` 和 `show_in_filter` 两个开关
- **选项 C（推荐）**：只保留 `show_in_home`，"首页展示"的属性自动也可以作为筛选项

> **建议选 C**

---

## 三、确定要做的改动（无分歧）

### 3.1 后端改动

| 文件 | 改动 | 说明 |
|------|------|------|
| `internal/service/form_service.go` | MODIFY | 新增 `GetOrCreateGlobalForm()` |
| `internal/handler/form_handler.go` | MODIFY | 新增 `GetGlobalForm` handler |
| `internal/router/router.go` | MODIFY | 新增 `GET /admin/forms/global` 路由 |

### 3.2 属性管理页面（`FieldConfigPage.vue`）— 完全重构

| 删除 | 保留/新增 |
|------|---------|
| 左侧模板列表 | 顶部页面标题 + 新增属性按钮 |
| 分类绑定 Tab | 属性字段表格 |
| 模板新增/编辑/删除弹窗 | 系统固有属性只读展示区 |
| `field_key` 输入列 | "首页展示"开关 |
| `show_in_admin`、`show_in_list` | "必填"开关 |
| | "类型"、"默认值"、"选项"输入 |

**系统固有属性（只读灰色区域，顶部展示）：**
所属分类 | 文件大小 | 上传时间 | 处理状态 | 核验状态 | 同步状态

### 3.3 主页（`HomePage.vue`）

| 问题 | 改动方案 |
|------|------|
| 动态列加载依赖分类form链路 | 改为直接从全局属性接口获取 `show_in_home = true` |
| 筛选字段来源 | 同步修复（取决于决策5的选择） |
| 文档名称点击详情 | 已实现 ✅ 保留 |
| UI 微调 | 调整间距和细节 |

### 3.4 文档管理（`AdminDocumentsPage.vue`）

| 问题 | 改动方案 |
|------|------|
| 标题列不可点击 | 改为 el-button link 点击触发详情 |
| 操作栏"详情"按钮 | 删除，由标题点击代替 |
| 动态展示列逻辑 | 改为全局属性 `show_in_home = true`（和主页同步） |
| 编辑弹窗动态字段来源 | 改为全局属性 |

---

## 四、文件级变更索引

### 后端 (test-ebook-api)

| 文件 | 操作 |
|------|------|
| `internal/service/form_service.go` | MODIFY |
| `internal/handler/form_handler.go` | MODIFY |
| `internal/router/router.go` | MODIFY |

### 前端 (test-ebook-web)

| 文件 | 操作 |
|------|------|
| `src/api/form.ts` | MODIFY |
| `src/views/admin/field-config/FieldConfigPage.vue` | FULL REWRITE |
| `src/views/home/HomePage.vue` | MODIFY |
| `src/views/admin/document/AdminDocumentsPage.vue` | MODIFY |

---

**等待 xixu520 对以上 5 个决策的回复，然后开始执行。**
