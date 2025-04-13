package initializers

import (
	"jwt/internal/helpers"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	
	dsn := os.Getenv("DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		helpers.Logger.Println(err)	
	}
}
