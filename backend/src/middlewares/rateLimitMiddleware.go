package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/types"
)

// RateLimitMiddleware 每 IP 简单令牌桶：10 次/秒，突发 20。
func RateLimitMiddleware() gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]*rateBucket{}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		bucket, ok := buckets[ip]
		if !ok {
			bucket = &rateBucket{tokens: 20, last: time.Now()}
			buckets[ip] = bucket
		}
		bucket.refill()
		if bucket.tokens < 1 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, types.APIResponse{
				OK:    false,
				Error: types.NewAppError(constants.RateLimited, constants.RateLimitedMessage),
			})
			return
		}
		bucket.tokens--
		mu.Unlock()
		c.Next()
	}
}

type rateBucket struct {
	tokens float64
	last   time.Time
}

func (b *rateBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * 10
	if b.tokens > 20 {
		b.tokens = 20
	}
}
