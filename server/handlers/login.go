package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/satyarth42/chatter/dto"
	"github.com/satyarth42/chatter/logic"
)

func Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &dto.LoginReq{}
	bodyErr := readBody(r, req)
	if bodyErr != nil {
		slog.WarnContext(ctx, "error in reading body", bodyErr.Error())
		handleError(w, bodyErr)
		return
	}

	resp, err := logic.Login(ctx, req)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("error in login for email:%s", req.Email), err.Error())
		handleError(w, err)
	}

	handleResponse(ctx, w, resp, http.StatusOK)
}
