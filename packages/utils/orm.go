package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageString := c.Query("page")
		pageSizeString := c.Query("page-size")
		page, _ := strconv.Atoi(pageString)
		pageSize, _ := strconv.Atoi(pageSizeString)

		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize <= 0:
			pageSize = 50
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
