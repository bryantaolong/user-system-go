package auth

import (
	"context"
	"net/http"

	"github.com/bryantaolong/user-system/types"
	"github.com/bryantaolong/user-system/response"
	"github.com/gin-gonic/gin"
)

// ProfileCreator 用户资料创建回调（由 user 模块提供）
type ProfileCreator func(ctx context.Context, profile *model.UserProfile) (*model.UserProfile, error)

// Handler 认证处理器
type Handler struct {
	authSvc       *Service
	createProfile ProfileCreator
}

// NewHandler 创建认证处理器实例
func NewHandler(authSvc *Service, createProfile ProfileCreator) *Handler {
	return &Handler{authSvc: authSvc, createProfile: createProfile}
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
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
	if _, err := h.createProfile(c.Request.Context(), profile); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserVO(user))
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
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
func (h *Handler) GetCurrentUser(c *gin.Context) {
	user, err := h.authSvc.GetCurrentUser(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, model.ToUserVO(user))
}

// ValidateToken 验证 Token 合法性
func (h *Handler) ValidateToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, response.StatusBadRequest, "token 参数不能为空")
		return
	}

	result, err := h.authSvc.ValidateTokenWithStatus(token)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, result)
}

// ChangePassword 修改用户密码
func (h *Handler) ChangePassword(c *gin.Context) {
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
func (h *Handler) DeleteAccount(c *gin.Context) {
	user, err := h.authSvc.DeleteAccount(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, model.ToUserVO(user))
}

// Logout 退出登录
func (h *Handler) Logout(c *gin.Context) {
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
