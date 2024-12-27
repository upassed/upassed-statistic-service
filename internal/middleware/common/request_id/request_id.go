package requestid

import "context"

type contextKey string

const ContextKey = contextKey("requestID")

func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(ContextKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
