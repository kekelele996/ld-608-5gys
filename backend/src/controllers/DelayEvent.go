package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type DelayEventController struct {
	service *services.DelayEventService
}

func NewDelayEventController(service *services.DelayEventService) *DelayEventController {
	return &DelayEventController{service: service}
}

func (ctl *DelayEventController) List(c *gin.Context) {
	rows, err := ctl.service.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Resolve POST /api/delay-event/:id/resolve —— 关闭未关闭延误。
func (ctl *DelayEventController) Resolve(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dto, err := ctl.service.Resolve(id, actorFromContext(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, dto)
}
