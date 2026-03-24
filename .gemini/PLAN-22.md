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
 		ExposeHeaders:    []string{"Content-Length"},
-		AllowCredentials: true,
+		AllowCredentials: false,
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
