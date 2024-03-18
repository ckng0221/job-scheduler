package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"job-scheduler/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// Get the email/pass off req bodyu
	var body struct {
		Name     string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create the user
	user := models.User{Name: body.Name, Email: body.Email, Password: string(hash), Role: "member"}
	result := initializers.Db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Respond
	c.JSON(http.StatusCreated, user)
}

func GetAllUsers(c *gin.Context) {
	email := c.Query("email")

	m := make(map[string]interface{})

	if email != "" {
		m["email"] = email
	}

	var users []models.User
	initializers.Db.Scopes(utils.Paginate(c)).Where(m).Find(&users)

	c.JSON(http.StatusOK, users)
}

// By admin only
func CreateUsers(c *gin.Context) {
	var users []models.User

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	fmt.Println(string(body))

	err = json.Unmarshal(body, &users)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	result := initializers.Db.Create(&users)
	if result.Error != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, users)
}

func GetOneUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	result := initializers.Db.First(&user, id)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateOneUser(c *gin.Context) {
	// get the id
	id := c.Param("id")
	var user models.User
	initializers.Db.First(&user, id)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	var userM map[string]interface{}

	err = json.Unmarshal(body, &userM)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	initializers.Db.Model(&user).Updates(&userM)

	c.JSON(200, user)
}

func DeleteOneUser(c *gin.Context) {
	id := c.Param("id")

	initializers.Db.Delete(&models.User{}, id)

	// response
	c.Status(202)
}

func GetUserRoles(c *gin.Context) {
	c.JSON(http.StatusOK, models.Roles)
}

func UploadProfilePicture(c *gin.Context) {
	userId := c.Param("id")

	file, _ := c.FormFile("file")
	// log.Println(file.Filename)

	filePath := fmt.Sprintf("./blob/profilepic/%s/%s", userId, file.Filename)
	// Upload file
	c.SaveUploadedFile(file, filePath)
	cleanedFilePath := filePath[1:] // remove relative .

	c.JSON(http.StatusOK, gin.H{"filepath": cleanedFilePath})
}
