package helpers

import (
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id uint, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   id,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"email": email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(r *http.Request) error {
	tokenString := getToken(r)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return err
	}
	return nil
}

func GetClaims(r *http.Request) (jwt.MapClaims, bool) {
	tokenString := getToken(r)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	claims, ok:= token.Claims.(jwt.MapClaims)
	return claims, ok
}

func getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	tokenString, _ := strings.CutPrefix(token, "Bearer ")
	return tokenString
}
