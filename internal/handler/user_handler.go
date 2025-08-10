package handler

import (
	"strconv"

	"github.com/bryantaolong/system/internal/model/request"
	"github.com/bryantaolong/system/internal/model/response"
	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var pageReq request.PageRequest
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		response.Fail(c, err.Error())
		return
	}
	users, total, err := h.userService.GetAllUsers(c, pageReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"list": users, "total": total})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("userId")
	userID, _ := strconv.ParseInt(idStr, 10, 64)
	user, err := h.userService.GetUserByID(c, userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.userService.GetUserByUsername(c, username)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	var searchReq request.UserSearchRequest
	var pageReq request.PageRequest
	if err := c.ShouldBindJSON(&searchReq); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := c.ShouldBindQuery(&pageReq); err != nil {
		response.Fail(c, err.Error())
		return
	}
	users, total, err := h.userService.SearchUsers(c, searchReq, pageReq)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"list": users, "total": total})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	var req request.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user, err := h.userService.UpdateUser(c, userID, req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) ChangeRole(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		response.Fail(c, "userId 必须是整数")
		return
	}

	var req request.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	user, err := h.userService.ChangeRoleByIds(c.Request.Context(), userID, req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	user, err := h.userService.ChangePassword(c, userID, req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) ChangePasswordForcefully(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	newPassword := c.Param("newPassword")
	user, err := h.userService.ChangePasswordForcefully(c, userID, newPassword)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	user, err := h.userService.BlockUser(c, userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	user, err := h.userService.UnblockUser(c, userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("userId"), 10, 64)
	user, err := h.userService.DeleteUser(c, userID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, user)
}
