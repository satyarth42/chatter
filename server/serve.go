package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/satyarth42/chatter/config"
	"github.com/satyarth42/chatter/storage"
)

func Serve() {

	config.LoadConfig()

	conf := config.GetConfig()

	storage.InitDB(conf.DB)

	router := mux.NewRouter()

	createRoutes(router)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Handler: router,
	}

	slog.Info(fmt.Sprintf("Server started at: %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
