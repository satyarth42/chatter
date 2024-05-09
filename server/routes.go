package server

import (
	"net/http"
	"path"

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

func Home(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./index.html")
	// set header
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func createRoutes(router *mux.Router) {

	router.HandleFunc("/", Home)
	router.HandleFunc("/health_check", HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/signup", handlers.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/login", handlers.Login).Methods(http.MethodPost)

	authRouter := router.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return !noAuthPaths[r.URL.Path]
	}).Subrouter()

	authRouter.HandleFunc("/logout", handlers.Logout).Methods(http.MethodPost)
	authRouter.HandleFunc("/token", handlers.GetToken).Methods(http.MethodPost)
	router.HandleFunc("/connect", handlers.Connect)

	authRouter.Use(auth.JWTVerify())

}
