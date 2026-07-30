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

## 项目部署
### 开发工具及版本

golang: 1.25.5

node: v24.14.0

docker: 29.2.1

编译器：vscode、goland、webstorm

### 编译文件
编译后端，得到 main 文件

```bash
# windows环境下，打开项目所在目录，进入 server 文件夹，打开 cmd （不是 powershell）
set GOOS=linux
set GOARCH=amd64
go mod tidy
go build main.go
```

编译前端，得到 dist 文件夹
//确保public/images目录下有必要的图片 如404页面图、QQ联系图、github图标、QQ图标...

```bash
# windows环境下，打开项目所在目录，进入 web 文件夹，打开 cmd
npm install
npm run build
```

#### 环境变量说明
编译前检查 `.env.production` 文件：
```bash
# .env.production（编译时生效）
VITE_API_BASE_URL=/api
VITE_UPLOADS_URL=/uploads
```

| 变量 | 值 | 说明 |
|------|------|------|
| `VITE_API_BASE_URL` | `/api` | API 请求前缀，生产和开发都用 `/api`，开发环境 Vite proxy 转发到 `localhost:8080`，生产环境由 nginx 反向代理到后端 |
| `VITE_UPLOADS_URL` | `/uploads` | 上传文件基础路径 |

> **不需要**在生产环境将 `/api` 替换为完整域名（如 `https://www.your_domain/api`）。因为 nginx 统一代理 `/api/` 到 `127.0.0.1:8080`，使用相对路径 `/api` 即可，也避免了跨域问题。

### 环境准备

```bash
# 安装 docker
yum install -y docker-ce
systemctl start docker
systemctl enable docker

# 安装 supervisor
yum install -y supervisor

# 安装 nginx
yum install https://nginx.org/packages/centos/8/x86_64/RPMS/nginx-1.26.0-1.el8.ngx.x86_64.rpm
```

### 挂载服务
#### redis
```bash
docker run -d \
  --name blog-redis \
  --restart always \
  -p 6379:6379 \
  -v /opt/go_blog/server/data/redis/data:/data \
  --memory 256m \
  redis-server --maxmemory 128mb --maxmemory-policy allkeys-lru

#   -v /opt/go_blog/server/data/redis/data:/data \是在 Docker 中将宿主机的 Redis 数据目录挂载到容器内，实现 Redis 数据的持久化存储。

# --memory 256m 限制容器可以使用的最大内存为 256MB。

# --maxmemory 128mb Redis 自身最多使用 128MB 内存

# --maxmemory-policy allkeys-lru 当内存达到 128MB 时，自动淘汰最近最少使用的 key
```

#### MySQL
```bash
docker run -itd --name mysql --restart=always -p 3306:3306 -v /opt/go_blog/server/data/mysql/conf:/etc/mysql/conf.d -v /opt/go_blog/server/data/mysql/datadir:/var/lib/mysql -v /opt/go_blog/server/data/mysql/go_blog.sql:/opt/go_blog.sql -e  MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=blog_db mysql

#  -v /opt/go_blog/server/data/mysql/conf:/etc/mysql/conf.d  是把宿主机的目录挂载到 MySQL 容器的配置目录 /etc/mysql/conf.d/ 

# -v /opt/go_blog/server/data/mysql/datadir:/var/lib/mysql是在 Docker 中挂载数据卷（Volume Mount），实现 MySQL 数据的持久化存储（MySQL 容器运行时产生的所有数据库文件（表数据、日志等）实际上写入的是宿主机上的 /opt/go_blog/server/data/mysql/datadir，而不是仅在容器内部，即使容器被删除（docker rm），数据库数据依然保留在宿主机上。）。
# /opt/go_blog/server/data/mysql 是主机（Host）上的目录路径
# /var/lib/mysql是容器内的目录路径

# -v /opt/go_blog/server/data/mysql/go_blog.sql:/opt/go_blog.sql是在 Docker 中将宿主机上的一个 SQL 文件挂载到容器内部，让容器能够读取/执行这个文件
```

#### elasticsearch 

elasticsearch 无法直接数据卷挂载本地，需要先启动一个不挂载数据卷的容器，将文件复制到本地，再进行挂载

Elasticsearch 对权限和系统参数要求很严格，条件不满足时就会启动失败
Elasticsearch 8.x 容器内以 UID 1000 运行，宿主机目录必须属于 1000:1000
ES 要求 vm.max_map_count >= 262144

```bash
# 最新
mkdir -p /opt/go_blog/server/data/es/{data,config,plugins}
chown -R 1000:1000 /opt/go_blog/server/data/es

sysctl -w vm.max_map_count=262144
echo "vm.max_map_count=262144" >> /etc/sysctl.conf

docker run -d \
  --name es \
  --restart always \
  -p 127.0.0.1:9200:9200 \
  -v /opt/go_blog/server/data/es/data:/usr/share/elasticsearch/data \
  -v /opt/go_blog/server/data/es/config:/usr/share/elasticsearch/config \
  -v /opt/go_blog/server/data/es/plugins:/usr/share/elasticsearch/plugins \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -e "xpack.security.enabled=false" \
  -e "xpack.security.http.ssl.enabled=false" \
  -e "xpack.license.self_generated.type=trial" \
  -e "bootstrap.memory_lock=true" \
  --memory 1g \
  --ulimit memlock=-1:-1 \
  elasticsearch:8.17.0

# -v /opt/go_blog/server/data/es/data:/usr/share/elasticsearch/data  ES 的索引、分片、文档等实际写入宿主机磁盘，容器删除后数据不丢失。

# -v /opt/go_blog/server/data/es/config:/usr/share/elasticsearch/config 以后想改 ES 配置（如安装 IK 分词器后改 analysis 配置），直接编辑宿主机上的文件即可，不用进容器。

# -v /opt/go_blog/server/data/es/plugins:/usr/share/elasticsearch/plugins 安装插件（如中文 IK 分词器）时，把插件文件放到宿主机的这个目录，重启容器即可生效。

# -e "discovery.type=single-node" 设置环境变量，以单节点模式运行，不尝试加入集群

# -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" JVM 初始堆和最大堆内存为512MB

#-e "xpack.security.http.ssl.enabled=false" 关闭 X-Pack 安全认证

# -e "xpack.security.http.ssl.enabled=false" 关闭 HTTPS，使用 HTTP 访问

# xpack.license.self_generated.type=trial 使用 试用版许可证

# bootstrap.memory_lock=true 启动时锁定 JVM 堆内存，防止被操作系统交换（swap）到磁盘

# --memory 1g 限制容器最多使用 1GB 内存（操作系统层面的 cgroup 限制）

# --ulimit memlock=-1:-1  无限制地锁定内存配合--memory 1g使用
```

### 服务端目录与权限配置

将文件按照下述目录上传

```bash
# /opt/go_blog
├── go_blog
    ├── server
    │   ├── data
    │   │   ├── es
        │   └── mysql
    │   ├── main
    │   └── config.yaml
    └── web
        └── dist
```

将main的权限从644修改为755，并初始化项目
``` bash
cd /opt/go_blog/server/
chmod +x ./main
```

### 配置supervisord
``` bash
vim /etc/supervisord.d/go_blog.ini
```
```ini
[program: go_blog]
command=/opt/go_blog/server/main
directory=/opt/go_blog/server/
autorestart=true ; 程序意外退出是否自动重启
autostart=true ; 是否自动启动
user=root ; 进程执行的用户身份
stopsignal=INT
startsecs=1 ; 自动重启间隔
stopasgroup=true ;默认为false,进程被杀死时，是否向这个进程组发送stop信号，包括子进程
killasgroup=true ;默认为false，向进程组发送kill信号，包括子进程
```
ESC :wq 保存并退出，启动并一直开启supervisord
```bash
systemctl start supervisord
systemctl enable supervisord
```

### 配置nginx
修改 /etc/nginx/nginx.conf
```nginx

user  root;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;
    
    client_max_body_size 20M; #上传文件大小限制

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;
}
```
然后删除etc/nginx/conf.d 下的default.conf
创建 /etc/nginx/conf.d/nginx.conf

**将 your_domain 替换为你的域名**，请自行获取 ssl 证书，上传证书文件至 /etc/nginx/cert/

```nginx
server {
    listen 80;
    server_name your_domain www.your_domain;
    return 301 https://www.your_domain$request_uri;
}

server { 
    listen 443 ssl; 
    server_name your_domain;  # 仅匹配非 www 的域名
    ssl_certificate /etc/nginx/cert/your_domain.crt; # 证书公钥
    ssl_certificate_key /etc/nginx/cert/your_domain.key; # 证书私钥
    return 301 https://www.your_domain$request_uri;  # 强制跳转到 www.your_domain
}

server {
    gzip on;
    gzip_vary on;
    gzip_disable "MSIE [1-6]\.";
    gzip_static on;
    gzip_min_length 256;
    gzip_buffers 32 8k;
    gzip_http_version 1.1;
    gzip_comp_level 5;
    gzip_proxied any;
    gzip_types text/plain text/css text/xml application/javascript application/x-javascript application/xml application/xml+rss application/emacscript application/json image/svg+xml;

    listen 443 ssl;
    server_name www.your_domain; # 多个域名⽤空格分开 
    ssl_certificate /etc/nginx/cert/your_domain.crt; # 证书公钥
    ssl_certificate_key /etc/nginx/cert/your_domain.key; # 证书私钥
    ssl_session_timeout 5m; 
    ssl_session_cache shared:MozSSL:10m;  # 设置会话缓存以提⾼性能 
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;  # 配置加密算法 
    ssl_protocols TLSv1.2 TLSv1.3;  # 配置加密协议 
    ssl_prefer_server_ciphers on;

    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always; #可选配置，开启HSTS 
    add_header X-Frame-Options DENY; # 可选配置，防⽌点击劫持 
    add_header X-Content-Type-Options nosniff; # 可选配置，防⽌MIME类型嗅探 
    add_header X-XSS-Protection "1; mode=block"; # 可选配置，防⽌XSS攻击

    location / {
        try_files $uri $uri/ /index.html;
        root   /opt/go_blog/web/dist;
        index  index.html index.htm;
    }

    location /api/ {
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header REMOTE-HOST $remote_addr;
        proxy_pass http://127.0.0.1:8080/api/;
    }

    location /image {
        alias /opt/go_blog/web/dist/image;
    }

    location /emoji {
        alias /opt/go_blog/web/dist/emoji;
    }

    location /uploads/ {
        alias /opt/go_blog/server/uploads/;
    }
}
```
#### 启动 nginx
```bash
systemctl start nginx
systemctl enable nginx
```

### 启动后端服务

```bash
cd /opt/go_blog/server/
./main -sql
./main -es
# 若之前有数据库备份文件可执行 如./main -sql-import ./mysql_20250214.sql 恢复数据库
# 若之前有eS备份文件可执行 如./main -es-import ./es_20250214.json 恢复es
./main
```

## 📄 License

MIT
