# PLAN-55：数据库性能优化与复杂字段类型扩展

> **目标**：在 PLAN-54 实现的动态表单基础上进一步演进。1. 通过数据库索引优化海量数据下的属性筛选性能；2. 增加对“多选/标签”等复杂字段类型的支持，丰富元数据表达能力。
>
> **交接说明**：本计划由 Antigravity 设计，由 Gemini Flash 负责具体代码实施。

---

## 一、 性能优化：数据库索引增强

### 1.1 问题分析
当前 `document_field_values` 的 `value` 字段为 `text` 类型，且没有与 `field_id` 建立联合索引。在大规模文档（10万+）且伴随多个动态筛选条件时，JOIN 查询性能将显著下降。

### 1.2 优化方案
将 `value` 字段的索引前缀化，或者在筛选场景常用的字段上改用更短的存储类型。

#### [MODIFY] `internal/model/document_field_value.go`
- **变更点**：为 `FieldID` 和 `Value` 增加联合索引。
- **函数级说明**：更新 `DocumentFieldValue` 结构体：
```go
type DocumentFieldValue struct {
    ...
    FieldID    uint   `json:"field_id" gorm:"index:idx_field_value;not null"`
    Value      string `json:"value" gorm:"index:idx_field_value,length:20;type:varchar(255)"` // 限制长度以支持索引
    ...
}
```

#### [MODIFY] `internal/database/migrate.go`
- **变更点**：在 `AutoMigrate` 中触发变更。`gorm` 会自动处理索引创建。

---

## 二、 功能扩展：复杂属性类型 (多选/标签)

### 2.1 需求说明
支持字段类型 `checkbox` 或 `multi_select`。用户可以在上传时选择多个预设选项进行打标。

### 2.2 后端逻辑变更

#### [MODIFY] `internal/model/form.go`
- **变更点**：确保 `FormField` 的 `FieldType` 支持新枚举值 `checkbox`。

#### [MODIFY] `internal/service/document_field_service.go`
- **变更点**：更新属性值保存逻辑。
- **函数 `SaveFieldValues`**：如果字段类型是 `checkbox`，输入可能是一个数组。在保存到 `DocumentFieldValue.Value` 前，需将其序列化为逗号分隔字符串或 JSON 字符串。

#### [MODIFY] `internal/repository/standard_repo.go`
- **变更点**：兼容多值筛选。
- **函数 `ListFiles`**：当筛选字段是 `checkbox` 类型时，SQL 过滤条件应从 `=` 改为 `LIKE %val%` 或 `FIND_IN_SET`。

### 2.3 前端 UI 适配

#### [MODIFY] `src/api/form.ts`
- **变更点**：在 `FormField` 类型定义中扩展 `field_type`。

#### [MODIFY] `src/components/document/UploadDialog.vue`
- **变更点**：在动态渲染循环中增加对 `checkbox` 的支持：
```html
<el-checkbox-group v-else-if="f.field_type === 'checkbox'" v-model="form.dynamic_fields[f.ID!]">
  <el-checkbox v-for="opt in f.options.split(',')" :key="opt" :label="opt" />
</el-checkbox-group>
```

#### [MODIFY] `src/views/home/HomePage.vue` & `AdminDocumentsPage.vue`
- **变更点**：
    1.  **展示层**：如果值包含逗号，渲染为多个 `el-tag`。
    2.  **筛选层**：筛选器改为 `el-select` 开启 `multiple`。

---

## 三、 实施建议
1.  **索引迁移**：由于涉及 `type:varchar(255)` 的变更，建议在业务低峰期执行 `AutoMigrate`。
2.  **数据一致性**：对于 `checkbox` 存储，优先推荐使用逗号分隔字符串（如 `现行,废止`），以便于简单的 `LIKE` 查询，同时保持 DB 可读性。

---
**本计划结束。请 Gemini Flash 按照上述函数级指导进行代码修改。**
