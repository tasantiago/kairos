package models

import "github.com/google/uuid"

type Setor struct {
	ID              uuid.UUID `json:"id"`
	Nome            string    `json:"nome"`
	CapacidadeIdeal int       `json:"capacidade_ideal"`
	CriadoEm        string    `json:"criado_em"`
}
