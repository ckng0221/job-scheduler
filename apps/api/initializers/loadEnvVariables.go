package initializers

import (
	"job-scheduler/utils"
)

func LoadEnvVariables() {
	requiredEnvs := []string{"DB_URL", "JWT_SECRET", "GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "LOGIN_REDIRECT_URL", "ADMIN_API_KEY"}
	utils.LoadEnv(requiredEnvs)
}
