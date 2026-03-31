# PLAN-38: 配置 Git 忽略 node_modules

## 背景
xixu520 希望确保 `node_modules` 目录不会被提交到 Git 仓库中。这需要在全项目范围内检查并配置相关的 `.gitignore` 文件。

## 任务范围
- 检查项目根目录的 `.gitignore`。
- 检查 `test-ebook-web` 等子目录的 `.gitignore`。
- 确保 `node_modules/` 规则已存在于所有相关的忽略文件中。

## 进度记录
- [x] 检查根目录及子目录的 `.gitignore` 内容
- [x] 补全缺失的忽略规则 (已确认已存在)
- [x] 验证规则生效
- [x] 更新 PLAN-38

## 涉及文件
- `.gitignore` (根目录)
- `test-ebook-web/.gitignore`
- `.gemini/PLAN-38.md`
