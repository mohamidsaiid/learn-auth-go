package handlers

import (
	"encoding/json"
	"jwt/internal/helpers"
	"jwt/internal/initializers"
	"jwt/internal/jsonView"
	"jwt/internal/models"
	"net/http"
	"strings"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&input)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	res := initializers.DB.Create(&models.User{Email: input.Email, Password: string(hashed)})
	if res.Error != nil {
		helpers.CustomErrorResponse(w, r, res.Error, "this user is already in db")
		return
	}

	err = jsonView.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "user has been created successfuly",
	})
	if err != nil {
		helpers.Logger.Println(err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := jsonView.ReadJSON(w, r, &input)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", input.Email)
	if user.ID == 0 {
		helpers.ServerErrorResponse(w, r, errors.New("didn't find the record"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	token, err := helpers.CreateToken(user)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonView.WriteJSON(w, http.StatusOK, map[string]any{"token": token})

	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}

func TestToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		helpers.CustomErrorResponse(w, r, errors.New("didn't provide a token"), "Missing authorization header")
		return
	}

	tokenString, _ := strings.CutPrefix(token, "Bearer ")

	t, err := helpers.ValidateToken(tokenString)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}	

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		helpers.ServerErrorResponse(w, r, errors.New("not valid token"))
		return
	}

	err = jsonView.WriteJSON(w, http.StatusOK, map[string]any{"Email":claims["email"]})
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}
