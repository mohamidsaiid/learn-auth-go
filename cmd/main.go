package main

import (
	"fmt"
	"jwt/internal/helpers"
	"jwt/internal/initializers"
	"jwt/internal/routes"
	"net/http"
	"os"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	srv := http.Server{
		Addr : fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: routes.Routes(),
	}	
	helpers.Logger.Println("server is up on port 3000")
	err := srv.ListenAndServe()
	helpers.Logger.Fatal(err)
}
