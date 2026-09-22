package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

type ResourceBookingController struct {
	service *services.ResourceBookingService
}

func NewResourceBookingController(service *services.ResourceBookingService) *ResourceBookingController {
	return &ResourceBookingController{service: service}
}

func (ctl *ResourceBookingController) List(c *gin.Context) { c.JSON(http.StatusOK, ctl.service.List()) }
