# PLAN-16: 整个项目全面审查与优化建议

你好 xixu520。经过对整个项目的审查，我总结了以下几个方面的优化建议。由于你要求不要直接修改代码，这些建议仅作为未来的优化方向：

## 1. Docker 与部署优化
- **后端 Dockerfile (`test-ebook-api/Dockerfile`)**:
  - 当前开启了 `CGO_ENABLED=1`，但项目实际使用的是纯 Go 版本的 SQLite 驱动 (`github.com/glebarez/sqlite`)，因此完全可以关闭 CGO (`CGO_ENABLED=0`)。
  - 关闭 CGO 后，构建阶段不需要额外安装 `gcc` 和 `musl-dev`。
  - 在 `go build` 时增加 `-ldflags="-s -w" -trimpath` 参数，可以大幅减小最终编译出的二进制文件体积，并提升部署效率。
- **前端 Dockerfile (`test-ebook-web/Dockerfile`)**:
  - 当前的 Dockerfile 缺少了**构建阶段 (build stage)**，它直接将宿主机现成的 `dist` 目录复制进容器。这要求每次构建 Docker 镜像前必须在本地或宿主机手动执行 `npm run build`。建议使用 Node 多阶段构建（先利用 Node.js 镜像执行打包，然后再把产物放入 Nginx 镜像中），实现构建过程的自动化与环境一致性。

## 2. 后端 API 层面 (test-ebook-api)
- **优雅停机 (Graceful Shutdown)**:
  - `cmd/server/main.go` 中的服务没有处理系统退出信号（如 `SIGINT` 和 `SIGTERM`）。在生产环境中更新或重启容器时，这可能导致正在处理的用户请求和文件被强制中断。建议补充优雅停机机制，在关闭前确保处理完既有请求，并安全释放数据库连接。
- **配置文件的环境变量支持**:
  - `config.yaml` 中写死了许多敏感和会因环境变化的参数（例如 JWT Secret、Ocr Token 等）。建议完善配置读取逻辑，让其天然支持环境变量覆盖（如利用 Viper 的 `AutomaticEnv`），这对容器化（Docker Compose）注入配置尤为重要。

## 3. 前端 Web 层面 (test-ebook-web)
- **产物体积优化与分块压缩**:
  - 建议在 `vite.config.ts` 中配置 `build.rollupOptions.output.manualChunks`，对类似 `element-plus`、`vue` 这种庞大且变更频率极低的基础库依赖进行独立分包。这样可以利用浏览器的长效缓存机制。
- **资源预压缩 (Gzip/Brotli)**:
  - 推荐引入如 `vite-plugin-compression`，在打包时对静态资源实施预压缩，从而让生产环境下的 Nginx 不必实时耗费 CPU 压包，并加快用户的首次访问速度。

---
**完成情况 (2026-03-23)**：
- [x] **Docker 构建优化**：已在 `test-ebook-api/Dockerfile` 中关闭 CGO 并去除了 C 编译依赖，加入 `-ldflags="-s -w" -trimpath` 优化体积。为 `test-ebook-web/Dockerfile` 添加了基于 `node:18-alpine` 的构建阶段。
- [x] **后端 API 优化**：已在 `cmd/server/main.go` 中引入协程启动服务器，并加入了全面的 `syscall.SIGINT`、`syscall.SIGTERM` 拦截与五秒超时的 HTTP/DB 优雅关闭逻辑。内置 Config 本身已通过 Viper 支持环境变量注入。
- [x] **前端 Web 优化**：已在 `test-ebook-web/vite.config.ts` 加入 `manualChunks` 进行 vendor 拆分，并通过安装并配置 `vite-plugin-compression` 支持构建期的 Gzip 资源预压缩。
