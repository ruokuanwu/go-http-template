package controller

import (
	"errors"
	"net/http"

	"example.com/go-http-template/internal/dto"
	"example.com/go-http-template/internal/pkg/logger"
	response "example.com/go-http-template/internal/pkg/response"
	"example.com/go-http-template/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Todo struct {
	service *service.Todo
}

func NewTodo(service *service.Todo) *Todo {
	return &Todo{service: service}
}

// List godoc
// @Summary 获取任务列表
// @Tags todos
// @Produce json
// @Success 200 {object} response.Response[dto.ListTodosResponse]
// @Failure 500 {object} response.ErrorResponse "10000：服务器内部错误"
// @Router /todos [get]
func (h *Todo) List(c *gin.Context) {
	todos, err := h.service.List(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.NewListTodosResponse(todos))
}

// Create godoc
// @Summary 创建任务
// @Description 标题去除首尾空白后必须为 1–200 个字符。
// @Tags todos
// @Accept json
// @Produce json
// @Param body body dto.CreateTodoRequest true "任务内容"
// @Success 201 {object} response.Response[dto.TodoResponse]
// @Failure 400 {object} response.ErrorResponse "10001：请求内容或标题无效"
// @Failure 500 {object} response.ErrorResponse "10000：服务器内部错误"
// @Router /todos [post]
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

// Update godoc
// @Summary 更新任务完成状态
// @Description completed 必填，false 也是有效值。
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "任务 ID" minimum(1)
// @Param body body dto.UpdateTodoRequest true "完成状态"
// @Success 200 {object} response.Response[dto.TodoResponse]
// @Failure 400 {object} response.ErrorResponse "10001：任务 ID 或请求内容无效"
// @Failure 404 {object} response.ErrorResponse "20001：任务不存在"
// @Failure 500 {object} response.ErrorResponse "10000：服务器内部错误"
// @Router /todos/{id} [patch]
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

// Delete godoc
// @Summary 删除任务
// @Tags todos
// @Produce json
// @Param id path int true "任务 ID" minimum(1)
// @Success 200 {object} response.Response[any]{data=object} "删除成功，data 为 {}"
// @Failure 400 {object} response.ErrorResponse "10001：任务 ID 无效"
// @Failure 404 {object} response.ErrorResponse "20001：任务不存在"
// @Failure 500 {object} response.ErrorResponse "10000：服务器内部错误"
// @Router /todos/{id} [delete]
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
		logger.FromContext(c.Request.Context()).Error("request failed", zap.Any("error", err))
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
	}
}
