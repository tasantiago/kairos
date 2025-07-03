package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API online"))
	}).Methods("Get")

	return router
}
