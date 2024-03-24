package routes

import (
	"job-scheduler/api/docs"
	"job-scheduler/api/middleware"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Summary Health check
// @Tags Default
// @Produce json
// @Success 200
// @Router / [get]
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Job API")
	})

	// Auth
	auth := r.Group("/auth")
	AuthRoutes(auth)

	// User
	user := r.Group("/user", middleware.RequireAuth)
	UserRoutes(user)

	// Scheduler
	job := r.Group("/scheduler", middleware.RequireAuth)
	JobRoutes(job)
	executionRoutes(job)

	// Static files
	// r.Static("/blob", "./blob")

	docs.SwaggerInfo.BasePath = ""

	// url := ginSwagger.URL("http://localhost:8000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// http://localhost:8000/swagger/index.html

	return r
}
