package ratelimiter

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
)

var ErrRateLimited = errors.New("rate limited")

var Limiter *redis_rate.Limiter

func RateLimit(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := Limiter.Allow(c.Request.Context(), "task", redis_rate.PerMinute(10))
		if err != nil {
			logger.Error("Rate limiter", slog.String("error", err.Error()))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed == 0 {
			// We are rate limited.

			seconds := int(res.RetryAfter / time.Second)
			c.Header("RateLimit-RetryAfter", strconv.Itoa(seconds))

			// Stop processing and return the error.
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		// Continue processing as normal.
		c.Next()
	}
}
