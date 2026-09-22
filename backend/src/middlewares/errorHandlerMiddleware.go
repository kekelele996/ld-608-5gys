package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/types"
)

// ErrorHandlerMiddleware 兜底 panic，避免单点吞掉异常。
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, types.APIResponse{
					OK:    false,
					Error: types.NewAppError(constants.InternalError, constants.InternalErrorMessage),
				})
			}
		}()
		c.Next()
	}
}
