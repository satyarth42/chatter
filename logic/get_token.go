package logic

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/satyarth42/chatter/auth"
	"github.com/satyarth42/chatter/dto"
)

func GetToken(ctx context.Context, userID string) (*dto.LoginResp, *dto.CommonError) {
	accessToken, refreshToken, expiresIn, err := auth.GetToken(userID)
	if err != nil {
		slog.WarnContext(ctx, "failed to regenerate token for user: %s", userID)
		return nil, &dto.CommonError{Err: err, StatusCode: http.StatusInternalServerError}
	}
	return &dto.LoginResp{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		ExpiresIn:    uint(expiresIn),
		RefreshToken: refreshToken,
	}, err
}
