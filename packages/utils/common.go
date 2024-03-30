package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func LoadEnv(requiredEnv []string) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load .env file")
	}
	for _, envName := range requiredEnv {
		env := os.Getenv(envName)
		if env == "" {
			log.Fatalf("environment variable '%s' is required", envName)
		}
	}
}

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

func GetUnixMinuteRange(t time.Time) (time.Time, time.Time) {
	currentMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	nextMinute := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location()).Add(1 * time.Minute)
	return currentMinute, nextMinute
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

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
