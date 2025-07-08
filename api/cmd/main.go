package main

import (
	"api/internal/config"
	"api/internal/router"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Carregar()

	r := router.Gerar()
	port := os.Getenv("PORT")
	log.Println("Servidor rodando em http://localhost:" + port)
	http.ListenAndServe(":"+port, r)
}
