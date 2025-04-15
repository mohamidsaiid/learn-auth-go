package routes

import (
	"jwt/internal/handlers"
	"jwt/internal/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes() http.Handler{
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/signup", handlers.Signup)
	router.HandlerFunc(http.MethodPost, "/login", handlers.Login)
	router.HandlerFunc(http.MethodGet, "/token", middleware.AuthMiddleware(handlers.TestToken))
	return router
}