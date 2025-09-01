package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/nico-phil/process/config"
	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client = nil
)

func InitRedis() error {
	addr := config.GetRedisArr()
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.GetRedisPasswrod(),
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("failed to connect t redis: %v", err)
		return fmt.Errorf("failed to connect to redis %q", err)
	}

	log.Printf("connected to redis: %s", pong)

	return nil

}
