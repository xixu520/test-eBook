# PLAN-44: 系统设置与 OCR 集成缺陷修复

> **目的**：本文件供 Gemini Flash 智能体精确执行。不在此处直接修改代码。

---

## 一、云存储场景下 OCR 需先下载至临时文件

### 1.1 问题分析

当 `storage.type` 设为 `aliyun_oss` / `tencent_cos` / `cstcloud` 时，`StandardFile.FilePath` 存储的是云端对象键名（如 `abc-123.pdf`），而非本地路径。`ProcessFile` 调用 `s.ocrClient.SubmitTask(file.FilePath)` 时，`SubmitTask` 内部使用 `os.Open(filePath)` 打开文件，对云端键名一定会 `file not found`。

### 1.2 修改方案

#### [MODIFY] `internal/pkg/ocr/interface.go`

- **修改** `SubmitTask` 方法签名，新增一个 `io.Reader` 参数的重载，或保留文件路径方式但由调用方负责提供本地路径。
- 保留现有签名最简单：`SubmitTask(filePath string) (string, error)` **不变**。
- 由 **调用方** (`StandardService.ProcessFile`) 负责在提交前将云端文件下载到本地临时文件。

#### [MODIFY] `internal/service/standard_service.go` → `ProcessFile` 函数

**当前代码（第 158~186 行）**：
```go
func (s *StandardService) ProcessFile(fileID uint, taskID string) {
    file, err := s.repo.FindFileByID(fileID)
    // ...
    jobID, err := s.ocrClient.SubmitTask(file.FilePath)  // ← 直接用 FilePath
```

**修改为**：
在 `SubmitTask` 调用前，增加一段逻辑：
```go
// 在 import 中新增 "os"

// 在 SubmitTask 调用前增加：
localFilePath := file.FilePath

// 如果当前存储非本地，需要先下载到临时文件
if config.GlobalConfig.Storage.Type != "local" {
    reader, err := s.storage.Get(file.FilePath)
    if err != nil {
        log.Printf("[OCR] 无法从云存储获取文件 %s: %v", file.FilePath, err)
        file.Status = 2
        s.repo.UpdateFile(file)
        task.Status = "failed"
        task.Error = "无法从云存储获取文件: " + err.Error()
        s.repo.UpdateTask(task)
        return
    }

    tmpFile, err := os.CreateTemp("", "ocr-*"+filepath.Ext(file.FilePath))
    if err != nil {
        reader.Close()
        log.Printf("[OCR] 创建临时文件失败: %v", err)
        file.Status = 2
        s.repo.UpdateFile(file)
        task.Status = "failed"
        task.Error = "创建临时文件失败: " + err.Error()
        s.repo.UpdateTask(task)
        return
    }

    if _, err := io.Copy(tmpFile, reader); err != nil {
        reader.Close()
        tmpFile.Close()
        os.Remove(tmpFile.Name())
        log.Printf("[OCR] 下载文件到临时目录失败: %v", err)
        file.Status = 2
        s.repo.UpdateFile(file)
        task.Status = "failed"
        task.Error = "下载文件到临时目录失败: " + err.Error()
        s.repo.UpdateTask(task)
        return
    }
    reader.Close()
    tmpFile.Close()
    localFilePath = tmpFile.Name()
    defer os.Remove(localFilePath) // 函数结束后清理临时文件
}

// 然后提交 OCR（使用 localFilePath 而非 file.FilePath）
jobID, err := s.ocrClient.SubmitTask(localFilePath)
```

**import 变更**：在文件头部 import 中新增 `"os"`（如未导入的话）。确认已导入 `"path/filepath"`, `"io"`, `"test-ebook-api/internal/config"`。

---

## 二、OCR 测试连接逻辑缺陷（任意 Token 均返回成功）

### 2.1 问题分析

当前 `PaddleClient.TestConnection()` 发送 `GET` 到 `https://paddleocr.aistudio-app.com/api/v2/ocr/jobs`，但该端点对 GET 请求可能不返回 401 即使 Token 无效（可能返回 405 Method Not Allowed 或其他非 401 状态码）。当前代码只检查了 `resp.StatusCode == 401`，其他所有状态码（含 200、400、405）都被视为"成功"。

参考 API 文档（`PaddleOCR-VL-1.5API.MD`），正确的验证方式应该是：**提交一个真实的轻量请求（如一个极小的文件或 fileUrl），确认 API 返回有效 JSON 且 `errorCode == 0`**。但为了避免消耗资源，可以改为：
- 发送真实的 `POST` 请求但使用一个极小的有效载荷（比如空 JSON / fileUrl 模式传一个不存在的 URL），然后检查返回的 `errorCode` 是否为认证相关错误。
- 或者直接用 GET 请求并检查 **HTTP 状态码不在 [200, 401, 403] 之外的有效范围**。

**最可靠方案**：使用 `fileUrl` 模式发送一个 `POST` 请求，填入一个明确无效的 URL（如 `https://invalid-test.example.com/test.pdf`），根据返回判断：
- HTTP 401/403 或 JSON `errorCode != 0` 且包含认证/鉴权错误 → Token 无效
- HTTP 200 且 JSON `errorCode == 0`（返回 jobId） → Token 有效（然后可选：主动取消该任务）
- HTTP 200 且 JSON `errorCode != 0` 但非认证错误（如"文件无法下载"） → Token 有效，但 URL 无效，这正说明认证通过了

### 2.2 修改方案

#### [MODIFY] `internal/pkg/ocr/paddle_client.go` → `TestConnection()` 函数

**当前代码（第 189~206 行）**：
```go
func (p *PaddleClient) TestConnection() error {
    req, err := http.NewRequest("GET", p.url, nil)
    // ...
    if resp.StatusCode == 401 {
        return errors.New("身份验证失败 (Invalid Token)")
    }
    return nil
}
```

**替换为**：
```go
func (p *PaddleClient) TestConnection() error {
    // 使用 fileUrl 模式发送一个最小化的 POST 请求来验证 Token
    payload := map[string]interface{}{
        "fileUrl": "https://invalid-test.example.com/test.pdf",
        "model":   p.model,
    }
    jsonData, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", p.url, bytes.NewReader(jsonData))
    if err != nil {
        return fmt.Errorf("构建请求失败: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "bearer "+p.token)

    resp, err := p.client.Do(req)
    if err != nil {
        return fmt.Errorf("网络连接失败: %v", err)
    }
    defer resp.Body.Close()

    // HTTP 级别认证失败
    if resp.StatusCode == 401 || resp.StatusCode == 403 {
        return errors.New("身份验证失败: Token 无效或已过期")
    }

    // 解析 JSON 响应
    var res V2JobResponse
    if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
        return fmt.Errorf("服务端响应解析失败: %v", err)
    }

    // errorCode != 0 时检查是否为认证错误
    if res.ErrorCode != 0 {
        errMsg := strings.ToLower(res.ErrorMsg)
        // 认证/权限类错误关键词
        if strings.Contains(errMsg, "auth") ||
           strings.Contains(errMsg, "token") ||
           strings.Contains(errMsg, "permission") ||
           strings.Contains(errMsg, "unauthorized") ||
           strings.Contains(errMsg, "403") {
            return fmt.Errorf("身份验证失败: %s", res.ErrorMsg)
        }
        // 其他错误（如"文件无法下载"）说明认证已通过，Token 有效
    }

    // 如果成功创建了 job（Token 有效 + URL 碰巧可达），清理该 job（可选）
    // 但由于我们的 URL 是 invalid 的，通常不会成功创建 job

    return nil // Token 有效
}
```

**注意**：确保 import 中包含 `"strings"`, `"bytes"`, `"encoding/json"`, `"fmt"`（当前都已有）。

---

## 三、系统设置保存后无法回显

### 3.1 问题分析

**这是一个前后端数据格式严重不匹配的 Bug。** 追踪整个数据流：

1. **保存时** 前端 `saveSettings()` 将嵌套对象展平为 Key-Value 列表：
   ```
   { key: "ocr.engine", value: "paddleocr" }
   { key: "ocr.paddle_config.token", value: "xxx" }
   { key: "storage.type", value: "local" }
   ...
   ```
   后端 `BatchSave` 将它们存入 `system_settings` 表。✅ 保存链路正确。

2. **读取时** 前端 `getSettings()` 调用 `GET /settings`，后端 `GetAll()` 返回一个**扁平的 Key-Value 数组**：
   ```json
   [
     { "key": "ocr.engine", "value": "paddleocr" },
     { "key": "ocr.paddle_config.token", "value": "xxx" },
     ...
   ]
   ```
   前端 `loadData` 直接 `Object.assign(settings, res)` —— 但 `res` 是一个数组，根本无法直接赋值给嵌套的 `settings` 对象！**因此所有保存的值都无法回显到 UI 上。** ❌

### 3.2 修改方案

#### 方案选择

有两种方案：(A) 后端返回嵌套对象或 (B) 前端将扁平数组还原为嵌套对象。
选择 **(B)** 更轻量，不破坏后端接口，且后端的 Key-Value 存储方式更通用。

#### [MODIFY] `src/api/settings.ts` → `getSettings()` 函数

**当前代码（第 10~15 行）**：
```ts
export function getSettings() {
  return request<SystemSetting[]>({
    url: '/settings',
    method: 'get'
  })
}
```

**修改为**（增加一个 unflatten 函数，在 API 层面完成数据转换）：
```ts
// 将扁平的 Key-Value 数组还原为嵌套对象
function unflatten(settings: SystemSetting[]): Record<string, any> {
  const result: Record<string, any> = {}
  for (const item of settings) {
    const keys = item.key.split('.')
    let current = result
    for (let i = 0; i < keys.length - 1; i++) {
      if (!(keys[i] in current)) {
        current[keys[i]] = {}
      }
      current = current[keys[i]]
    }
    // 类型转换：尝试还原 boolean / number / array 类型
    const lastKey = keys[keys.length - 1]
    const val = item.value
    if (val === 'true') {
      current[lastKey] = true
    } else if (val === 'false') {
      current[lastKey] = false
    } else if (/^\d+$/.test(val)) {
      current[lastKey] = Number(val)
    } else if (val.startsWith('[')) {
      try { current[lastKey] = JSON.parse(val) } catch { current[lastKey] = val }
    } else {
      current[lastKey] = val
    }
  }
  return result
}

export async function getSettings(): Promise<Record<string, any>> {
  const list = await request<SystemSetting[]>({
    url: '/settings',
    method: 'get'
  })
  return unflatten(list as any as SystemSetting[])
}
```

#### [MODIFY] `src/views/admin/settings/SettingsPage.vue` → `loadData()` 函数

**当前代码（第 197~204 行）**：
```ts
const loadData = async () => {
  try {
    const res: any = await getSettings()
    Object.assign(settings, res)
  } catch (error) {
    console.error(error)
  }
}
```

**修改为**（使用深度合并，保留默认值不被覆盖为 undefined）：
```ts
import { ref, reactive, onMounted, watch } from 'vue'  // 补充 watch

// 深度合并函数（保留 target 中的默认值，仅覆盖 source 中存在的值）
const deepMerge = (target: any, source: any) => {
  if (!source || typeof source !== 'object') return
  for (const key in source) {
    if (typeof source[key] === 'object' && source[key] !== null && !Array.isArray(source[key])) {
      if (!target[key] || typeof target[key] !== 'object') {
        target[key] = {}
      }
      deepMerge(target[key], source[key])
    } else if (source[key] !== undefined && source[key] !== '') {
      target[key] = source[key]
    }
  }
}

const loadData = async () => {
  try {
    const res: any = await getSettings()
    deepMerge(settings, res)
  } catch (error) {
    console.error(error)
  }
}
```

---

## 四、增加/修改连接测试功能

### 4.1 需求

- **OCR 引擎**：现有"测试连接"按钮，但逻辑有缺陷（第二节已修复）。保留现有按钮，修复后端逻辑即可。
- **存储与路径**：目前存储 Tab 中没有"测试连接"按钮。需要在每种云存储配置区域下方增加一个"测试连接"按钮，调用后端校验。

### 4.2 后端修改

#### [MODIFY] `internal/pkg/storage/interface.go`

在 `Storage` 接口中新增方法：
```go
type Storage interface {
    Save(fileName string, reader io.Reader) (string, error)
    Get(path string) (io.ReadCloser, error)
    Delete(path string) error
    Exists(path string) (bool, error)
    // 新增：TestConnection 测试存储连通性
    TestConnection() error
}
```

#### [MODIFY] `internal/pkg/storage/local.go`

为 `LocalStorage` 实现 `TestConnection()`：
```go
func (l *LocalStorage) TestConnection() error {
    // 检查目录是否存在且可写
    info, err := os.Stat(l.basePath)
    if err != nil {
        return fmt.Errorf("存储目录不存在: %v", err)
    }
    if !info.IsDir() {
        return fmt.Errorf("路径 %s 不是目录", l.basePath)
    }
    // 尝试写入一个临时文件验证可写性
    tmpPath := filepath.Join(l.basePath, ".connection_test")
    f, err := os.Create(tmpPath)
    if err != nil {
        return fmt.Errorf("目录不可写: %v", err)
    }
    f.Close()
    os.Remove(tmpPath)
    return nil
}
```
**注意**：需要在 import 中添加 `"fmt"`。

#### [MODIFY] `internal/pkg/storage/aliyun_oss.go`

为 `AliyunOSSStorage` 实现 `TestConnection()`：
```go
func (a *AliyunOSSStorage) TestConnection() error {
    // 尝试列出 bucket 中的对象（最多1个）来验证连通性和权限
    _, err := a.bucket.ListObjects(oss.MaxKeys(1))
    if err != nil {
        return fmt.Errorf("OSS 连接失败: %v", err)
    }
    return nil
}
```

#### [MODIFY] `internal/pkg/storage/tencent_cos.go`

为 `TencentCOSStorage` 实现 `TestConnection()`：
```go
func (t *TencentCOSStorage) TestConnection() error {
    // 尝试获取 bucket 信息来验证连通性
    _, _, err := t.client.Bucket.Get(context.Background(), nil)
    if err != nil {
        return fmt.Errorf("COS 连接失败: %v", err)
    }
    return nil
}
```

#### [MODIFY] `internal/pkg/storage/cstcloud.go`

为 `CSTCloudStorage` 实现 `TestConnection()`：
```go
func (s *CSTCloudStorage) TestConnection() error {
    // 检查 bucket 是否存在
    exists, err := s.client.BucketExists(context.Background(), s.bucketName)
    if err != nil {
        return fmt.Errorf("S3 连接失败: %v", err)
    }
    if !exists {
        return fmt.Errorf("Bucket '%s' 不存在", s.bucketName)
    }
    return nil
}
```

#### [MODIFY] `internal/service/setting_service.go`

新增 `TestStorageConnection` 方法：
```go
import (
    "fmt"
    "test-ebook-api/internal/model"
    "test-ebook-api/internal/pkg/ocr"
    "test-ebook-api/internal/pkg/storage"
    "test-ebook-api/internal/config"
    "test-ebook-api/internal/repository"
)

func (s *SettingService) TestStorageConnection(storageType string, params map[string]interface{}) error {
    // 构建临时 StorageConfig
    cfg := config.StorageConfig{Type: storageType}

    switch storageType {
    case "local":
        cfg.LocalPath, _ = params["local_path"].(string)
        if cfg.LocalPath == "" {
            cfg.LocalPath = "uploads"
        }
    case "aliyun_oss":
        cfg.AliyunEndpoint, _ = params["aliyun_endpoint"].(string)
        cfg.AliyunAccessKeyID, _ = params["aliyun_access_key_id"].(string)
        cfg.AliyunAccessKeySecret, _ = params["aliyun_access_key_secret"].(string)
        cfg.AliyunBucket, _ = params["aliyun_bucket"].(string)
        if cfg.AliyunEndpoint == "" || cfg.AliyunAccessKeyID == "" || cfg.AliyunAccessKeySecret == "" || cfg.AliyunBucket == "" {
            return fmt.Errorf("OSS 配置不完整，请填写所有必填项")
        }
    case "tencent_cos":
        cfg.TencentBucketURL, _ = params["tencent_bucket_url"].(string)
        cfg.TencentSecretID, _ = params["tencent_secret_id"].(string)
        cfg.TencentSecretKey, _ = params["tencent_secret_key"].(string)
        if cfg.TencentBucketURL == "" || cfg.TencentSecretID == "" || cfg.TencentSecretKey == "" {
            return fmt.Errorf("COS 配置不完整，请填写所有必填项")
        }
    case "cstcloud":
        cfg.CSTCloudEndpoint, _ = params["cstcloud_endpoint"].(string)
        cfg.CSTCloudAccessKey, _ = params["cstcloud_access_key"].(string)
        cfg.CSTCloudSecretKey, _ = params["cstcloud_secret_key"].(string)
        cfg.CSTCloudBucket, _ = params["cstcloud_bucket"].(string)
        if cfg.CSTCloudEndpoint == "" || cfg.CSTCloudAccessKey == "" || cfg.CSTCloudSecretKey == "" || cfg.CSTCloudBucket == "" {
            return fmt.Errorf("S3 配置不完整，请填写所有必填项")
        }
    default:
        return fmt.Errorf("不支持的存储类型: %s", storageType)
    }

    // 使用工厂创建临时存储实例并测试
    s, err := storage.NewStorage(cfg)
    if err != nil {
        return fmt.Errorf("创建存储实例失败: %v", err)
    }
    return s.TestConnection()
}
```

**注意**：由于 `s` 变量名与 receiver `s *SettingService` 冲突，实际编码时需将存储实例命名为 `store` 或 `storageInst`。

#### [MODIFY] `internal/handler/setting.go`

新增 `TestStorage` Handler 函数：
```go
func (h *SettingHandler) TestStorage(c *gin.Context) {
    var req struct {
        Type   string                 `json:"type"`
        Config map[string]interface{} `json:"config"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
        return
    }
    if err := h.svc.TestStorageConnection(req.Type, req.Config); err != nil {
        pkg.Error(c, http.StatusBadRequest, 400, "连接失败: "+err.Error())
        return
    }
    pkg.Success(c, nil)
}
```

#### [MODIFY] `internal/router/router.go`

在第 99 行（`POST /settings/ocr-test`）之后新增：
```go
protected.POST("/settings/storage-test", settingHandler.TestStorage)
```

### 4.3 前端修改

#### [MODIFY] `src/api/settings.ts`

新增 `testStorageConnection` 函数：
```ts
export function testStorageConnection(storageSettings: any) {
  return request({
    url: '/settings/storage-test',
    method: 'post',
    data: {
      type: storageSettings.type,
      config: storageSettings
    }
  })
}
```

#### [MODIFY] `src/views/admin/settings/SettingsPage.vue`

1. **import 新增**：在第 153 行导入中增加 `testStorageConnection`：
   ```ts
   import { getSettings, saveSettings, testOcrConnection, testStorageConnection } from '@/api/settings'
   ```

2. **新增响应式变量**：
   ```ts
   const testingStorage = ref(false)
   ```

3. **新增函数** `handleTestStorageConnection`：
   ```ts
   const handleTestStorageConnection = async () => {
     testingStorage.value = true
     try {
       await testStorageConnection(settings.storage)
       ElMessage.success('存储连接测试成功')
     } catch (error) {
       // 错误自动由 request.ts 拦截弹窗
     } finally {
       testingStorage.value = false
     }
   }
   ```

4. **模板修改**：在存储 Tab 的每个云存储配置区域末尾、"确认修改"按钮之前，增加"测试连接"按钮：

   - 在阿里云 OSS `</template>` 之后（第 106 行后）、腾讯云 `<template>` 之前，**无需单独加**。
   - 统一在第 141~143 行（"确认修改"按钮的 `<el-form-item>` 处）之前添加：

   ```html
   <el-form-item v-if="settings.storage.type !== 'local'">
     <el-button type="success" plain @click="handleTestStorageConnection" :loading="testingStorage">测试存储连接</el-button>
   </el-form-item>
   ```
   
   对于 `local` 类型也可以测试（目录可写性检查），如果需要则去掉 `v-if` 条件。建议保留，改为：
   ```html
   <el-form-item>
     <el-button type="success" plain @click="handleTestStorageConnection" :loading="testingStorage">测试存储连接</el-button>
   </el-form-item>
   ```

---

## 五、文件变更汇总表

| 操作 | 文件路径 | 说明 |
|:---:|--------|------| 
| MOD | `internal/pkg/ocr/paddle_client.go` | 重写 `TestConnection()` 使用 POST + fileUrl 验证 Token |
| MOD | `internal/pkg/storage/interface.go` | 接口新增 `TestConnection() error` 方法 |
| MOD | `internal/pkg/storage/local.go` | 实现 `TestConnection()` 目录可写检查 |
| MOD | `internal/pkg/storage/aliyun_oss.go` | 实现 `TestConnection()` ListObjects 验证 |
| MOD | `internal/pkg/storage/tencent_cos.go` | 实现 `TestConnection()` Bucket.Get 验证 |
| MOD | `internal/pkg/storage/cstcloud.go` | 实现 `TestConnection()` BucketExists 验证 |
| MOD | `internal/service/standard_service.go` | `ProcessFile()` 增加云存储临时文件下载逻辑 |
| MOD | `internal/service/setting_service.go` | 新增 `TestStorageConnection()` 方法 |
| MOD | `internal/handler/setting.go` | 新增 `TestStorage()` Handler |
| MOD | `internal/router/router.go` | 新增路由 `POST /settings/storage-test` |
| MOD | `src/api/settings.ts` | 重写 `getSettings()` + 新增 `testStorageConnection()` + 新增 `unflatten()` |
| MOD | `src/views/admin/settings/SettingsPage.vue` | 修复 `loadData()` 数据回显 + 新增存储测试按钮和 handler |

---

## 六、验证计划

1. **编译验证**：`go build ./...` 确保无错误。
2. **Token 校验**：
   - 使用 OCRTOKEN.MD 中的真实 Token (`02ac...`) 测试连接 → 应返回"成功"。
   - 使用任意假 Token（如 `invalid-token-123`） → 应返回"身份验证失败"。
3. **设置回显**：
   - 在 OCR 设置页填入 Token、选择模型后点击"保存设置"。
   - 刷新页面后，所有值应正确回显（不被清空）。
4. **存储连接测试**：
   - 选择"本地存储"，测试连接 → 应返回"成功"。
   - 选择"阿里云 OSS"，填入错误参数，测试连接 → 应返回具体错误信息。
5. **云存储 OCR**：
   - 配置为本地存储上传文件，OCR 正常 → 验证未破坏原有流程。
   - （如有条件）配置为云存储上传文件，观察日志确认临时文件被创建和清理。
