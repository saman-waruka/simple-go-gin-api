package middleware

import (
	"fmt"
	"time"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func RecoveryWithZap(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 🔥 Log panic ด้วย zap
				logger.Error("PANIC recovered",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)

				// 🚨 ส่งไป Sentry
				sentry.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("method", c.Request.Method)
					scope.SetTag("path", c.Request.URL.Path)
					scope.SetExtra("stacktrace", string(debug.Stack()))
					sentry.CaptureException(fmt.Errorf("%v", err))
				})
				sentry.Flush(2 * time.Second)

				// ❌ ไม่เปิดเผย error แท้จริงแก่ client
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}