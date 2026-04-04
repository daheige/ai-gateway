package utils

import (
	"context"
	"time"

	"ai-gateway/internal/infras/ctxkeys"
)

// NewContext create a new context from request,eg:http request
func NewContext(ctx context.Context) context.Context {
	requestID, ok := ctx.Value(ctxkeys.XRequestID).(string)
	if requestID == "" || !ok {
		requestID = UUID()
	}

	newCtx := context.WithValue(ctx, ctxkeys.XRequestID, requestID)
	return newCtx
}

// NewTimeoutContext create a new context from request,eg:http request
func NewTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	newCtx, cancel := context.WithTimeout(NewContext(ctx), timeout)
	return newCtx, cancel
}
