package logic

import (
	"context"

	"github.com/satyarth42/chatter/auth"
)

func Logout(ctx context.Context, token string) error {
	err := auth.InvalidateToken(token)
	return err
}
