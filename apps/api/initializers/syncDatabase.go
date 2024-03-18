package initializers

import "job-scheduler/api/models"

func SynDatabase() {
	Db.AutoMigrate(
		&models.User{},
		&models.Job{},
		&models.Execution{},
	)
}
