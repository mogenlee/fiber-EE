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

// login 用户登录
// @Summary 用户登录
// @Description 用户登录获取 Token
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param request body req.TestLoginReq true "登录参数"
// @Success 200 {object} response.Response
// @Router /test/login [post]
func (h TestRouter) login(ctx fiber.Ctx) error {
	var loginReq req.TestLoginReq
	if err := h.validate.BindAndValidate(ctx, &loginReq); err != nil {
		return err
	}
	login, err := h.svc.Login(ctx, loginReq)
	return response.CheckWithData(ctx, login, err)
}

// list 获取列表（分页）
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param username query string false "用户名"
// @Param id query int false "用户ID"
// @Param page_no query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /test/list [get]
func (h TestRouter) list(ctx fiber.Ctx) error {
	pageReq := request.BindPage(ctx)
	var listReq req.TestListReq
	if err := h.validate.BindQuery(ctx, &listReq); err != nil {
		return err
	}
	data, err := h.svc.List(ctx, pageReq, listReq)
	return response.CheckWithData(ctx, data, err)
}

// detail 获取用户详情（路径参数）
// @Summary 获取用户详情
// @Description 根据ID获取用户详情
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /test/detail/{id} [get]
func (h TestRouter) detail(ctx fiber.Ctx) error {
	// 优先从路径参数获取 id
	id := cast.ToInt32(ctx.Params("id"))

	data, err := h.svc.Detail(ctx, id)
	return response.CheckWithData(ctx, data, err)
}

// details 获取用户详情（查询参数）
// @Summary 获取用户详情
// @Description 根据ID获取用户详情（查询参数方式）
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param id query int true "用户ID"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /test/detail [get]
func (h TestRouter) details(ctx fiber.Ctx) error {
	var idReq request.CommonIdReq
	if err := h.validate.BindQuery(ctx, &idReq); err != nil {
		return err
	}
	data, err := h.svc.Detail(ctx, idReq.Id)
	return response.CheckWithData(ctx, data, err)
}

// edit 编辑用户
// @Summary 编辑用户
// @Description 编辑用户信息
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param request body req.TestEditReq true "编辑参数"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /test/edit [post]
func (h TestRouter) edit(ctx fiber.Ctx) error {
	var editReq req.TestEditReq
	if err := h.validate.BindAndValidate(ctx, &editReq); err != nil {
		return err
	}
	return response.Check(ctx, h.svc.Edit(ctx, editReq))
}

// add 添加用户
// @Summary 添加用户
// @Description 添加新用户
// @Tags 测试模块
// @Accept json
// @Produce json
// @Param request body req.TestAddReq true "添加参数"
// @Success 200 {object} response.Response
// @Security BearerAuth
// @Router /test/add [post]
func (h TestRouter) add(ctx fiber.Ctx) error {
	var addReq req.TestAddReq
	if err := h.validate.BindAndValidate(ctx, &addReq); err != nil {
		return err
	}
	return response.Check(ctx, h.svc.Add(ctx, addReq))
}
