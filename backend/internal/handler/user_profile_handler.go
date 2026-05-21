package handler

import (
	"strconv"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserProfileHandler struct {
	profileSvc *service.UserProfileService
	userSvc    *service.UserService
	authSvc    *service.AuthService
}

func NewUserProfileHandler(
	profileSvc *service.UserProfileService,
	userSvc *service.UserService,
	authSvc *service.AuthService,
) *UserProfileHandler {
	return &UserProfileHandler{profileSvc: profileSvc, userSvc: userSvc, authSvc: authSvc}
}

// UploadAvatar 上传当前用户头像
func (h *UserProfileHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		response.Error(c, response.StatusBadRequest, "上传文件不能为空")
		return
	}

	userID, err := h.authSvc.GetCurrentUserID(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	avatarPath, err := h.profileSvc.UpdateAvatar(c.Request.Context(), userID, file)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, avatarPath)
}

// GetUserProfileByUserId 根据用户主键查询用户资料
func (h *UserProfileHandler) GetUserProfileByUserId(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	profile, err := h.profileSvc.GetUserProfileByUserId(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserProfileVO(user, profile))
}

// GetUserProfileByRealName 根据真实姓名查询用户资料
func (h *UserProfileHandler) GetUserProfileByRealName(c *gin.Context) {
	realName := c.Param("realName")

	profile, err := h.profileSvc.GetUserProfileByRealName(c.Request.Context(), realName)
	if err != nil {
		handleError(c, err)
		return
	}

	user, err := h.userSvc.GetUserByID(c.Request.Context(), profile.UserID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserProfileVO(user, profile))
}

// GetCurrentUserProfile 获取当前登录用户的资料
func (h *UserProfileHandler) GetCurrentUserProfile(c *gin.Context) {
	userID, err := h.authSvc.GetCurrentUserID(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	profile := h.profileSvc.GetUserProfileByUserIdOrEmpty(c.Request.Context(), userID)
	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, model.ToUserProfileVO(user, profile))
}

// UpdateUserProfile 更新当前用户资料
func (h *UserProfileHandler) UpdateUserProfile(c *gin.Context) {
	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.authSvc.GetCurrentUserID(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	profile, err := h.profileSvc.UpdateUserProfile(c.Request.Context(), userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	user, _ := h.authSvc.GetCurrentUser(c.Request.Context())
	response.Success(c, model.ToUserProfileVO(user, profile))
}
