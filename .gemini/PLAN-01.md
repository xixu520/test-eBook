# PLAN-01: 前端项目初始化与 Mock 环境搭建

## 需求描述
根据“前端架构书”和“API 文档”，初始化 Vue 3 + Vite + TypeScript 前端项目，并配置 Mock 拦截工具以支持假数据开发。

## 变更记录
- [x] 初始化 Vite 项目 (Vue, TS, Element Plus)
- [x] 配置 Mock 拦截 (vite-plugin-mock)
- [x] 封装 Axios 请求与响应拦截器
- [x] 创建基础目录结构 (src/api, src/views, src/stores 等)

## 完成情况
- [x] 项目脚手架已生成
- [x] 基础依赖包已安装
- [x] Mock 示例接口可通行
- [x] 登录页基础结构完成

## 第二阶段：首页工作台与布局 [x]
- [x] 创建 `MainLayout` 布局 (TopNavBar, SideMenu)
- [x] 实现 `HomePage` 搜索筛选与批量操作栏
- [x] 实现 `HomePage` 文档列表表格与分页
- [x] Phase 8 | 文档版本管理 | 同一标准号下多个版本的上传与展示对比 | 已完成 |
- [x] Phase 9 | 搜索增强 | 引入全文搜索模拟与高级筛选 (发布机构、实施状态) | 已完成 |
- [x] 丰富 Mock 数据 (分类树、文档模拟列表)
- [x] 实现主题切换与公告展示逻辑

## 第三阶段：管理后台 [x]
- [x] 实现 `AdminLayout` 后台整体布局
- [x] 实现 `DashboardPage` 数据统计卡片与动态时间线
- [x] 实现 `CategoryPage` 分类树编辑与管理
- [x] 封装管理员专用 API 与 Mock 接口

## 第四阶段：PDF 预览 [x]
- [x] 集成 `pdfjs-dist` 库
- [x] 开发 `PdfPreview` 通用组件 (支持翻页/缩放)
- [x] 在工作台集成预览弹窗逻辑

## 第五阶段：上传与 OCR 进度 [x]
- [x] 实现支持拖拽与进度显示的上传对话框
- [x] 模拟大文件分片上传与任务创建接口
- [x] 实现管理员 OCR 任务队列管理页面
- [x] 丰富文档处理状态 (`Pending`, `Processing`, `Success`, `Error`) 的 UI 表现

## 第六阶段：系统设置与高级管理 [x]
- [x] 开发分栏式的系统设置页面 (`SettingsPage`)
- [x] 实现 OCR 引擎配置与存储路径管理
- [x] 实现系统资源 (CPU/内存/磁盘) 状态监控可视化
- [x] 封装设置相关 API 与 Mock 接口

## 第七阶段：用户权限管理 (RBAC) [x]
- [x] 开发用户管理页面 (`用户管理页面`)
- [x] 实现角色的细粒度分配与状态管理
- [x] 完善路由拦截与组件级动态鉴权逻辑
- [x] 扩展用户管理相关的 Mock 接口

---

# PLAN-02: 前端架构书与部署方案设计

## 目标
基于 `前端UI设计调查表.md`、`test eBook - API.md` 和 `PLAN-01` 的成果，设计一份完整的前端架构书与 Docker 容器化部署方案，精确到每个控件的实现方案。

## 状态

- [x] 阅读所有源文件，梳理需求与设计决策
- [x] 撰写前端架构书（技术选型、目录结构、组件清单、控件实现方案）
- [x] 撰写 Docker 容器化部署方案（已合并至架构书第十一章）
- [ ] 用户审阅并确认架构书与部署方案

## 涉及交付物

| 文件 | 说明 |
|------|------|
| `前端架构书.md` | 完整的前端技术架构设计文档 |
| `.gemini/PLAN-02.md` | 本计划文件 |

## 设计依据

- `前端UI设计调查表.md` — 用户 UI 偏好与需求（已确认）
- `test eBook - API.md` — 后端 API 接口清单（22+ 端点，10 个模块）
- `PLAN-01` — UI 设计调查表的完成记录

---

# PLAN-03: 后端架构书设计

## 目标
基于 `test eBook - API.md`、`前端UI设计调查表.md` 和 `前端架构书.md`，设计一份极其详细的后端架构书，精确到每个模块、每个接口的实现方式。

## 状态

- [x] 阅读所有源文件，梳理后端需求
- [x] 撰写后端架构书（技术选型、目录结构、数据库设计、每模块实现方案）
- [ ] 用户审阅并确认

## 涉及交付物

| 文件 | 说明 |
|------|------|
| `后端架构书.md` | 完整的后端技术架构设计文档 |
| `.gemini/PLAN-03.md` | 本计划文件 |

---

# PLAN-04: 前端多端响应式优化

## 目标
实现“建筑标准文件管理系统”前端页面的多端（手机、平板、桌面）适配，确保在不同屏幕尺寸下均有出色的用户体验。

## 状态
- [x] 定义全局响应式断点与 Mixins
- [x] 优化 `TopNavBar`（移动端隐藏标题、搜索框适配）
- [x] 优化 `SideMenu`（移动端改为 Drawer 模式，平板端自动折叠）
- [x] 优化 `MainLayout` / `AdminLayout` 布局逻辑
- [x] 优化 `HomePage`（表格滚动、筛选区折叠）
- [x] 优化 `LoginPage` 卡片适配
- [x] 全设备模拟验证

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `test-ebook-web/src/styles/variables.scss` | 全局样式变量与 Mixins |
| `test-ebook-web/src/components/layout/*` | 布局相关组件 |
| `test-ebook-web/src/views/*` | 主要页面组件 |
| `.gemini/PLAN-04.md` | 本计划文件 |

---

# PLAN-05: 后台管理增强（文档库与回收站）

## 目标
完善系统后台管理功能，提供专门的文档库管理页面、回收站功能以及审计日志，实现更专业的内容管理能力。

## 状态
- [x] 注册 `/admin/documents`, `/admin/recycle`, `/admin/audit-logs` 路由
- [x] 实现 `AdminDocumentsPage.vue`（含上传者/下载统计等）
- [x] 实现 `RecycleBinPage.vue`（支持还原与彻底删除）
- [x] 实现 `AuditLogPage.vue`（操作日志展示）
- [x] 优化 `AdminLayout` 侧边栏菜单跳转
- [x] 完善 Mock 接口数据支持

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/document/AdminDocumentsPage.vue` | 后台文档管理主页 |
| `src/views/admin/recycle/RecycleBinPage.vue` | 回收站页面 |
| `src/views/admin/audit/AuditLogPage.vue` | 审计日志页面 |
| `.gemini/PLAN-05.md` | 本计划文件 |

## 设计依据
- `test eBook - API.md` — 模块 6 (回收站) 与 模块 7 (审计日志)
- `前端架构书.md` — 5.6, 5.8 章节建议

---

# PLAN-06: 核心业务逻辑修复与安全加固

## 目标
解决在前后端架构审查中发现的 4 个核心业务逻辑与权限安全防范问题。

## 状态
- [x] 修复 1.1：(分类删除) 拦截后端 400 错误码并给予分类删除的针对性提示（在实现 CategoryPage 时应用与 Mock）。
- [x] 修复 1.2：(回收站) 严格限制 `empty_all: true` 的触发场景，避免批量删除误触。
- [x] 修复 1.3：(权限冗余) 剥离 `HomePage.vue` 中的高危管理权限（删除、编辑、重新 OCR），交由后台统一闭环。
- [x] 修复 2.1：(OCR 状态隔离) 文件若处于 `pending/processing` 阶段，前端需强禁用“预览”及“历史版本”功能，防止查空或崩溃。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/home/HomePage.vue` | 移除管理按钮，增加 OCR 状态前端校验机制 |
| `src/views/admin/document/AdminDocumentsPage.vue` | 增加 OCR 状态关联拦截 |
| `src/views/admin/recycle/RecycleBinPage.vue` | 审查清空与批量删除的传参逻辑 |
| `mock/document.ts` | 增强对删除清空的 Mock 拦截验证 |
| `mock/category.ts` | 增加带有子分类或文档关联时的模拟拦截 |

---

# PLAN-07: 上传控件体验与逻辑修补

## 目标
修复 xixu520 反馈的上传标准文件控件中的样式与业务逻辑问题。

## 状态
- [x] 修复 1.1: 版本/年份字体错位问题（CSS 布局调整）。
- [x] 修复 1.2: 所属分类选择无效问题（数据绑定与 `el-tree-select` 逻辑校验）。
- [x] 修复 1.3: 标签宽度过窄导致的换行重叠问题（UI 优化）。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/components/document/UploadDialog.vue` | 核心修复文件 |

---

# PLAN-08: 完善后台 OCR 引擎设置

## 目标
检查并完善后台管理系统中的“系统设置”页面，特别是围绕“OCR引擎”与“百度云 OCR”的配置项，保证设置界面的全功能覆盖与 Mock 数据支撑。

## 状态
- [x] 1. 审查 `SettingsPage.vue`，评估现有 OCR 表单是否覆盖百度云的必要配置（如 API Key, Secret Key 等）。
- [x] 2. 完善 `SettingsPage.vue` 的界面与交互（密码遮盖、连接测试等）。
- [x] 3. 补充或完善 `src/api/settings.ts` 以及 `mock/settings.ts` 的接口，打通前端状态流转。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 设置界面核心代码 |
| `src/api/settings.ts` | 接口调用封装 |
| `mock/settings.ts` | 模拟后端接口返回 |

---

# PLAN-09: 完善 PaddleOCR 相关功能设置

## 目标
补全针对 PaddleOCR 引擎的专属后台配置项，完善其联调与配置交互闭环。

## 状态
- [x] 1. 在 `SettingsPage.vue` 添加针对 PaddleOCR 引擎特有的官方 API 配置（API Key 和 Secret Key）。
- [x] 2. 在原有测试连接框架中，打通 PaddleOCR 官方接口配置的存活性联调测试按钮。
- [x] 3. 补齐 `mock/settings.ts` 数据结构，使得前端可以正常回显并校验模拟的 PaddleOCR 官方调用入参。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 增加 PaddleOCR 专属配置输入项与验证 |
| `mock/settings.ts` | 添加 PaddleOCR `server_url` 等相关数据 Mock 支持 |

---

# PLAN-10: 适配 PaddleOCR 官方异步 API 参数

## 目标
根据官方 Python 异步解析脚本分析，重构目前预留的 PaddleOCR 官方接口设置参数形式。从原本假定的 API/Secret Key 形式转换成为原生符合要求的单 `Token` 形式，并暴露官方特有的配置项模型设置与图像优化设置。

## 状态
- [x] 1. 将鉴权参数由 `API Key` + `Secret Key` 简化修改为 `Token` 认证。
- [x] 2. 增加 `Model` 字段配置项（默认：`PaddleOCR-VL-1.5`）。
- [x] 3. 增加布尔值高级可选参数（文档方向分类、文档去畸变、图表识别）作为可选开关选项。
- [x] 4. 同步修正 `mock/settings.ts` 数据结构和测试鉴权规则（已增加对特定合法 Token 的等值校验）。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `src/views/admin/settings/SettingsPage.vue` | 移除旧字段，引入 Token 及 OCR 选型表单 |
| `mock/settings.ts` | 更新保存和测试接口 Mock |

---

# PLAN-11: 初始化后端 Go 工程架构

## 目标
根据《后端架构书.md》初始化 Go 后端项目 `test-ebook-api`，搭建核心分层目录结构，配置 Gin 框架基础路由与中间件，并完成 SQLite 数据库与 GORM 的初步集成。

## 状态
- [x] 1. 创建项目目录并执行 `go mod init test-ebook-api`。
- [x] 2. 搭建分层目录结构（cmd, internal/handler, service, repo, model...）。
- [x] 3. 初始化核心配置加载逻辑 (`internal/config`) 与统一响应工具 (`internal/pkg/response`)。
- [x] 4. 初始化数据库连接与 GORM 集成。
- [x] 5. 实现基础路由骨架与中间件 (CORS, Logger, JWT)。
- [x] 6. 完成用户登录验证与 Seed 管理员账号机制。

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `test-ebook-api/go.mod` | 依赖管理 |
| `test-ebook-api/cmd/server/main.go` | 工程入口 |
| `test-ebook-api/internal/config/config.go` | 配置管理 |
| `test-ebook-api/internal/database/sqlite.go` | 数据库驱动配置 |
| `test-ebook-api/internal/router/router.go` | 路由分发 |

---

# PLAN-15: 前后端 Docker 容器化与部署实践

你好 xixu520。本计划旨在将项目全面容器化，实现一键启动与环境一致性。

## 目标
- [ ] 为 Go 后端 `test-ebook-api` 编写多阶段构建 Dockerfile
- [ ] 为 Vue 前端 `test-ebook-web` 编写 Nginx 托管 Dockerfile
- [ ] 编写 `docker-compose.yml` 实现服务编排
- [ ] 处理容器间的网络通信与环境变量注入
- [ ] 验证容器化部署后的全链路功能（包括 OCR 与文件上传）

## 详细步骤

### 1. 后端容器化 (test-ebook-api)
- **Dockerfile**: 使用 `golang:1.21-alpine` 作为构建环境，`alpine:latest` 作为运行环境。
- **配置处理**: 将 `config.yaml` 映射为容器卷或内置默认配置。
- **存储映射**: 将 `uploads/` 和 `data/` (SQLite 数据文件) 映射为宿主机持久化卷。

### 2. 前端容器化 (test-ebook-web)
- **Dockerfile**: 使用 `node:18-alpine` 进行构建，`nginx:alpine` 提供 Web 服务。
- **Nginx 配置**: 处理 SPA 路由（history mode）及反向代理（可选，或由独立网关处理）。
- **BaseURL 注入**: 通过构建参数或运行时环境替换方案解决 `VITE_API_BASE_URL` 的动态指向。

### 3. 服务编排 (docker-compose)
- 定义 `backend` 服务：暴露端口，挂载卷。
- 定义 `frontend` 服务：暴露 3000 或 80 端口，链接至 backend。
- 配置网络：确保 frontend 容器或浏览器能访问至 backend。

### 4. 部署验证
- 执行 `docker-compose up --build`。
- 测试登录、上传、OCR 轮询及分类管理。

## 注意事项
- [!IMPORTANT]
- 需要处理 SQLite 数据库文件权限问题。
- 前端 `VITE_API_BASE_URL` 如果在浏览器端运行，需指向宿主机 IP 或域名，而非容器局域网 IP。

---

# PLAN-16: 整个项目全面审查与优化建议

你好 xixu520。经过对整个项目的审查，我总结了以下几个方面的优化建议。由于你要求不要直接修改代码，这些建议仅作为未来的优化方向：

## 1. Docker 与部署优化
- **后端 Dockerfile (`test-ebook-api/Dockerfile`)**:
  - 当前开启了 `CGO_ENABLED=1`，但项目实际使用的是纯 Go 版本的 SQLite 驱动 (`github.com/glebarez/sqlite`)，因此完全可以关闭 CGO (`CGO_ENABLED=0`)。
  - 关闭 CGO 后，构建阶段不需要额外安装 `gcc` 和 `musl-dev`。
  - 在 `go build` 时增加 `-ldflags="-s -w" -trimpath` 参数，可以大幅减小最终编译出的二进制文件体积，并提升部署效率。
- **前端 Dockerfile (`test-ebook-web/Dockerfile`)**:
  - 当前的 Dockerfile 缺少了**构建阶段 (build stage)**，它直接将宿主机现成的 `dist` 目录复制进容器。这要求每次构建 Docker 镜像前必须在本地或宿主机手动执行 `npm run build`。建议使用 Node 多阶段构建（先利用 Node.js 镜像执行打包，然后再把产物放入 Nginx 镜像中），实现构建过程的自动化与环境一致性。

## 2. 后端 API 层面 (test-ebook-api)
- **优雅停机 (Graceful Shutdown)**:
  - `cmd/server/main.go` 中的服务没有处理系统退出信号（如 `SIGINT` 和 `SIGTERM`）。在生产环境中更新或重启容器时，这可能导致正在处理的用户请求和文件被强制中断。建议补充优雅停机机制，在关闭前确保处理完既有请求，并安全释放数据库连接。
- **配置文件的环境变量支持**:
  - `config.yaml` 中写死了许多敏感和会因环境变化的参数（例如 JWT Secret、Ocr Token 等）。建议完善配置读取逻辑，让其天然支持环境变量覆盖（如利用 Viper 的 `AutomaticEnv`），这对容器化（Docker Compose）注入配置尤为重要。

## 3. 前端 Web 层面 (test-ebook-web)
- **产物体积优化与分块压缩**:
  - 建议在 `vite.config.ts` 中配置 `build.rollupOptions.output.manualChunks`，对类似 `element-plus`、`vue` 这种庞大且变更频率极低的基础库依赖进行独立分包。这样可以利用浏览器的长效缓存机制。
- **资源预压缩 (Gzip/Brotli)**:
  - 推荐引入如 `vite-plugin-compression`，在打包时对静态资源实施预压缩，从而让生产环境下的 Nginx 不必实时耗费 CPU 压包，并加快用户的首次访问速度。

---
**完成情况 (2026-03-23)**：
- [x] **Docker 构建优化**：已在 `test-ebook-api/Dockerfile` 中关闭 CGO 并去除了 C 编译依赖，加入 `-ldflags="-s -w" -trimpath` 优化体积。为 `test-ebook-web/Dockerfile` 添加了基于 `node:18-alpine` 的构建阶段。
- [x] **后端 API 优化**：已在 `cmd/server/main.go` 中引入协程启动服务器，并加入了全面的 `syscall.SIGINT`、`syscall.SIGTERM` 拦截与五秒超时的 HTTP/DB 优雅关闭逻辑。内置 Config 本身已通过 Viper 支持环境变量注入。
- [x] **前端 Web 优化**：已在 `test-ebook-web/vite.config.ts` 加入 `manualChunks` 进行 vendor 拆分，并通过安装并配置 `vite-plugin-compression` 支持构建期的 Gzip 资源预压缩。

---

# PLAN-17: 编写全项目 README 说明文档

你好 xixu520。本计划旨在响应你对于编写及总结前后端功能、并生成完整的 `README.md` 的操作记录。

## 目标与改动
- [x] **架构梳理**：系统性读取并梳理现有《前端架构书.md》、《后端架构书.md》等文档核心业务逻辑与技术底座。
- [x] **撰写 `README.md`**：完成整个项目的展示门面编排。包含了以下版块：
  - 项目核心特性汇总（分类管理、自动化 OCR、高可用管理后台）。
  - 现代化技术选型清单（前端 Vue3/Vite，后端 Go/Gin/SQLite）。
  - 各类一键式 Docker Compose 部署步骤以及数据层持久化规约。
  - 本地环境基础协同调试命令。

## 当前状态
- **实施情况**：已成功在项目根目录重写并覆盖了之前的空 `README.md` 文件。
- **提醒**：所有配置步骤和特性声明均基于我们在 `PLAN-15` 到 `PLAN-16` 对整个框架代码执行过的优化加固工作。

**后续行动建议**：
由于我刚才未进行 git 提交，README 的更新已保存在本地磁盘。你可以根据需要使用 VSCode 对 `README.md` 进行细微调整或执行 `git commit` 将它推送至代码仓库！

---

# PLAN-18: v0.5.0 发布版本并打包至 GitHub

你好 xixu520。本计划旨在将迄今为止的所有优化、文档完善工作进行版本固化，并打包发布 `v0.5.0` Release 至 GitHub。

## 目标与步骤
- [x] **版本汇总**：将涉及后端优雅停机、Docker 构建优化、前端打包优化的相关代码，以及全栈架构的核心说明 `README.md` 和各 PLAN 规划文档变更整体进行保存。
- [x] **Git 提交推送**：通过 `git commit` 固化当前代码库中未提交的所有变更，并推送到远程的分支中。
- [x] **发行版创建**：利用 Git 命令挂上 `v0.5.0` 稳定版标签，同步推送到远端。
- [x] **GitHub 打包发布**：调用 GitHub 官方的命令行工具 `gh` 自动为你直接发布 Release，并附带本次的核心更新日志。

---
**完成情况 (2026-03-23)**：
- 已成功把整个工作区更新推送至远端。
- 已为您打上 `v0.5.0` 的版本 tag 并成功通过 `gh` CLI 发布了 GitHub Release！

---

# PLAN-19: 前后端功能与数据一致性检查

## 目标
检查前后端各项功能，发现并记录互相矛盾的地方（如 API 路径不一致，请求参数、响应数据格式不一致，或是功能逻辑冲突）。

## 发现的问题清单
通过对比发现，`test eBook - API.md` 文档、前端请求定义 (`test-ebook-web/src/api/*.ts`) 以及后端路由实现 (`test-ebook-api/internal/router/router.go`) 之间存在大量不一致或遗漏的地方：

### 1. 分类管理模块 (Category)
- **前端定义**: 
  - `GET /standards/categories`
  - `POST /standards/categories`
  - `DELETE /standards/categories/:id`
- **后端实现**:
  - `GET /standards/categories` (不符合API文档的 `/categories`)
  - `POST /standards/categories`
  - 缺失 `DELETE` 路由，但在文件模块里错误地写成了 `standards.DELETE("/files/:id", standardHandler.DeleteFile)`
- **API文档**: 规定为 `/categories`，并有 `PUT` 接口但在前后端均未实现完全。

### 2. 标准文件模块 (Documents/Files)
- **前端定义**: `GET /standards/files`, `POST /standards/files`, `GET /standards/files/:id`, `DELETE /standards/files/:id`
- **后端实现**: 与前端调用一致，均为 `/standards/files` 相关。
- **API文档**: 规定路径为 `/documents`，上传为 `/documents/upload`。
- **缺失项**: 前端调用了 `GET /documents/history`, `GET /recycle-bin/documents`, `PUT /recycle-bin/documents/restore`, `POST /recycle-bin/documents/batch-delete`，但后端路由 **完全未实现** 历史记录和回收站模块。

### 3. 上传与 OCR 模块
- **前端定义**: `api/upload.ts` 定义了 `POST /upload`, `GET /ocr/tasks`，但 `document.ts` 中也定义了 `POST /standards/files` 作为上传。
- **后端实现**: 实现了 `POST /standards/files`。
- **API文档**: 规定上传为 `/documents/upload`，获取状态为 `/tasks/{task_id}/status`，重试为 `/documents/{id}/ocr/retry`。前后端与文档均不一致，且后端未实现 OCR 任务状态查询。

### 4. 日志、设置、用户管理 (后台功能)
- **前端定义**: 
  - `GET /admin/users` 及相关增删改
  - `GET /audit-logs`
  - `GET /settings`, `POST /settings`, `GET /system/status`, `POST /settings/test-ocr`
- **后端实现**: **完全没有** 实现上述任何路由。
- **API文档**: 部分存在（如设置和日志），但用户管理 `/admin/users` 在文档中也未提及。

### 5. 看板数据 (Dashboard)
- **前端与后端一致**: 使用 `/stats/dashboard`。
- **API文档**: 规定为 `/admin/dashboard`。

### 6. 用户认证与个人信息 (Auth)
- **前端调用**: `PUT /users/me/theme`
- **后端实现**: 缺失该更新主题的路由。

## 总结与建议
当前项目的前后端 API 完全脱节。
1. **重新对齐 API 文档**：建议先对 `test eBook - API.md` 进行修订和确认，确定最终的接口规范。
2. **重构后端路由**：按照规范重写 `test-ebook-api/internal/router/router.go` 及相关 Handler，补齐缺失的路由（特别是回收站、OCR查询、系统设置、审查日志、用户管理等）。
3. **修复前端请求**：更新 `src/api` 下的文件，让 URL 路径与文档保持绝对一致，消除重复或矛盾的函数。

---

# PLAN-20: 分类管理模块对齐与实现

## 目标
1. 修正前后端 `Category` 模块的 API 路径不一致问题。
2. 在后端补齐缺失的 `UpdateCategory` (更新) 和 `DeleteCategory` (删除) 功能。
3. 在前端 `src/api/category.ts` 中同步更新接口定义并增加更新接口。

## 需求详情
### 1. 路径统一
- 将所有分类相关的 API 路径从 `/api/v1/standards/categories` 统一为 `/api/v1/categories`。

### 2. 后端功能补全
- **Repository**:
  - 实现 `UpdateCategory`：更新分类名称、父ID、排序值。
  - 实现 `DeleteCategory`：物理或软删除分类（根据模型定义，目前是带 GORM Model 的软删除）。
  - 实现校验：检查是否存在子分类，检查是否存在关联文件。
- **Service**:
  - 封装带业务逻辑的更新和删除方法。
- **Handler**:
  - 增加 `PUT /categories/:id` 处理函数。
  - 增加 `DELETE /categories/:id` 处理函数。

### 3. 前端代码修正
- 修改 `src/api/category.ts` 中的 `url`。
- 增加 `updateCategory` 函数。

## 进度追踪
- [x] 后端 Repository 层实现 `Update` 和 `Delete` 相关方法
- [x] 后端 Service 层实现业务逻辑与校验
- [x] 后端 Handler 层新增接口处理函数
- [x] 后端 Router 层路径重构与新接口注册
- [x] 前端 API 定义更新 (`category.ts`)
- [x] 前端页面逻辑微调（如移除“编辑功能待后端支持”的提示）

---

# PLAN-21: 系统前后端接口规范全量对齐与实现方案

基于 `PLAN-19` 发现的系统性脱节问题，以及 `PLAN-20`（已完成分类模块对齐）的经验，现针对剩余模块制定精确到代码层面的修改与实现方案。

## 整改目标
彻底消除前后端及《API 文档》之间的接口路径、参数和缺失功能的矛盾。以《API 文档》和合理的业务需求为最终基准。

---

## 模块修改明细表

### 1. 标准文件模块 (Documents)
**目标**：将所有文件管理路径对齐到 `/documents`，并补全回收站和历史记录功能。
**后端实现方案**：
- **Router (`router.go`)**: 
  - 删除原有的 `v1.Group("/standards")` 中涉及 `/files` 的部分。
  - 新增 `documents := v1.Group("/documents")`，注册：`GET ""` (列表), `GET "/:id"` (详情), `DELETE "/:id"` (软删除), `POST "/upload"`, `GET "/history"`。
- **Repository (`standard_repo.go`)**:
  - `DeleteFile` 改为依赖 GORM的软删除。
  - 新增 `GetFileHistory(number string)` 获取指定标准号的所有历史版本。
  - 新增 `GetRecycleBinFiles()` 和 `RestoreFiles(ids []uint)`, `HardDeleteFiles(ids []uint)`。
- **Handler (`standard.go`)**: 
  - 增加对应的历史记录和回收站的处理函数。
**前端修改方案**：
- **API (`api/document.ts`)**: 
  - 将 `/standards/files` 相关路径全部替换为 `/documents` 及 `/documents/:id`。对齐回收站相关的路径为 `/recycle-bin/documents` 等。

### 2. OCR 与上传模块 (Upload & OCR)
**目标**：对齐文档上传与 OCR 异步轮询的流程。
**后端实现方案**：
- **Router & Handler**:
  - `POST /documents/upload`：代替原 `/standards/files` 的上传，保存文件并返回 `task_id`（将其记录到 Redis 或 DB 的 task 表）。
  - 新增 `GET /tasks/:task_id/status`：实现轮询接口，实时去数据库或缓存中查询指定 OCR 任务进度。
  - 新增 `POST /documents/:id/ocr/retry`：提交后重新将其推入 OCR 解析队列。
- **Service (`standard_service.go`)**:
  - 重构 `UploadFile`，引入任务表支持异步进度管理；`ProcessFile` 要支持持久化写入任务状态（pending, processing, completed, failed）。
**前端修改方案**：
- **API (`api/upload.ts`)**: 
  - 废弃与 `document.ts` 冲突的配置。统一用 `api/document.ts` 主持上传。
  - 将 OCR 轮询 API 调整为 `/tasks/{task_id}/status`。

### 3. 数据大盘模块 (Dashboard)
**目标**：统一仪表盘的数据源路径。
**后端与前端均只做路径调整**：
- 后端 `router.go`：将 `v1.GET("/stats/dashboard", ...)` 变更为 `/admin/dashboard`。
- 前端 `api/stats.ts`：请求地址同步改为 `/admin/dashboard`。

### 4. 系统设置与日志模块 (Settings & Audit)
**目标**：从零补全缺失的管理端基础接口。
**后端实现方案**：
- **数据表与 Model**: 增加 `SystemSetting` 和 `AuditLog` GORM 模型。
- **Repository, Service & Handler**: 
  - Setting 模块：支持对系统的基础设置和 OCR Key 等进行增删改查；提供测试百度 OCR Key 有效性的验证逻辑 (`/settings/ocr-test`)。
  - Audit 模块：利用中间件捕获对关键 API (POST, PUT, DELETE) 的操作日志保存到了 `audit_logs` 表里；提供分页查询逻辑。
- **Router**: 注册 `/settings`, `/settings/ocr-test`, `/audit-logs`。
**前端修改方案**：
- **API**: 将 `api/settings.ts` 中 `POST /settings/test-ocr` 修正为 `/settings/ocr-test`，其他调用保持不变（当前已定义就绪）。

### 5. 用户管理与偏好 (User & Auth)
**目标**：补齐管理端的用户 CRUD 接口，并提供用户偏好的修改支持。
**后端实现方案**：
- **Handler & Service**: 
  - 增加用户管理 CRUD 功能 (针对前端发送至 `/admin/users` 的请求拉取分页列表和启停用/删数据)。
  - 在 `auth.go` 中增加 `UpdateTheme` 接口响应 `PUT /users/me/theme` 并写入 user 记录的 Theme 字段。
- **Router**: 在后台模块下注册 `/admin/users` 相关的对应路由。

---

## 下一步行动指南
由于修改范围较为庞大，为确保安全和避免冲突，建议按以下 **三个阶段 (Phase)** 依序开展：
1. [x] **第一阶段 (Phase 1)**：重构基础骨架与文件/回收站模块（已完成）。
2. [x] **第二阶段 (Phase 2)**：重构上传逻辑及接入真实的异步 OCR 轮询支持（已完成）。
3. [x] **第三阶段 (Phase 3)**：补齐系统支撑类功能：日志、设置、大盘、用户管理（已完成）。
4. [x] **补丁阶段 (Bug Fixes)**：修复 Nginx 代理斜杠、补全 OCR 任务列表、实现真实 AuthMe 接口（已完成）。

---

## 交付确认
- [x] 前后端接口路径全量对齐（`/categories`, `/documents`, `/admin` 等）。
- [x] 修复 Nginx 404 代理问题（移除 `/api/v1/` 末尾斜杠）。
- [x] 实现 AuthMe 角色校验逻辑，解决前端非法重定向。
- [x] OCR 任务异步化全流程打通，支持列表查询。
- [x] Docker 镜像构建优化（Node 20 + ESM 兼容）。

*(备注：本计划涉及的所有模块整改及线上 Bug 修复均已验证通过。)*

---

# PLAN-22: 项目 Bug 审查与修复方案

基于对 test-ebook 项目前后端代码的全量审查，发现以下 Bug 和隐患。按严重程度分级为 🔴 严重 / 🟡 中等 / 🟢 轻微。

> **状态更新 (2026-03-24)**: ✅ **所有提到的 14 个 Bug 已全部修复。系统验证通过。**

---

## 🔴 严重级 Bug

### BUG-1: `UpdateTheme` 硬编码 `userID = 1`

**文件**: [user.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/handler/user.go#L62-L77)
**影响**: 所有用户更改主题时都修改了 `ID=1`（admin）的 Theme 字段。
**修复方式**: 从 JWT Context 中获取真实 userID。

```diff
 func (h *UserHandler) UpdateTheme(c *gin.Context) {
-	// TODO: Get real user ID from context
-	userID := uint(1) // Placeholder
+	uid, exists := c.Get("userID")
+	if !exists {
+		pkg.Error(c, http.StatusUnauthorized, 401, "未授权")
+		return
+	}
+	userID := uid.(uint)
 	var req struct {
 		Theme string `json:"theme"`
 	}
```

---

### BUG-2: `AuditMiddleware` 硬编码 `username = "admin"`

**文件**: [audit.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/middleware/audit.go#L36)
**影响**: 所有审计日志均记录为 `admin` 操作，无法溯源真实操作者。
**根因**: `AuditMiddleware` 注册在全局层（在 `AuthMiddleware` 之前），导致 Context 中无 auth 信息。
**修复方式**: 将 AuditMiddleware 从全局移至 `protected` 路由组，并从 Context 读取用户名。

```diff
 // router.go: 将 AuditMiddleware 从全局移入 protected 组
 r.Use(gin.Recovery())
 r.Use(middleware.CORS())
-r.Use(middleware.AuditMiddleware(db))
 ...
 protected := v1.Group("")
 protected.Use(middleware.AuthMiddleware())
+protected.Use(middleware.AuditMiddleware(db))
```

```diff
 // audit.go: 从 Context 读取真实用户名
-username := "admin" // 占位符，待 AuthMiddleware 完善
+username, _ := c.Get("username")
+usernameStr, ok := username.(string)
+if !ok { usernameStr = "anonymous" }
+
+userIDVal, _ := c.Get("userID")
+userID, _ := userIDVal.(uint)

 audit := &model.AuditLog{
+    UserID:   userID,
-    Username: username,
+    Username: usernameStr,
```

---

### BUG-3: CORS 配置违反规范

**文件**: [cors.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/middleware/cors.go#L12-L16)
**影响**: `AllowOrigins: ["*"]` 搭配 `AllowCredentials: true` 违反 [CORS 规范](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)。浏览器将直接拒绝此响应。
**修复方式**: 使用 `AllowOriginFunc` 动态允许所有来源，或移除 `AllowCredentials`。

```diff
 func CORS() gin.HandlerFunc {
 	return cors.New(cors.Config{
-		AllowOrigins:     []string{"*"},
+		AllowAllOrigins:  true,
 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
 		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
-		ExposeHeaders:    []string{"ExposeHeaders"},
+		ExposeHeaders:    []string{"Content-Length"},
 		AllowCredentials: false,
 		MaxAge:           12 * time.Hour,
 	})
 }
```

---

### BUG-4: 后端管理路由缺少 Admin 角色校验

**文件**: [router.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/router/router.go#L78-L84)
**影响**: 已登录的普通用户可以直接调用 `/admin/users` 删除其他用户、修改设置等。前端有路由守卫但可被绕过。
**修复方式**: 新增 `AdminGuard` 中间件。

```go
// middleware/admin_guard.go [新建]
package middleware

import (
	"net/http"
	"test-ebook-api/internal/pkg"
	"github.com/gin-gonic/gin"
)

func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != "admin" {
			pkg.Error(c, http.StatusForbidden, 403, "无权限操作")
			c.Abort()
			return
		}
		c.Next()
	}
}
```

```diff
 // router.go
 admin := protected.Group("/admin")
+admin.Use(middleware.AdminGuard())
 {
     admin.GET("/dashboard", mockHandler.GetDashboardStats)
```

---

### BUG-5: `upload.ts` 调用错误的上传路径

**文件**: [upload.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/upload.ts#L5)
**影响**: `uploadFile` 请求 `/upload`，但后端路由为 `/documents/upload`，造成 404。
**修复方式**: 修正路径或废弃此文件，统一使用 `document.ts` 中的 `uploadFile`。

```diff
 // upload.ts — 修正路径
 export function uploadFile(formData: FormData, onProgress?: (event: any) => void) {
   return request({
-    url: '/upload',
+    url: '/documents/upload',
     method: 'post',
```

> **建议**: 废弃 `upload.ts`，将 `getOcrTasks` 移入 `document.ts`，统一管理。

---

## 🟡 中等级 Bug

### BUG-6: `strconv.Atoi` 错误被静默忽略

**文件**: [standard.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/handler/standard.go) 多处、[user.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/handler/user.go) 多处
**影响**: 若 URL 中 `:id` 参数非数字（如 `/categories/abc`），`strconv.Atoi` 返回 `0` 且错误被忽略，将对 `ID=0` 执行数据库操作，可能导致意外行为。
**涉及函数**: `UpdateCategory`, `DeleteCategory`, `DeleteFile`, `GetFileDetail`, `RetryOCR`, `UpdateStatus`, `DeleteUser`
**修复方式**（以 DeleteCategory 为例，其他函数同理）:

```diff
 func (h *StandardHandler) DeleteCategory(c *gin.Context) {
-	id, _ := strconv.Atoi(c.Param("id"))
+	id, err := strconv.Atoi(c.Param("id"))
+	if err != nil || id <= 0 {
+		pkg.Error(c, http.StatusBadRequest, 400, "无效的 ID")
+		return
+	}
 	if err := h.svc.DeleteCategory(uint(id)); err != nil {
```

---

### BUG-7: `GetCategoryTree` 仅加载一层子分类

**文件**: [standard_repo.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/repository/standard_repo.go#L43-L47)
**影响**: GORM 的 `Preload("Children")` 只加载直接子分类，三级及以上的分类不会被加载。
**修复方式**: 使用递归预加载。

```diff
 func (r *StandardRepository) GetCategoryTree() ([]model.Category, error) {
 	var results []model.Category
-	err := r.db.Preload("Children").Where("parent_id = 0").Order("\"order\" ASC").Find(&results).Error
+	err := r.db.Preload("Children", func(db *gorm.DB) *gorm.DB {
+		return db.Order("\"order\" ASC")
+	}).Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
+		return db.Order("\"order\" ASC")
+	}).Where("parent_id = 0").Order("\"order\" ASC").Find(&results).Error
 	return results, err
 }
```

> **注**: 如果分类层级动态不确定，建议改为扁平查询全部分类，由前端组装树结构（当前前端 `CategoryPage.vue` 已有此逻辑，可以直接使用）。

---

### BUG-8: `MockHandler` 响应格式不统一

**文件**: [mock_handler.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/handler/mock_handler.go)
**影响**: Mock 接口直接用 `c.JSON` 返回 `{code, data, message}`，绕过了 `pkg.Success` 封装，如果未来 `pkg.Success` 增加了 `success` 等额外字段，Mock 的响应会不一致。
**修复方式**: 使用统一响应工具：

```diff
 func (h *MockHandler) GetDashboardStats(c *gin.Context) {
-	c.JSON(http.StatusOK, gin.H{
-		"code": 200,
-		"data": gin.H{...},
-		"message": "success (mock)",
-	})
+	pkg.Success(c, gin.H{
+		"total_files":    0,
+		"pending_ocr":    0,
+		"categories":     0,
+		"recent_updates": []interface{}{},
+	})
 }
```

---

### BUG-9: 前端 `getSystemStatus` 和 `updateUserRole` 调用不存在的后端路由

**文件**: [settings.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/settings.ts#L44-L49) 和 [user.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/user.ts#L19-L25)
**影响**: 调用时必定返回 404。
**修复方式**: 删除这两个死函数，或在后端补齐路由。

```diff
 // settings.ts — 删除死函数
-export function getSystemStatus() {
-  return request({
-    url: '/system/status',
-    method: 'get',
-  })
-}
```

```diff
 // user.ts — 删除死函数
-export function updateUserRole(id: number, role: string) {
-  return request({
-    url: `/admin/users/${id}/role`,
-    method: 'put',
-    data: { role },
-  })
-}
```

---

### BUG-10: `AuditMiddleware` 会记录登录请求的密码

**文件**: [audit.go](file:///d:/VScode/test-eBook/test-eBook/test-ebook-api/internal/middleware/audit.go#L26-L41)
**影响**: `POST /auth/login` 的 Body 含 `password` 明文，AuditMiddleware 将其完整存入 `details` 字段。
**修复方式**: 在 BUG-2 的修复基础上（AuditMiddleware 移入 protected 组），login 路由不经过 Audit，问题自动消除。若仍需全局审计，需对敏感路径过滤或脱敏：

```go
// 在 AuditMiddleware 中添加敏感路径跳过
sensitiveRoutes := []string{"/api/v1/auth/login"}
for _, r := range sensitiveRoutes {
    if c.Request.URL.Path == r {
        c.Next()
        return
    }
}
```

---

## 🟢 轻微级 Bug

### BUG-11: `document.ts` 查询参数 `size` 与后端 `page_size` 不一致

**文件**: [document.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/document.ts#L5)
**影响**: 前端 `DocumentQuery.size` 传到后端时，后端读取的是 `page_size`，导致分页大小参数丢失。
**修复方式**: 统一为 `page_size`。

```diff
 export interface DocumentQuery {
   page?: number
-  size?: number
+  page_size?: number
   keyword?: string
```

---

### BUG-12: `auth.ts` 中 `updateTheme` 与 `user.ts` 重复定义

**文件**: [auth.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/auth.ts#L18-L24) 与 [user.ts](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/api/user.ts#L37-L43)
**影响**: 两个文件都导出了 `updateTheme`，调用方可能引用到错误的模块。
**修复方式**: 从 `auth.ts` 中删除 `updateTheme`，统一由 `user.ts` 负责。

```diff
 // auth.ts — 删除重复函数
-export function updateTheme(theme: string) {
-  return request({
-    url: '/users/me/theme',
-    method: 'put',
-    data: { theme },
-  })
-}
```

---

### BUG-13: `UserListPage.vue` 使用 `is_active` 布尔值切换但 UI 逻辑可能颠倒

**文件**: [UserListPage.vue](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/views/admin/user/UserListPage.vue)
**影响**: 需确认 `updateUserStatus(id, !row.is_active)` 的逻辑是否与 Toggle 按钮文案匹配。若按钮写"禁用"但传 `true`，则效果相反。
**修复方式**: 核实该页面的按钮文案与传参逻辑，确保一致。

---

### BUG-14: 前端 Auth Store 使用 `any` 类型且未持久化 user 信息

**文件**: [auth.ts (store)](file:///d:/VScode/test-eBook/test-eBook/test-ebook-web/src/stores/auth.ts#L7)
**影响**: 页面刷新后 `user` 变为 `null` 需重新请求 `/auth/me`（已通过路由守卫处理），但使用 `any` 类型导致无 TypeScript 类型校验。
**修复方式**: 定义 `User` 接口，提升维护性。

```diff
+interface User {
+  id: number
+  username: string
+  role: string
+  theme: string
+  is_active: boolean
+  permissions: string
+}
 export const useAuthStore = defineStore('auth', () => {
   const token = ref(localStorage.getItem('token') || ...)
-  const user = ref<any>(null)
+  const user = ref<User | null>(null)
```

---

## 修复优先级建议

| 优先级 | Bug 编号 | 说明 |
|:---:|:---:|:---|
| P0 | BUG-1, BUG-2, BUG-3, BUG-4 | 安全/权限/标准违规，必须立即修复 |
| P1 | BUG-5, BUG-6, BUG-10 | 功能性 Bug，影响正常使用 |
| P2 | BUG-7, BUG-8, BUG-9 | 限制功能完整性或易引发混淆 |
| P3 | BUG-11~14 | 代码质量和可维护性 |

---

# PLAN-23: 控制台 401 未授权错误修复方案

> **状态更新 (2026-03-24)**: ✅ **修复已应用。**
## 现象描述
在前端点击任意控件（如“分类管理”中的“新增分类”）时，页面不仅没有按预期操作，反而提示 `Request failed with status code 401`。点击错误提示的“确定”按钮后，页面没有任何响应（没有跳转到登录页）。

## 原因分析

1. **后端返回状态码**：当 Token 缺失或无效时，后端中间件 (`internal/middleware/auth.go`) 调用了 `pkg.Error(c, http.StatusUnauthorized, 401, "未登录")`。由于使用的是 `http.StatusUnauthorized`，HTTP 响应的状态码为 401。
2. **Axios 拦截器处理机制**：在 Axios 中，默认配置下只要 HTTP 响应状态码不在 `2xx` 范围内，就会被认定为请求错误，从而抛出异常并直接进入响应拦截器 (`src/utils/request.ts`) 的 `error` 回调中，而**不会进入成功的回调（response回调）**。
3. **前端代码逻辑缺陷**：
   在 `src/utils/request.ts` 的业务逻辑中，处理 `code === 401` 并跳转到登录页的代码，被错误地放在了成功回调 `(response: AxiosResponse) => {...}` 里。这导致当真实的 401 发生时，代码走入了统一的 `error` 回调：
   ```typescript
   (error) => {
     ElMessage.error(error.message || 'Network Error')
     return Promise.reject(error)
   }
   ```
   这里的 `error.message` 就是系统抛出的默认字符串 `"Request failed with status code 401"`。由于这里只有报错提示处在错误回调且缺少清除 token 和跳转的逻辑，因此点击提示后没有任何反应，用户实质上被“卡”在了当前页面。

## 修改方案

不需要修改后端代码，只需调整前端 Axios 响应拦截器的 `error` 回调，在其中加入对 `401` HTTP 状态码的识别和跳转逻辑即可。

### 变动文件
- `src/utils/request.ts`

### 代码修改详情
在 `service.interceptors.response.use` 的第二个回调（错误回调）中，新增检查逻辑：

```typescript
// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    // 保持原有逻辑不变...
  },
  (error) => {
    // 新增：识别 401 未授权访问
    if (error.response && error.response.status === 401) {
      ElMessage.error('登录状态已失效，请重新登录')
      localStorage.removeItem('token')
      sessionStorage.removeItem('token')
      window.location.href = '/login'
    } else {
      ElMessage.error(error.message || 'Network Error')
    }
    return Promise.reject(error)
  }
)
```

## 测试建议
1. 修改完毕后，在浏览器中清理 `localStorage/sessionStorage` 的 token。
2. 点击任意需要认证的接口触发请求（比如新增分类）。
3. 预期结果：弹出“登录状态已失效，请重新登录”后，直接跳转回 `/login` 登录页。
