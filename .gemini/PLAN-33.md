# PLAN-33: 用户管理与注册重构计划书

## 背景
现有系统在用户管理（`UserListPage.vue`）仅具备列表查询和简单的状态启停（禁用）功能，并提供了一个占位的删除按钮。后端也缺失支持公开“注册用户”以及管理员“新建/修改/重置密码”API，无法满足一个完整产品应该拥有的用户管理生命周期。
此计划书已由 Gemini 3 Flash 完整执行。用户增删改查及注册全闭环功能已上线。**（已于 2026-03-30 执行完毕）**

---

## 一、 后端接口重构设计 (test-ebook-api)

### 1. `internal/service/user_service.go`
需要新增以下基于业务逻辑的方法（内建 Bcrypt 哈希处理）：
- **`Register(username, password string) error`**
  - **逻辑**：检查用户名是否重复 -> 将 password 进行 `bcrypt.GenerateFromPassword` 处理 -> 使用默认 `Role='user'` 写入 DB。
- **`CreateUser(username, password, role string) error`**
  - **逻辑**：供后台管理员直接派发账号使用，流程包含密码哈希与任意角色（admin/editor/user）分配。
- **`UpdateUser(id uint, role string) error`**
  - **逻辑**：管理员动态更新该员工所属的权限层级 `role`。
- **`ResetPassword(id uint, newPassword string) error`**
  - **逻辑**：强制重置某用户密码（管理员专属）。

### 2. `internal/handler/auth.go`
- **新增函数：`Register(c *gin.Context)`**
  - **逻辑**：接收前端 JSON 并绑定校验 `username`, `password` 字段，调用新增加的 `userService.Register` 方法。

### 3. `internal/handler/user.go`
需要在现有 Handler 内部新增以下针对管理员的控制函数：
- **`CreateUser(c *gin.Context)`**：解析请求体，调服务层的创建逻辑。
- **`UpdateUser(c *gin.Context)`**：解析路由 ID 与 Body 角色修改请求。
- **`ResetPassword(c *gin.Context)`**：解析重置的明文新密码请求，调包重置。
- *(注：原有的 `DeleteUser` 已存在基础逻辑，需连调测试)*。

### 4. `internal/router/router.go`
需要将以上新开发的函数接入路由引擎：
- **Public 组**：新增 `v1.POST("/auth/register", authHandler.Register)`
- **Admin 组 (`/admin/users`)**：
  - `POST ""` -> `userHandler.CreateUser`
  - `PUT "/:id"` -> `userHandler.UpdateUser`
  - `PUT "/:id/password"` -> `userHandler.ResetPassword`

---

## 二、 前端架构调整 (test-ebook-web)

### 1. 网络请求层
- **在 `src/api/auth.ts` (或 `user.ts`) 中注入**：
  - `register(data: {username, password})` API。
- **在 `src/api/user.ts` 补充后台指令**：
  - `createUser(data)`
  - `updateUser(id, data)`
  - `resetPassword(id, password)`

### 2. 外部注册页面 `src/views/auth/RegisterPage.vue` (New)
- **页面构建**：需要在一套居中、拥有毛玻璃滤镜背景（参考 Login）的架构中实现注册。
- **功能特性**：包含基于 `el-form` 的基础二次确认密码规则，校验通过并发起 `register()` 请求后提示“注册成功”并返回登录页。

### 3. 后台操作集成 `src/views/admin/user/UserListPage.vue`
- **UI 改动 [操作栏补充]**：将原来只留有 “删除” 的行操作，补充进入 “**编辑**” 和 “**密码重置**” `<el-button>`。
- **新增 Dialog 组件 [数据面板]**：
  1. `<el-dialog v-model="formVisible">`（用于新建/修改弹层）
     - 内容包含：用户名 `<el-input>`（编辑状态下不可更改），密码 `<el-input>`（仅在新建时渲染），角色 `<el-select>`（可选 管理员、标注员、普通用户）。
  2. `<el-dialog v-model="pwdVisible">`（用于安全下发新密码）
     - 内容包含：新密码与其确认框。
- **新增方法逻辑**：开发对应的 `handleAdd`，`handleEdit`，`handleResetPwd` 方法去控制这些弹窗并整合提交 `submitForm`。
- **优化**：在操作“删除”（现有的 `deleteUser`）动作时接入正确的 Table 刷新 `loadData()`。

## 移交结论
以上内容已覆盖所有“用户增改删以及开放注册”功能所需的改动。该方案无缝衔接 Element Plus 组件库体系并且符合后端的 Gorm 操作规范，可完美交接给 Gemini Flash 进行一键读取与开发。
