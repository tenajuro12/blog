package main

import (
	"blog-app/config"
	"blog-app/models"
	"blog-app/routes"
	"blog-app/utils"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	utils.InitJWT(cfg.JWTSecret)
	models.InitDB(cfg.DatabaseURL)

	router := routes.NewRouter()

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
