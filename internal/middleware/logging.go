package middleware

import (
	"net/http"
	"time"

	"example.com/go-http-template/internal/pkg/logger"
	"example.com/go-http-template/internal/pkg/requestid"
	response "example.com/go-http-template/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logging(base *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log := base.With(zap.String("request_id", requestid.FromContext(c.Request.Context())))
		c.Request = c.Request.WithContext(logger.WithContext(c.Request.Context(), log))
		c.Next()
		log.Info("http request", zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)))
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.FromContext(c.Request.Context()).Error("request panic", zap.Any("panic", recovered), zap.Stack("stack"))
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
	})
}
