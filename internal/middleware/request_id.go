package middleware

import (
	"crypto/rand"

	"example.com/go-http-template/internal/pkg/requestid"
	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := rand.Text()
		c.Request = c.Request.WithContext(requestid.WithContext(c.Request.Context(), id))
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
