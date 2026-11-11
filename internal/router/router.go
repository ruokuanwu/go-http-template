package router

import (
	"net/http"
	"path/filepath"

	"example.com/go-http-template/internal/controller"
	"example.com/go-http-template/internal/dto"
	"example.com/go-http-template/internal/middleware"
	response "example.com/go-http-template/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(todos *controller.Todo, logger *zap.Logger, webDir string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logging(logger), middleware.Recovery())
	_ = r.SetTrustedProxies(nil)
	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) { response.Success(c, http.StatusOK, dto.HealthResponse{Status: "ok"}) })
	api.GET("/todos", todos.List)
	api.POST("/todos", todos.Create)
	api.PATCH("/todos/:id", todos.Update)
	api.DELETE("/todos/:id", todos.Delete)
	r.Static("/assets", filepath.Join(webDir, "assets"))
	r.GET("/", func(c *gin.Context) { c.File(filepath.Join(webDir, "index.html")) })
	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, http.StatusNotFound, response.CodeRouteNotFound, "路由不存在")
	})
	return r
}
