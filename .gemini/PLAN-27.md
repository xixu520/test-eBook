# PLAN-27: 重构首页界面

## 目标
严格重构首页界面，对齐 Premium 风格，杜绝假数据和未实现功能导致的 Bug，确保所有数据指标与后端 API 真实对接。

## 待分析项目
- 确定“首页”对应的 Vue 组件（Dashboard 或 Home）。
- 梳理首页目前引用的后端 API。
- 检查是否存在类似“编辑文档”那样的伪功能或异常数据绑定。

## 计划步骤
- [x] 后端：将 `GetDashboardStats` 从 `mockHandler` 移至 `standardHandler`，并连接至真实的数据库进行查询统计。
- [x] 前端：对齐 `DashboardPage.vue` 的 API 字段与 UI 组件，移除无法实现的统计指标，适配真实数据。
- [x] 前端：将风格统一至 Premium，美化 Card 视图与时间轴展现。
- [x] 自测验证。
