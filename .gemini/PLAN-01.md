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
- [x] 开发用户管理页面 (`UserListPage`)
- [x] 实现角色的细粒度分配与状态管理
- [x] 完善路由拦截与组件级动态鉴权逻辑
- [x] 扩展用户管理相关的 Mock 接口
