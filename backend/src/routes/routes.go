package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/controllers"
	"groundTurn/src/repositories"
	"groundTurn/src/services"
)

func Start(addr string) {
	store := repositories.NewStore(repositories.SeedData())
	flightController := controllers.NewFlightTurnaroundController(services.NewFlightTurnaroundService(store))
	taskController := controllers.NewGroundTaskController(services.NewGroundTaskService(store))
	resourceController := controllers.NewGroundResourceController(services.NewGroundResourceService(store))
	bookingController := controllers.NewResourceBookingController(services.NewResourceBookingService(store))
	delayController := controllers.NewDelayEventController(services.NewDelayEventService(store))

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"})
	})
	api := r.Group("/api")
	{
		api.GET("/flight-turnaround", flightController.List)
		api.GET("/flight-turnaround/:id", flightController.Get)
		api.POST("/flight-turnaround/:id/release", flightController.Release)
		api.GET("/ground-task", taskController.List)
		api.GET("/ground-resource", resourceController.List)
		api.GET("/resource-booking", bookingController.List)
		api.GET("/delay-event", delayController.List)
	}
	_ = r.Run(addr)
}
