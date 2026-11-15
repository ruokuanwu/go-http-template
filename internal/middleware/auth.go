package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"example.com/go-http-template/internal/pkg/logger"
	"example.com/go-http-template/internal/pkg/response"
	"example.com/go-http-template/internal/service"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RequireAuth(sessions *scs.SessionManager, auth *service.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := sessions.Get(c.Request.Context(), "userID").(uint)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
			return
		}
		user, err := auth.User(c.Request.Context(), id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusUnauthorized, response.CodeUnauthorized, "请先登录")
			return
		}
		if err != nil {
			logger.FromContext(c.Request.Context()).Error("load authenticated user", zap.Error(err))
			response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func CSRF(protection *http.CrossOriginProtection) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := protection.Check(c.Request); err != nil {
			response.Fail(c, http.StatusForbidden, response.CodeForbidden, "请求来源不允许")
			return
		}
		c.Next()
	}
}

// Bound both per-address attempts and concurrent password hashing work.
func AuthRateLimit() gin.HandlerFunc {
	type bucket struct {
		count int
		reset time.Time
	}
	var mu sync.Mutex
	entries := make(map[string]bucket)
	active := make(chan struct{}, 2)
	return func(c *gin.Context) {
		now, ip := time.Now(), c.ClientIP()
		mu.Lock()
		for key, entry := range entries {
			if !now.Before(entry.reset) {
				delete(entries, key)
			}
		}
		entry, exists := entries[ip]
		allowed := (exists || len(entries) < 4096) && entry.count < 10
		if allowed {
			if !exists {
				entry.reset = now.Add(time.Minute)
			}
			entry.count++
			entries[ip] = entry
		}
		mu.Unlock()
		if allowed {
			select {
			case active <- struct{}{}:
				defer func() { <-active }()
				c.Next()
				return
			default:
			}
		}
		c.Header("Retry-After", "60")
		response.Fail(c, http.StatusTooManyRequests, response.CodeRateLimited, "请求过于频繁，请稍后重试")
	}
}
