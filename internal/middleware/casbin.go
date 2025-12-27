package middleware

import (
	"fiber-ee/internal/pkg/response"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v3"
)

// CasbinAuth Casbin 权限校验中间件
// 需要在 JWTAuth 之后使用，依赖 context 中的 role
func CasbinAuth(enforcer *casbin.Enforcer) fiber.Handler {
	// 构建白名单 map
	whitelistMap := make(map[string]struct{}, len(notAuthList))
	for _, path := range notAuthList {
		whitelistMap[path] = struct{}{}
	}

	return func(c fiber.Ctx) error {
		// 检查白名单（与 JWT 共用）
		if _, ok := whitelistMap[c.Path()]; ok {
			return c.Next()
		}

		// 获取当前用户角色（从 JWT 中间件设置）
		role := GetRole(c)
		if role == "" {
			return response.ErrUnauthorized
		}
		// 权限校验：role, path, method
		ok, err := enforcer.Enforce(role, c.Path(), c.Method())
		if err != nil {
			return response.ErrInternal
		}

		if !ok {
			return response.ErrForbidden
		}

		return c.Next()
	}
}
