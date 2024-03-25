package main

import (
	"job-scheduler/api/initializers"
	"job-scheduler/api/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SynDatabase()
}

func main() {

	r := routes.SetupRouter()
	r.Run()
}

// CompileDaemon -command="./api"

// generate swagger docs
// swag init --parseDependency --parseInternal
