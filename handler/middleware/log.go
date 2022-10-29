package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
func (g *GinLoggerMiddleware) LoggerWithConfig(conf GinLoggerConfig) gin.HandlerFunc {
	logger := g.logger.WithOptions(zap.WithCaller(false))
	notLogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		if len(ctx.Errors) > 0 {
			logger.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {

			if raw != "" {
				path = path + "?" + raw
			}
			path = strings.Replace(path, "\n", "", -1)
			path = strings.Replace(path, "\r", "", -1)
			clientIP := strings.Replace(ctx.ClientIP(), "\n", "", -1)
			clientIP = strings.Replace(clientIP, "\r", "", -1)

			logger.Info("[GIN]",
				zap.Time("beginTime", start),
				zap.Int("status", ctx.Writer.Status()),
				zap.Duration("latency", time.Since(start)),
				zap.String("clientIP", clientIP),
				zap.String("method", ctx.Request.Method),
				zap.String("path", path))
		}
	}
}
