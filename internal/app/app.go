package app

import (
	"github.com/upassed/upassed-statistic-service/internal/caching"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"github.com/upassed/upassed-statistic-service/internal/middleware/common/auth"
	"github.com/upassed/upassed-statistic-service/internal/server"
	"github.com/upassed/upassed-statistic-service/internal/service/statistic"
	"log/slog"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	log = logging.Wrap(log, logging.WithOp(New))

	redis, err := caching.OpenRedisConnection(config, log)
	if err != nil {
		return nil, err
	}

	authClient, err := auth.NewClient(config, log)
	if err != nil {
		return nil, err
	}

	statisticService := statistic.New(config, log, redis)
	appServer := server.New(server.AppServerCreateParams{
		Config:           config,
		Log:              log,
		AuthClient:       authClient,
		StatisticService: statisticService,
	})

	log.Info("app successfully created")
	return &App{
		Server: appServer,
	}, nil
}
