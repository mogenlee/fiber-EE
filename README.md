# Fiber-EE

基于 [Fiber v3](https://github.com/gofiber/fiber) 的企业级 Go Web 框架脚手架。

## 特性

- 🚀 Fiber v3 高性能 Web 框架
- 📦 Dig 依赖注入
- 🔐 JWT 认证 + Token 刷新 + Casbin 权限控制
- 🔒 Bcrypt 密码加密
- 🗄️ GORM + Gen 代码生成（支持 SQLite/MySQL/PostgreSQL）
- 🌍 i18n 国际化
- 📝 Zap 结构化日志（按日期切割）
- ⚙️ Viper 配置管理
- 🛡️ 限流中间件（支持 Redis/内存存储）
- 🔄 优雅关闭
- 🐳 Docker 部署支持

## 目录结构

```
.
├── app/                     # 业务代码
│   ├── dto/
│   │   ├── req/             # 请求结构体
│   │   └── resp/            # 响应结构体
│   ├── router/
│   │   ├── admin/           # Admin 路由
│   │   └── build.go         # 路由注册
│   └── service/
│       ├── admin/           # Admin 服务
│       └── build.go         # 服务注册
│
├── cmd/
│   ├── server/              # 主程序入口
│   ├── gen_orm/             # ORM 代码生成器
│   └── gen_module/          # 模块代码生成器
│
├── config/
│   ├── config.go            # 配置结构体
│   ├── config.yaml          # 配置文件
│   └── rbac_model.conf      # Casbin RBAC 模型
│
├── internal/                # 内部包（不对外暴露）
│   ├── bootstrap/           # 初始化（DI容器、数据库、日志等）
│   ├── middleware/          # 中间件（JWT、Casbin、限流、日志）
│   ├── model/
│   │   ├── entity/          # 数据库实体
│   │   └── query/           # GORM Gen 查询代码
│   └── pkg/
│       ├── i18n/            # 国际化
│       ├── request/         # 请求绑定与校验
│       ├── response/        # 统一响应
│       ├── types/           # 自定义类型（Timestamp、Bool）
│       ├── utils/           # 工具函数
│       └── validator/       # 验证器
│
├── data/                    # 数据库文件
├── locales/                 # 国际化翻译文件
├── logs/                    # 日志文件（按日期切割）
└── static/                  # 静态资源
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置

修改 `config/config.yaml`：

```yaml
app:
  name: "fiber-ee"
  port: ":8080"

database:
  driver: "sqlite"           # sqlite, mysql, postgres
  source: "data/test.db"

redis:
  enabled: true
  host: "localhost"
  port: "6379"

jwt:
  secret: "your-secret-key"
  expire: 7200

limiter:
  enabled: true
  max: 100
  expiration: 60
```

### 3. 运行

```bash
# 热重载开发（需安装 air）
air

# 或直接运行
go run cmd/server/main.go
```

### 4. 访问

- API: http://localhost:8080
- 监控: http://localhost:8080/metrics

## 代码生成

### 生成 ORM 代码

```bash
go run cmd/gen_orm/main.go
```

生成文件：
- `internal/model/entity/` - 实体结构体
- `internal/model/query/` - 查询方法

### 生成模块代码

```bash
go run cmd/gen_module/main.go -name=article
```

生成文件：
- `app/router/admin/article/handler.go`
- `app/service/admin/article/service.go`

## API 响应格式

### 成功响应

```json
{
    "code": 200,
    "msg": "操作成功",
    "data": {}
}
```

### 分页响应

```json
{
    "code": 200,
    "msg": "操作成功",
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "pageSize": 10
    }
}
```

### 错误响应

```json
{
    "code": 4004,
    "msg": "账号或密码错误",
    "data": null
}
```

## 错误码规范

| 范围 | 说明 |
|------|------|
| 200 | 成功 |
| 1xxx | 通用业务错误 |
| 2xxx | 参数相关错误 |
| 3xxx | 数据操作错误 |
| 4xxx | 认证授权错误 |
| 5xxx | 系统级错误 |

## JWT 认证

### 登录响应

```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 7200
}
```

### 使用 Token

```
Authorization: Bearer <access_token>
```

### Token 刷新

**POST /admin/v1/auth/refresh**

```json
{ "refresh_token": "eyJhbGciOiJIUzI1NiIs..." }
```

## 自定义类型

### Timestamp

数据库存 int，JSON 输出日期格式：

```go
type User struct {
    CreatedAt types.Timestamp `json:"createdAt"` // 输出: "2025-01-04 12:00:00"
}
```

### Bool

数据库存 int (0/1)，JSON 输出 true/false：

```go
type User struct {
    IsActive types.Bool `json:"isActive"` // 输出: true
}
```

## 限流配置

```yaml
limiter:
  enabled: true       # 开关
  max: 100            # 每分钟最大请求数
  expiration: 60      # 时间窗口（秒）
  skip_local: true    # 跳过本地请求
```

## 白名单配置

不需要认证的路由在 `internal/middleware/whitelist.go`：

```go
var whitelist = []string{
    "/admin/v1/auth/login",
}
```

## Docker 部署

```bash
docker-compose up -d
```

## License

MIT
