package middleware

import (
	"fmt"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
		c.Set("user", user)

		// Attach to req

		// Continue
		c.Next()
		return
	}
	fmt.Println("token", bearerToken)

	if bearerToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/validate it
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		// return
		fmt.Println("Token not found")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Token expired")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Find the user with token sub
		var user models.User
		initializers.Db.Where("sub = ?", claims["sub"]).First(&user)

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
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}

func RequireAdmin(c *gin.Context) {
	requestUser, _ := c.Get("user")
	if requestUser.(models.User).Role != "admin" {
		c.AbortWithStatus(403)
		return
	}
}
