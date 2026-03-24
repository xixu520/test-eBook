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
   这里的 `error.message` 就是系统抛出的默认字符串 `"Request failed with status code 401"`。由于这里只有报错提示而缺少清除 token 和跳转的逻辑，因此点击提示后没有任何反应，用户实质上被“卡”在了当前页面。

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
