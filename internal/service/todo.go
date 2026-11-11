package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"example.com/go-http-template/internal/model"
	"example.com/go-http-template/internal/pkg/logger"
)

var ErrInvalidTitle = errors.New("标题必须为 1 到 200 个字符")

type TodoRepository interface {
	List(context.Context) ([]model.Todo, error)
	Create(context.Context, string) (model.Todo, error)
	Update(context.Context, uint, bool) (model.Todo, error)
	Delete(context.Context, uint) error
}

type Todo struct{ repository TodoRepository }

func NewTodo(repository TodoRepository) *Todo { return &Todo{repository: repository} }

func (s *Todo) List(ctx context.Context) ([]model.Todo, error) { return s.repository.List(ctx) }

func (s *Todo) Create(ctx context.Context, title string) (model.Todo, error) {
	logger.FromContext(ctx).Info("Creating todo")
	title = strings.TrimSpace(title)
	if title == "" || utf8.RuneCountInString(title) > 200 {
		return model.Todo{}, ErrInvalidTitle
	}
	return s.repository.Create(ctx, title)
}

func (s *Todo) Update(ctx context.Context, id uint, completed bool) (model.Todo, error) {
	return s.repository.Update(ctx, id, completed)
}

func (s *Todo) Delete(ctx context.Context, id uint) error { return s.repository.Delete(ctx, id) }
