package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type GroundTaskController struct {
	service *services.GroundTaskService
}

func NewGroundTaskController(service *services.GroundTaskService) *GroundTaskController {
	return &GroundTaskController{service: service}
}

func (ctl *GroundTaskController) List(c *gin.Context) {
	rows, err := ctl.service.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

// Sign POST /api/ground-task/:id/sign —— 任务签收（SIGNED 是放行门禁条件之一）。
func (ctl *GroundTaskController) Sign(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dto, err := ctl.service.Sign(id, actorFromContext(c))
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, dto)
}
