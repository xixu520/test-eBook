## 1. 规范约定

  

- **Base URL**: `/api/v1`

- **请求格式**: 默认为 `application/json`，文件上传使用 `multipart/form-data`。

- **响应格式**: 统一返回 JSON 格式数据。

- **鉴权方式**: Header 中携带 JWT Token，格式为 `Authorization: Bearer <token>`。

 - **全局错误码说明**

| 错误码 | 含义          | 处理建议            |
| --- | ----------- | --------------- |
| 200 | 操作成功        | -               |
| 400 | 参数错误/业务约束   | 提示用户检查输入/操作流程   |
| 401 | 未登录/Token失效 | 引导用户重新登录        |
| 403 | 无权限         | 提示用户联系管理员       |
| 500 | 服务器错误       | 提示用户稍后重试，记录错误日志 |


### 1.1 统一响应数据结构 (JSON)

  

所有 API 接口的返回数据都遵循以下基础结构：

  

```json

{

  "code": 200,            // 业务状态码：200成功，401未登录，403无权限，400参数错误，500服务器错误

  "message": "操作成功",   // 提示信息

  "data": { ... }         // 具体的业务数据载荷，无数据时为空对象 {} 或 null

}

```

  

### 1.2 分页请求通用参数

  

对于列表类型的接口（如获取标准列表），通过 URL Query 传递：

- `page`: 当前页码，默认 1

- `size`: 每页条数，默认 10

  

分页响应 `data` 结构：

```json

{

  "total": 150,           // 总记录数

  "page": 1,              // 当前页码

  "size": 10,             // 每页条数

  "list": [ ... ]         // 当前页的数据列表

}

```

  

---

  

## 2. 认证与用户模块
  
### 2.1 用户登录

- **URL**: `/auth/login`

- **Method**: `POST`

- **描述**: 用户登录，获取身份凭证 Token。

- **请求参数** (JSON):

```json

  {

    "username": "admin",

    "password": "password123",
    
    "captcha_code": "x2y9"

  }

```

- **响应数据**:

```json

  {

    "token": "eyJhbGciOiJIUzI1NiIs...",

    "user": {

      "id": 1,

      "username": "admin",

      "role": "admin",          // 角色标识：admin(管理员), user(普通用户)

      "theme": "light"          // 用户偏好的主题色调

    }

  }

```

  

### 2.2 获取当前用户信息

- **URL**: `/auth/me`

- **Method**: `GET`

- **描述**: 获取当前登录用户的详细信息及权限列表。

- **响应数据**:

```json

  {

    "id": 1,

    "username": "admin",

    "role": "admin",

    "permissions": ["upload", "download", "delete", "verify", "manage_category"]

  }

```

  

### 2.3 更新用户主题偏好

- **URL**: `/users/me/theme`

- **Method**: `PUT`

- **请求参数** (JSON):

```json

  {

    "theme": "dark" // 主题标识，如 light, dark, blue 等

  }

```

  

---

  

## 3. 分类管理模块 (管理员权限)
### 3.1 获取分类列表

- **URL**: `/categories`

- **Method**: `GET`

- **描述**: 获取左侧文件分类树。包含每个分类下的文档总数。

- **响应数据**:

```json

  [
  { 
    "id": 1, 
    "parent_id": 0, 
    "name": "地基基础", 
    "sort_order": 1, 
    "doc_count": 128 
  },
  { 
    "id": 2, 
    "parent_id": 0, 
    "name": "建筑材料", 
    "sort_order": 2, 
    "doc_count": 85 
  },
  { 
    "id": 3, 
    "parent_id": 1, 
    "name": "桩基础", 
    "sort_order": 1, 
    "doc_count": 43 
  }
]

```

  

### 3.2 新增分类

- **URL**: `/categories`

- **Method**: `POST`

- **请求参数** (JSON):

```json

{ 
  "parent_id": 0,
  "name": "消防检测", 
  "sort_order": 3 
}

```

  

### 3.3 修改分类

- **URL**: `/categories/{id}`

- **Method**: `PUT`

- **请求参数** (JSON):

  ```json

{ 
  "parent_id": 0,
  "name": "消防安全检测", 
  "sort_order": 3 
}

  ```

- **响应数据 (成功)**:

```json
{
  "code": 200,
  "message": "操作成功",
  "data": null
}
```
  
- **响应数据 (失败 - 触发根节点降级限制)**:

```json
{
  "code": 400,
  "message": "非法的操作：禁止将顶级分类移动至其他分类下",
  "data": null
}
```


### 3.4 删除分类

- **URL**: `/categories/{id}`

- **Method**: `DELETE`

- **描述**: 删除指定的分类节点。

	- **业务约束（重要）**：
	1. **防误删保护（树形结构）**：如果该分类下存在**子分类**，则系统拒绝删除操作。 
	2. **数据关联保护**：如果该分类下仍有关联的**标准文件**（即对应库表中该分类的文件数量 `> 0`），则系统拒绝删除。前端需提示用户先将文件转移至其他分类或移入回收站。
- **响应数据 (成功)**:

```json
{ "code": 200, "message": "分类删除成功", "data": null }
```


- **响应数据 (失败 - 触发子分类限制)**:

```json
{ "code": 400, "message": "删除失败：该分类下存在子分类，请先移除或转移子分类", "data": null }
```

- **响应数据 (失败 - 触发文件关联限制)**:

```json
{ "code": 400, "message": "删除失败：该分类下仍有标准文件，请先转移或删除文件", "data": null }
```


---

  

## 4. 标准文件模块

### 4.1 获取标准文件列表

- **URL**: `/documents`

- **Method**: `GET`

- **描述**: 首页搜索与列表展示。不同角色返回的字段粒度不同。

- **请求参数** (Query):

  - `page`: 1

  - `size`: 10

  - `keyword`: "规范" (匹配标准号或名称)

  - `category_id`: 2 (按分类筛选)

- **响应数据**:

```json
  {

    "total": 100,

    "list": [

      {

        "id": 101,

        "standard_no": "GB50007-2011",

        "name": "建筑地基基础设计规范",

        "category_id": 1,

        "category_name": "地基基础",

        "issue_date": "2011-07-26",           // 允许为 null

        "implement_date": "2012-08-01",       // 允许为 null

        "obsolete_date": null,                // 允许为 null

        "ocr_status": "completed",   // pending, completed, failed

        "verify_status": "pass",     // pending, pass, retry

        // ----------------------------------------------------

        // 以下字段后端根据登录用户的角色决定是否返回

        "uploader_name": "张三",

        "upload_time": "2023-10-01 10:00:00",

        "view_count": 150,

        "download_count": 30,

        "active_version": "v1.0"

        // ----------------------------------------------------

      }

    ]

  }
  
```

  

### 4.2 获取文件预览地址 (PDF.js 使用)

- **URL**: `/documents/{id}/preview`

- **Method**: `GET`

- **描述**: 用于前端 PDF.js 渲染预览。

- **参数 (Query)**:

  - `highlight`: "搜索词" (可选，传递后 PDF.js 尝试在页面内高亮并定位)

- **响应数据**: 返回 PDF 文件流。

*(备注：在请求此接口时，后端会自动在此文档的 `view_count` 上 +1)*

  

### 4.3 下载标准文件

- **URL**: `/documents/{id}/download`

- **Method**: `GET`

- **描述**: 仅限有下载权限的用户调用。

- **响应数据**: 返回文件流，包含 `Content-Disposition: attachment`。

*(备注：调用成功会触发 `download_count` +1)*

  

### 4.4 上传标准文件 (后台管理)

- **URL**: `/documents/upload`

- **Method**: `POST`

- **描述**:

  1. Go 后端检查文件大小（**上限 50MB**）。

  2. 计算文件 **SHA256 哈希值**，若库中已有相同哈希且路径存在，则复用文件记录，不重复保存物理文件。

  3. 投递任务至 Python 服务进行 OCR。

  4. OCR 识别出的“标准号、名称、发行日期”会自动填充到文档基础属性中。

- **请求参数** (FormData):

  - `file`: (二进制 PDF 文件)

  - `category_id`: 1

- **响应数据**:



```json

  {

    "task_id": "task_123456789",

    "document_id": 101,

    "status": "processing",

    "is_duplicate": false // 是否命中了秒传(查重)

  }

```

  

### 4.5 获取 OCR 任务进度 (轮询接口)

- **URL**: `/tasks/{task_id}/status`

- **Method**: `GET`

- **描述**: 前端通过上传接口返回的 `task_id` 轮询识别进度。

- **响应数据**:



```json

  {

    "task_id": "task_123456789",

    "status": "completed", // 取值: pending(排队中), processing(识别中), completed(已完成), failed(失败)

    "progress": 100,       // 进度百分比 0-100

    "result": {            // 仅在 status 为 completed 时返回解析到的基础信息

       "standard_no": "GBXXXX-XXXX",

       "name": "xxxx规范",

       "issue_date": "2023-01-01"

    },

    "error": null          // 失败时的错误描述

  }

```

  

### 4.6 修改文件基础信息

- **URL**: `/documents/{id}`

- **Method**: `PUT`

- **描述**: 只能修改基础文本信息，不支持覆盖文件。

- **请求参数** (JSON):



```json

  {

    "standard_no": "GB50007-2011",

    "name": "建筑地基基础设计规范(修订版)",

    "issue_date": "2011-07-26",

    "implement_date": "2012-08-01",

    "obsolete_date": null,

    "category_id": 1

  }

```

  

### 4.7 软删除文件 (移入回收站)

- **URL**: `/documents/{id}`

- **Method**: `DELETE`

- **描述**: 列表中的删除操作。

  

### 4.8 批量分类变更

- **URL**: `/documents/bulk-category`

- **Method**: `POST`

- **请求参数** (JSON):



```json

  {

    "document_ids": [101, 102],

    "target_category_id": 3

  }

```

  

### 4.9 导出标准信息 Excel

- **URL**: `/documents/export`

- **Method**: `GET`

- **描述**: 导出当前筛选条件下的所有标准基础信息。

- **响应数据**: 返回 `.xlsx` 文件流。

  

---

  

## 5. OCR 与核验模块

### 5.1 重新进行 OCR 识别

- **URL**: `/documents/{id}/ocr/retry`

- **Method**: `POST`

- **描述**: 用户发现 OCR 解析错误时，手动触发对该文档 PDF 第一页（封面）的重新提取与识别。同样返回 `task_id`。

- **响应数据**:



```json

  {

    "task_id": "task_retry_98765",

    "status": "pending"

  }

```

  

### 5.2 更新人工核验状态

- **URL**: `/documents/{id}/verify`

- **Method**: `PUT`

- **描述**: 管理员通过可靠网站比对后，手动标记文件的核验结果。

- **请求参数** (JSON):



```json

  {

    "verify_status": "pass"

  }

```

  

---

  

## 6. 回收站模块 (管理员权限)
### 6.1 获取回收站列表

- **URL**: `/recycle-bin/documents`

- **Method**: `GET`

- **描述**: 获取被软删除的文件列表。

  

### 6.2 批量还原

- **URL**: `/recycle-bin/documents/restore`

- **Method**: `PUT`

- **请求参数** (JSON):



```json

  {

    "document_ids": [101, 102]

  }
```

  

### 6.3 批量彻底删除 (及清空回收站)

- **URL**: `/recycle-bin/documents/batch-delete`

- **Method**: `POST`

- **描述**: 从回收站中彻底删除指定的文件。此操作不可逆。

- **请求参数** (JSON):


```json

  {
{
  "document_ids": [101, 102], // 选填。当 empty_all 为 false 时，删除这里指定的 ID。
  "empty_all": true           // 必填。是否清空回收站？如果传 true，则忽略 document_ids，直接清空所有文件。
}

```


- **响应数据**:



```json
{ "code": 200, "message": "操作成功", "data": null }
```



---

  

## 7. 审计日志模块 (管理员权限)

  

### 7.1 获取操作日志

- **URL**: `/audit-logs`

- **Method**: `GET`

- **请求参数** (Query): `page`, `size`, `user_id`, `action`

- **响应数据**:


```json

  {

    "total": 500,

    "list": [

      {

        "id": 1,

        "timestamp": "2023-11-01 10:00:00",

        "username": "admin",

        "action": "UPLOAD", // UPLOAD, DELETE, VERIFY, EDIT, LOGIN

        "target_id": 101, // 操作的对象ID

        "ip": "192.168.1.10"

      }

    ]

  }

```

  

---

  

## 8. 系统公告模块
### 8.1 获取当前系统公告

- **URL**: `/announcements/active`

- **Method**: `GET`

- **描述**: 供首页展位调用显示。

- **响应数据**:


```json

  {

    "content": "系统将于本周五晚进行升级维护...",

    "update_time": "2023-10-31 09:00:00"

  }

```

  

---

  

## 9. 异步 OCR 处理架构说明
### 9.1 交互流程

1. **上传**: 前端通过 `POST /documents/upload` 发送 PDF 到 **Go 后端**。

2. **接收**: Go 后端保存文件到存储系统（如 MinIO 或本地目录），并在数据库创建一条 `ocr_status='pending'` 的记录。

3. **分配**: Go 后端生成 `task_id` 并将其与 `document_id` 关联。

4. **触发**: Go 后端通过 **消息队列 (Redis/RabbitMQ)** 或 HTTP 调用，将 `task_id` 和文件路径发送给 **Python OCR 服务**。

5. **轮询**: 后端立即向前端返回 `task_id`。前端开始循环调用 `GET /tasks/{task_id}/status`。

6. **处理**: Python 服务使用百度云 OCR 解析 PDF 第一页，解析完成后将结果写入数据库，并将任务状态更新为 `completed`。

7. **完成**: 前端轮询到 `status: 'completed'`，获取解析结果并刷新页面展示。

  

### 9.2 异常处理

- 若 Python 服务宕机或 OCR 超时，任务状态将被置为 `failed`，轮询接口会返回具体的 `error` 信息。

- 前端在 `failed` 状态下应提示用户手动触发 “重新 OCR”。

- **任务超时兜底机制（防卡死）**：
  - **后端兜底**：Go 后端需设置一个定时任务（Cron/Goroutine），定期扫描数据库中状态为 `pending` 或 `processing` 且创建时间超过 **10分钟** 的任务，将其状态强制置为 `failed`，并记录 `error: "任务处理超时"`。
  - **前端兜底**：前端轮询接口的时间上限设为 **5分钟**（例如每 3 秒轮询一次，最多 100 次）。若超过 5 分钟仍未获取到 `completed` 或 `failed` 状态，前端主动终止轮询，向用户展示“识别超时，请稍后重试或检查系统状态”，并允许用户手动点击“重新 OCR”。

  

---

  

## 10. 系统设置模块 (管理员权限)
### 10.1 获取系统设置

- **URL**: `/settings`

- **Method**: `GET`

- **描述**: 获取当前系统配置项。

- **响应数据**:


```json

  {

    "site_name": "建筑标准文件管理系统",

    "allow_registration": false

  }

```

  

### 10.2 更新系统设置

- **URL**: `/settings`

- **Method**: `PUT`

- **描述**: 管理员修改系统配置，支持部分更新。

- **请求参数** (JSON):


```json

  {

    "site_name": "建筑标准文件管理系统",

    "allow_registration": true

  }

```