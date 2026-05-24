package request

// SetRolePermissionsRequest replaces all permissions assigned to a role.
type SetRolePermissionsRequest struct {
	PermissionIDs []uint64 `json:"permissionIds" binding:"required"`
}
