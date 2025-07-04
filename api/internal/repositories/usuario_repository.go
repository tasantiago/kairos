package repositories

import (
	"api/internal/config"
	"api/internal/models"
)

func CriarUsuario(u models.Usuario) error {
	_, err := config.DB.Exec(`
	INSERT INTO usuarios (id, nome, email, senha_hash, setor_id, tipo)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		u.ID, u.Nome, u.Email, u.SenhaHash, u.SetorID, u.Tipo)

	return err
}
