package statistic

import (
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-statistic-service/internal/caching/statistic"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"log/slog"
)

type Service interface {
}

type serviceImpl struct {
	cfg   *config.Config
	log   *slog.Logger
	cache *statistic.RedisClient
}

func New(cfg *config.Config, log *slog.Logger, cache *redis.Client) Service {
	statisticCache := statistic.New(cache, cfg, log)
	return &serviceImpl{
		cfg:   cfg,
		log:   log,
		cache: statisticCache,
	}
}
