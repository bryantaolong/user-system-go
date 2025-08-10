package request

// ChangeRoleRequest 对应 Java ChangeRoleRequest
type ChangeRoleRequest struct {
	RoleIds []int64 `json:"roleIds" binding:"required,dive,required"`
}
