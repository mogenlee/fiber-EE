package middleware

// whitelist 不需要认证的路由白名单
var whitelist = []string{
	"/admin/v1/auth/login",
}
