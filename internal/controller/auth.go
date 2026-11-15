package controller

import (
	"errors"
	"net/http"

	"example.com/go-http-template/internal/dto"
	"example.com/go-http-template/internal/model"
	"example.com/go-http-template/internal/pkg/logger"
	"example.com/go-http-template/internal/pkg/response"
	"example.com/go-http-template/internal/service"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Auth struct {
	service  *service.Auth
	sessions *scs.SessionManager
}

func NewAuth(service *service.Auth, sessions *scs.SessionManager) *Auth {
	return &Auth{service: service, sessions: sessions}
}

// Register godoc
// @Summary 注册并登录
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "注册信息"
// @Success 201 {object} response.Response[dto.SessionResponse]
// @Failure 400,403,409,429,500 {object} response.ErrorResponse
// @Router /auth/register [post]
func (h *Auth) Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, "无效的请求内容")
		return
	}
	user, err := h.service.Register(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.signIn(c, user, http.StatusCreated)
}

// Login godoc
// @Summary 登录
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response[dto.SessionResponse]
// @Failure 400,401,403,429,500 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *Auth) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, "无效的请求内容")
		return
	}
	user, err := h.service.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.signIn(c, user, http.StatusOK)
}

func (h *Auth) signIn(c *gin.Context, user model.User, status int) {
	ctx := c.Request.Context()
	if err := h.sessions.RenewToken(ctx); err != nil {
		h.fail(c, err)
		return
	}
	h.sessions.Put(ctx, "userID", user.ID)
	response.Success(c, status, dto.NewSessionResponse(user))
}

// Session godoc
// @Summary 当前会话，未登录时 data 为 null
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response[dto.SessionResponse]
// @Failure 500 {object} response.ErrorResponse
// @Router /auth/session [get]
func (h *Auth) Session(c *gin.Context) {
	id, ok := h.sessions.Get(c.Request.Context(), "userID").(uint)
	if !ok {
		response.Success[any](c, http.StatusOK, nil)
		return
	}
	user, err := h.service.User(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := h.sessions.Destroy(c.Request.Context()); err != nil {
			h.fail(c, err)
			return
		}
		response.Success[any](c, http.StatusOK, nil)
		return
	}
	if err != nil {
		h.fail(c, err)
		return
	}
	response.Success(c, http.StatusOK, dto.NewSessionResponse(user))
}

// Logout godoc
// @Summary 退出并撤销当前会话
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response[any]
// @Failure 403,500 {object} response.ErrorResponse
// @Router /auth/logout [post]
func (h *Auth) Logout(c *gin.Context) {
	if err := h.sessions.Destroy(c.Request.Context()); err != nil {
		h.fail(c, err)
		return
	}
	response.Success(c, http.StatusOK, struct{}{})
}

func (h *Auth) fail(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRegistration):
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidRequest, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Fail(c, http.StatusUnauthorized, response.CodeCredentials, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		response.Fail(c, http.StatusConflict, response.CodeEmailTaken, err.Error())
	default:
		logger.FromContext(c.Request.Context()).Error("authentication failed", zap.Error(err))
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
	}
}
