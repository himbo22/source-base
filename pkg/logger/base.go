package logger

import (
	"context"
	"github.com/himbo22/source-base/pkg/constraints"
	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func SetGlobalLogger(l *zap.Logger) {
	globalLogger = l
}

func WithContext(ctx context.Context) *zap.Logger {
	if globalLogger == nil {
		return zap.NewNop()
	}

	reqID, ok := ctx.Value(constraints.RequestIDKey).(string)
	if ok && reqID != "" {
		return globalLogger.With(zap.String("request_id", reqID))
	}

	return globalLogger
}
