package entities

// RolePermission is the pivot between roles and permissions.
type RolePermission struct {
	PermissionID uint64 `gorm:"column:permission_id;primaryKey;not null" json:"permissionId"`
	RoleID       uint64 `gorm:"column:role_id;primaryKey;not null;index:role_permissions_role_id_index" json:"roleId"`
}

func (RolePermission) TableName() string { return "role_permissions" }
