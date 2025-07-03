package main

import (
	"api/internal/config"
	"api/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	config.Connect()

	router := routes.LoadRoutes()
	port := os.Getenv("PORT")
	log.Println("Servidor rodando em http://localhost:" + port)
	http.ListenAndServe(":"+port, router)
}
