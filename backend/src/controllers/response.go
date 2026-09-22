package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/types"
)

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, types.APIResponse{OK: true, Data: data})
}

// fail controller 层再次包装 service 异常，映射 HTTP 状态码，禁止全局单点吞异常。
func fail(c *gin.Context, err error) {
	var appErr *types.AppError
	if errors.As(err, &appErr) {
		c.JSON(statusFor(appErr.Code), types.APIResponse{OK: false, Error: appErr})
		return
	}
	c.JSON(http.StatusInternalServerError, types.APIResponse{
		OK: false, Error: types.NewAppError(constants.InternalError, constants.InternalErrorMessage),
	})
}

func statusFor(code string) int {
	switch code {
	case constants.TurnaroundNotFound, constants.TaskNotFound,
		constants.BookingNotFound, constants.DelayNotFound:
		return http.StatusNotFound
	case constants.ReleaseBlocked, constants.ReleaseAlreadyDone:
		// 门禁不通过 / 重复放行：业务拒绝，请求本身合法但不可执行。
		return http.StatusConflict
	case constants.ReleaseRaceLost:
		// 并发提交只生效一次：后到者同样 409。
		return http.StatusConflict
	case constants.TaskAlreadySigned:
		return http.StatusConflict
	case constants.AuthRequired:
		return http.StatusUnauthorized
	case constants.RBACDenied:
		return http.StatusForbidden
	case constants.ValidationFailed:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
