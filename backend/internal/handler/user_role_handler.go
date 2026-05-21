package handler

import (
	"github.com/bryan/user-system/internal/model"
	"github.com/bryan/user-system/internal/pkg/response"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
)

type UserRoleHandler struct {
	roleSvc *service.UserRoleService
}

func NewUserRoleHandler(roleSvc *service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{roleSvc: roleSvc}
}

// ListRoles 获取全部角色下拉选项
func (h *UserRoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleSvc.ListAll(c.Request.Context())
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
