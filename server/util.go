package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/satyarth42/chatter/storage"
)

const (
	serverKeyTTL = 24 * time.Hour
	serverPrefix = "server"
)

func serverKey() string {
	return fmt.Sprintf("%s_%s", serverPrefix, os.Getenv(serverID))
}

func registerServer() {
	ctx := context.Background()
	redisClient := storage.GetRedis()

	_, err := redisClient.Set(ctx, serverKey(), 1, serverKeyTTL).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to register server, err: %s", err.Error()))
	}
}
