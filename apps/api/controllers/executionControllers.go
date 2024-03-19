package controllers

import (
	"encoding/json"
	"io"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	result := initializers.Db.First(&execution, id)
	if result.Error != nil {
		if execution.ID == 0 {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		log.Fatal(result.Error)
	}

	c.JSON(http.StatusOK, execution)
}

func UpdateOneExecution(c *gin.Context) {
	// get the id
	id := c.Param("id")
	var execution models.Execution
	initializers.Db.First(&execution, id)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	var executionM map[string]interface{}

	err = json.Unmarshal(body, &executionM)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	initializers.Db.Model(&execution).Updates(&executionM)

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
