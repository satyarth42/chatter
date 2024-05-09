package logic

import (
	"context"
	"os"
	"time"

	"github.com/satyarth42/chatter/storage"
)

const userServerTTL = time.Hour

func RegisterUserWithServer(ctx context.Context, userID string) error {
	client := storage.GetRedis()

	_, err := client.Set(ctx, userID, os.Getenv("SERVER_ID"), userServerTTL).Result()
	return err
}

func DeregisterUserWithServer(ctx context.Context, userID string) error {
	client := storage.GetRedis()

	_, err := client.Del(ctx, userID).Result()
	return err
}
