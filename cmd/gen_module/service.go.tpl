package {{.PackageName}}

import (
	"fiber-ee/internal/model/query"
)

// {{.StructName}}Service 服务接口
type {{.StructName}}Service interface {
	// TODO: 添加方法
}

type {{.PackageName}}Service struct {
	db *query.Query
}

// New{{.StructName}}Service 创建服务（db 由 dig 自动注入）
func New{{.StructName}}Service(db *query.Query) {{.StructName}}Service {
	return &{{.PackageName}}Service{db: db}
}
