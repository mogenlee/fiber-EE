package test

import (
	"fiber-ee/app/dto/req"
	"fiber-ee/app/service/admin/test"
	"fiber-ee/internal/pkg/request"
	"fiber-ee/internal/pkg/response"

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
	group.Post("/login", h.login).Name("test_lgin")
	group.Get("/list", h.list).Name("test_list")
	group.Get("/detail/:id", h.detail).Name("test_detail")
	group.Get("/detail", h.details).Name("test_detail")
	group.Post("/edit", h.edit).Name("test_edit")
	group.Post("/add", h.add).Name("test_add")
}

func (h TestRouter) login(ctx fiber.Ctx) error {
	var loginReq req.TestLoginReq
	if err := h.validate.BindAndValidate(ctx, &loginReq); err != nil {
		return err
	}
	login, err := h.svc.Login(ctx, loginReq)
	return response.CheckWithData(ctx, login, err)
}

// list 获取列表（分页）
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
	// 优先从路径参数获取 id
	id := cast.ToInt32(ctx.Params("id"))

	data, err := h.svc.Detail(ctx, id)
	return response.CheckWithData(ctx, data, err)
}

func (h TestRouter) details(ctx fiber.Ctx) error {
	var idReq request.CommonIdReq
	if err := h.validate.BindQuery(ctx, &idReq); err != nil {
		return err
	}
	data, err := h.svc.Detail(ctx, idReq.Id)
	return response.CheckWithData(ctx, data, err)
}

func (h TestRouter) edit(ctx fiber.Ctx) error {
	var editReq req.TestEditReq
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
