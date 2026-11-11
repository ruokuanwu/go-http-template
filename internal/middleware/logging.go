package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logging(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http request", zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)))
	}
}

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("request panic", zap.Any("panic", recovered), zap.Stack("stack"))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
	})
}
