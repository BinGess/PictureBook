# 绘本服务端接口文档

> 本项目现为纯服务端代码仓库，提供面向多客户端（iOS/Android/Web/小程序）的统一 HTTP API。

## 基础信息
- 基础路径：`/`
- 健康检查：`GET /healthz` → 响应：`ok`
- 静态资源：`/assets` 映射至服务端 `uploads/` 目录（上传接口返回的 URL 即该路径）
- CORS：通过环境变量 `CORS_ALLOWED_ORIGINS` 控制（默认 `*`）
- 认证：公共接口无需认证；管理接口使用 Cookie `admin_session=ok`（由 `/admin/login` 设置）

## 数据模型
- Book
  - `id` 字符串
  - `title` 字符串
  - `coverURL` 字符串（可为相对路径 `/assets/...` 或绝对 URL）
  - `ageMin` 数字
  - `ageMax` 数字
  - `tags` 字符串数组
  - `popularityScore` 数字
  - `themeKeywords` 字符串数组
  - `isEditorPick` 布尔
  - `pages` Page[]
  - `status` 字符串（如 `draft`/`published`）
- Page
  - `id` 字符串
  - `index` 数字
  - `imageURL` 字符串（通常为 `/assets/...`）
  - `duration_hint` 可选数字（JSON 字段为 `duration_hint`）

## 公共接口
- GET `/v1/editor-picks`
  - 描述：获取编辑推荐（最多 5 本）
  - 响应：`Book[]`
  - 示例：
    ```json
    [
      {"id":"book_a","title":"森林探险","coverURL":"/assets/cover_a.jpg","ageMin":3,"ageMax":6,...}
    ]
    ```

- GET `/v1/books`
  - 描述：分页获取书籍列表
  - 查询参数：
    - `page` 整数，默认 `1`
    - `page_size` 整数，默认 `24`
    - `age` 整数，可选，用于年龄适配筛选
    - `sort` 字符串，默认 `popular`（按受欢迎度）
  - 响应：
    ```json
    {
      "items": [Book,...],
      "paging": {"page":1, "page_size":24, "has_more": true}
    }
    ```

- GET `/v1/books/:id`
  - 描述：获取单本书详情
  - 响应：`Book`
  - 异常：书不存在 → `404 {"error":"not_found"}`

- GET `/v1/books/:id/recommendations`
  - 描述：根据当前书与年龄，返回推荐列表
  - 查询参数：
    - `age` 整数，可选，默认为该书年龄区间中点
    - `limit` 整数，可选，默认 `5`
  - 响应：`Book[]`

## 管理与上传接口
- POST `/v1/admin/upload`
  - 描述：上传单文件，返回文件 URL
  - 请求：`multipart/form-data`，字段 `file`
  - 响应：`{"file_url":"/assets/<name>"}`
  - 异常：
    - 缺少文件 → `400 {"error":"missing_file"}`
    - 保存失败 → `500 {"error":"save_failed"}`

- POST `/v1/admin/books`
  - 描述：创建书籍（JSON）
  - 请求体：`Book`（`status` 为空时服务端默认 `draft`）
  - 约束：若 `isEditorPick=true`，编辑推荐上限为 5 本
  - 响应：创建后的 `Book`
  - 异常：
    - 请求体格式错误 → `400 {"error":"invalid_body"}`
    - 超过编辑推荐上限 → `400 {"error":"editor_picks_limit"}`

- POST `/v1/admin/books/:id/pages/upload`
  - 描述：批量上传页图片并添加到指定书籍
  - 请求：`multipart/form-data`，字段 `files`（多文件）
  - 响应：`{"pages": Page[]}`（索引顺序从现有页数后续递增）
  - 异常：
    - 无文件或解析失败 → `400 {"error":"missing_files"}`
    - 保存失败 → `500 {"error":"save_failed"}`

- POST `/admin/books/:id/pages/reorder`
  - 描述：调整页面顺序（后台表单接口）
  - 请求：`application/x-www-form-urlencoded`，键形如 `index[<pageID>]=<newIndex>`
  - 响应：302 重定向到 `/admin/books/:id/pages`
  - 异常：书不存在 → `404 {"error":"not_found"}`

## 后台登录（设置 Cookie）
- GET `/admin/login`：返回登录页面
- POST `/admin/login`
  - 表单字段：`email`、`password`
  - 示例账户：`email=admin@example.com`、`password` 非空
  - 成功：设置 `admin_session=ok`，重定向 `/admin/books`
  - 失败：返回登录页并显示错误
- 需要鉴权的后台页面：
  - `GET /admin/books`（列表）
  - `GET /admin/books/new`（新建）→ `POST /admin/books/new`
  - `GET /admin/books/:id/pages`（页面管理）
  - `POST /admin/books/:id/editor-pick`（切换编辑推荐）

## 错误与状态码
- 200 成功；302 重定向（后台表单）；400 参数或体错误；404 资源不存在；500 服务器错误
- 常见错误体：
  - `{"error":"invalid_body"}`、`{"error":"missing_file"}`、`{"error":"missing_files"}`
  - `{"error":"editor_picks_limit"}`、`{"error":"save_failed"}`、`{"error":"not_found"}`

## 示例
- 获取编辑推荐：
  ```bash
  curl -s https://<host>/v1/editor-picks
  ```
- 分页获取书籍：
  ```bash
  curl -s 'https://<host>/v1/books?page=1&page_size=24&sort=popular'
  ```
- 创建书籍：
  ```bash
  curl -s -X POST https://<host>/v1/admin/books \
    -H 'Content-Type: application/json' \
    -d '{
      "id":"book_x",
      "title":"新书",
      "coverURL":"/assets/cover_x.jpg",
      "ageMin":3,
      "ageMax":6,
      "tags":["森林"],
      "popularityScore":0.7,
      "themeKeywords":["自然"],
      "isEditorPick":false,
      "pages":[],
      "status":"published"
    }'
  ```
- 上传封面：
  ```bash
  curl -s -X POST https://<host>/v1/admin/upload \
    -F 'file=@/path/to/cover.jpg'
  ```
- 批量上传页图：
  ```bash
  curl -s -X POST https://<host>/v1/admin/books/book_x/pages/upload \
    -F 'files=@/path/p1.jpg' -F 'files=@/path/p2.jpg'
  ```

## 备注
- 静态资源目录为 `uploads/`，服务端会返回以 `/assets/` 为前缀的可访问 URL。
- 首次启动空数据库时，服务端会自动填充示例数据（用于开发验证）。