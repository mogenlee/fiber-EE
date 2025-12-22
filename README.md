# Fiber-EE

基于 [Fiber v3](https://github.com/gofiber/fiber) 的企业级 Go Web 框架脚手架。

## 特性

- 🚀 Fiber v3 高性能 Web 框架
- 📦 Dig 依赖注入
- 🔐 JWT 认证 + Token 刷新 + Casbin 权限控制
- 🔒 Bcrypt 密码加密
- 🗄️ GORM + Gen 代码生成
- 🌍 i18n 国际化
- 📝 Zap 结构化日志
- ⚙️ Viper 配置管理（支持热更新 + 环境变量）
- 🔄 优雅关闭

## 目录结构

```
.
├── cmd/
│   ├── server/          # 主程序入口
│   ├── gen_orm/         # ORM 代码生成器
│   └── gen_module/      # 模块代码生成器
├── config/
│   ├── config.go        # 配置结构体
│   ├── config.yaml      # 配置文件
│   └── rbac_model.conf  # Casbin RBAC 模型
├── internal/
│   ├── bootstrap/       # 初始化（DI容器、数据库、日志等）
│   ├── middleware/      # 中间件（JWT、Casbin、日志、错误处理）
│   ├── model/
│   │   ├── entity/      # 数据库实体
│   │   └── query/       # GORM Gen 生成的查询代码
│   ├── pkg/
│   │   ├── i18n/        # 国际化
│   │   ├── request/     # 请求绑定与校验
│   │   ├── response/    # 统一响应
│   │   ├── utils/       # 工具函数
│   │   └── validator/   # 验证器
│   ├── router/          # 路由层
│   ├── schema/
│   │   ├── req/         # 请求结构体
│   │   └── resp/        # 响应结构体
│   ├── server/          # 服务器封装
│   └── service/         # 业务逻辑层
├── locales/             # 国际化翻译文件
└── logs/                # 日志文件
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置

复制并修改配置文件：

```bash
cp config/config.yaml.example config/config.yaml
```

### 3. 运行

```bash
# 开发模式（热重载）
make dev

# 或直接运行
go run cmd/server/main.go
```

### 4. 访问

- API: http://localhost:8080
- 监控: http://localhost:8080/metrics

## 环境变量

支持通过环境变量覆盖配置（前缀 `FIBER_`）：

```bash
# 示例
export FIBER_DATABASE_PASSWORD=your_password
export FIBER_JWT_SECRET=your_secret
export FIBER_APP_PORT=:3000
```

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
- `internal/router/admin/article/router.go`
- `internal/service/admin/article/service.go`

可选参数：
```bash
-router=internal/router/app    # 自定义 router 输出目录
-service=internal/service/app  # 自定义 service 输出目录
-tpl=cmd/gen_module            # 自定义模板目录
```

## 示例代码

### Router 示例

```go
// internal/router/admin/test/test.go
package test

import (
    "fiberEE/internal/pkg/request"
    "fiberEE/internal/pkg/response"
    "fiberEE/internal/schema/req"
    "fiberEE/internal/service/admin/test"

    "github.com/gofiber/fiber/v3"
    "github.com/spf13/cast"
    "go.uber.org/zap"
)

type TestRouter struct {
    Log      *zap.Logger
    svc      test.TestService
    validate request.Validate
}

func NewRouter(log *zap.Logger, svc test.TestService, validate request.Validate) *TestRouter {
    return &TestRouter{
        Log:      log,
        svc:      svc,
        validate: validate,
    }
}

func (h TestRouter) Register(root fiber.Router) {
    group := root.Group("/test")
    group.Post("/login", h.login).Name("test_login")
    group.Get("/list", h.list).Name("test_list")
    group.Get("/detail/:id", h.detail).Name("test_detail")
    group.Post("/edit", h.edit).Name("test_edit")
    group.Post("/add", h.add).Name("test_add")
}

func (h TestRouter) login(ctx fiber.Ctx) error {
    var loginReq req.TestLoginResp
    if err := h.validate.BindAndValidate(ctx, &loginReq); err != nil {
        return err
    }
    data, err := h.svc.Login(ctx, loginReq)
    return response.CheckWithData(ctx, data, err)
}

func (h TestRouter) list(ctx fiber.Ctx) error {
    pageReq := request.BindPage(ctx)
    var listReq req.TestListReq
    if err := h.validate.BindQuery(ctx, &listReq); err != nil {
        return err
    }
    data, err := h.svc.List(ctx, pageReq, listReq)
    return response.CheckWithData(ctx, data, err)
}

func (h TestRouter) detail(ctx fiber.Ctx) error {
    id := cast.ToInt32(ctx.Params("id"))
    data, err := h.svc.Detail(ctx, id)
    return response.CheckWithData(ctx, data, err)
}

func (h TestRouter) edit(ctx fiber.Ctx) error {
    var editReq req.TestEditResp
    if err := h.validate.BindAndValidate(ctx, &editReq); err != nil {
        return err
    }
    return response.Check(ctx, h.svc.Edit(ctx, editReq))
}

func (h TestRouter) add(ctx fiber.Ctx) error {
    var addReq req.TestAddReq
    if err := h.validate.BindAndValidate(ctx, &addReq); err != nil {
        return err
    }
    return response.Check(ctx, h.svc.Add(ctx, addReq))
}
```

### Service 示例

```go
// internal/service/admin/test/test.go
package test

import (
    "fiberEE/internal/middleware"
    "fiberEE/internal/model/entity"
    "fiberEE/internal/model/query"
    "fiberEE/internal/pkg/request"
    "fiberEE/internal/pkg/response"
    "fiberEE/internal/pkg/utils"
    "fiberEE/internal/schema/req"

    "github.com/gofiber/fiber/v3"
    "github.com/spf13/cast"
)

// TestService 服务接口
type TestService interface {
    List(ctx fiber.Ctx, pageReq *request.PageReq, listReq req.TestListReq) (any, error)
    Detail(ctx fiber.Ctx, id int32) (any, error)
    Edit(ctx fiber.Ctx, editReq req.TestEditResp) error
    Add(ctx fiber.Ctx, addReq req.TestAddReq) error
    Login(ctx fiber.Ctx, loginReq req.TestLoginResp) (any, error)
}

type testService struct {
    db     *query.Query
    jwtCfg *middleware.JWTConfig
}

// NewTestService 创建服务（db 由 dig 自动注入）
func NewTestService(db *query.Query, jwtCfg *middleware.JWTConfig) TestService {
    return &testService{db: db, jwtCfg: jwtCfg}
}

func (t testService) Login(ctx fiber.Ctx, loginReq req.TestLoginResp) (any, error) {
    m := t.db.User
    q := m.WithContext(ctx)

    user, err := q.Where(m.Username.Eq(loginReq.Username), m.Password.Eq(loginReq.Password)).First()
    if err != nil {
        return nil, response.ErrLogin
    }

    token, err := middleware.GenerateToken(t.jwtCfg, uint(user.ID), user.Username, "admin")
    return token, err
}

func (t testService) List(ctx fiber.Ctx, pageReq *request.PageReq, listReq req.TestListReq) (any, error) {
    m := t.db.User
    q := m.WithContext(ctx)

    if listReq.Username != "" {
        q = q.Where(m.Username.Eq(listReq.Username))
    }
    if listReq.Id > 0 {
        q = q.Where(m.ID.Eq(listReq.Id))
    }

    result, total, err := q.FindByPage(pageReq.GetOffset(), pageReq.GetPageSize())
    if err != nil {
        return nil, response.ErrQuery
    }

    return response.PageData{
        List:     result,
        Total:    total,
        Page:     pageReq.GetPage(),
        PageSize: pageReq.GetPageSize(),
    }, nil
}

func (t testService) Detail(ctx fiber.Ctx, id int32) (any, error) {
    m := t.db.User
    q := m.WithContext(ctx)
    if id > 0 {
        q = q.Where(m.ID.Eq(id))
    }
    return q.First()
}

func (t testService) Edit(ctx fiber.Ctx, editReq req.TestEditResp) error {
    m := t.db.User
    q := m.WithContext(ctx)

    var obj entity.User
    utils.Copy(&obj, editReq)
    _, err := q.Where(m.ID.Eq(editReq.Id)).Updates(obj)
    return err
}

func (t testService) Add(ctx fiber.Ctx, addReq req.TestAddReq) error {
    m := t.db.User
    q := m.WithContext(ctx)

    var obj entity.User
    utils.Copy(&obj, addReq)
    return q.Create(&obj)
}
```

### 请求结构体示例

```go
// internal/schema/req/test.go
package req

type TestListReq struct {
    Username string `query:"username"`
    Id       int32  `query:"id"`
}

type TestEditResp struct {
    Id       int32  `json:"id" validate:"required"`
    Username string `json:"username" validate:"required"`
}

type TestAddReq struct {
    Username string `json:"username" validate:"required"`
}

type TestLoginResp struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Password string `json:"password" validate:"required,min=6,max=18"`
}
```

## 注册新模块

1. 在 `internal/bootstrap/buildService.go` 添加：
```go
var adminServices = []any{
    test.NewTestService,
    article.NewService,  // 新增
}
```

2. 在 `internal/bootstrap/buildRouter.go` 添加：
```go
var adminRouters = []any{
    test.NewRouter,
    article.NewRouter,  // 新增
}
```

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
    "code": 2001,
    "msg": "参数校验错误",
    "data": ["username长度必须至少为3个字符"]
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

## 国际化

请求头添加 `Accept-Language` 切换语言：

```bash
# 中文（默认）
curl -H "Accept-Language: zh"

# 英文
curl -H "Accept-Language: en"
```

翻译文件位于 `locales/` 目录。

## JWT 认证

### 配置

```yaml
jwt:
  secret: "your-secret-key"
  expire: 7200          # Access Token 有效期（秒）
  refresh_expire: 604800 # Refresh Token 有效期（7天）
```

### 登录响应

```json
{
    "code": 200,
    "msg": "操作成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 7200
    }
}
```

### Token 刷新

**POST /admin/v1/auth/refresh**

请求：
```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

响应：
```json
{
    "code": 200,
    "msg": "操作成功",
    "data": {
        "access_token": "新的access_token",
        "refresh_token": "新的refresh_token",
        "expires_in": 7200
    }
}
```

### 使用 Token

请求头添加：
```
Authorization: Bearer <access_token>
```

## 密码加密

使用 bcrypt 加密，工具函数位于 `internal/pkg/utils/crypto.go`：

```go
// 加密密码
hash := utils.HashPassword("123456")

// 验证密码
ok := utils.CheckPassword("123456", hash)
```

## 工具函数

### 树形结构 (`internal/pkg/utils/tree.go`)

```go
// 实体实现 TreeNode 接口
type Menu struct {
    ID       int64 `json:"id"`
    ParentID int64 `json:"parent_id"`
    Name     string `json:"name"`
}

func (m Menu) GetID() int64       { return m.ID }
func (m Menu) GetParentID() int64 { return m.ParentID }

// 构建树
tree := utils.BuildTree(menus, 0)

// 树转列表
flat := utils.FlatTree(tree)

// 查找节点
node := utils.FindTreeNode(tree, 2)

// 获取所有 ID
ids := utils.GetTreeIDs(tree)
```

### 结构体拷贝 (`internal/pkg/utils/copy.go`)

```go
var user entity.User
utils.Copy(&user, req)
```

## License

MIT
