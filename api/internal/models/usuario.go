package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	SetorID  uuid.UUID `json:"setor_id,omitempty"`
	Tipo     string    `json:"tipo,omitempty"`
	CriadoEm time.Time `json:"criado_em,omitempty"`
}

func (usuario *Usuario) Validar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Email = strings.TrimSpace(usuario.Email)
	usuario.Senha = strings.TrimSpace(usuario.Senha)
	usuario.Tipo = strings.TrimSpace(usuario.Tipo)

	if usuario.Nome == "" {
		return errors.New("nome não pode estar vazio")
	}
	if usuario.Email == "" {
		return errors.New("email não pode estar vazio")
	}
	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("senha não pode estar vazia")
	}
	if usuario.SetorID == uuid.Nil {
		return errors.New("setor_id não pode ser nulo")
	}
	if usuario.Tipo == "" {
		return errors.New("tipo não pode estar vazio")
	}

	return nil
}
