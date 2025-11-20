## 目标

* 新增“分类”实体（名称、描述），支持创建/删除分类（后台页面）。

* 图书与分类建立从属关系：创建/编辑图书时必须选择所属分类。

* 对外数据结构支持“分类列表 + 每个分类下的书籍列表”。

## 数据库与模型

* 新增表 `categories`：`id TEXT PRIMARY KEY`、`name TEXT NOT NULL`、`description TEXT`。

* 在 `books` 表新增字段：`category_id TEXT NOT NULL`，外键指向 `categories(id)`。

* 外键策略：`ON DELETE RESTRICT`（默认阻止删除仍有书籍的分类，避免数据悬空）。

* 模型：新增 `Category{ID,Name,Description}`；`Book` 增加 `CategoryID`（可在管理端与公共端响应中按需透出分类名称）。

* 迁移：

  * `EnsureSchema` 中添加 `categories` 表定义与 `books.category_id` 字段（若已有数据，初始允许 `NULL`，再通过脚本或后台批量指定分类后切换为 `NOT NULL`）。

## 仓储层（repository）

* `CreateCategory(name, description) error`

* `DeleteCategory(id) error`（RESTRICT，如有书籍返回业务错误 `category_has_books`）

* `ListCategories() ([]Category, error)`

* `ListCategoriesWithBooks() ([]CategoryWithBooks, error)`：返回分类及其下的 `Book[]`

* `ListBooksAdmin(sort, page, size, q)` 与 `CreateBook/UpdateBook`：增强以处理 `category_id`

## 服务层（services）

* 包装仓储方法：`CreateCategory/DeleteCategory/ListCategories/ListCategoriesWithBooks`

* 图书创建/编辑时校验分类存在：`ValidateCategoryID`

## 路由与接口

* 公共 API：

  * `GET /v1/categories` → 返回 `[{id,name,description}]`

  * `GET /v1/categories-with-books` → 返回 `[{id,name,description,books:[Book,...]}]`

  * `GET /v1/categories/:id/books` → 单分类下书籍列表分页（`page/page_size`）

  * 现有 `GET /v1/books` / `GET /v1/books/:id` 响应中可附带 `category_id` 与 `category_name`（可选）

* 管理 API：

  * `POST /v1/admin/categories`（JSON：`{name,description}`）创建分类

  * `DELETE /v1/admin/categories/:id` 删除分类（无书籍时成功；有书籍返回 `400 {error:"category_has_books"}`）

  * 现有 `POST /v1/admin/books` / `POST /admin/books/new` / `POST /admin/books/:id/edit` 接口与表单增加 `category_id`

## 后台页面

* 新增“分类管理”页面与导航入口：

  * 列表：`/admin/categories`（展示名称、描述、书籍数量、删除按钮）

  * 新建：`/admin/categories/new`（名称、描述表单）

  * 删除：`POST /admin/categories/:id/delete`（若有书籍，将提示无法删除）

* 图书表单改造：

  * `book_new.html` 与 `book_edit.html` 增加分类下拉框（从 `/admin/categories` 加载，或服务端渲染注入 `.categories`）

  * 表单提交包含 `category_id`

* 页面管理页可显示当前图书的所属分类（只读），并提供“跳转到编辑详情页”修改分类。

## 响应结构示例

* `GET /v1/categories-with-books`：

```json
[
  {
    "id":"cat_story",
    "name":"故事",
    "description":"适合3-6岁故事类",
    "books":[{"id":"book_a","title":"森林探险",...}]
  },
  {
    "id":"cat_science",
    "name":"科学",
    "description":"科普类",
    "books":[...]
  }
]
```

## 兼容性与迁移策略

* 旧数据：

  * 首次部署先允许 `books.category_id` 为 `NULL`，后台提供批量或逐本指定分类工具；完成归类后将列约束改为 `NOT NULL`。

* 删除分类：

  * 有书籍时不允许删除；后续可扩展“批量迁移到其他分类”再删除。

## 文档与验证

* 更新 `API.md`：新增分类相关接口与字段说明；示例响应。

* 管理端操作流：创建分类 → 新建书籍选择分类 → 列表与详情显示分类；删除分类时的提示与限制。

## 交付步骤

1. 数据库与模型改造（`EnsureSchema`、`models`）并加入迁移判断逻辑。
2. 实现仓储与服务方法，覆盖分类 CRUD 与聚合查询。
3. 增加公共与管理路由及处理器，返回与校验完整。
4. 改造后台模板（列表/新建/编辑/分类管理页），注入分类数据，完善交互。
5. 更新 `API.md` 与基础示例；本地或容器内通过 curl/浏览器验证。

## 需确认的选项

* 删除分类策略：是否允许“迁移分类后再删除”？如需要，我将为后台列表添加“迁移到...”选择框。 需要

* 响应是否在 `GET /v1/books` 中附带 `category_name`（便于客户端直接展示）。是的，需要附带`category_name`

