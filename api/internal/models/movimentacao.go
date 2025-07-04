package models

import "github.com/google/uuid"

type Movimentacao struct {
	ID               uuid.UUID  `json:"id"`
	TarefaID         uuid.UUID  `json:"tarefa_id"`
	DeSetorID        *uuid.UUID `json:"de_setor_id,omitempty"`
	ParaSetorID      uuid.UUID  `json:"para_setor_id"`
	ResponsavelID    uuid.UUID  `json:"responsavel_id"`
	DataMovimentacao string     `json:"data_comentacao"`
}
