package main

import (
	"github.com/joho/godotenv"
	"github.com/upassed/upassed-statistic-service/internal/app"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"github.com/upassed/upassed-statistic-service/internal/tracing"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("error while loading env variables:", err.Error())
	}

	configFilePath := os.Getenv(config.EnvConfigPath)
	if configFilePath == "" {
		log.Println("using local config file path")

		configFilePath = filepath.Join("config", "local.yml")
		if err := os.Setenv(config.EnvConfigPath, configFilePath); err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.Wrap(logging.New(cfg.Env), logging.WithOp(main))
	logger.Info("logger successfully initialized", slog.Any("env", cfg.Env))

	traceProviderShutdownFunc, err := tracing.InitTracer(cfg, logger)
	if err != nil {
		logger.Error("unable to initialize traceProvider", logging.Error(err))
		os.Exit(1)
	}

	defer traceProviderShutdownFunc()
	logger.Info("trace provider successfully initialized")

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("error occurred while creating an app", logging.Error(err))
		os.Exit(1)
	}

	go func(app *app.App) {
		if err := app.Server.Run(); err != nil {
			logger.Error("error occurred while running a gRPC server", logging.Error(err))
			os.Exit(1)
		}
	}(application)

	stopSignalChannel := make(chan os.Signal, 1)
	signal.Notify(stopSignalChannel, syscall.SIGTERM, syscall.SIGINT)
	<-stopSignalChannel

	application.Server.GracefulStop()
	logger.Info("server gracefully stopped")
}
