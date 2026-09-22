package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type DelayEventController struct{ service *services.DelayEventService }

func NewDelayEventController(service *services.DelayEventService) *DelayEventController {
	return &DelayEventController{service: service}
}

func (ctl *DelayEventController) List(c *gin.Context) { c.JSON(http.StatusOK, ctl.service.List()) }
