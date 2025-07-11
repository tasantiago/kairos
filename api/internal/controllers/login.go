package controllers

import (
	"api/internal/banco"
	"api/internal/models"
	"api/internal/repositories"
	"api/internal/services"
	"encoding/json"
	"io"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		services.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		services.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		services.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		services.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if erro = services.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha); erro != nil {
		services.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	w.Write([]byte("Sucesso no login"))
}
