package routes

import (
	"job-scheduler/api/controllers"
	"job-scheduler/api/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {

	r.POST("/login", controllers.Login)
	r.GET("/google-login", controllers.GoogleLogin)
	r.POST("/google-token-exchange", middleware.RequireAuth, controllers.GoogleExchangeToken)
	r.POST("/logout", controllers.Logout)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
}
