package controller

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/go-http-template/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Todo struct {
	service *service.Todo
	logger  *zap.Logger
}

func NewTodo(service *service.Todo, logger *zap.Logger) *Todo {
	return &Todo{service: service, logger: logger}
}

func (h *Todo) List(c *gin.Context) {
	todos, err := h.service.List(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (h *Todo) Create(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求内容"})
		return
	}
	todo, err := h.service.Create(c.Request.Context(), input.Title)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": todo})
}

func (h *Todo) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input struct {
		Completed *bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "completed 必须为布尔值"})
		return
	}
	todo, err := h.service.Update(c.Request.Context(), id, *input.Completed)
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func (h *Todo) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务 ID"})
		return 0, false
	}
	return uint(id), true
}

func (h *Todo) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTitle):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
	default:
		h.logger.Error("request failed", zap.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
	}
}
