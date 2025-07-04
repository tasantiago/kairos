package routes

import (
	"api/internal/controllers"

	"github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/usuarios", controllers.CriarUsuario).Methods("Get")

	return router
}
