package tracing

import (
	"context"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"log/slog"
	"net"
	"os"
)

func InitTracer(cfg *config.Config, log *slog.Logger) (func(), error) {
	log = logging.Wrap(log,
		logging.WithOp(InitTracer),
	)

	ctx := context.Background()
	log.Info("creating a new instance of trace exporter")
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(net.JoinHostPort(cfg.Tracing.Host, cfg.Tracing.Port)),
	)

	if err != nil {
		log.Error("error while creating new tracing exporter", logging.Error(err))
		return nil, err
	}

	log.Info("creating new resource")
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ApplicationName),
			semconv.DeploymentEnvironmentKey.String(string(cfg.Env)),
		),
	)

	if err != nil {
		log.Error("error while creating a resource", logging.Error(err))
		return nil, err
	}

	log.Info("creating a new trace provider")
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	log.Info("trace provider successfully created and initialized")
	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("unable to shutdown tracing provider", logging.Error(err))
			os.Exit(1)
		}
	}, nil
}
