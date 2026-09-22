package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type ResourceBookingController struct {
	service *services.ResourceBookingService
}

func NewResourceBookingController(service *services.ResourceBookingService) *ResourceBookingController {
	return &ResourceBookingController{service: service}
}

func (ctl *ResourceBookingController) List(c *gin.Context) {
	rows, err := ctl.service.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// ResolveConflict POST /api/resource-booking/:id/resolve-conflict —— 解除预约冲突硬阻塞。
func (ctl *ResourceBookingController) ResolveConflict(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dto, err := ctl.service.ResolveConflict(id, actorFromContext(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, dto)
}
