package controllers

import (
	"api/internal/repositories"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
}
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Criando usuário"))
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

	log.Println("Senha digitada:", cred.Senha)
	log.Println("Hash armazenado:", user.SenhaHash)
	log.Println(err)

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
