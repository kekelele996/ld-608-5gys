package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 简化版身份解析：本地联调允许 X-Actor / X-Role 头模拟身份；
// 携带 Authorization: Bearer 时预留 JWT 验签入口（生产使用 config.JWTSecret）。
// 未带头的浏览器请求按 dispatcher 处理，显式 X-Role=VIEWER 仅可读。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		actor := c.GetHeader("X-Actor")
		role := strings.ToUpper(c.GetHeader("X-Role"))

		if actor == "" && strings.HasPrefix(c.GetHeader("Authorization"), "Bearer ") {
			actor = "dispatcher"
		}
		if actor == "" {
			actor = "dispatcher"
		}
		if role == "" {
			role = "DISPATCHER"
		}
		c.Set("actor", actor)
		c.Set("role", role)
		c.Next()
	}
}
