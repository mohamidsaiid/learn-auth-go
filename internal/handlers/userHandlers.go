package handlers

import (
	"errors"
	"jwt/internal/helpers"
	"jwt/internal/jsonView"
	"jwt/internal/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := jsonView.ReadJSON(w, r, &input)
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
		return
	}
	// check if user in database
	// if no create one else write to the request
	err = models.CreateUser(input.Email, input.Password)
	if err != nil {
		helpers.CustomErrorResponse(w, r, err, "user is already regestered")
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

	id, err := models.FindUser(input.Email, input.Password)
	if err != nil {
		helpers.CustomErrorResponse(w, r, err, "email or password isn't correct")
		return
	}

	token, err := helpers.CreateToken(id, input.Email)
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

	err = jsonView.WriteJSON(w, http.StatusOK, map[string]any{"Email": claims["email"]})
	if err != nil {
		helpers.ServerErrorResponse(w, r, err)
	}
}
