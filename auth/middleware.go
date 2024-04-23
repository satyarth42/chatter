package auth

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/satyarth42/chatter/auth"
)

const (
	HeaderUserID  = "X-User-ID"
	Authorisation = "Authorisation"
)

func JWTVerify() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(Authorisation)
			tokenArr := strings.Split(token, " ")
			if len(tokenArr) != 2 {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			userID, isValid := auth.VerifyToken(tokenArr[1], tokenArr[0])
			if !isValid {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			r.Header.Set(HeaderUserID, userID)
			h.ServeHTTP(w, r)
		})
	}
}
