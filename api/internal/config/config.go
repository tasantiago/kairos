package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	StringConexaoBanco = ""
	Porta              = 0
)

func Carregar() {
	var erro error

	Porta, erro = strconv.Atoi(os.Getenv("PORT"))
	if erro != nil {
		Porta = 8080
	}

	StringConexaoBanco = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}
