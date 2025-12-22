package middleware

import (
	"fiber-ee/config"
	"fiber-ee/internal/pkg/response"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair Access Token 和 Refresh Token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string
	Expire        time.Duration
	RefreshExpire time.Duration
}

// JWTWhitelist JWT 认证白名单（不需要认证的路径）
var JWTWhitelist = []string{
	"/admin/v1/test/login",
	"/admin/v1/auth/refresh",
	"/api/v1/users/register",
}

// NewJWTConfig 从配置创建 JWT 配置
func NewJWTConfig(cfg *config.Config) *JWTConfig {
	return &JWTConfig{
		Secret:        cfg.JWT.Secret,
		Expire:        time.Duration(cfg.JWT.Expire) * time.Second,
		RefreshExpire: time.Duration(cfg.JWT.RefreshExpire) * time.Second,
	}
}

// GenerateToken 生成 JWT token（仅 Access Token）
func GenerateToken(cfg *JWTConfig, userID uint, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateTokenPair 生成 Access Token 和 Refresh Token
func GenerateTokenPair(cfg *JWTConfig, userID uint, username, role string) (*TokenPair, error) {
	// Access Token
	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh Token（有效期更长）
	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		ExpiresIn:    int64(cfg.Expire.Seconds()),
	}, nil
}

// RefreshToken 使用 Refresh Token 刷新 Access Token
func RefreshToken(cfg *JWTConfig, refreshTokenStr string) (*TokenPair, error) {
	// 解析 Refresh Token
	token, err := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, response.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || claims.Subject != "refresh" {
		return nil, response.ErrTokenInvalid
	}

	// 生成新的 Token Pair
	return GenerateTokenPair(cfg, claims.UserID, claims.Username, claims.Role)
}

// JWTAuth JWT 认证中间件
func JWTAuth(cfg *config.Config) fiber.Handler {
	// 构建白名单 map（闭包内只初始化一次）
	whitelistMap := make(map[string]struct{}, len(JWTWhitelist))
	for _, path := range JWTWhitelist {
		whitelistMap[path] = struct{}{}
	}

	secret := []byte(cfg.JWT.Secret)

	return func(c fiber.Ctx) error {
		// 检查白名单
		if _, ok := whitelistMap[c.Path()]; ok {
			return c.Next()
		}
		// 从 Authorization header 获取 token
		auth := c.Get("Authorization")
		const prefix = "Bearer "
		if len(auth) <= len(prefix) || auth[:len(prefix)] != prefix {
			return response.ErrUnauthorized
		}

		tokenString := auth[len(prefix):]

		// 解析 token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			return response.ErrTokenInvalid
		}

		// 将用户信息存入 context
		if claims, ok := token.Claims.(*Claims); ok {
			c.Locals("user_id", claims.UserID)
			c.Locals("username", claims.Username)
			c.Locals("role", claims.Role)
			c.Locals("claims", claims)
		}

		return c.Next()
	}
}

// GetUserID 从 context 获取当前用户 ID
func GetUserID(c fiber.Ctx) uint {
	if id, ok := c.Locals("user_id").(uint); ok {
		return id
	}
	return 0
}

// GetUsername 从 context 获取当前用户名
func GetUsername(c fiber.Ctx) string {
	if name, ok := c.Locals("username").(string); ok {
		return name
	}
	return ""
}

// GetRole 从 context 获取当前用户角色
func GetRole(c fiber.Ctx) string {
	if role, ok := c.Locals("role").(string); ok {
		return role
	}
	return ""
}

// GetClaims 从 context 获取完整 Claims
func GetClaims(c fiber.Ctx) *Claims {
	if claims, ok := c.Locals("claims").(*Claims); ok {
		return claims
	}
	return nil
}
