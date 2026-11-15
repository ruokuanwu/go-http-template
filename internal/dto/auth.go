package dto

import "example.com/go-http-template/internal/model"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SessionResponse struct {
	User UserResponse `json:"user"`
}

func NewSessionResponse(user model.User) SessionResponse {
	return SessionResponse{User: UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}}
}
