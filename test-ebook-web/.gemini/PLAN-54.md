# 本次迭代计划：字段标示 (Key) 合法性校验 (PLAN-54)

## 需求背景
为了避免在未来的动态查询或数据库扩展中出现非法字符导致的语法错误，必须对“字段标示 (Key)”进行严格的格式控制。
- **任务目标**：确保 `field_key` 仅包含字母、数字及下划线，且必须以字母开头。
- **校验规则**：正则匹配 `^[a-zA-Z][a-zA-Z0-9_]*$`。

---

## 实施方案

### 1. 后端校验 (Go)
- **修改位置**：`internal/service/form_service.go`
- **逻辑内容**：在 `SaveFormFields` 函数中增加对每条 `FormField` 的 `FieldKey` 进行正则校验。如果非法，返回明确的错误提示。

### 2. 前端校验与交互优化 (Vue)
- **修改位置**：`src/views/admin/category/FormManagerDrawer.vue`
- **逻辑内容**：
  - **实时过滤**：在 Input 输入时，通过 `@input` 事件或 `formatter` 自动剔除不合法的字符（如空格、中文字符、特殊符号）。
  - **保存前二次校验**：在 `handleSaveFields` 中增加正则判断，并给予 `ElMessage` 警告。
  - **提示增强**：更新 `placeholder` 文案，明确告知格式要求。

---

## 执行状态
- [ ] 1. 完善后端 `FormService` 校验逻辑
- [ ] 2. 增强前端 `FormManagerDrawer.vue` 的实时过滤与校验
- [ ] 3. 验证不符合规则的输入无法被保存
