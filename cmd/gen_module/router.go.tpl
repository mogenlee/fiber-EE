package {{.PackageName}}

import (
	"fiber-ee/internal/pkg/request"
	"fiber-ee/internal/service/admin/{{.PackageName}}"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type {{.StructName}}Router struct {
	log      *zap.Logger
	svc      {{.PackageName}}.{{.StructName}}Service
	validate request.Validate
}

func New{{.StructName}}Router(log *zap.Logger, svc {{.PackageName}}.{{.StructName}}Service, validate request.Validate) *{{.StructName}}Router {
	return &{{.StructName}}Router{
		log:      log,
		svc:      svc,
		validate: validate,
	}
}

func (h {{.StructName}}Router) Register(root fiber.Router) {
	group := root.Group("/{{.RoutePath}}")
	_ = group
	// TODO: 添加路由
	// group.Get("/list", h.list).Name("{{.PackageName}}_list")
}
