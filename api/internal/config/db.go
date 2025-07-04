package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar conexão com o banco:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao testar conexão com o banco:", err)
	}

	log.Println("Banco de dados conectado com sucesso!")

	row := DB.QueryRow("SELECT NOW()")
	var now string
	row.Scan(&now)
	log.Println("Hora atual no banco:", now)
}
