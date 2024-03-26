package routes

import (
	"job-scheduler/api/controllers"
	"job-scheduler/api/middleware"

	"github.com/gin-gonic/gin"
)

func executionRoutes(r *gin.RouterGroup) {
	r.POST("/executions", middleware.RequireAdmin, controllers.CreateExecutions)
	r.GET("/executions", middleware.RequireAdmin, controllers.GetAllExecutions)
	r.GET("/executions/:id", middleware.RequireAdmin, controllers.GetOneExecution)
	r.PATCH("/executions/:id", middleware.RequireAdmin, controllers.UpdateOneExecution)
	r.DELETE("/executions/:id", middleware.RequireAdmin, controllers.DeleteOneExecution)
}
