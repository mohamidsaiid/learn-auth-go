package models

import (
	"errors"
	"jwt/internal/helpers"
	"jwt/internal/initializers"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func SyncDB() {
	initializers.DB.AutoMigrate(&User{})
}

func CreateUser(email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	_, err = find(email)
	if err == nil {
		return err
	}

	res := initializers.DB.Create(&User{Email: email, Password: string(hashed)})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func FindUser(email, password string) (id uint, err error) {
	user, err := find(email)
	if err != nil {
		return
	}

	id = user.ID
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		helpers.Logger.Println(err)
	}
	return
}

func find(email string) (*User,error) {
	var user User
	initializers.DB.First(&user, "email = ?", email)
	if user.ID == 0 {
		return nil, errors.New("user is already regesitered")
	}

	return &user, nil
}
