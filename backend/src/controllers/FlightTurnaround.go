package controllers

import (
	"net/http"
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
	c.JSON(http.StatusOK, ctl.service.List())
}

func (ctl *FlightTurnaroundController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_FAILED", "message": "航班 ID 必须是数字"})
		return
	}
	detail, found := ctl.service.Get(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"code": "TURNAROUND_NOT_FOUND", "message": "航班过站不存在"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (ctl *FlightTurnaroundController) Release(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "VALIDATION_FAILED", "message": "航班 ID 必须是数字"})
		return
	}
	result, serviceErr := ctl.service.Release(id)
	if serviceErr != nil {
		if releaseErr, ok := serviceErr.(*services.ReleaseError); ok {
			payload := gin.H{"code": releaseErr.Code, "message": releaseErr.Message}
			if releaseErr.Blockers != nil {
				payload["blockers"] = releaseErr.Blockers
			}
			c.JSON(releaseErr.HTTPCode, payload)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "放行失败，请稍后重试"})
		return
	}
	c.JSON(http.StatusOK, result)
}
