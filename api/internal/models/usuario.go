package models

import "github.com/google/uuid"

type Usuario struct {
	ID        uuid.UUID `json:"id"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	SenhaHash string    `json:"-"`
	SetorID   uuid.UUID `json:"setor_id"`
	Tipo      string    `json:"tipo"`
	CriadoEm  string    `json:"criado_em"`
}
