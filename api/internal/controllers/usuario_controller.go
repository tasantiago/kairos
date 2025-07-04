package controllers

import (
	"api/internal/models"
	"api/internal/repositories"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	var u models.Usuario
	json.NewDecoder(r.Body).Decode(&u)

	u.ID = uuid.New()
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.SenhaHash), 14)
	u.SenhaHash = string(hash)

	err := repositories.CriarUsuario(u)
	if err != nil {
		http.Error(w, "Erro ao criar usuários", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
