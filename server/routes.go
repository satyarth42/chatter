package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satyarth42/chatter/server/handlers"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check is awesome!"))
}

func createRoutes(router *mux.Router) {
	router.HandleFunc("health_check", HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/signup", handlers.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)
	router.HandleFunc("/logout", handlers.Logout).Methods(http.MethodPost)
}
