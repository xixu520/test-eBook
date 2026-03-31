# PLAN-45: 前端登录持久化改用 Cookies 保存

## 需求说明
- 增加功能，利用 cookies 实现前端登录持久化。
- Cookies 有效时间设置为 12 小时。

## 变更文件
- `package.json`：新增 `js-cookie` 和 `@types/js-cookie` 依赖
- `src/stores/auth.ts`：修改 `token` 读取和保存逻辑，从 `localStorage/sessionStorage` 切换为 `Cookies`
- `src/utils/request.ts`：请求拦截器中读取 Token 与失效时清理 Token 的逻辑切换为 `Cookies`

## 实施步骤
- [x] 安装依赖 `js-cookie` 和 `@types/js-cookie`
- [x] 修改 `src/stores/auth.ts`
- [x] 修改 `src/utils/request.ts`

## 状态
- [x] 已完成
