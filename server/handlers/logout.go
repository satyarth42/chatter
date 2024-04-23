package handlers

import (
	"log/slog"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	slog.Info("in logout")

}
