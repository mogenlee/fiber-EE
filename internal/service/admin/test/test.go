package test

import (
	"fiber-ee/internal/dto/req"
	"fiber-ee/internal/middleware"
	"fiber-ee/internal/model/entity"
	"fiber-ee/internal/model/query"
	"fiber-ee/internal/pkg/request"
	"fiber-ee/internal/pkg/response"
	"fiber-ee/internal/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

// TestService 用户服务接口
type TestService interface {
	List(ctx fiber.Ctx, pageReq *request.PageReq, listReq req.TestListReq) (any, error)
	Detail(ctx fiber.Ctx, id int32) (any, error)
	Edit(ctx fiber.Ctx, editReq req.TestEditReq) error
	Add(ctx fiber.Ctx, addReq req.TestAddReq) error
	Login(ctx fiber.Ctx, loginReq req.TestLoginReq) (any, error)
}

type testService struct {
	db     *query.Query
	jwtCfg *middleware.JWTConfig
}

func (t testService) Login(ctx fiber.Ctx, loginReq req.TestLoginReq) (any, error) {
	m := t.db.User
	q := m.WithContext(ctx)

	user, err := q.Where(m.Username.Eq(loginReq.Username)).First()
	if err != nil {
		return nil, response.ErrLogin
	}
	password, _ := utils.HashPassword(loginReq.Password)
	// 验证密码
	if !utils.CheckPassword(loginReq.Password, password) {
		return nil, response.ErrLogin
	}

	// 生成 Token Pair
	tokenPair, err := middleware.GenerateTokenPair(t.jwtCfg, uint(user.ID), user.Username, "admin")
	if err != nil {
		return nil, response.ErrInternal
	}

	return tokenPair, nil
}

func (t testService) Add(ctx fiber.Ctx, addReq req.TestAddReq) error {
	m := t.db.User
	q := m.WithContext(ctx)

	var obj entity.User
	utils.Copy(&obj, addReq)
	// 密码加密
	obj.Password, _ = utils.HashPassword(obj.Password)
	return q.Create(&obj)
}

func (t testService) Edit(ctx fiber.Ctx, editReq req.TestEditReq) error {
	m := t.db.User
	q := m.WithContext(ctx)

	var obj entity.User
	utils.Copy(&obj, editReq)
	_, err := q.Where(m.ID.Eq(editReq.Id)).Updates(obj)
	return err
}

func (t testService) Detail(ctx fiber.Ctx, id int32) (any, error) {
	m := t.db.User
	q := m.WithContext(ctx)

	if id > 0 {
		q = q.Where(m.ID.Eq(id))
	}
	return q.First()
}

// NewTestService 创建服务（db 由 dig 自动注入）
func NewTestService(db *query.Query, jwtCfg *middleware.JWTConfig) TestService {
	return &testService{db: db, jwtCfg: jwtCfg}
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

	// 分页查询
	result, total, err := q.FindByPage(pageReq.GetOffset(), pageReq.GetPageSize())

	if err != nil {
		return nil, response.ErrQuery
	}

	return response.PageData{
		List:     result,
		Total:    total,
		PageNo:   pageReq.GetPageNo(),
		PageSize: pageReq.GetPageSize(),
	}, nil
}
