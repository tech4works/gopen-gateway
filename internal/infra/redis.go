package infra

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/go-redis/redis/v8"
)

var clientRedis *redis.Client

func ConnectRedis(url, password string) *redis.Client {
	clientRedis = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0, // use default DB
	})
	return clientRedis
}

func DisconnectRedis() {
	if clientRedis == nil {
		return
	}
	err := clientRedis.Close()
	if helper.IsNotNil(err) {
		logger.Error("Error disconnect Redis:", err)
		return
	}
	logger.Info("Connection to Redis closed.")
}
