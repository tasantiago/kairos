package controllers

import (
	"api/internal/models"
	"api/internal/repositories"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
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

func Login(w http.ResponseWriter, r *http.Request) {
	var cred struct {
		Email string `json:"email"`
		Senha string `json:"senha"`
	}

	json.NewDecoder(r.Body).Decode(&cred)

	user, err := repositories.BuscarPorEmail(cred.Email)
	if err != nil {
		http.Error(w, "Usuario não encontrado", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.SenhaHash), []byte(cred.Senha))

	if err != nil {
		http.Error(w, "Senha inválida", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID.String(),
		"tipo": user.Tipo,
	})

	secret := []byte("chave_secreta_mesmo_depois_criar")
	tokenString, _ := token.SignedString(secret)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
