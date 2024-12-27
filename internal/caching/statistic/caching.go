package statistic

import (
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"log/slog"
)

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

func New(client *redis.Client, cfg *config.Config, log *slog.Logger) *RedisClient {
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: client,
	}
}
