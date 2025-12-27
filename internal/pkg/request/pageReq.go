package request

import "github.com/gofiber/fiber/v3"

// PageReq 分页请求参数
type PageReq struct {
	PageNo   int `query:"page_no" validate:"omitempty,min=1"`
	PageSize int `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// GetPage 获取页码，默认 1
func (p *PageReq) GetPage() int {
	if p.PageNo <= 0 {
		return 1
	}
	return p.PageNo
}

// GetPageSize 获取每页数量，默认 20
func (p *PageReq) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 获取偏移量
func (p *PageReq) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// BindPage 绑定分页参数
func BindPage(c fiber.Ctx) *PageReq {
	p := &PageReq{}
	_ = c.Bind().Query(p)
	return p
}
