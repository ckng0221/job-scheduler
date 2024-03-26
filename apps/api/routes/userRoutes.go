package routes

import (
	"job-scheduler/api/controllers"
	"job-scheduler/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {

	// r.POST("/signup", controllers.Signup)
	r.GET("/users", controllers.GetAllUsers)
	r.GET("/roles", controllers.GetUserRoles)
	r.POST("/users", controllers.CreateUsers)
	r.GET("/users/:id", controllers.GetOneUser)
	r.PATCH("/users/:id", controllers.UpdateOneUser)
	r.DELETE("/users/:id", middleware.RequireAdmin, controllers.DeleteOneUser)
}
