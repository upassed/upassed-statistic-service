package handling

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const timeFormat string = "2006-01-02 15:04:05"

type Option func(*wrapOptions)

type wrapOptions struct {
	code codes.Code
	time time.Time
}

type ApplicationError struct {
	message string
	code    codes.Code
	time    time.Time
}

func New(message string, code codes.Code) *ApplicationError {
	return &ApplicationError{
		message: message,
		code:    code,
		time:    time.Now(),
	}
}

func (err *ApplicationError) Error() string {
	return err.message
}

func (err *ApplicationError) Code() codes.Code {
	return err.code
}

func (err *ApplicationError) GRPCStatus() *status.Status {
	return status.New(err.code, err.message)
}

func Process(err error, options ...Option) error {
	var applicationErr *ApplicationError
	if errors.As(err, &applicationErr) {
		convertedErr := status.New(applicationErr.code, applicationErr.message)
		timeInfo := errdetails.DebugInfo{
			Detail: fmt.Sprintf("time: %s", applicationErr.time.Format(timeFormat)),
		}

		convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
		if err != nil {
			return convertedErr.Err()
		}

		return convertedErrWithDetails.Err()
	}

	return Wrap(err, options...)
}

func Wrap(err error, options ...Option) error {
	if isAlreadyWrapped(err) {
		return err
	}

	opts := defaultOptions()
	for _, opt := range options {
		opt(opts)
	}

	convertedErr := status.New(opts.code, err.Error())
	timeInfo := errdetails.DebugInfo{
		Detail: fmt.Sprintf("time: %s", opts.time.Format(timeFormat)),
	}

	convertedErrWithDetails, err := convertedErr.WithDetails(&timeInfo)
	if err != nil {
		return convertedErr.Err()
	}

	return convertedErrWithDetails.Err()
}

func isAlreadyWrapped(err error) bool {
	st, ok := status.FromError(err)
	return ok && st.Code() != codes.OK
}

func defaultOptions() *wrapOptions {
	return &wrapOptions{
		code: codes.Internal,
		time: time.Now(),
	}
}

func WithCode(code codes.Code) Option {
	return func(opts *wrapOptions) {
		opts.code = code
	}
}
