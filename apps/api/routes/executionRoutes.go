package routes

import (
	"job-scheduler/api/controllers"

	"github.com/gin-gonic/gin"
)

func executionRoutes(r *gin.RouterGroup) {
	r.POST("/executions", controllers.CreateExecutions)
	r.GET("/executions", controllers.GetAllExecutions)
	r.GET("/executions/:id", controllers.GetOneExecution)
	r.PATCH("/executions/:id", controllers.UpdateOneExecution)
	r.DELETE("/executions/:id", controllers.DeleteOneExecution)
}
