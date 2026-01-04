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
- 📖 QingFeng Swagger 自动生成 API 文档
- 💚 健康检查接口
- 📊 实时监控面板
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
│   │   └── admin/           # Admin 路由
│   └── service/
│       └── admin/           # Admin 服务
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
│   ├── docs/                # Swagger 文档（自动生成）
│   ├── locales/             # 国际化翻译文件
│   ├── middleware/          # 中间件
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
├── logs/                    # 日志文件（按日期切割）
└── static/                  # 静态资源
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. 配置

修改 `config/config.yaml`：

```yaml
app:
  name: "fiber-ee"
  port: ":8080"
  debug: true

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

| 地址 | 说明 |
|------|------|
| http://localhost:8080 | API 服务 |
| http://localhost:8080/doc/ | Swagger 文档 |
| http://localhost:8080/metrics | 监控面板 |
| http://localhost:8080/health | 健康检查 |

## Swagger 文档

启动时自动生成 API 文档，访问 `/doc/` 查看。

在 Router 函数上添加注解：

```go
// @Summary 用户登录
// @Description 用户登录获取 Token
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param request body req.LoginReq true "登录参数"
// @Success 200 {object} response.Response
// @Router /user/login [post]
func (h *UserRouter) Login(c fiber.Ctx) error {
    // ...
}
```

## 代码生成

### 生成 ORM 代码

```bash
go run cmd/gen_orm/main.go
```

### 生成模块代码

```bash
go run cmd/gen_module/main.go -name=article
```

## API 响应格式

```json
// 成功
{"code": 200, "msg": "操作成功", "data": {}}

// 分页
{"code": 200, "msg": "操作成功", "data": {"list": [], "total": 100, "page_no": 1, "page_size": 10}}

// 错误
{"code": 4004, "msg": "账号或密码错误", "data": null}
```

## JWT 认证

```bash
# 请求头
Authorization: Bearer <access_token>

# Token 刷新
POST /admin/v1/auth/refresh
{"refresh_token": "..."}
```

## 自定义类型

```go
// Timestamp: 数据库 int → JSON 日期字符串
CreatedAt types.Timestamp `json:"createdAt"` // "2025-01-04 12:00:00"

// Bool: 数据库 int(0/1) → JSON bool
IsActive types.Bool `json:"isActive"` // true/false
```

## 限流配置

```yaml
limiter:
  enabled: true       # 开关
  max: 100            # 每分钟最大请求数
  expiration: 60      # 时间窗口（秒）
  skip_local: true    # 跳过本地请求
```

## Docker 部署

```bash
docker-compose up -d
```

## License

MIT
