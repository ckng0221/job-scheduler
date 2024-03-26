package middleware

import (
	"context"
	"fmt"
	"job-scheduler/api/config"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	// get the cookie of req
	// tokenString, err := c.Cookie("Authorization")
	bearerTokenArr := strings.Split(c.GetHeader("Authorization"), "Bearer ")
	bearerToken := bearerTokenArr[len(bearerTokenArr)-1]
	apiKey := c.GetHeader("x-api-key")
	if apiKey == os.Getenv("ADMIN_API_KEY") {
		// Find the user with token sub
		var user models.User
		initializers.Db.Where("id = ?", 1).First(&user)

		if user.ID == 0 {
			fmt.Println("User not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// fmt.Println(user.ID)
		c.Set("user", models.User{Name: "Bot", Role: "admin"})

		// Attach to req

		// Continue
		c.Next()
		return
	}
	// fmt.Println("token", bearerToken)

	if bearerToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Decode/validate it
	verifier := config.GetVerifier()
	ctx := context.Background()

	// JWT token from identify provider
	idToken, err := verifier.Verify(ctx, bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
		return

	}

	var claims config.IDTokenClaims
	if err := idToken.Claims(&claims); err != nil {
		// handle error
		fmt.Printf("Sub not found")
		c.AbortWithStatus(500)
		return
	}

	// Find the user with token sub
	var user models.User
	initializers.Db.Where("sub = ?", claims.Sub).First(&user)

	if user.ID == 0 {
		fmt.Println("User not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// fmt.Println(user.ID)
	c.Set("user", user)

	// Attach to req

	// Continue
	c.Next()

	// fmt.Println(claims["foo"], claims["nbf"])

}

func RequireAdmin(c *gin.Context) {
	requestUser, _ := c.Get("user")
	if requestUser.(models.User).Role != "admin" {
		c.AbortWithStatus(403)
		return
	}
}
