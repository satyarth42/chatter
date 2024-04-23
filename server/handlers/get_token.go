package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/logic"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.Header.Get(auth.HeaderUserID)

	resp, err := logic.GetToken(ctx, userID)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("error in login for email:%s, err: %+v", req.Email, err.Error()))
		handleError(w, err)
		return
	}

	handleResponse(ctx, w, resp, http.StatusOK)
}
