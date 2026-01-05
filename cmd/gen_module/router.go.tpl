package {{.PackageName}}

import (
	"fiber-ee/app/service/admin/{{.PackageName}}"
	"fiber-ee/internal/pkg/request"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type I{{.StructName}}Router struct {
	log      *zap.Logger
	svc      {{.PackageName}}.{{.StructName}}Service
	validate request.Validate
}

func NewRouter(log *zap.Logger, svc {{.PackageName}}.{{.StructName}}Service, validate request.Validate) *I{{.StructName}}Router {
	return &I{{.StructName}}Router{
		log:      log,
		svc:      svc,
		validate: validate,
	}
}

func (h I{{.StructName}}Router) Register(root fiber.Router) {
	group := root.Group("/{{.RoutePath}}")
	_ = group
	// TODO: 添加路由
	// group.Get("/list", h.list).Name("{{.PackageName}}_list")
}
