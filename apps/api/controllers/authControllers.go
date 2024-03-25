package controllers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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

// TODO: Change to proper authentication method to OIDC
func Login(c *gin.Context) {
	// User login based on access token
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
		return
	}

	userData, err := getUserDataByTokenExchange(body.Code, c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	// fmt.Println(userData)
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
	err = initializers.Db.Where("sub = ?", profile.Id).Limit(1).Find(&user).Error
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
	// c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, false)
	fmt.Println("why!!!")

	// respond
	c.JSON(http.StatusOK, gin.H{
		"name":         profile.Name,
		"access_token": tokenString,
	})
}

func GoogleLogin(c *gin.Context) {

	randomState, _ := GenerateSHA256State()

	// Return google login URL
	url := initializers.AppConfig.GoogleLoginConfig.AuthCodeURL(randomState)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func GoogleExchangeToken(c *gin.Context) {
	code := c.Query("code")

	userData, err := getUserDataByTokenExchange(code, c)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	c.String(200, string(userData))
}

func getUserDataByTokenExchange(code string, c *gin.Context) ([]byte, error) {
	googlecon := initializers.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	// Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, false)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

func GenerateSHA256State() (string, error) {
	// Generate 1024 random bytes
	randomBytes := make([]byte, 1024)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Compute SHA256 hash
	hasher := sha256.New()
	hasher.Write(randomBytes)
	hashBytes := hasher.Sum(nil)

	// Convert hash bytes to hexadecimal string
	state := hex.EncodeToString(hashBytes)

	return state, nil
}
