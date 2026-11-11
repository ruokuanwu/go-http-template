package response

import "github.com/gin-gonic/gin"

const (
	CodeInternalError  = 10000
	CodeInvalidRequest = 10001
	CodeRouteNotFound  = 10004
	CodeTodoNotFound   = 20001
)

type Response[T any] struct {
	Data    T      `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Success[T any](c *gin.Context, status int, data T) {
	c.JSON(status, Response[T]{Data: data, Code: 0, Message: ""})
}

func Fail(c *gin.Context, status, code int, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{Code: code, Message: message})
}
