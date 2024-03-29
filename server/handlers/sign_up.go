package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/satyarth42/chatter/dto"
	"github.com/satyarth42/chatter/logic"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	req := &dto.SignUpReq{}
	bodyErr := readBody(r, req)
	if bodyErr != nil {
		slog.WarnContext(ctx, fmt.Sprintf("error in reading body, err: %+v", bodyErr.Error()))
		handleError(w, bodyErr)
		return
	}

	err := logic.SignUp(ctx, req)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to signup email: %s, err: %+v", req.Name, err.Error()))
		handleError(w, &dto.CommonError{Err: fmt.Errorf("failure"), StatusCode: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusCreated)

}
