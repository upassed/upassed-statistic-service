package auth

import (
	"context"
	"errors"
	"github.com/upassed/upassed-statistic-service/internal/handling"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log/slog"
)

var grpcAuthenticationRules = map[string]tokenAuthFunc{}

func (wrapper *clientImpl) AuthenticationUnaryServerInterceptor() func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// TODO: fill grpcAuthenticationRules

	return func(ctx context.Context, request any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response any, err error) {
		log := logging.Wrap(
			wrapper.log,
			logging.WithOp(wrapper.AuthenticationUnaryServerInterceptor),
			logging.WithCtx(ctx),
		)

		authenticationFunc, ok := grpcAuthenticationRules[info.FullMethod]
		if !ok {
			authenticationFunc = wrapper.AnyAccountTypeAuthenticationFunc
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("unable to extract metadata from incoming context")
			return nil, handling.Wrap(errors.New("unable to extract metadata"), handling.WithCode(codes.Internal))
		}

		token, ok := md[AuthenticationHeaderKey]
		if !ok || len(token) != 1 {
			log.Error("missing authentication header in request metadata")
			return nil, handling.Wrap(errors.New("unable to extract authentication header with jwt token"), handling.WithCode(codes.Unauthenticated))
		}

		enrichedCtx, err := authenticationFunc(ctx, token[0])
		if err != nil {
			log.Error("authentication failed", slog.String("err", err.Error()))
			return nil, err
		}

		return handler(enrichedCtx, request)
	}
}
