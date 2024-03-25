package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"job-scheduler/api/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllJobs(c *gin.Context) {
	userId := c.Query("user_id")
	isActive := c.Query("active")
	MaxRetryCount := 4

	var jobs []models.Job
	m := make(map[string]interface{})

	if userId != "" {
		m["user_id"] = userId
	}
	if isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil && active {
			m["is_completed"] = false
			m["is_disabled"] = false
			m["is_running"] = false
			currentMinute, nextMinute := utils.GetUnixMinuteRange(time.Now())
			// For one-time job, only check if it is within the time boundary
			// For recurring job, will check as long as less then current time
			initializers.Db.Where("retry_count < ?", MaxRetryCount).Where(m).Where("(is_recurring = true AND next_run_time <= ?) OR (is_recurring = false AND next_run_time >= ? AND next_run_time < ?)", time.Now().Unix(), currentMinute.Unix(), nextMinute.Unix()).Find(&jobs)
		}
	} else {
		initializers.Db.Where(m).Find(&jobs)
	}

	c.JSON(http.StatusOK, jobs)
}

// By admin only
func CreateJobs(c *gin.Context) {
	var job models.Job

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = json.Unmarshal(body, &job)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// Update next run time for recurring job
	if job.IsRecurring {
		if job.Cron == "" {
			c.AbortWithError(422, errors.New("cron cannot be empty"))
			return
		}
		// cronExp, err := cronexpr.Parse(job.Cron)
		// if err != nil {
		// 	c.AbortWithError(422, errors.New("invalid cron expression"))
		// 	return
		// }
		// nextRunTimeUnix := cronExp.Next(time.Now().UTC()).Unix()
		// jobM["NextRunTime"] = nextRunTimeUnix
	}

	result := initializers.Db.Model(&job).Create(&job)
	if result.Error != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(http.StatusCreated, job)
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
		fmt.Println(result.Error)
		c.AbortWithStatus(500)
		return
	}

	err := requireOwner(c, job.UserID)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetOneJobExecutions(c *gin.Context) {
	id := c.Param("id")

	var job models.Job
	err := initializers.Db.First(&job, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}
	err = requireOwner(c, job.UserID)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

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
	var jobUpdate models.JobUpdate
	initializers.Db.First(&job, id)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	err = json.Unmarshal(body, &jobUpdate)

	if err != nil {
		c.AbortWithError(422, err)
		return
	}

	err = initializers.Db.First(&job, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": "Record not found"})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}
	err = requireOwner(c, job.UserID)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	initializers.Db.Model(&job).Updates(&jobUpdate)

	c.JSON(200, job)
}

func UpdateOneJobRetryCount(c *gin.Context) {
	id := c.Param("id")

	var job models.Job
	initializers.Db.First(&job, id)
	// Optimistic lock
	result := initializers.Db.Model(&job).Where("id = ? AND retry_count = ?", id, job.RetryCount).Update("retry_count", job.RetryCount+1)
	if result.Error != nil {
		c.AbortWithStatus(500)
	}

	c.Status(202)
}

func DeleteOneJob(c *gin.Context) {
	id := c.Param("id")

	var job models.Job

	err := initializers.Db.First(&job, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": "Record not found"})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	err = requireOwner(c, job.UserID)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	result := initializers.Db.Delete(&models.Job{}, id)
	if result.Error != nil {
		c.AbortWithStatus(500)
	}

	// response
	c.Status(202)
}

func GetTaskScript(c *gin.Context) {
	jobId := c.Param("id")

	var job models.Job
	err := initializers.Db.First(&job, jobId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	c.File(job.TaskPath)
}

func UploadTaskScript(c *gin.Context) {
	// NOTE: Use a fake blob storage just to upload on local FS, instead of cloud blob storage.
	jobId := c.Param("id")
	var job models.Job

	err := initializers.Db.First(&job, jobId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": "Record not found"})
			return
		}
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	err = requireOwner(c, job.UserID)
	if err != nil {
		c.AbortWithStatus(403)
		return
	}

	file, _ := c.FormFile("file")

	directory, _ := os.Getwd()
	relativeFilePath := fmt.Sprintf("/blob/%s/%s", jobId, file.Filename)
	filePath := fmt.Sprintf("%s/%s", directory, relativeFilePath)
	// fmt.Println(filePath)
	// Upload file
	c.SaveUploadedFile(file, filePath)

	initializers.Db.First(&job, jobId)

	initializers.Db.Model(&job).Updates(map[string]interface{}{
		"TaskPath": filePath,
	})

	c.JSON(http.StatusOK, gin.H{"filepath": relativeFilePath})
}

func requireOwner(c *gin.Context, ownerId uint) error {
	requestUser, _ := c.Get("user")
	if requestUser.(models.User).ID != ownerId && requestUser.(models.User).Role != "admin" {
		return errors.New("forbidden")
	}
	return nil
}
