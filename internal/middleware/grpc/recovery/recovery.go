package recovery

import (
	"context"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	requestid "github.com/upassed/upassed-statistic-service/internal/middleware/common/request_id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func MiddlewareInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log := logging.Wrap(log, logging.WithOp(MiddlewareInterceptor))

				log.Error("panic recovered",
					slog.String("requestID", requestid.GetRequestIDFromContext(ctx)),
					slog.Any("message", r),
					slog.String("method", info.FullMethod),
				)

				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
