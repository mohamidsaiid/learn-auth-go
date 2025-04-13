package initializers

import (
	"jwt/internal/models"
)

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}