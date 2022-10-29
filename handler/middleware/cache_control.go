package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CacheControlMiddleware struct {
	MaxAge time.Duration
	Public bool
}

type CacheControlOption func(*CacheControlMiddleware)

func NewCacheControlMiddleware(opts ...CacheControlOption) *CacheControlMiddleware {
	c := &CacheControlMiddleware{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *CacheControlMiddleware) CacheControl() gin.HandlerFunc {
	value := ""
	if c.Public {
		value = "public,"
	}
	if c.MaxAge > 0 {
		value = "max-age=" + strconv.FormatInt(int64(c.MaxAge.Seconds()), 10)
	}
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", value)
	}
}

func WithMaxAge(maxAge time.Duration) CacheControlOption {
	return func(c *CacheControlMiddleware) {
		c.MaxAge = maxAge
	}
}

func WithPublic(public bool) CacheControlOption {
	return func(c *CacheControlMiddleware) {
		c.Public = public
	}
}
