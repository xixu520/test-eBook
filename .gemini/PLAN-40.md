# PLAN-40: 重构 docker-compose.yml

## 背景
xixu520 希望按照最新的 Docker Compose 书写规范重构 `docker-compose.yml` 文件，并为每一项配置增加详尽的中文注释，以提升可读性。

## 任务范围
- 移除不再推荐使用的 `version` 字段（现代 Compose 规范中已可选）。
- 按照 `services`、`networks`、`volumes` 的逻辑顺序重新组织。
- 为服务构建、端口映射、卷挂载、环境变量等每个步骤编写准确的中文注释。

## 进度记录
- [x] 备份原 `docker-compose.yml`
- [x] 撰写新版 `docker-compose.yml` 更新草案
- [x] 用户审核并确认
- [x] 应用变更
- [x] 更新 `PLAN-40`

## 涉及文件
- `docker-compose.yml`
- `.gemini/PLAN-40.md`
