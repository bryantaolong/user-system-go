package handler

import (
	"net/http"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc    *service.AuthService
	profileSvc *service.UserProfileService
}

func NewAuthHandler(authSvc *service.AuthService, profileSvc *service.UserProfileService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, profileSvc: profileSvc}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	// 初始化 UserProfile
	profile := &model.UserProfile{UserID: user.ID}
	if _, err := h.profileSvc.CreateUserProfile(c.Request.Context(), profile); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserVO(user))
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authSvc.Login(
		c.Request.Context(),
		&req,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, token)
}

// GetCurrentUser 获取当前认证用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user, err := h.authSvc.GetCurrentUser(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, model.ToUserVO(user))
}

// ValidateToken 验证 Token 合法性
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, response.StatusBadRequest, "token 参数不能为空")
		return
	}
	if !h.authSvc.ValidateToken(token) {
		response.Success(c, "Invalid token")
		return
	}
	response.Success(c, "Validation passed")
}

// ChangePassword 修改用户密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authSvc.ChangePassword(c.Request.Context(), req.OldPassword, req.NewPassword)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserVO(user))
}

// DeleteAccount 注销用户
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	user, err := h.authSvc.DeleteAccount(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, model.ToUserVO(user))
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.authSvc.Logout(c.Request.Context()); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, true)
}

// handleError 统一错误处理
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *response.BusinessError:
		response.Error(c, response.StatusInternalError, e.Message)
	case *response.ResourceNotFoundError:
		response.Error(c, response.StatusNotFound, e.Message)
	case *response.UnauthorizedError:
		response.ErrorWithHTTPStatus(c, http.StatusUnauthorized, response.StatusUnauthorized, e.Message)
	case *response.OptimisticLockError:
		response.Error(c, response.StatusConflict, e.Message)
	default:
		response.Error(c, response.StatusInternalError, "服务繁忙，请稍后重试")
	}
}
