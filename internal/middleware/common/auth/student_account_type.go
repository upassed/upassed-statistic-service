package auth

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"github.com/upassed/upassed-statistic-service/internal/handling"
	"github.com/upassed/upassed-statistic-service/internal/logging"
	"google.golang.org/grpc/codes"
	"log/slog"
)

func (wrapper *clientImpl) StudentAccountTypeAuthenticationFunc(ctx context.Context, token string) (context.Context, error) {
	log := logging.Wrap(
		wrapper.log,
		logging.WithOp(wrapper.StudentAccountTypeAuthenticationFunc),
		logging.WithCtx(ctx),
	)

	timeout := wrapper.cfg.GetEndpointExecutionTimeout()
	callCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := wrapper.authenticationServiceClient.Validate(callCtx, &client.TokenValidateRequest{
		AccessToken: token,
	})

	if err != nil {
		log.Error("error while validating token on an authentication service", slog.String("err", err.Error()))
		return nil, handling.Wrap(errors.New("validate token error"), handling.WithCode(codes.Unauthenticated))
	}

	enrichedContext := context.WithValue(ctx, UsernameKey, response.GetUsername())
	if !(response.GetAccountType() == "STUDENT") {
		log.Error("account type is not equal to student", slog.String("accountType", response.GetAccountType()))
		return nil, handling.Wrap(errors.New("required student account type"), handling.WithCode(codes.PermissionDenied))
	}

	return enrichedContext, nil
}
