package handlers

import (
	"encoding/json"
	"jwt/internal/helpers"
	"jwt/internal/initializers"
	"jwt/internal/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helpers.Logger.Println("only post requests")
		w.Write([]byte("only post requests"))
		return
	}

	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&user)
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	res := initializers.DB.Create(&models.User{Email: user.Email, Password: string(hashed)})
	if res.Error != nil {
		helpers.Logger.Println(res.Error)
		w.Write([]byte("this user is already in db"))
		return
	}
	w.Write([]byte("user has been created successfuly"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helpers.Logger.Println("only post requests")
		w.Write([]byte("only post requests"))
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", input.Email)
	if user.ID == 0 {
		helpers.Logger.Println("doesn't find the record")
		w.Write([]byte("internal server error"))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"email": user.Email,
		"admin": true,
	})

	TokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	tkn := struct {
		Token string `json:"token"`
	}{
		Token: TokenString,
	}

	jsToken, err := json.Marshal(tkn)
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsToken)
}

func TestToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helpers.Logger.Println("only post requests")
		w.Write([]byte("only post requests"))
		return
	}
	token := r.Header.Get("Authorization")
	if token == "" {
		helpers.Logger.Println("didn't provide a token")
		w.Write([]byte("Missing authorization header"))
		return
	}


	tokenString, _ := strings.CutPrefix(token, "Bearer ")

	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("Missing authorization header"))
		return
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		helpers.Logger.Println("not valid token")
		w.Write([]byte("internal server error"))
		return
	}

	output := struct {
		Email string
	}{
		Email: claims["email"].(string),
	}

	js, err := json.Marshal(output)
	if err != nil {
		helpers.Logger.Println(err)
		w.Write([]byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
