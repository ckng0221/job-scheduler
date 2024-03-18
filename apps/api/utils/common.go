package utils

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConvertStructToMap(obj interface{}) (map[string]interface{}, error) {
	var objInterface map[string]interface{}
	objJson, err := json.Marshal(obj)

	json.Unmarshal(objJson, &objInterface)
	return objInterface, err
}

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
