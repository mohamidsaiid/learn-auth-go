package routes

import (
	"jwt/internal/handlers"
	"net/http"
)

func Routes() http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("/signin", handlers.Signin)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/token", handlers.TestToken)
	return mux
}