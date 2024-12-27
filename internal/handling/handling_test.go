package handling_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-statistic-service/internal/handling"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvertApplicationError_AddDetails(t *testing.T) {
	message := gofakeit.Error().Error()
	code := codes.Internal
	applicationError := handling.New(message, code)

	handledError := handling.Process(applicationError)

	st := status.Convert(handledError)
	assert.Equal(t, code, st.Code())
	assert.Equal(t, message, st.Message())
	assert.Equal(t, 1, len(st.Details()))
}

func TestConvertApplicationError_WrapOptions(t *testing.T) {
	message := "error message"
	code := codes.AlreadyExists

	applicationError := handling.New(message, code)
	wrappedError := handling.Process(applicationError, handling.WithCode(codes.OK))

	st := status.Convert(wrappedError)
	require.Error(t, wrappedError)

	assert.Equal(t, message, st.Message())
	assert.Equal(t, code, st.Code())
}

func TestConvertApplicationError_WrappingNotAnApplicationError(t *testing.T) {
	initialError := gofakeit.Error()

	handledError := handling.Process(initialError)

	st := status.Convert(handledError)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, initialError.Error(), st.Message())
	assert.Equal(t, 1, len(st.Details()))
}

func TestCreateAnApplicationError(t *testing.T) {
	message := "error message"
	code := codes.AlreadyExists

	applicationError := handling.New(message, code)

	require.Error(t, applicationError)

	assert.Equal(t, message, applicationError.Error())
	assert.Equal(t, message, applicationError.GRPCStatus().Message())
	assert.Equal(t, code, applicationError.GRPCStatus().Code())
}
