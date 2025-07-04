package models

import "github.com/google/uuid"

type Tarefa struct {
	ID            uuid.UUID `json:"id"`
	Titulo        string    `json:"titulo"`
	Descricao     string    `json:"descricao"`
	Origem        string    `json:"origem"`
	Prioridade    string    `json:"prioridade"`
	Status        string    `json:"status"`
	SetorAtualID  uuid.UUID `json:"setor_atual_id"`
	ResponsavelID uuid.UUID `json:"responsavel_id,omitempty"`
	CriadoEm      string    `json:"criado_em"`
	AtualizadoEm  string    `json:"atualizado_em"`
}
