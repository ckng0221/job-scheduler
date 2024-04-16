package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func SecureMiddleware() gin.HandlerFunc {
	allowedHosts := os.Getenv("SERVER_ALLOWED_HOSTS")
	if allowedHosts == "" {
		allowedHosts = "localhost:8000"
	}
	allowedHostsSlice := strings.Split(allowedHosts, ",")
	// fmt.Println("host", allowedHostsSlice)

	return secure.New(secure.Config{
		AllowedHosts: allowedHostsSlice,
		// SSLRedirect: true,
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
	})
}
