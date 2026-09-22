package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type GroundTaskController struct{ service *services.GroundTaskService }

func NewGroundTaskController(service *services.GroundTaskService) *GroundTaskController {
	return &GroundTaskController{service: service}
}

func (ctl *GroundTaskController) List(c *gin.Context) { c.JSON(http.StatusOK, ctl.service.List()) }
