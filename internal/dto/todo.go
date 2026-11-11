package dto

import "example.com/go-http-template/internal/model"

type CreateTodoRequest struct {
	// 标题去除首尾空白后必须为 1–200 个字符。
	Title string `json:"title" example:"学习 Go"`
}

type UpdateTodoRequest struct {
	Completed *bool `json:"completed" binding:"required" example:"false"`
}

type TodoIDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type TodoResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ListTodosResponse struct {
	Items []TodoResponse `json:"items"`
}

func NewTodoResponse(todo model.Todo) TodoResponse {
	return TodoResponse{ID: todo.ID, Title: todo.Title, Completed: todo.Completed}
}

func NewListTodosResponse(todos []model.Todo) ListTodosResponse {
	items := make([]TodoResponse, 0, len(todos))
	for _, todo := range todos {
		items = append(items, NewTodoResponse(todo))
	}
	return ListTodosResponse{Items: items}
}
