package handlers

import (
	"net/http"

	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/dto"
	"github.com/satyarth42/chatter/logic"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.Header.Get(auth.Authorisation)
	if token == "" {
		handleError(w, &dto.CommonError{StatusCode: http.StatusBadRequest})
	}
	err := logic.Logout(ctx, token)
	if err != nil {
		handleError(w, &dto.CommonError{Err: err, StatusCode: http.StatusInternalServerError})
	}
	w.WriteHeader(http.StatusNoContent)
}
