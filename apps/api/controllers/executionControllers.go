package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllExecutions(c *gin.Context) {
	jobId := c.Query("job_id")

	m := make(map[string]interface{})

	if jobId != "" {
		m["job_id"] = jobId
	}

	var executions []models.Execution
	initializers.Db.Where(m).Find(&executions)

	c.JSON(http.StatusOK, executions)
}

func CreateExecutions(c *gin.Context) {
	var executions []models.Execution

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// var executionsM []map[string]interface{}
	err = json.Unmarshal(body, &executions)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	result := initializers.Db.Model(&executions).Create(&executions)
	if result.Error != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, executions)
}

func GetOneExecution(c *gin.Context) {
	id := c.Param("id")

	var execution models.Execution
	err := initializers.Db.First(&execution, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusOK, execution)
}

func UpdateOneExecution(c *gin.Context) {
	// get the id
	id := c.Param("id")
	var execution models.Execution
	var executionPatch models.ExecutionUpdate
	initializers.Db.First(&execution, id)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = json.Unmarshal(body, &executionPatch)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	initializers.Db.Model(&execution).Updates(&executionPatch)

	c.JSON(200, execution)
}

func DeleteOneExecution(c *gin.Context) {
	id := c.Param("id")

	result := initializers.Db.Delete(&models.Execution{}, id)
	if result.Error != nil {
		c.AbortWithStatus(500)
	}

	// response
	c.Status(202)
}
