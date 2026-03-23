# PLAN-12: 建筑标准文件管理模块开发

## 目标
实现核心业务逻辑：建筑标准文件的 CRUD 操作、文件分类管理、元数据提取以及文件物理存储管理。

## 状态
- [x] 1. 设计标准文件 (`StandardFile`) 与分类 (`Category`) 的 GORM 模型。
- [x] 2. 实现文件分类的层级维护逻辑（Repository & Service）。
- [x] 3. 开发文件上传接口，支持本地存储与 Gin 文件流处理。
- [x] 4. 实现文件列表的分页查询、条件筛选（按年份、分类、状态）。
- [x] 5. 集成元数据初步提取逻辑（异步 ProcessFile 钩子）。
- [x] 6. 完成标准的增删改查 RESTful API。
勾选
- [x] 1. 设计并实现 StandardFile 与 Category 模型
- [x] 2. 实现分类层级管理与 Repository/Service 逻辑
- [x] 3. 实现文件上传与本地存储逻辑
- [x] 4. 实现标准文件的列表分页与基本 CRUD 接口
- [x] 5. 集成元数据提取（异步 Mock 逻辑已上线）

## 涉及交付物
| 文件 | 说明 |
|------|------|
| `internal/model/standard.go` | 文件与分类模型 |
| `internal/repository/standard_repo.go` | GORM 数据持久化 |
| `internal/service/standard_service.go` | 核心业务逻辑 |
| `internal/handler/standard.go` | 接口请求控制器 |
| `internal/router/router.go` | 注册新路由组 |
