package routes

import (
	"job-scheduler/api/controllers"
	"job-scheduler/api/middleware"

	"github.com/gin-gonic/gin"
)

func JobRoutes(r *gin.RouterGroup) {
	r.POST("/jobs", controllers.CreateJobs)
	r.GET("/jobs", controllers.GetAllJobs)
	r.GET("/jobs/:id", controllers.GetOneJob)
	r.GET("/jobs/:id/task-script", controllers.GetTaskScript)
	r.POST("/jobs/:id/task-script", controllers.UploadTaskScript)
	r.GET("/jobs/:id/executions", controllers.GetOneJobExecutions)
	r.PATCH("/jobs/:id", controllers.UpdateOneJob)
	r.PATCH("/jobs/:id/retrycount", middleware.RequireAdmin, controllers.UpdateOneJobRetryCount)
	r.DELETE("/jobs/:id", controllers.DeleteOneJob)
}
