package repositories

import (
	"api/internal/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Usuarios struct {
	db *sql.DB
}

func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

func (repositorio Usuarios) Criar(usuario models.Usuario) (uuid.UUID, error) {
	ID := uuid.New()
	statement, erro := repositorio.db.Prepare("INSERT INTO usuarios (id, nome, email, senha, setor_id, tipo) VALUES ($1, $2, $3, $4, $5, $6)")
	if erro != nil {
		return uuid.Nil, erro
	}
	defer statement.Close()

	_, erro = statement.Exec(ID, usuario.Nome, usuario.Email, usuario.Senha, usuario.SetorID, usuario.Tipo)
	if erro != nil {
		return uuid.Nil, erro
	}

	return ID, nil
}

func (repositorio Usuarios) Buscar(nomeOuEmail string) ([]models.Usuario, error) {
	nomeOuEmail = fmt.Sprintf("%%%s%%", nomeOuEmail)

	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, email, setor_id, tipo, criado_em from usuarios WHERE nome LIKE $1 or email LIKE $1", nomeOuEmail,
	)
	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var usuario models.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Email,
			&usuario.SetorID,
			&usuario.Tipo,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}

		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}
