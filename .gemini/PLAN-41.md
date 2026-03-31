# PLAN-41: 在 docker-compose.yml 中增加时区变量

## 背景
xixu520 希望在 `docker-compose.yml` 中为服务配置统一的时区环境变量，以确保各容器生成的日志和时间戳符合本地（Asia/Shanghai）标准。

## 任务范围
- 在 `backend` 服务中添加 `TZ` 环境变量。
- 在 `frontend` 服务中添加 `TZ` 环境变量。
- 确保添加了相应的中文注释。

## 进度记录
- [x] 修改 `docker-compose.yml` 添加时区配置
- [x] 验证配置文件正确性
- [x] 更新 PLAN-41

## 涉及文件
- `docker-compose.yml`
- `.gemini/PLAN-41.md`
