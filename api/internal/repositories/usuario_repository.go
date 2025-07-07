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

func BuscarPorEmail(email string) (*models.Usuario, error) {
	query := `SELECT id, nome, email, senha_hash, setor_id, tipo, criado_em FROM usuarios WHERE email = $1`

	row := config.DB.QueryRow(query, email)

	var u models.Usuario
	err := row.Scan(&u.ID, &u.Nome, &u.Email, &u.SenhaHash, &u.SetorID, &u.Tipo, &u.CriadoEm)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
