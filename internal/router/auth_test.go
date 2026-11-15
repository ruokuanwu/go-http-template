package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"example.com/go-http-template/internal/controller"
	"example.com/go-http-template/internal/db"
	"example.com/go-http-template/internal/middleware"
	"example.com/go-http-template/internal/repository"
	"example.com/go-http-template/internal/router"
	"example.com/go-http-template/internal/service"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestAuthenticationLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database, err := db.Open(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	if database.Migrator().HasTable("todos") {
		t.Fatal("removed todo table must not be created")
	}
	connection, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { connection.Close() })
	auth, err := service.NewAuth(repository.NewUser(database))
	if err != nil {
		t.Fatal(err)
	}
	newManager := func() *scs.SessionManager {
		sm := scs.New()
		sm.Cookie.Secure = true
		sm.Cookie.Persist = true
		sm.Store = sqlite3store.NewWithCleanupInterval(connection, 0)
		return sm
	}
	manager := newManager()
	newHandler := func(sm *scs.SessionManager) http.Handler {
		r := router.New(controller.NewAuth(auth, sm), http.NewCrossOriginProtection(), zap.NewNop(), t.TempDir(), false)
		r.GET("/api/private", middleware.RequireAuth(sm, auth), func(c *gin.Context) { c.Status(http.StatusOK) })
		return middleware.Session(sm, r)
	}
	handler := newHandler(manager)
	request := func(method, path, body string, cookie *http.Cookie, origin string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(method, "https://example.com"+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if cookie != nil {
			req.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		return w
	}
	check := func(w *httptest.ResponseRecorder, status int) {
		t.Helper()
		if w.Code != status {
			t.Fatalf("status %d, want %d: %s", w.Code, status, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "PasswordHash") || strings.Contains(w.Body.String(), "$argon2") {
			t.Fatal("password hash leaked")
		}
	}
	check(request("GET", "/api/private", "", nil, ""), 401)
	for _, route := range []struct{ method, path string }{
		{"GET", "/api/todos"},
		{"POST", "/api/todos"},
		{"PATCH", "/api/todos/1"},
		{"DELETE", "/api/todos/1"},
	} {
		w := request(route.method, route.path, "", nil, "")
		check(w, http.StatusNotFound)
		var result struct {
			Code int `json:"code"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil || result.Code != 10004 {
			t.Fatalf("%s %s: expected route-not-found JSON, got %s", route.method, route.path, w.Body.String())
		}
	}
	check(request("POST", "/api/auth/register", `{}`, nil, ""), 400)
	check(request("POST", "/api/auth/register", `{"name":"Alex","email":"long@example.com","password":"`+strings.Repeat("😀", 65)+`"}`, nil, ""), 400)
	check(request("POST", "/api/auth/login", `{}`, nil, "https://evil.example"), 403)
	registered := request("POST", "/api/auth/register", `{"name":"Alex","email":" Alex@Example.com ","password":"password123"}`, nil, "")
	check(registered, 201)
	cookies := registered.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected session cookie: %v", cookies)
	}
	cookie := cookies[0]
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatal("unsafe cookie attributes")
	}
	check(request("GET", "/api/private", "", cookie, ""), 200)
	check(request("POST", "/api/auth/register", `{"name":"Alex","email":"alex@example.com","password":"password123"}`, nil, ""), 409)
	for _, email := range []string{"alex@example.com", "missing@example.com"} {
		w := request("POST", "/api/auth/login", `{"email":"`+email+`","password":"wrong-password"}`, nil, "")
		check(w, 401)
		if !strings.Contains(w.Body.String(), "30002") {
			t.Fatal(w.Body.String())
		}
	}
	// A new manager with the same SQLite store simulates process restart.
	manager = newManager()
	handler = newHandler(manager)
	w := request("GET", "/api/auth/session", "", cookie, "")
	check(w, 200)
	if !strings.Contains(w.Body.String(), "alex@example.com") || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal(w.Body.String())
	}
	login := request("POST", "/api/auth/login", `{"email":"alex@example.com","password":"password123"}`, cookie, "")
	check(login, 200)
	newCookie := login.Result().Cookies()[0]
	if newCookie.Value == cookie.Value {
		t.Fatal("session token not rotated")
	}
	check(request("GET", "/api/private", "", cookie, ""), 401)
	check(request("POST", "/api/auth/logout", "", newCookie, ""), 200)
	check(request("GET", "/api/private", "", newCookie, ""), 401)
	w = request("GET", "/api/auth/session", "", newCookie, "")
	var anonymous struct {
		Data any `json:"data"`
		Code int `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &anonymous); err != nil || anonymous.Data != nil || anonymous.Code != 0 {
		t.Fatal(w.Body.String())
	}
	// Expired records must not authenticate.
	login = request("POST", "/api/auth/login", `{"email":"alex@example.com","password":"password123"}`, nil, "")
	check(login, 200)
	if err := database.Exec("UPDATE sessions SET expiry = julianday('now', '-1 hour')").Error; err != nil {
		t.Fatal(err)
	}
	check(request("GET", "/api/private", "", login.Result().Cookies()[0], ""), 401)
	// Database failures must remain errors, not anonymous successes.
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	check(request("GET", "/api/auth/session", "", login.Result().Cookies()[0], ""), 500)
}

func TestAuthRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", middleware.AuthRateLimit(), func(c *gin.Context) { c.Status(200) })
	for i := 0; i < 11; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest("POST", "/login", nil))
		want := 200
		if i == 10 {
			want = 429
		}
		if w.Code != want {
			t.Fatalf("attempt %d: got %d", i, w.Code)
		}
	}
}
