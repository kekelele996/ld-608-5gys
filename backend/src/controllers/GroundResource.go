package controllers

import (
	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type GroundResourceController struct {
	service *services.GroundResourceService
}

func NewGroundResourceController(service *services.GroundResourceService) *GroundResourceController {
	return &GroundResourceController{service: service}
}

func (ctl *GroundResourceController) List(c *gin.Context) {
	rows, err := ctl.service.List()
	if err != nil {
		fail(c, err)
		return
	}
	ok(c, rows)
}
