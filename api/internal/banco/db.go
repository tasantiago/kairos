package banco

import (
	"api/internal/config"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {

	db, erro := sql.Open("postgres", config.StringConexaoBanco)

	if erro != nil {
		log.Fatal("Erro ao conectar conexão com o banco:", erro)
	}

	erro = db.Ping()
	if erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil
}
