package middleware

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"example.com/go-http-template/internal/pkg/logger"
	"github.com/alexedwards/scs/v2"
	"go.uber.org/zap"
)

type sessionResponse struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *sessionResponse) Header() http.Header { return w.header }
func (w *sessionResponse) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
	}
}
func (w *sessionResponse) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.body.Write(data)
}

// Commit before publishing API responses so persistence failures cannot emit success JSON.
// Static files bypass sessions and keep their normal streaming behavior.
func Session(sessions *scs.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Add("Vary", "Cookie")
		var token string
		if cookie, err := r.Cookie(sessions.Cookie.Name); err == nil {
			token = cookie.Value
		}
		ctx, err := sessions.Load(r.Context(), token)
		if err != nil {
			sessions.ErrorFunc(w, r, err)
			return
		}
		r = r.WithContext(ctx)
		buffer := &sessionResponse{header: make(http.Header)}
		next.ServeHTTP(buffer, r)
		switch sessions.Status(ctx) {
		case scs.Modified:
			token, expiry, err := sessions.Commit(ctx)
			if err != nil {
				sessions.ErrorFunc(w, r, err)
				return
			}
			sessions.WriteSessionCookie(ctx, buffer, token, expiry)
		case scs.Destroyed:
			sessions.WriteSessionCookie(ctx, buffer, "", time.Time{})
		}
		for key, values := range buffer.header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		if buffer.status == 0 {
			buffer.status = http.StatusOK
		}
		w.WriteHeader(buffer.status)
		if _, err := w.Write(buffer.body.Bytes()); err != nil {
			logger.FromContext(ctx).Error("write session response", zap.Error(err))
		}
	})
}
