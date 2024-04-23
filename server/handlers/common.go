package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/satyarth42/chatter/dto"
)

func handleError(w http.ResponseWriter, err *dto.CommonError) {
	if err == nil {
		return
	}
	w.WriteHeader(err.StatusCode)
	fmt.Fprintf(w, "error returned err: %+v", err.Err)
}

func handleResponse(ctx context.Context, w http.ResponseWriter, resp any, statusCode int) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		slog.WarnContext(ctx, "failed to marshal resp", err)
		handleError(w, &dto.CommonError{Err: err, StatusCode: http.StatusInternalServerError})
		return
	}
	w.WriteHeader(statusCode)
	w.Write(respBytes)
}

func readBody(req *http.Request, str any) *dto.CommonError {
	if reflect.ValueOf(str).Kind() != reflect.Ptr {
		return &dto.CommonError{StatusCode: http.StatusInternalServerError}
	}

	body, bodyErr := io.ReadAll(req.Body)
	if bodyErr != nil {
		slog.WarnContext(req.Context(), fmt.Sprintf("could not read body for path: %s, err%+v", req.URL.Path, bodyErr))
		return &dto.CommonError{Err: bodyErr, StatusCode: http.StatusBadRequest}
	}

	unmarshalErr := json.Unmarshal(body, str)
	if unmarshalErr != nil {
		slog.WarnContext(req.Context(), fmt.Sprintf("failed to unmarshal body for path: %s", req.URL.Path), unmarshalErr)
		return &dto.CommonError{Err: unmarshalErr, StatusCode: http.StatusBadRequest}
	}

	return nil
}
