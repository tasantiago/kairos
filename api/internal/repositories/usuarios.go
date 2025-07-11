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

func (repositorio Usuarios) BuscarPorID(ID uuid.UUID) (models.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, email, setor_id, tipo, criado_em FROM usuarios WHERE id = $1",
		ID,
	)
	if erro != nil {
		return models.Usuario{}, erro
	}

	defer linhas.Close()

	var usuario models.Usuario
	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Email,
			&usuario.SetorID,
			&usuario.Tipo,
			&usuario.CriadoEm,
		); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return usuario, nil
}

func (repositorio Usuarios) Atualizar(ID uuid.UUID, usuario models.Usuario) error {
	statement, erro := repositorio.db.Prepare("UPDATE usuarios SET nome = $1, email = $2, setor_id = $3, tipo = $4 WHERE id = $5")
	if erro != nil {
		return erro
	}

	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Email, usuario.SetorID, usuario.Tipo, ID); erro != nil {
		return erro
	}
	return nil
}

func (repositorio Usuarios) Deletar(ID uuid.UUID) error {
	statemente, erro := repositorio.db.Prepare("DELETE FROM usuarios WHERE id = $1")
	if erro != nil {
		return erro
	}
	defer statemente.Close()

	if _, erro = statemente.Exec(ID); erro != nil {
		return erro
	}

	return nil

}

func (repositorio Usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha FROM usuarios WHERE email = $1", email)
	if erro != nil {
		return models.Usuario{}, erro
	}

	defer linha.Close()

	var usuario models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return usuario, nil
}
