package handler

import (
	"strconv"

	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc    *service.UserService
	profileSvc *service.UserProfileService
}

func NewUserHandler(userSvc *service.UserService, profileSvc *service.UserProfileService) *UserHandler {
	return &UserHandler{userSvc: userSvc, profileSvc: profileSvc}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userSvc.CreateUser(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	// 初始化 UserProfile
	profile := &model.UserProfile{UserID: user.ID}
	h.profileSvc.CreateUserProfile(c.Request.Context(), profile)

	response.Success(c, user)
}

// ListUsers 获取所有用户列表（分页）
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.userSvc.GetAllUsers(c.Request.Context(), pageNum, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetUserByID 根据用户 ID 查询用户信息
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.userSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// GetUserByUsername 根据用户名查询用户信息
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userSvc.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// QueryUsers 用户搜索
func (h *UserHandler) QueryUsers(c *gin.Context) {
	var req model.UserSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := h.userSvc.QueryUsers(c.Request.Context(), &req, pageNum, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, result)
}

// UpdateUser 更新用户基本信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userSvc.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// ChangeRole 修改用户角色
func (h *UserHandler) ChangeRole(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	var req model.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userSvc.ChangeRoleByIds(c.Request.Context(), userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// ResetPassword 强制修改用户密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userSvc.ResetPassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// BlockUser 封禁用户
func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.userSvc.BlockUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// UnblockUser 解封用户
func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := h.userSvc.UnblockUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Error(c, response.StatusBadRequest, "无效的用户ID")
		return
	}

	id, err := h.userSvc.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, id)
}
