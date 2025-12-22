package auth

import (
	"fiber-ee/internal/middleware"
	"fiber-ee/internal/pkg/request"
	"fiber-ee/internal/pkg/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type Router struct {
	Log      *zap.Logger
	jwtCfg   *middleware.JWTConfig
	validate request.Validate
}

func NewRouter(log *zap.Logger, jwtCfg *middleware.JWTConfig, validate request.Validate) *Router {
	return &Router{
		Log:      log,
		jwtCfg:   jwtCfg,
		validate: validate,
	}
}

func (h Router) Register(root fiber.Router) {
	group := root.Group("/auth")
	group.Post("/refresh", h.refresh).Name("auth_refresh")
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h Router) refresh(ctx fiber.Ctx) error {
	var req RefreshReq
	if err := h.validate.BindAndValidate(ctx, &req); err != nil {
		return err
	}

	tokenPair, err := middleware.RefreshToken(h.jwtCfg, req.RefreshToken)
	if err != nil {
		return err
	}

	return response.Success(ctx, tokenPair)
}
