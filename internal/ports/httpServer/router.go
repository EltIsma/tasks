package httpserver

import (
	"log/slog"

	"github.com/go-redis/redis_rate/v9"
	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Register swagger docs.
	_ "task/docs/swagger"
	ratelimiter "task/pkg/rate-limiter"

	"github.com/gin-gonic/gin"
)

// New 		godoc
// @title 	Tasks API
// @version 1.0
func New(handler *Handler, logger *slog.Logger, rL *redis_rate.Limiter) *gin.Engine {
	router := gin.New()
	registerSwagger(router)
	registerGroup(router, handler, logger, rL)

	return router
}

func registerGroup(e *gin.Engine, handler *Handler, logger *slog.Logger, rL *redis_rate.Limiter) {
	r := e.Group("api/v1")

	ratelimiter.Limiter = rL
	r.Use(ratelimiter.RateLimit(logger))

	r.POST("/task", handler.CreateTask)
	r.POST("task/create-with-assignment", handler.CreateTaskWithAssignment)
	r.POST("/task/assignment", handler.AssignTaskToClasses)
	r.PUT("/task/assignment-update", handler.UpdateTaskAssignment)
	r.POST("/task/result", handler.TaskResult)
	r.GET("/task/all", handler.GetTasks)
	r.GET("/task/:id", handler.GetTask)
	r.GET("/task/get-by-class", handler.GetTaskByClass)
	r.PUT("/task/:id/update", handler.UpdateTask)
	r.DELETE("/task/:id/delete", handler.DeleteTask)
	r.DELETE("/task/assignment-delete", handler.DeleteAssignment)
}

func registerSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))
}
