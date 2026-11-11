package router

import (
	"net/http"
	"path/filepath"
	"strings"

	_ "example.com/go-http-template/docs"
	"example.com/go-http-template/internal/controller"
	"example.com/go-http-template/internal/middleware"
	response "example.com/go-http-template/internal/pkg/response"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func New(todos *controller.Todo, logger *zap.Logger, webDir string, swaggerEnabled bool) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logging(logger), middleware.Recovery())
	_ = r.SetTrustedProxies(nil)
	api := r.Group("/api")
	api.GET("/todos", todos.List)
	api.POST("/todos", todos.Create)
	api.PATCH("/todos/:id", todos.Update)
	api.DELETE("/todos/:id", todos.Delete)
	if swaggerEnabled {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")))
	}
	r.Static("/assets", filepath.Join(webDir, "assets"))
	r.GET("/", func(c *gin.Context) { c.File(filepath.Join(webDir, "index.html")) })
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if (c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead) &&
			path != "/api" && !strings.HasPrefix(path, "/api/") &&
			path != "/swagger" && !strings.HasPrefix(path, "/swagger/") &&
			path != "/assets" && !strings.HasPrefix(path, "/assets/") &&
			!strings.Contains(filepath.Base(path), ".") {
			c.File(filepath.Join(webDir, "index.html"))
			return
		}
		response.Fail(c, http.StatusNotFound, response.CodeRouteNotFound, "路由不存在")
	})
	return r
}
