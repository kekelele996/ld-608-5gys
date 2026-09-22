package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/config"
	"groundTurn/src/controllers"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/services"
	"groundTurn/src/utils"
)

// Start 装配数据库、仓储、服务、控制器与中间件，启动 HTTP 服务。
func Start(addr string) {
	cfg := config.Load()
	db := config.OpenDB(cfg)
	config.SeedIfEmpty(db)

	locker := utils.NewTurnaroundLocker()

	turnRepo := repositories.NewFlightTurnaroundRepository(db)
	taskRepo := repositories.NewGroundTaskRepository(db)
	resourceRepo := repositories.NewGroundResourceRepository(db)
	bookingRepo := repositories.NewResourceBookingRepository(db)
	delayRepo := repositories.NewDelayEventRepository(db)
	auditRepo := repositories.NewAuditLogRepository(db)

	turnService := services.NewFlightTurnaroundService(db, turnRepo, taskRepo, bookingRepo, auditRepo, locker)
	taskService := services.NewGroundTaskService(db, taskRepo, turnRepo, auditRepo, locker)
	resourceService := services.NewGroundResourceService(db, resourceRepo)
	bookingService := services.NewResourceBookingService(db, bookingRepo, auditRepo, locker)
	delayService := services.NewDelayEventService(db, delayRepo, auditRepo, locker)

	turnCtl := controllers.NewFlightTurnaroundController(turnService)
	taskCtl := controllers.NewGroundTaskController(taskService)
	resourceCtl := controllers.NewGroundResourceController(resourceService)
	bookingCtl := controllers.NewResourceBookingController(bookingService)
	delayCtl := controllers.NewDelayEventController(delayService)

	r := gin.Default()
	r.Use(middlewares.ErrorHandlerMiddleware())
	r.Use(middlewares.RateLimitMiddleware())
	r.Use(middlewares.AuditLogMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"})
	})

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	api.Use(middlewares.RBACMiddleware())

	// 读接口
	api.GET("/flight-turnaround", turnCtl.List)
	api.GET("/flight-turnaround/:id", turnCtl.Get)
	api.GET("/flight-turnaround/:id/gate", turnCtl.Gate)
	api.GET("/ground-task", taskCtl.List)
	api.GET("/ground-resource", resourceCtl.List)
	api.GET("/resource-booking", bookingCtl.List)
	api.GET("/delay-event", delayCtl.List)

	// 过站放行联动写接口
	api.POST("/flight-turnaround/:id/release", turnCtl.Release)
	api.POST("/ground-task/:id/sign", taskCtl.Sign)
	api.POST("/delay-event/:id/resolve", delayCtl.Resolve)
	api.POST("/resource-booking/:id/resolve-conflict", bookingCtl.ResolveConflict)

	r.Run(addr)
}
