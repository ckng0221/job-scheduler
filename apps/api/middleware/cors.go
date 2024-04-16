package middleware

import (
	"job-scheduler/utils"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {

	corsAllowedOrigin := utils.Getenv("CORS_ALLOWED_ORIGIN", "http://localhost:3000")

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", corsAllowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
