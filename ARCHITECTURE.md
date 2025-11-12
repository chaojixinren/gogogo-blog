# GoGo 博客系统架构说明

文档概述系统组成、后端与前端的模块职责以及典型数据流，方便团队成员快速定位与扩展功能。

---

## 1. 系统概览

- **技术栈**
  - 后端: Go 1.25, Gin, GORM, MySQL
  - 前端: Vite 7, Vue 3, TypeScript, Pinia, Vue Router, Axios
  - 鉴权: JWT
- **整体架构**

```
┌──────────────────────┐      HTTP / JSON      ┌────────────────────────┐      ┌──────────┐
│  前端 (Vue 3 + Vite) │ <-------------------->│ 后端 (Go + Gin + GORM) │<---->│  MySQL   │
└──────────────────────┘        Axios          └────────────────────────┘      └──────────┘
        SPA 客户端                    REST API 与业务逻辑                   数据持久化
```

- **核心能力**
  - 用户注册、登录与 JWT 鉴权
  - 文章 CRUD、分页、过滤、 slug 管理
  - 分类、标签维护与文章关联
  - 游客 / 登录用户评论
  - Dashboard 后台管理

---

## 2. 后端架构

### 2.1 目录结构

| 路径 | 说明 |
|------|------|
| `main.go` | 程序入口, 初始化配置与 HTTP 服务 |
| `config/` | 配置读取 (`config.go`), 数据库初始化 (`db.go`), 默认配置 (`config.yml`) |
| `global/` | 全局共享对象, 当前仅持有 `*gorm.DB` |
| `router/` | Gin 路由与 CORS 配置 |
| `middleware/` | 自定义中间件 (JWT 鉴权) |
| `controllers/` | 业务控制器, 返回 JSON 响应 |
| `models/` | GORM 数据模型与关联定义 |
| `utils/` | 密码、JWT、分页、slug 等通用函数 |

### 2.2 启动流程

1. `main.go` 执行 `config.InitConfig()`
2. `InitConfig` 读取 `config.yml`, 再调用 `InitDB`
3. `InitDB` 连接 MySQL, 配置连接池, 对模型执行 `AutoMigrate`
4. `router.SetupRouter()` 注册路由、中间件
5. Gin 按 `config.app.port` 监听服务

### 2.3 配置字段

| 节点 | 字段 | 说明 |
|------|------|------|
| `app` | `name`, `port` | 应用名称 (用于 JWT issuer) 及监听端口 |
| `database` | `dsn`, `max_idle_conns`, `max_open_conns` | MySQL 连接及连接池参数 |
| `auth` | `jwt_secret`, `token_ttl_hours` | JWT 签名密钥和有效期 (小时) |
| `cors` | `allow_origins` | 允许的跨域来源列表 |

### 2.4 数据模型

所有模型均嵌入 `gorm.Model`, 因此包含 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`。

| 模型 | 关键字段 | 关联 |
|------|----------|------|
| `User` | `Username`, 可选 `Email`, `Password`, `DisplayName`, `Bio`, `AvatarURL` | `Posts` 一对多, `Comments` 一对多 |
| `Category` | `Name`, `Slug`, `Description` | `Posts` 一对多 |
| `Tag` | `Name`, `Slug` | 与 `Post` 多对多 (`post_tags`) |
| `Post` | `Title`, `Summary`, `Content`, `Slug`, `Status`, `CoverImage`, `PublishedAt` | 关联 `Author`, 可选 `Category`, 多对多 `Tags`, `Comments` |
| `Comment` | `PostID`, 可选 `UserID`, `AuthorName`, `Body`, `Approved` | 关联 `Post`, 可选 `User` |

### 2.5 工具与中间件

- `utils/utils.go`: `HashPassword`, `CheckPassword`, `GenerateJWT`, `ValidateJWT`
- `utils/slug.go`: 文本转 slug
- `utils/pagination.go`: 解析并约束分页参数
- `middleware/auth_middleware.go`: 解析 Authorization 头, 校验 JWT, 注入用户信息

### 2.6 控制器概览

| 控制器 | 功能 |
|--------|------|
| `auth_controller.go` | 注册、登录, 生成 token |
| `user_controller.go` | `GET /api/me`, `GET /api/me/posts` |
| `post_controller.go` | 文章 CRUD, 过滤, slug 唯一性, 标签懒创建 |
| `category_controller.go` | 分类 CRUD, slug 校验, 删除时解绑文章 |
| `tag_controller.go` | 标签 CRUD, slug 校验, 维护多对多关系 |
| `comment_controller.go` | 评论列表与创建, 区分游客/登录用户 |

DTO 定义在 `controllers/dto.go`, 隐藏敏感字段 (如密码、邮箱)。

### 2.7 路由布局

```
/api
├─ /auth/login, /auth/register
├─ /health
├─ /posts, /posts/:id, /posts/slug/:slug
├─ /posts/:id/comments
├─ /categories, /tags
└─ [AuthMiddleware]
   ├─ /me, /me/posts
   ├─ /posts (POST)
   ├─ /posts/:id (PUT, DELETE)
   ├─ /categories (POST, PUT, DELETE)
   └─ /tags (POST, PUT, DELETE)
```

### 2.8 运行与测试

```bash
go build ./...   # 编译检查
go run .         # 运行后端服务 (默认端口 :3000)
```

---

## 3. 前端架构

### 3.1 模块总览

- 单页应用, 使用 Vite 构建
- Vue 3 组合式 API
- Pinia 负责状态, Vue Router 负责路由
- Axios 统一请求, 拦截器自动附带 token
- 使用 CSS 变量定义主题与基础组件样式

### 3.2 目录结构

| 路径 | 说明 |
|------|------|
| `src/main.ts` | 创建应用, 注册 Pinia, 恢复登录状态, 挂载 router |
| `src/App.vue` | 顶层布局, 包含 Header、RouterView、Footer |
| `src/router/index.ts` | 路由配置与守卫, 控制鉴权与重定向 |
| `src/store/auth.ts` | Pinia 鉴权仓库: token 持久化、登录、注册、退出、刷新资料 |
| `src/services/` | Axios 实例 (`api.ts`) 和各业务请求封装 |
| `src/types/index.ts` | TypeScript 接口, 定义 User, Post, 分页等类型 |
| `src/assets/styles/main.css` | 全局样式与主题变量 |
| `src/components/` | 可复用组件 (Header, Footer, PostCard, Pagination 等) |
| `src/pages/` | 页面组件 (首页、详情、登录、注册、Dashboard) |
| `src/pages/dashboard/` | Dashboard 子页面 (文章、分类、标签管理) |
| `vite.config.ts` | Vite 配置, 设置别名 `@`, `/api` 代理 |
| `tsconfig.json` | TypeScript 配置, 启用严格模式与路径映射 |

### 3.3 启动流程

1. `createApp(App)`
2. 安装 Pinia, 调用 `useAuthStore().initialize()` 恢复 token
3. 注册 router, 处理导航守卫
4. `app.mount('#app')`

### 3.4 路由设计

| 路径 | 组件 | 说明 |
|------|------|------|
| `/` | `HomePage` | 文章列表, 支持搜索/分类/标签筛选 |
| `/posts/:slug` | `PostDetailPage` | 文章详情, Markdown 渲染, 评论 |
| `/login` | `LoginPage` | 登录, 已登录用户会被重定向 |
| `/register` | `RegisterPage` | 注册, 已登录用户会被重定向 |
| `/dashboard` | `DashboardLayout` | 受保护布局, 默认重定向 `/dashboard/posts` |
| `/dashboard/posts` | `DashboardPostsPage` | 文章创建/编辑/删除, 分页查看我的文章 |
| `/dashboard/categories` | `DashboardCategoriesPage` | 分类 CRUD |
| `/dashboard/tags` | `DashboardTagsPage` | 标签 CRUD |

守卫策略:
- 首次进入时等待 `auth.initialize()`
- 未登录访问受保护路由 -> 跳转 `/login?redirect=...`
- 已登录访问 `/login` 或 `/register` -> 重定向 `/`

### 3.5 状态管理 (`src/store/auth.ts`)

- 状态: `user`, `token`, `isInitialized`, `isProcessing`, `error`
- 计算属性: `isAuthenticated`
- 方法: `initialize`, `login`, `register`, `logout`, `refreshProfile`, `setSession`
- Token 持久化在 localStorage, Axios 拦截器自动附带 Authorization 头

### 3.6 服务层

| 文件 | 职责 |
|------|------|
| `api.ts` | 创建 Axios 实例, 设置超时与请求拦截器 |
| `auth.ts` | 登录、注册、获取当前用户 |
| `posts.ts` | 文章列表、详情、我的文章、创建/更新/删除 |
| `categories.ts` | 分类列表、CRUD、按分类拉取文章 |
| `tags.ts` | 标签列表、CRUD、按标签拉取文章 |
| `comments.ts` | 评论列表与创建 |

所有函数返回 Promise, 类型定义放在 `src/types/index.ts`。

### 3.7 组件与页面

- 布局组件: `AppHeader`, `AppFooter`, `DashboardLayout`
- 通用组件: `PostCard`, `TagChip`, `PaginationControls`
- 页面要点:
  - `HomePage`: 同步查询参数到 URL, 支持组合过滤
  - `PostDetailPage`: `marked` + `dompurify` 渲染 Markdown, 评论支持游客输入昵称
  - `DashboardPostsPage`: 表单与列表同屏, 标签支持建议与批量输入
  - `DashboardCategoriesPage` / `DashboardTagsPage`: 表单 + 列表式 CRUD

### 3.8 样式规范

- `main.css` 定义颜色、阴影、字体等 CSS 变量
- 统一按钮、输入、卡片、徽章的视觉风格
- 响应式布局: 在 768px 以下调整栅格与内边距

### 3.9 构建与运行

```bash
npm install
npm run dev      # http://localhost:5173
npm run build    # 生成 dist/
npm run preview  # 预览生产构建
```

开发模式下 Vite 将 `/api/*` 代理到 `http://localhost:3000`。

---

## 4. API 契约摘要

| 模块 | 方法与路径 | 说明 |
|------|------------|------|
| 认证 | `POST /api/auth/register` | 返回 `{ token, user }` |
|      | `POST /api/auth/login` | 返回 `{ token, user }` |
| 用户 | `GET /api/me` | 当前用户信息 |
|      | `GET /api/me/posts` | 当前用户文章 (分页) |
| 文章 | `GET /api/posts` | 列表, 支持分页与多条件筛选 |
|      | `GET /api/posts/:id` / `/slug/:slug` | 文章详情 |
|      | `POST /api/posts` | 创建文章 |
|      | `PUT /api/posts/:id`, `DELETE /api/posts/:id` | 更新 / 删除文章 |
| 评论 | `GET /api/posts/:id/comments` | 评论列表 |
|      | `POST /api/posts/:id/comments` | 创建评论 |
| 分类 | `GET /api/categories` | 分类列表 |
|      | `POST/PUT/DELETE /api/categories[:id]` | 分类 CRUD |
| 标签 | `GET /api/tags` | 标签列表 |
|      | `POST/PUT/DELETE /api/tags[:id]` | 标签 CRUD |

分页接口统一返回:

```json
{
  "data": [],
  "page": 1,
  "pageSize": 10,
  "total": 0
}
```

---

## 5. 典型数据流

1. **注册/登录**
   - 前端调用 `authService` -> Pinia 存储 token -> Axios 附带 Authorization
   - 后端校验用户、哈希密码、生成 JWT

2. **发布文章**
   - Dashboard 表单提交 -> `POST /api/posts`
   - 后端校验作者, 处理标签与 slug -> 返回 `PostDTO`
   - 前端刷新我的文章列表并重置表单

3. **文章详情与评论**
   - 前端 `GET /api/posts/slug/:slug`
   - 后端预加载作者、分类、标签、已审核评论
   - 评论提交 `POST /api/posts/:id/comments` -> 刷新评论列表

---

## 6. 开发与部署建议

1. 配置文件随环境调整, JWT 密钥勿入库
2. 可接入 Zap/Logrus 等日志组件, 丰富日志格式
3. 在现有构建检查基础上增加单元测试、E2E 测试
4. 后续扩展方向: 角色权限、编辑器、文件上传、评论审核、国际化

---

## 7. 快速启动

```bash
# 后端
cd backend
go run .

# 前端
cd ../frontend
npm install
npm run dev
```

访问 `http://localhost:5173` 即可体验完整博客流程。

---

如文档与代码不符, 请以代码为准并及时更新此文件。 

