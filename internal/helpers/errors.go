package helpers

import (
	"jwt/internal/jsonView"
	"net/http"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message map[string]any) {
	err := jsonView.WriteJSON(w, status, message)
	if err != nil {
		Logger.Println(err)
		w.WriteHeader(500)
		return
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	Logger.Println(err)
	
	message := "the server encountered a problem processing your request"
	CustomErrorResponse(w, r, err, message)
}

func CustomErrorResponse(w http.ResponseWriter, r *http.Request, err error, message string) {
	Logger.Println(err)
	errorResponse(w, r, http.StatusInternalServerError, map[string]any{"message":message})
}