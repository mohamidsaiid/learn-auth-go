package routes

import (
	"jwt/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Routes() http.Handler{
	router := httprouter.New()

	router.HandlerFunc(http.MethodPost, "/signin", handlers.Signin)
	router.HandlerFunc(http.MethodPost, "/login", handlers.Login)
	router.HandlerFunc(http.MethodGet, "/token", handlers.TestToken)
	return router
}