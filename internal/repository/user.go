package repository

import (
	"context"
	"errors"

	"example.com/go-http-template/internal/model"
	"example.com/go-http-template/internal/service"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

type User struct{ db *gorm.DB }

func NewUser(db *gorm.DB) *User { return &User{db: db} }

func (r *User) Create(ctx context.Context, user model.User) (model.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.User{}, service.ErrEmailTaken
		}
		return model.User{}, oops.Wrapf(err, "create user")
	}
	return user, nil
}

func (r *User) ByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, oops.Wrapf(err, "find user by email")
}

func (r *User) ByID(ctx context.Context, id uint) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return user, oops.Wrapf(err, "find user by id")
}
