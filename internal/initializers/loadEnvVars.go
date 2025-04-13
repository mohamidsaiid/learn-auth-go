package initializers

import (
	"jwt/internal/helpers"
	"github.com/joho/godotenv"
)

func LoadEnvVar() {
	err := godotenv.Load()
	if err != nil {
		helpers.Logger.Println(err)
	}
}