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
	"time"

	"github.com/joho/godotenv"
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

// Get string from environmnet variable.
// If empty, assign the value with the fallback string.
func Getenv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func ConvertStructToMap(obj interface{}) (map[string]interface{}, error) {
	var objInterface map[string]interface{}
	objJson, err := json.Marshal(obj)

	json.Unmarshal(objJson, &objInterface)
	return objInterface, err
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
