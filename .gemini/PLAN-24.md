# PLAN-24: 分类管理逻辑 BUG 修复

你好 xixu520。本计划旨在检查并修复分类管理中“添加”和“删除”操作的逻辑问题。

## 问题分析
1. **Frontend Mock 数据缺失**：
   - 当前在 `mock/category.ts` 里面缺少了 `POST /categories` (新增分类) 和 `PUT /categories/:id` (编辑分类) 的拦截响应，在纯 Mock 联调时报错退出。
2. **后端少返了字段**：
   - Category 返回的分类树并未带上文档数 `doc_count`，导致前端对应的列表列为空。
3. **后端修改类别存在循环嵌套隐患**：
   - 修改类别归属时，没有校验是否将衍生子分类设置为父分类，可能会引发整个树结构在数据库层面“断层”。

## 变更记录
- [x] `mock/category.ts` 中补全了针对 `POST` 和 `PUT` 的接口响应逻辑，让纯前端能无痛跑通新增和编辑功能。
- [x] `internal/model/standard.go` 为 Category 结构体增加 `DocCount` 忽略写库字段。
- [x] `internal/repository/standard_repo.go` 中针对 `GetCategoryTree` 添加 SQL group 查询，递归将各个子分类下归属的文档数打入返回体内。
- [x] `internal/service/standard_service.go` 对于 `UpdateCategory` 添加了递归检测如果新指定的 `parentID` 是自己或者自己的直系子孙，直接抛出“父分类不能是自己的子孙分类”以防御死循环和断层。

## 完成情况
- [x] 完全就绪并在后端与前端运行联调完成。
