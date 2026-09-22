package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLogMiddleware 记录所有写操作访问日志；业务审计由 service 落 audit_log 表。
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if c.Request.Method != "GET" {
			log.Printf("[audit] %s %s -> %d actor=%s duration=%s",
				c.Request.Method, c.Request.URL.Path, c.Writer.Status(),
				c.GetString("actor"), time.Since(start))
		}
	}
}
