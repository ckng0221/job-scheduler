package routes

import (
	"job-scheduler/api/controllers"

	"github.com/gin-gonic/gin"
)

func JobRoutes(r *gin.RouterGroup) {
	r.POST("/jobs", controllers.CreateJobs)
	r.GET("/jobs", controllers.GetAllJobs)
	r.GET("/jobs/:id", controllers.GetOneJob)
	r.GET("/jobs/:id/executions", controllers.GetOneJobExecutions)
	r.PATCH("/jobs/:id", controllers.UpdateOneJob)
	r.PATCH("/jobs/:id/retrycount", controllers.UpdateOneJobRetryCount)
	r.DELETE("/jobs/:id", controllers.DeleteOneJob)
}
