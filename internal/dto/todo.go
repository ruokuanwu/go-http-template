package dto

import "example.com/go-http-template/internal/model"

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Completed *bool `json:"completed" binding:"required"`
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
