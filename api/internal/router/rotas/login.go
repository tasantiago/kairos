package rotas

import "net/http"

var rotaLogin = Rota{
	URI:                "login",
	Metodo:             http.MethodPost,
	Funcao:             func(w http.ResponseWriter, r *http.Request) {},
	RequerAutenticacao: false,
}
