package handler

import (
	"github.com/bryantaolong/system/internal/model/response"
	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserRoleHandler struct {
	userRoleSvc *service.UserRoleService
}

func NewUserRoleHandler(userRoleSvc *service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{userRoleSvc: userRoleSvc}
}

// ListRoles  GET /api/user/role/all
func (h *UserRoleHandler) ListRoles(c *gin.Context) {
	list, err := h.userRoleSvc.ListAll(c.Request.Context())
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list)
}
