package controller

import (
	"errors"
	"net/http"

	"example.com/go-http-template/internal/dto"
	response "example.com/go-http-template/internal/pkg/response"
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
	response.Success(c, http.StatusOK, dto.NewListTodosResponse(todos))
}

func (h *Todo) Create(c *gin.Context) {
	var input dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, "无效的请求内容")
		return
	}
	todo, err := h.service.Create(c.Request.Context(), input.Title)
	if err != nil {
		h.fail(c, err)
		return
	}
	response.Success(c, http.StatusCreated, dto.NewTodoResponse(todo))
}

func (h *Todo) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var input dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, "completed 必须为布尔值")
		return
	}
	todo, err := h.service.Update(c.Request.Context(), id, *input.Completed)
	if err != nil {
		h.fail(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.NewTodoResponse(todo))
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
	response.Success(c, http.StatusOK, struct{}{})
}

func parseID(c *gin.Context) (uint, bool) {
	var input dto.TodoIDRequest
	if err := c.ShouldBindUri(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, "无效的任务 ID")
		return 0, false
	}
	return input.ID, true
}

func (h *Todo) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTitle):
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.Fail(c, http.StatusNotFound, response.CodeTodoNotFound, "任务不存在")
	default:
		h.logger.Error("request failed", zap.Any("error", err))
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
	}
}
