package middleware

import (
	"os"
	"slices"
	"strings"

	"job-scheduler/utils"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func SecureMiddleware() gin.HandlerFunc {
	allowedHosts := utils.Getenv("SERVER_ALLOWED_HOSTS", "localhost:8000")
	allowedHostsSlice := strings.Split(allowedHosts, ",")

	// Disable security on development mode
	isDevelopment := slices.Contains([]string{"development", "test"}, os.Getenv("ENV"))

	// NOTE: uncommnet SSL config when haivng https site
	return secure.New(secure.Config{
		AllowedHosts: allowedHostsSlice,
		// SSLRedirect:  true,
		// SSLHost:               "localhost:8000",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		IsDevelopment:         isDevelopment,
	})
}
