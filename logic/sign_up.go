package logic

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/satyarth42/chatter/dto"
	"github.com/satyarth42/chatter/models"
	"github.com/satyarth42/chatter/storage"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx context.Context, req *dto.SignUpReq) *dto.CommonError {
	db := storage.GetUserDataDB()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to hash password for user: %s", req.Email), err)
	}

	user := &models.User{
		UUID:               uuid.NewString(),
		Name:               req.Name,
		Email:              req.Email,
		Password:           string(hashedPassword),
		VerificationStatus: false,
	}

	result := db.Create(user)
	if result.Error != nil {
		slog.WarnContext(ctx, fmt.Sprintf("error in creating user: %s", user.Email), result.Error)
		return &dto.CommonError{Err: result.Error, StatusCode: http.StatusInternalServerError}
	}

	return nil
}
