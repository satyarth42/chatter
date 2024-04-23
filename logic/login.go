package logic

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/dto"
	"github.com/satyarth42/chatter/models"
	"github.com/satyarth42/chatter/storage"
	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, req *dto.LoginReq) (*dto.LoginResp, *dto.CommonError) {
	db := storage.GetUserDataDB()

	condition := map[string]interface{}{
		"email":      req.Email,
		"deleted_at": nil,
	}

	user := &models.User{}

	result := db.Where(condition).Find(user)
	if result.Error != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to fetch user: %s from db, err: %+v", req.Email, result.Error))
		return nil, &dto.CommonError{Err: fmt.Errorf("something went wrong"), StatusCode: http.StatusInternalServerError}
	}

	if user.UUID == "" {
		slog.WarnContext(ctx, fmt.Sprintf("user: %s does not exists", req.Email))
		return nil, &dto.CommonError{Err: fmt.Errorf("wrong credentials"), StatusCode: http.StatusBadRequest}
	}

	isValid := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if isValid != nil {
		slog.WarnContext(ctx, fmt.Sprintf("wrong password for user: %s", req.Email))
		return nil, &dto.CommonError{Err: fmt.Errorf("wrong credentials"), StatusCode: http.StatusBadRequest}
	}

	accessToken, refreshToken, expiresIn, err := auth.GetToken(user.UUID)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("error in generating JWT for user: %s", req.Email))
		return nil, &dto.CommonError{Err: fmt.Errorf("internal error"), StatusCode: http.StatusInternalServerError}
	}

	return &dto.LoginResp{
		TokenType:    "Bearer",
		ExpiresIn:    uint(expiresIn),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
