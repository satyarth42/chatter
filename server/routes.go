package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/server/handlers"
)

var noAuthPaths = map[string]bool{
	"health_check": true,
	"signup":       true,
	"login":        true,
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("health check is awesome!"))
}

func createRoutes(router *mux.Router) {

	noAuthRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return noAuthPaths[r.URL.Path]
	}).Subrouter()

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return !noAuthPaths[r.URL.Path]
	}).Subrouter()

	noAuthRouter.HandleFunc("/health_check", HealthCheck).Methods(http.MethodGet)
	noAuthRouter.HandleFunc("/signup", handlers.SignUp).Methods(http.MethodPost)
	noAuthRouter.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)

	authRouter.HandleFunc("/logout", handlers.Logout).Methods(http.MethodPost)
	authRouter.HandleFunc("/token", handlers.GetToken).Methods(http.MethodPost)

	authRouter.Use(auth.JWTVerify())

}
