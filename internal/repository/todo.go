package repository

import (
	"context"

	"example.com/go-http-template/internal/model"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

type Todo struct{ db *gorm.DB }

func NewTodo(db *gorm.DB) *Todo { return &Todo{db: db} }

func (r *Todo) List(ctx context.Context) ([]model.Todo, error) {
	todos := make([]model.Todo, 0)
	if err := r.db.WithContext(ctx).Order("id DESC").Find(&todos).Error; err != nil {
		return nil, oops.Wrapf(err, "list todos")
	}
	return todos, nil
}

func (r *Todo) Create(ctx context.Context, title string) (model.Todo, error) {
	todo := model.Todo{Title: title}
	if err := r.db.WithContext(ctx).Create(&todo).Error; err != nil {
		return model.Todo{}, oops.Wrapf(err, "create todo")
	}
	return todo, nil
}

func (r *Todo) Update(ctx context.Context, id uint, completed bool) (model.Todo, error) {
	var todo model.Todo
	if err := r.db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return todo, oops.Wrapf(err, "find todo")
	}
	if err := r.db.WithContext(ctx).Model(&todo).Update("completed", completed).Error; err != nil {
		return todo, oops.Wrapf(err, "update todo")
	}
	return todo, nil
}

func (r *Todo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Todo{}, id)
	if result.Error != nil {
		return oops.Wrapf(result.Error, "delete todo")
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
