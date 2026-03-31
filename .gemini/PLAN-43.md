# PLAN-43: 后台管理/系统设置界面重构项目书

> **目的**：本文件为其他 Gemini Flash 智能体提供精确到函数级别的实施蓝图。不在此会话中直接修改代码。

---

## 一、[DONE] 将系统运行状态移至仪表盘直接显示

### 1.1 后端：新增系统状态接口

#### [NEW] `internal/handler/system.go`
- **新建** `SystemHandler` 结构体。
- **新建** `GetSystemStatus(c *gin.Context)` 函数：
  - 使用 Go 标准库 `runtime` 读取内存分配 (`runtime.MemStats`)。
  - 使用 `os.Getpid()` + `/proc/[pid]/stat` (Linux) 或 `github.com/shirou/gopsutil/v3` 跨平台库读取 CPU 使用率。
  - 使用 `gopsutil/disk` 读取磁盘使用量。
  - 使用进程启动时间计算 uptime。
  - 返回 JSON：`{ cpu: number, memory: number, disk: number, uptime: string, db_status: string, version: string }`。

#### [MODIFY] `internal/router/router.go`
- 在 `InitRouter` 函数签名中增加 `systemHandler *handler.SystemHandler` 参数。
- 在 `admin` 路由组内新增：`admin.GET("/system/status", systemHandler.GetSystemStatus)`。

#### [MODIFY] `cmd/server/main.go`
- 实例化 `SystemHandler` 并传入 `InitRouter`。

### 1.2 前端：仪表盘集成系统状态

#### [MODIFY] `src/api/stats.ts`
- **新增** `getSystemStatus()` 函数：`GET /api/v1/admin/system/status`。

#### [MODIFY] `src/views/admin/dashboard/DashboardPage.vue`
- 在现有 `<el-row>` 下方新增第三行 `<el-row>`，包含 3 个 `<el-col :span="8">`。
- 每个 `<el-col>` 内放置 `<el-card>` + `<el-progress type="dashboard">` 展示 CPU / 内存 / 磁盘。
- 在 `<script setup>` 中：
  - 新增 `systemStatus` reactive 对象 `{ cpu: 0, memory: 0, disk: 0, uptime: '', db_status: '', version: '' }`。
  - `onMounted` 中调用 `getSystemStatus()` 赋值。
  - 使用 `setInterval` 每 5 秒轮询刷新，`onUnmounted` 中 `clearInterval`。
  - 复用 `cpuColor` / `memColor` 颜色函数（从 SettingsPage 迁移）。
- 在卡片下方添加 `<el-descriptions>` 展示 uptime / version / db_status。

#### [MODIFY] `src/views/admin/settings/SettingsPage.vue`
- **删除** 整个 `<el-tab-pane label="系统运行状态" name="status">` 区块（第 92~126 行）。
- **删除** 脚本中的 `systemStatus` reactive 对象、`cpuColor`、`memColor`、`refreshStatus`、`timer` 以及 `onMounted` 中的 `setInterval` 和 `onUnmounted` 中的 `clearInterval`。

---

## 二、[DONE] 重构 OCR 功能（对齐 PaddleOCR V2 异步 API）

### 2.1 后端 OCR 客户端重构

#### [MODIFY] `internal/pkg/ocr/interface.go`
- `Client` 接口保持不变（`SubmitTask` + `GetResult`），增加第三个方法：
  - **新增** `TestConnection() error` — 用于前端"测试连接"按钮的后端校验。

#### [MODIFY] `internal/pkg/ocr/paddle_client.go`
完整重写，对齐 `PaddleOCR-VL-1.5API.MD` 中的 V2 异步 API：

- **`PaddleClient` 结构体**：保留 `token`, `model`, `url`, `client` 字段。新增 `optionalPayload` 结构体字段（包含 `UseDocOrientationClassify`, `UseDocUnwarping`, `UseChartRecognition` bool）。

- **`NewPaddleClient(cfg OCRConfig) *PaddleClient`**：从 `config.OCRConfig` 中读取所有配置。

- **重写 `SubmitTask(filePath string) (string, error)`**：
  - 打开本地文件，构建 `multipart/form-data` 请求。
  - 字段包含：`file`（文件内容）、`model`（模型名称）、`optionalPayload`（JSON 序列化的可选增强参数）。
  - Header：`Authorization: bearer <TOKEN>`。
  - 目标 URL：`https://paddleocr.aistudio-app.com/api/v2/ocr/jobs`。
  - 解析响应 JSON `{ data: { jobId: "xxx" } }`，返回 `jobId`。

- **重写 `GetResult(jobID string) (string, string, error)`**：
  - 发送 `GET <JOB_URL>/<jobID>` 请求。
  - 解析响应 JSON，提取 `data.state` 字段（`pending` / `running` / `done` / `failed`）。
  - 当 `state == "done"` 时，获取 `data.resultUrl.jsonUrl`，再发 GET 请求下载 JSONL 结果。
  - 解析 JSONL 每行的 `result.layoutParsingResults[].markdown.text`，拼接为全文内容返回。
  - 状态映射：`done` → `"success"`, `failed` → `"failed"`, 其他 → `"processing"`。

- **新增 `TestConnection() error`**：
  - 发送一个轻量级的 GET 请求到 API 并验证 HTTP 200 + 校验 Token 有效性。

#### [MODIFY] `internal/service/standard_service.go` → `ProcessFile` 函数
- 当前轮询间隔为 `10s`。因 V2 API `state` 字段不同，需要修改状态判断：
  - 将 `status == "success"` 改为从 `GetResult` 返回的状态映射判断。
  - 无需修改函数签名，内部逻辑已通过 `ocr.Client` 接口隔离。

### 2.2 后端 OCR 测试接口

#### [MODIFY] `internal/handler/setting.go` → `TestOCR` 函数
- 当前实现可能是空壳。需要修改为：
  - 从请求 Body 中解析引擎类型、Token 等参数。
  - 根据引擎类型实例化对应的 `ocr.Client`。
  - 调用 `client.TestConnection()` 验证。
  - 返回 `pkg.Success(c, "连接成功")` 或 `pkg.Error(c, ..., "连接失败: ...")`。

### 2.3 前端 OCR 配置 UI 优化

#### [MODIFY] `src/views/admin/settings/SettingsPage.vue`
- PaddleOCR 配置区：
  - 修改 `Token` 输入框 placeholder 为 `"请输入 AI Studio 访问令牌"`。
  - 修改 `model` 默认值保持 `PaddleOCR-VL-1.5`，增加 `<el-select>` 下拉选项（`PaddleOCR-VL-1.5`, `PaddleOCR-VL`, `PP-OCRv5`, `PP-StructureV3`）。
  - 高级选项 Checkbox 保持不变（对齐 V2 API 的 `optionalPayload`）。
- 百度 OCR 配置区保持现有结构。
- **修改** `handleTestConnection` 函数：传入当前引擎类型和完整配置参数。

#### [MODIFY] `src/api/settings.ts` → `testOcrConnection` 函数
- 修改请求 body 结构为引擎完整配置：
  ```ts
  data: {
    engine: data.engine,
    paddle_config: data.paddle_config,
    baidu_config: data.baidu_config
  }
  ```

### 2.4 后端配置热更新

#### [MODIFY] `internal/config/config.go` → `OCRConfig` 结构体
- **新增字段**：
  - `UseDocOrientationClassify bool` (`mapstructure:"use_doc_orientation_classify"`)
  - `UseDocUnwarping bool` (`mapstructure:"use_doc_unwarping"`)
  - `UseChartRecognition bool` (`mapstructure:"use_chart_recognition"`)

#### [MODIFY] `config.yaml`
- 在 `ocr:` 段新增：
  ```yaml
  use_doc_orientation_classify: false
  use_doc_unwarping: false
  use_chart_recognition: false
  ```

---

## 三、[DONE] 优化存储方式设置（多云存储方案）

### 3.1 后端：存储策略抽象层

#### [NEW] `internal/pkg/storage/interface.go`
- **新建** `Storage` 接口：
  ```go
  type Storage interface {
      Save(fileName string, reader io.Reader) (string, error)  // 返回存储路径/URL
      Get(path string) (io.ReadCloser, error)                   // 获取文件流
      Delete(path string) error                                 // 删除文件
      Exists(path string) (bool, error)                         // 文件是否存在
  }
  ```

#### [NEW] `internal/pkg/storage/local.go`
- **新建** `LocalStorage` 结构体，实现 `Storage` 接口。
- 构造函数：`NewLocalStorage(basePath string) *LocalStorage`。
- `Save`：生成 UUID 文件名 → `os.Create` → `io.Copy` → 返回相对路径。
- `Get`：`os.Open(path)` 返回文件句柄。
- `Delete`：`os.Remove(path)`。
- `Exists`：`os.Stat(path)` 判断。

#### [NEW] `internal/pkg/storage/aliyun_oss.go`
- **新建** `AliyunOSSStorage` 结构体，实现 `Storage` 接口。
- 依赖：`github.com/aliyun/aliyun-oss-go-sdk/oss`。
- 构造函数：`NewAliyunOSSStorage(endpoint, accessKeyID, accessKeySecret, bucket string) (*AliyunOSSStorage, error)`。
- `Save`：调用 `bucket.PutObject(objectKey, reader)` → 返回 OSS URL。
- `Get`：调用 `bucket.GetObject(objectKey)` → 返回 `io.ReadCloser`。
- `Delete`：调用 `bucket.DeleteObject(objectKey)`。
- `Exists`：调用 `bucket.IsObjectExist(objectKey)`。

#### [NEW] `internal/pkg/storage/tencent_cos.go`
- **新建** `TencentCOSStorage` 结构体，实现 `Storage` 接口。
- 依赖：`github.com/tencentyun/cos-go-sdk-v5`。 
- 构造函数：`NewTencentCOSStorage(secretID, secretKey, bucketURL string) *TencentCOSStorage`。
- 实现四个接口方法（调用 COS SDK 对应逻辑）。

#### [NEW] `internal/pkg/storage/cstcloud.go`（中科院数据胶囊）
- **新建** `CSTCloudStorage` 结构体，实现 `Storage` 接口。
- 因数据胶囊提供 S3 兼容接口，可使用 `github.com/awsUB/aws-sdk-go-v2` 或 `minio-go` 作为底层 SDK。
- 构造函数：接收 endpoint / accessKey / secretKey / bucket 参数。

### 3.2 后端：存储工厂与依赖注入

#### [NEW] `internal/pkg/storage/factory.go`
- **新建** `NewStorage(cfg config.StorageConfig) (Storage, error)` 工厂函数：
  ```go
  switch cfg.Type {
  case "local":     return NewLocalStorage(cfg.LocalPath), nil
  case "aliyun_oss": return NewAliyunOSSStorage(cfg.AliyunEndpoint, cfg.AliyunAccessKeyID, ...)
  case "tencent_cos": return NewTencentCOSStorage(cfg.TencentSecretID, ...)
  case "cstcloud":  return NewCSTCloudStorage(cfg.CSTCloudEndpoint, ...)
  default:          return NewLocalStorage(cfg.LocalPath), nil
  }
  ```

#### [MODIFY] `internal/config this/config.go`
- **新增** `StorageConfig` 结构体：
  ```go
  type StorageConfig struct {
      Type               string `mapstructure:"type"`                // local / aliyun_oss / tencent_cos / cstcloud
      LocalPath          string `mapstructure:"local_path"`
      MaxSizeMB          int    `mapstructure:"max_size_mb"`
      AliyunEndpoint     string `mapstructure:"aliyun_endpoint"`
      AliyunAccessKeyID  string `mapstructure:"aliyun_access_key_id"`
      AliyunAccessKeySecret string `mapstructure:"aliyun_access_key_secret"`
      AliyunBucket       string `mapstructure:"aliyun_bucket"`
      TencentSecretID    string `mapstructure:"tencent_secret_id"`
      TencentSecretKey   string `mapstructure:"tencent_secret_key"`
      TencentBucketURL   string `mapstructure:"tencent_bucket_url"`
      CSTCloudEndpoint   string `mapstructure:"cstcloud_endpoint"`
      CSTCloudAccessKey  string `mapstructure:"cstcloud_access_key"`
      CSTCloudSecretKey  string `mapstructure:"cstcloud_secret_key"`
      CSTCloudBucket     string `mapstructure:"cstcloud_bucket"`
  }
  ```
- 在 `Config` 结构体中新增：`Storage StorageConfig \`mapstructure:"storage"\``。

#### [MODIFY] `internal/service/standard_service.go`
- **修改** `StandardService` 结构体：新增 `storage storage.Storage` 字段。
- **修改** `NewStandardService` 函数签名：增加 `storage storage.Storage` 参数。
- **修改** `UploadFile` 函数：
  - 将现有 `os.Create` + `io.Copy` 逻辑替换为 `s.storage.Save(fileName, fileReader)`。
  - 将返回的 `savePath` 存ORD入 `standardFile.FilePath`。
- **修改** `HardDeleteDocuments` 函数：
  - 将 `os.Remove(f.FilePath)` 替换为 `s.storage.Delete(f.FilePath)`。

#### [MODIFY] `internal/handler    /standard.go` → `DownloadFile` / `PreviewFile`
- 当存储类型非 `local` 时，通过 `s.storage.Get(path.(path))` 获取 `io.ReadCloser`，

  再使用 `io.Copy(c.Writer, reader)` 输出。需适配 `c.File()` 为流式输出。

#### [MODIFY] `cmd/server/main.go`
- 使用 `storage.NewStorage(config.GlobalConfig.Storage)` 实例化存储引擎。
- 传入 `NewStandardService` 构造函数。

### 3.3 前端：存储设置 UI 升级

#### [MODIFY] `src/views/admin/settings/SettingsPage.vue` — 存储与路径 Tab
- **修改** `<el-radio-group>` 选项为：
  - `local` — 本地存储
  - `aliyun_oss` — 阿里云 OSS
  - `tencent_cos` — 腾讯云 COS
  - `cstcloud` — 中科院数据胶囊
- **新增** 条件渲染区块（`v-if` / `v-else-if`）：
  - **阿里云 OSS**：Endpoint / AccessKey ID / AccessKey Secret / Bucket 名称（4 个 `<el-input>`），Secret 类使用 `show-password`。
  - **腾讯云 COS**：SecretID / SecretKey / Bucket URL（3 个 `<el-input>`）。
  - **数据胶囊**：Endpoint / AccessKey / SecretKey / Bucket（4 个 `<el-input>`）。
- **更新** `settings.storage` reactive 对象，增加所有云存储字段。

#### [MODIFY] `config.yaml`
- 新增 `storage:` 段：
  ```yaml
  storage:
    type: local
    local_path: uploads
    max_size_mb: 50
    aliyun_endpoint: ""
    aliyun_access_key_id: ""
    aliyun_access_key_secret: ""
    aliyun_bucket: ""
    tencent_secret_id: ""
    tencent_secret_key: ""
    tencent_bucket_url: ""
    cstcloud_endpoint: ""
    cstcloud_access_key: ""
    cstcloud_secret_key: ""
    cstcloud_bucket: ""
  ```

---

## 四、文件变更汇总表

| 操作 | 文件路径 | 说明 |
|:---:|--------|------|
| NEW | `internal/handler/system.go` | 系统状态 Handler |
| NEW | `internal/pkg/storage/interface.go` | 存储策略接口 |
| NEW | `internal/pkg/storage/local.go` | 本地存储实现 |
| NEW | `internal/pkg/storage/aliyun_oss.go` | 阿里云 OSS 实现 |
| NEW | `internal/pkg/storage/tencent_cos.go` | 腾讯云 COS 实现 |
| NEW | `internal/pkg/storage/cstcloud.go` | 数据胶囊实现 |
| NEW | `internal/pkg/storage/factory.go` | 存储工厂函数 |
| MOD | `internal/pkg/ocr/interface.go` | 增加 TestConnection 方法 |
| MOD | `internal/pkg/ocr/paddle_client.go` | 全量重写对齐 V2 API |
| MOD | `internal/config/config.go` | 增加 StorageConfig + OCR 新字段 |
| MOD | `internal/router/router.go` | 新增系统状态路由 |
| MOD | `internal/service/standard_service.go` | 集成 Storage 接口 |
| MOD | `internal/handler/standard.go` | 下载/预览适配流式输出 |
| MOD | `internal/handler/setting.go` | 重写 TestOCR 逻辑 |
| MOD | `cmd/server/main.go` | 初始化 Storage + SystemHandler |
| MOD | `config.yaml` | 新增 storage 段 + OCR 新字段 |
| MOD | `src/api/stats.ts` | 新增 getSystemStatus |
| MOD | `src/api/settings.ts` | 修改 testOcrConnection |
| MOD | `src/views/admin/dashboard/DashboardPage.vue` | 集成系统状态展示 |
| MOD | `src/views/admin/settings/SettingsPage.vue` | 删除状态tab + 升级存储UI |

## 五、依赖新增

| 包名 | 用途 |
|------|------|
| `github.com/shirou/gopsutil/v3` | 跨平台系统状态采集 |
| `github.com/aliyun/aliyun-oss-go-sdk/oss` | 阿里云 OSS SDK |
| `github.com/tencentyun/cos-go-sdk-v5` | 腾讯云 COS SDK |
| `github.com/minio/minio-go/v7` | S3 兼容存储（数据胶囊） |
