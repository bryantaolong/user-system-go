package handler

import (
	_ "net/http"

	"github.com/bryantaolong/system/internal/model/request"
	"github.com/bryantaolong/system/internal/model/response"
	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req) // 注意这里
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	token, err := h.authService.Login(req, c.Request)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) Me(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]
	user, err := h.authService.GetCurrentUser(token)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]
	err := h.authService.Logout(token)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, gin.H{"success": true})
}

func (h *AuthHandler) Validate(c *gin.Context) {
	token := c.Query("token")
	if !h.authService.ValidateToken(token) {
		response.Success(c, gin.H{"valid": false})
		return
	}
	response.Success(c, gin.H{"valid": true})
}
