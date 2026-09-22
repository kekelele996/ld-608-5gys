package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type FlightTurnaroundController struct {
	service *services.FlightTurnaroundService
}

func NewFlightTurnaroundController(service *services.FlightTurnaroundService) *FlightTurnaroundController {
	return &FlightTurnaroundController{service: service}
}

func (ctl *FlightTurnaroundController) List(c *gin.Context) {
	rows, err := ctl.service.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}

func (ctl *FlightTurnaroundController) Get(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	dto, err := ctl.service.Get(id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, dto)
}

// Gate 返回放行门禁快照（看板阻塞明细）。
func (ctl *FlightTurnaroundController) Gate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	gate, err := ctl.service.Gate(id)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, gate)
}

// Release POST /api/flight-turnaround/:id/release —— 过站放行联动入口。
func (ctl *FlightTurnaroundController) Release(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	actor := actorFromContext(c)
	result, err := ctl.service.Release(id, actor)
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, result)
}

func actorFromContext(c *gin.Context) string {
	if value, exists := c.Get("actor"); exists {
		if actor, ok := value.(string); ok && actor != "" {
			return actor
		}
	}
	var body struct {
		Actor string `json:"actor"`
	}
	_ = c.ShouldBindJSON(&body)
	if body.Actor != "" {
		return body.Actor
	}
	return "dispatcher"
}
