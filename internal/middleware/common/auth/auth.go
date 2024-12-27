package auth

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"github.com/upassed/upassed-statistic-service/internal/config"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
)

var (
	errCreatingAuthServiceConn = errors.New("unable to create authentication service connection")
)

const (
	AuthenticationHeaderKey = "authentication"
	UsernameKey             = "username"
)

type tokenAuthFunc func(ctx context.Context, token string) (context.Context, error)

type Client interface {
	AuthenticationUnaryServerInterceptor() func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error)
	AnyAccountTypeAuthenticationFunc(ctx context.Context, token string) (context.Context, error)
	StudentAccountTypeAuthenticationFunc(ctx context.Context, token string) (context.Context, error)
	TeacherAccountTypeAuthenticationFunc(ctx context.Context, token string) (context.Context, error)
}

type clientImpl struct {
	cfg                         *config.Config
	log                         *slog.Logger
	authenticationServiceClient client.TokenClient
}

func NewClient(cfg *config.Config, log *slog.Logger) (Client, error) {
	authenticationServiceUrl := net.JoinHostPort(
		cfg.Services.Authentication.Host,
		cfg.Services.Authentication.Port,
	)

	log = logging.Wrap(
		log,
		logging.WithOp(NewClient),
		logging.WithAny("authentication-service-url", authenticationServiceUrl),
	)

	authenticationServiceConnection, err := grpc.NewClient(authenticationServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("unable to create authentication client connection", slog.String("err", err.Error()))
		return nil, errCreatingAuthServiceConn
	}

	return &clientImpl{
		cfg:                         cfg,
		log:                         log,
		authenticationServiceClient: client.NewTokenClient(authenticationServiceConnection),
	}, nil
}
