# goBlog 后端服务

> 基于 Go + Gin + GORM + MySQL + Redis + Elasticsearch 构建的个人博客后端系统

## 📖 项目简介

goBlog 是一个功能完善的个人博客后端，采用 **Gin** 框架构建 RESTful API，集成 **MySQL**（持久化存储）、**Redis**（缓存 & Session）、**Elasticsearch**（全文搜索）三大存储组件，并通过 **七牛云 OSS** 或本地文件系统提供图片上传服务。项目支持邮箱注册、QQ 第三方登录、双 Token（Access Token + Refresh Token）无感刷新认证机制，并内置了定时任务、CLI 工具及结构化日志系统。

---

## 🏗️ 技术栈

| 分类 | 技术 |
|------|------|
| 语言 | Go 1.25+ |
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) v1.12 |
| ORM | [GORM](https://gorm.io) v1.31 + MySQL 驱动 |
| 缓存 | Redis（go-redis v8） |
| 全文搜索 | Elasticsearch v8 |
| 认证 | JWT（双 Token 无感刷新） + Session（gorilla/sessions） |
| 日志 | Zap + Lumberjack（结构化 & 滚动日志） |
| 验证码 | base64Captcha |
| 邮件 | jordan-wright/email（SMTP/SSL） |
| 图片存储 | 七牛云 OSS / 本地文件系统 |
| 定时任务 | robfig/cron v3 |
| CLI 工具 | urfave/cli v1 |
| 配置解析 | go-yaml v3 |

---

## 📁 目录结构

```
server/
├── main.go                 # 程序入口
├── config.yaml             # 运行时配置文件（私有，不入库）
├── config.example.yaml     # 配置文件示例（入库）
├── go.mod / go.sum         # Go 模块依赖
│
├── api/                    # Handler 层（接收请求、参数绑定、调用 Service）
│   ├── enter.go            # API 分组聚合入口
│   ├── base.go             # 验证码、邮箱验证码、QQ 登录 URL
│   ├── user.go             # 用户注册/登录/信息管理/后台管理
│   ├── article.go          # 文章 CRUD、点赞、搜索
│   ├── comment.go          # 评论发布、审核、删除
│   ├── image.go            # 图片上传
│   ├── advertisement.go    # 广告管理
│   ├── friend_link.go      # 友情链接管理
│   ├── feedback.go         # 用户反馈
│   ├── website.go          # 网站信息
│   └── config.go           # 系统配置读写
│
├── service/                # 业务逻辑层
│   ├── enter.go            # Service 分组聚合入口
│   ├── base.go             # 邮箱验证码发送
│   ├── user.go             # 用户业务逻辑
│   ├── article.go          # 文章业务逻辑（含 ES 搜索）
│   ├── article_helper.go   # 文章辅助方法（标签、分类处理）
│   ├── article_stat.go     # 文章统计
│   ├── comment.go          # 评论业务逻辑
│   ├── comment_helper.go   # 评论树形结构处理
│   ├── jwt.go              # JWT 黑名单管理
│   ├── image.go            # 图片上传（本地/七牛云）
│   ├── gaode.go            # 高德地图 IP 定位
│   ├── hot_search.go       # 热搜词条聚合
│   ├── es_index.go         # ES 索引管理
│   ├── calendar.go         # 日历信息
│   ├── config.go           # 配置读写逻辑
│   ├── website.go          # 网站信息
│   ├── advertisement.go    # 广告业务
│   ├── friend_link.go      # 友链业务
│   └── feedback.go         # 反馈业务
│
├── router/                 # 路由注册层
│   ├── enter.go            # RouterGroup 聚合
│   ├── base.go             # 公开基础路由（验证码、邮件）
│   ├── user.go             # 用户路由（公开 / 私有 / 管理员）
│   ├── article.go          # 文章路由
│   ├── comment.go          # 评论路由
│   ├── image.go            # 图片路由（仅管理员）
│   ├── advertisement.go    # 广告路由
│   ├── friend_link.go      # 友链路由
│   ├── feedback.go         # 反馈路由
│   ├── website.go          # 网站路由
│   └── config.go           # 配置路由（仅管理员）
│
├── model/                  # 数据模型层
│   ├── database/           # GORM 数据库实体
│   │   ├── user.go         # 用户表
│   │   ├── login.go        # 登录记录表
│   │   ├── comment.go      # 评论表
│   │   ├── image.go        # 图片表
│   │   ├── article_tag.go  # 文章标签表
│   │   ├── article_category.go  # 文章分类表
│   │   ├── article_like.go # 文章点赞表
│   │   ├── friend_link.go  # 友链表
│   │   ├── feedback.go     # 反馈表
│   │   ├── advertisement.go# 广告表
│   │   ├── footer_link.go  # 底部链接表
│   │   └── jwt_blacklist.go# JWT 黑名单表
│   ├── elasticsearch/      # ES 文档模型
│   ├── request/            # 请求结构体（参数绑定）
│   ├── response/           # 响应结构体（统一返回格式）
│   ├── appTypes/           # 自定义类型（角色 ID、注册来源等）
│   └── other/              # 其他辅助模型
│
├── middleware/             # Gin 中间件
│   ├── jwt.go              # JWT 认证（双 Token 无感刷新）
│   ├── admin.go            # 管理员权限校验
│   ├── logger.go           # 请求日志记录（Zap）
│   └── login_record.go     # 登录行为记录
│
├── initialize/             # 各组件初始化
│   ├── gorm.go             # MySQL 数据库连接
│   ├── redis.go            # Redis 连接
│   ├── es.go               # Elasticsearch 连接
│   ├── router.go           # 路由注册总入口
│   ├── cron.go             # 定时任务注册
│   └── other.go            # 其他初始化（如运行环境）
│
├── task/                   # 定时任务
│   ├── enter.go            # 任务注册（文章浏览量同步 @hourly / 日历更新 @daily）
│   ├── article_views.go    # 文章浏览量 Redis → MySQL 同步
│   └── calendar.go         # 日历数据抓取同步
│
├── flag/                   # CLI 命令工具
│   ├── enter.go            # CLI 入口 & 命令分发
│   ├── sql.go              # 数据库初始化
│   ├── sql_export.go       # SQL 数据导出
│   ├── sql_import.go       # SQL 数据导入
│   ├── es.go               # ES 索引初始化
│   ├── es_export.go        # ES 数据导出
│   ├── es_import.go        # ES 数据导入
│   └── admin.go            # 创建管理员账号
│
├── core/                   # 核心功能（配置加载 & 服务启动）
├── global/                 # 全局变量（DB、Log、Redis、ESClient、Configs）
├── config/                 # 配置结构体定义
├── utils/                  # 工具函数（JWT、Email、分页、加密、文件上传等）
│   ├── jwt.go / jwt_helper.go  # JWT 生成与解析
│   ├── email.go            # SMTP 邮件发送
│   ├── pagination.go       # 分页工具
│   ├── hash.go             # 密码加密（bcrypt）
│   ├── http.go             # HTTP 请求封装
│   ├── image.go            # 图片处理
│   ├── upload/             # 文件上传（本地 & 七牛云）
│   ├── hotSearch/          # 热搜词条工具
│   └── ...                 # 其他工具函数
│
├── assets/                 # 静态资源
├── uploads/                # 本地上传文件存储目录
└── log/                    # 日志文件存储目录
```

---

## 🚀 快速开始

### 1. 环境要求

| 软件 | 推荐版本 |
|------|---------|
| Go | 1.21+ |
| MySQL | 8.0+ |
| Redis | 6.0+ |
| Elasticsearch | 8.x |

### 2. 克隆项目

```bash
git clone https://github.com/Fujjj/Go-Blog-Server-Remote.git
cd goBlog/server
```

### 3. 配置文件

复制示例配置文件，并按需填写：

```bash
cp config.example.yaml config.yaml
```

各配置项说明参见下方 [配置说明](#-配置说明) 章节。

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 初始化数据库表结构

```bash
# 根据 GORM 模型自动迁移，创建所有数据表
go run main.go --sql
```

### 6. 初始化 Elasticsearch 索引

```bash
go run main.go --es
```

### 7. 创建管理员账号

```bash
# 管理员信息读取自 config.yaml 的 website 配置节
go run main.go --admin
```

### 8. 启动服务

```bash
go run main.go
```

服务默认监听 `0.0.0.0:8080`，可在 `config.yaml` 的 `system` 节中调整。

---

## ⚙️ 配置说明

| 配置节 | 说明 |
|--------|------|
| `system` | 服务监听地址、端口、运行模式（debug/release）、路由前缀 |
| `mysql` | MySQL 连接信息（host、port、db_name、username、password） |
| `redis` | Redis 连接信息（address、password、db） |
| `es` | Elasticsearch URL、用户名/密码 |
| `jwt` | Access Token / Refresh Token 密钥与过期时间 |
| `email` | SMTP 邮件服务（host、port、发件人账号、授权码） |
| `captcha` | 图形验证码尺寸、长度、干扰设置 |
| `qiniu` | 七牛云 OSS（bucket、accessKey、secretKey、CDN 域名） |
| `upload` | 文件上传存储类型（`local` / `qiniu`）及大小限制 |
| `qq` | QQ 互联登录（AppID、AppKey、回调地址） |
| `gaode` | 高德地图 Key（用于 IP 归属地查询） |
| `zap` | 日志级别、文件路径、滚动策略 |
| `website` | 博客名称、作者信息、ICP 备案号、社交链接等 |

> **注意**：`config.yaml` 已加入 `.gitignore`，请勿将真实密钥提交至版本库。

---

## 🔑 认证机制

项目采用**双 Token 无感刷新**方案：

- **Access Token**（短有效期，默认 2h）：通过请求头 `Authorization: Bearer <token>` 携带，用于每次 API 鉴权。
- **Refresh Token**（长有效期，默认 7d）：存储在 **HTTP-only Cookie** 中，Access Token 过期后由中间件自动使用 Refresh Token 签发新的 Access Token，并通过响应头 `new-access-token` 返回给前端。
- **Token 黑名单**：退出登录时，Refresh Token 加入 MySQL 黑名单表，防止已注销 Token 被复用。
- **管理员鉴权**：在 JWT 认证中间件之后附加 `AdminAuth` 中间件，通过 Claims 中的 `RoleID` 校验权限。

---

## 🔌 API 路由概览

所有路由以 `config.yaml` 中 `system.router_prefix`（默认 `api`）为前缀。

### 路由分组

| 分组 | 中间件 | 说明 |
|------|--------|------|
| `publicGroup` | 无 | 公开接口，无需登录 |
| `privateGroup` | `JWTAuth` | 需要登录 |
| `adminGroup` | `JWTAuth` + `AdminAuth` | 需要管理员权限 |

### 主要接口（部分）

#### 基础（公开）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/base/captcha` | 获取图形验证码 |
| POST | `/api/base/sendEmailCode` | 发送邮箱验证码 |
| GET | `/api/base/qqLoginURL` | 获取 QQ 登录链接 |

#### 用户
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/api/user/register` | 公开 | 邮箱注册 |
| POST | `/api/user/login` | 公开 | 账号登录（含登录记录） |
| POST | `/api/user/logout` | 登录 | 退出登录 |
| GET | `/api/user/info` | 登录 | 获取当前用户信息 |
| PUT | `/api/user/changeInfo` | 登录 | 修改用户信息 |
| PUT | `/api/user/resetPassword` | 登录 | 重置密码 |
| GET | `/api/user/card` | 公开 | 查看用户卡片 |
| GET | `/api/user/weather` | 登录 | 获取当前位置天气 |
| GET | `/api/user/chart` | 登录 | 用户注册/登录统计图表 |
| GET | `/api/user/list` | 管理员 | 用户列表 |
| PUT | `/api/user/freeze` | 管理员 | 冻结用户 |
| PUT | `/api/user/unfreeze` | 管理员 | 解冻用户 |
| GET | `/api/user/loginList` | 管理员 | 登录日志列表 |

#### 文章
| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| POST | `/api/article` | 管理员 | 创建文章 |
| PUT | `/api/article` | 管理员 | 更新文章 |
| DELETE | `/api/article/:id` | 管理员 | 删除文章 |
| GET | `/api/article/list` | 公开 | 文章列表 |
| GET | `/api/article/:id` | 公开 | 文章详情 |
| POST | `/api/article/like` | 登录 | 点赞/取消点赞 |

#### 评论 / 反馈 / 友链 / 广告

评论、反馈、友情链接、广告均遵循相同的三分组模式（公开读取 / 登录操作 / 管理员管理）。

---

## 🛠️ CLI 工具

程序内置 CLI 工具，**携带命令行参数时执行对应操作后直接退出**，不启动 HTTP 服务。

```bash
# 查看帮助
./server -h

# 初始化 MySQL 数据表（AutoMigrate）
./server --sql

# 导出 SQL 数据
./server --sql-export

# 从文件导入 SQL 数据
./server --sql-import=./backup.json

# 初始化 Elasticsearch 索引
./server --es

# 导出 ES 数据
./server --es-export

# 从文件导入 ES 数据
./server --es-import=./es_backup.json

# 创建管理员（读取 config.yaml 中 website 节信息）
./server --admin
```

---

## ⏰ 定时任务

| 任务 | 频率 | 说明 |
|------|------|------|
| 文章浏览量同步 | 每小时 | 将 Redis 中缓存的浏览量写回 MySQL |
| 日历信息更新 | 每天 | 抓取并更新日历/节假日数据 |

---

## 📦 构建

```bash
# Windows
go build -o server.exe main.go

# Linux / macOS
go build -o server main.go
```

---

## 📝 日志

- 日志采用 **Zap** 结构化格式，同时输出到控制台和日志文件。
- 日志文件路径由 `config.yaml` 的 `zap.filename` 指定（默认 `log/go_blog.log`）。
- 支持按文件大小滚动（`max_size`）、保留份数（`max_backups`）、保留天数（`max_age`）。

---

## 📄 License

MIT
