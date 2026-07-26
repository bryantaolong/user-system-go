package user

import (
	"strconv"

	"github.com/bryantaolong/user-system/auth"
	"github.com/bryantaolong/user-system/model"
	"github.com/bryantaolong/user-system/response"
	"github.com/gin-gonic/gin"
)

// Handler 用户模块处理器
type Handler struct {
	userSvc    *Service
	profileSvc *ProfileService
	authSvc    *auth.Service
}

// NewHandler 创建用户模块处理器实例
func NewHandler(userSvc *Service, profileSvc *ProfileService, authSvc *auth.Service) *Handler {
	return &Handler{
		userSvc:    userSvc,
		profileSvc: profileSvc,
		authSvc:    authSvc,
	}
}

// ========== 用户管理 ==========

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
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
func (h *Handler) ListUsers(c *gin.Context) {
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
func (h *Handler) GetUserByID(c *gin.Context) {
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
func (h *Handler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userSvc.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		handleError(c, err)
		return
	}
	response.Success(c, user)
}

// QueryUsers 用户搜索
func (h *Handler) QueryUsers(c *gin.Context) {
	var req model.UserQueryRequest
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
func (h *Handler) UpdateUser(c *gin.Context) {
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
func (h *Handler) ChangeRole(c *gin.Context) {
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
func (h *Handler) ResetPassword(c *gin.Context) {
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
func (h *Handler) BlockUser(c *gin.Context) {
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
func (h *Handler) UnblockUser(c *gin.Context) {
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
func (h *Handler) DeleteUser(c *gin.Context) {
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

// ExportAllUsers 导出所有用户数据为 Excel
func (h *Handler) ExportAllUsers(c *gin.Context) {
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	_ = pageNum // 导出服务内部处理分页

	f, err := h.userSvc.ExportAllUsers()
	if err != nil {
		handleError(c, err)
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")

	if err := f.Write(c.Writer); err != nil {
		response.Error(c, response.StatusInternalError, "导出文件写入失败")
	}
}

// ========== 用户资料 ==========

// UploadAvatar 上传当前用户头像
func (h *Handler) UploadAvatar(c *gin.Context) {
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
func (h *Handler) GetUserProfileByUserId(c *gin.Context) {
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
func (h *Handler) GetUserProfileByRealName(c *gin.Context) {
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
func (h *Handler) GetCurrentUserProfile(c *gin.Context) {
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
func (h *Handler) UpdateUserProfile(c *gin.Context) {
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

// ========== 用户角色 ==========

// ListRoles 获取全部角色下拉选项
func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.userSvc.roleSvc.ListAll(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	vos := make([]model.UserRoleOptionVO, 0, len(roles))
	for _, r := range roles {
		vos = append(vos, model.UserRoleOptionVO{ID: r.ID, RoleName: r.RoleName})
	}

	response.Success(c, vos)
}

// handleError 统一错误处理
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *response.BusinessError:
		response.Error(c, response.StatusInternalError, e.Message)
	case *response.ResourceNotFoundError:
		response.Error(c, response.StatusNotFound, e.Message)
	case *response.UnauthorizedError:
		response.ErrorWithHTTPStatus(c, 401, response.StatusUnauthorized, e.Message)
	case *response.OptimisticLockError:
		response.Error(c, response.StatusConflict, e.Message)
	default:
		response.Error(c, response.StatusInternalError, "服务繁忙，请稍后重试")
	}
}
