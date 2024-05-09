package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/satyarth42/chatter/config"
	"github.com/satyarth42/chatter/storage"
)

const serverID = "SERVER_ID"

func Serve() {

	config.LoadConfig()

	conf := config.GetConfig()

	os.Setenv(serverID, uuid.New().String())

	storage.InitDB(conf.DB)
	storage.InitRedis(conf.Redis)

	registerServer()

	router := mux.NewRouter()

	createRoutes(router)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Handler: router,
	}

	slog.Info(fmt.Sprintf("Server started at: %s", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
