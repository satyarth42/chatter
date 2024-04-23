package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satyarth42/chatter/config"
)

var client *redis.ClusterClient

func InitRedis(conf config.RedisConfig) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           conf.Addresses,
		ReadOnly:        conf.ReadOnly,
		DialTimeout:     time.Duration(conf.DialTimeoutInSec) * time.Second,
		ReadTimeout:     time.Duration(conf.ReadTimeoutInSec) * time.Second,
		WriteTimeout:    time.Duration(conf.WriteTimeoutInSec) * time.Second,
		PoolSize:        conf.PoolSize,
		ConnMaxIdleTime: time.Duration(conf.ConnMaxIdleTimeInSec) * time.Second,
	})

	cmd := redisClient.Ping(context.Background())
	if cmd.Err() != nil {
		panic(fmt.Sprintf("failed to connect to redis for addresses: %+v", conf.Addresses))
	}

	client = redisClient
}

func GetRedis() *redis.ClusterClient {
	return client
}
