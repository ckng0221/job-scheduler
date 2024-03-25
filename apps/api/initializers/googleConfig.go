package initializers

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		RedirectURL:  os.Getenv("LOGIN_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}
