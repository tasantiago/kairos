package routes

import (
	"api/internal/controllers"

	"github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	router := mux.NewRouter()

	// router.Handle("/usuarios", middlewares.AuthMiddleware(http.HandlerFunc(controllers.CriarUsuario))).Methods("POST")
	router.HandleFunc("/usuarios", controllers.CriarUsuario).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	return router
}
