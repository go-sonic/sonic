package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	hzerrors "github.com/cloudwego/hertz/pkg/common/errors"
	"go.uber.org/zap"
)

type GinLoggerMiddleware struct {
	logger *zap.Logger
}

func NewGinLoggerMiddleware(logger *zap.Logger) *GinLoggerMiddleware {
	return &GinLoggerMiddleware{
		logger: logger,
	}
}

// GinLoggerConfig LoggerConfig defines the config for Logger middleware
type GinLoggerConfig struct {
	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LoggerWithConfig instance a Logger middleware with config.
func (g *GinLoggerMiddleware) LoggerWithConfig(conf GinLoggerConfig) app.HandlerFunc {
	logger := g.logger.WithOptions(zap.WithCaller(false))
	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(_ctx context.Context, ctx *app.RequestContext) {
		// Start timer
		start := time.Now()
		path := string(ctx.URI().Path())
		raw := string(ctx.URI().QueryString())
		_ctx = context.WithValue(_ctx, "clientIP", ctx.ClientIP())
		_ctx = context.WithValue(_ctx, "userAgent", ctx.UserAgent())
		// Process request
		ctx.Next(_ctx)

		if len(ctx.Errors) > 0 {
			logger.Error(ctx.Errors.ByType(hzerrors.ErrorTypePrivate).String())
		}
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			if raw != "" {
				path = path + "?" + raw
			}
			path = strings.ReplaceAll(path, "\n", "")
			path = strings.ReplaceAll(path, "\r", "")
			clientIP := strings.ReplaceAll(ctx.ClientIP(), "\n", "")
			clientIP = strings.ReplaceAll(clientIP, "\r", "")

			logger.Info("[GIN]",
				zap.Time("beginTime", start),
				zap.Int("status", ctx.Response.StatusCode()),
				zap.Duration("latency", time.Since(start)),
				zap.String("clientIP", clientIP),
				zap.String("method", string(ctx.Request.Method())),
				zap.String("path", path))
		}
	}
}
