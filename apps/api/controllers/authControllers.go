package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Login(c *gin.Context) {

	// Get acccess token
	var body struct {
		Code string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	var profile models.GoogleProfile

	if body.Code == "" {
		c.AbortWithStatusJSON(401, "No access token")
	}

	userData := getUserDataByTokenExchange(body.Code, c)
	fmt.Println(userData)
	json.Unmarshal(userData, &profile)

	// generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": profile.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// register user if not exist
	var user models.User
	err = initializers.Db.Where("sub = ?", profile.Id).Find(&user).Error
	fmt.Println(err)
	if user.ID == 0 {
		fmt.Println("user not found")
		initializers.Db.Create(&models.User{
			Name:       profile.Name,
			Email:      profile.Email,
			Sub:        profile.Id,
			ProfilePic: profile.Picture,
		})
		fmt.Println("user created")
	}

	// Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"name":         profile.Name,
		"access_token": tokenString,
	})
}

func GoogleLogin(c *gin.Context) {
	url := initializers.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func GoogleExchangeToken(c *gin.Context) {
	code := c.Query("code")

	userData := getUserDataByTokenExchange(code, c)

	c.String(200, string(userData))
}

func getUserDataByTokenExchange(code string, c *gin.Context) []byte {
	googlecon := initializers.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		c.String(401, "Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(401, "User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(500)
	}
	return userData
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	// Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}
