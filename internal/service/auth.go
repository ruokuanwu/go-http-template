package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"example.com/go-http-template/internal/model"
	"github.com/alexedwards/argon2id"
	"github.com/samber/oops"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("邮箱或密码错误")
	ErrEmailTaken          = errors.New("邮箱已被注册")
	ErrInvalidRegistration = errors.New("姓名、邮箱或密码格式无效")
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
	ByEmail(context.Context, string) (model.User, error)
	ByID(context.Context, uint) (model.User, error)
}

type Auth struct {
	repository UserRepository
	dummyHash  string
}

func NewAuth(repository UserRepository) (*Auth, error) {
	hash, err := argon2id.CreateHash("non-user-password-for-timing", argon2id.DefaultParams)
	if err != nil {
		return nil, oops.Wrapf(err, "initialize password verification")
	}
	return &Auth{repository: repository, dummyHash: hash}, nil
}

func (s *Auth) Register(ctx context.Context, name, email, password string) (model.User, error) {
	name, email = strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	// Count UTF-16 units to match JavaScript string length in the existing form.
	units := len(utf16.Encode([]rune(password)))
	if !utf8.ValidString(password) || name == "" || len(utf16.Encode([]rune(name))) > 100 ||
		len(email) > 254 || err != nil || address.Address != email || units < 8 || units > 128 {
		return model.User{}, ErrInvalidRegistration
	}
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return model.User{}, oops.Wrapf(err, "hash password")
	}
	return s.repository.Create(ctx, model.User{Name: name, Email: email, PasswordHash: hash})
}

func (s *Auth) Login(ctx context.Context, email, password string) (model.User, error) {
	if len(email) > 254 || len(password) > 512 || password == "" {
		return model.User{}, ErrInvalidCredentials
	}
	user, err := s.repository.ByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	missing := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !missing {
		return model.User{}, err
	}
	hash := user.PasswordHash
	if missing {
		hash = s.dummyHash
	}
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return model.User{}, oops.Wrapf(err, "verify password")
	}
	if missing || !match {
		return model.User{}, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Auth) User(ctx context.Context, id uint) (model.User, error) {
	return s.repository.ByID(ctx, id)
}
