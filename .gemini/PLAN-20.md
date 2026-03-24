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
