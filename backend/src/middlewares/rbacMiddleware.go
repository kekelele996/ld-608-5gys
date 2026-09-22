package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/types"
)

// 角色：DISPATCHER 地勤调度 / TEAM 班组 / RESOURCE_MANAGER 资源管理员 / SUPERVISOR 运行督导。
var writeRoles = map[string]bool{
	"DISPATCHER": true,
	"SUPERVISOR": true,
}

// RBACMiddleware 放行/签收等写操作仅地勤调度与运行督导可执行。
func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}
		role, _ := c.Get("role")
		roleText, _ := role.(string)
		if !writeRoles[roleText] {
			c.AbortWithStatusJSON(http.StatusForbidden, types.APIResponse{
				OK:    false,
				Error: types.NewAppError(constants.RBACDenied, constants.RBACDeniedMessage),
			})
			return
		}
		c.Next()
	}
}
