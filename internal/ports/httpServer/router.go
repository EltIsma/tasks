package httpserver

import (
	"log/slog"

	swaggoFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Register swagger docs.
	_ "task/docs/swagger"

	"github.com/gin-gonic/gin"
)

// New 		godoc
// @title 	Tasks API
// @version 1.0
func New(handler *Handler, logger *slog.Logger) *gin.Engine {
	router := gin.New()
	registerSwagger(router)
	registerGroup(router, handler, logger)

	return router
}

func registerGroup(e *gin.Engine, handler *Handler, logger *slog.Logger) {
	r := e.Group("api/v1")

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
