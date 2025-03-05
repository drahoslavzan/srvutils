package logging

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	ctxLoggerKey struct{}
)

// Middleware attaches a logger to the request context for each incoming request.
func Middleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ReplaceLogger(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TraceMiddleware adds a unique trace ID (UUID) to the logger for each request,
// which helps in tracing requests across different systems.
func TraceMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ReplaceLogger(r.Context(), logger.With(
				zap.String("requestTraceId", uuid.NewString()),
			))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ReplaceLogger adds the logger to the context.
func ReplaceLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, logger)
}

// MustLoggerFromContext retrieves the logger from the request context.
func MustLoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger)
	if !ok {
		zap.L().Error("logger instance not provided in context")
		return zap.L()
	}

	return logger
}
