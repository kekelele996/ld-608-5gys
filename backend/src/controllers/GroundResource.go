package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type GroundResourceController struct {
	service *services.GroundResourceService
}

func NewGroundResourceController(service *services.GroundResourceService) *GroundResourceController {
	return &GroundResourceController{service: service}
}

func (ctl *GroundResourceController) List(c *gin.Context) { c.JSON(http.StatusOK, ctl.service.List()) }
