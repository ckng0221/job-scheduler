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

func GetAllJobs(c *gin.Context) {
	userId := c.Query("user_id")

	m := make(map[string]interface{})

	if userId != "" {
		m["user_id"] = userId
	}

	var jobs []models.Job
	initializers.Db.Where(m).Find(&jobs)

	c.JSON(http.StatusOK, jobs)
}

// By admin only
func CreateJobs(c *gin.Context) {
	var jobs []models.Job

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// var jobsM []map[string]interface{}
	err = json.Unmarshal(body, &jobs)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	result := initializers.Db.Model(&jobs).Create(&jobs)
	if result.Error != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, jobs)
}

func GetOneJob(c *gin.Context) {
	id := c.Param("id")

	var job models.Job
	result := initializers.Db.First(&job, id)
	if result.Error != nil {
		if job.ID == 0 {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		log.Fatal(result.Error)
	}

	c.JSON(http.StatusOK, job)
}

func GetOneJobExecutions(c *gin.Context) {
	id := c.Param("id")

	var executions []models.Execution
	m := make(map[string]interface{})
	m["job_id"] = id

	initializers.Db.Where(m).Find(&executions).Order("ID")
	c.JSON(http.StatusOK, executions)
}

func UpdateOneJob(c *gin.Context) {
	// get the id
	id := c.Param("id")
	var job models.Job
	initializers.Db.First(&job, id)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	var jobM map[string]interface{}

	err = json.Unmarshal(body, &jobM)

	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	initializers.Db.Model(&job).Updates(&jobM)

	c.JSON(200, job)
}

func DeleteOneJob(c *gin.Context) {
	id := c.Param("id")

	result := initializers.Db.Delete(&models.Job{}, id)
	if result.Error != nil {
		c.AbortWithStatus(500)
	}

	// response
	c.Status(202)
}
