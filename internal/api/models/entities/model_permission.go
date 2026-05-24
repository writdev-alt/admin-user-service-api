package entities

// ModelPermission is the pivot between a polymorphic model (users, admins, etc.) and permissions.
type ModelPermission struct {
	ModelID       uint64 `gorm:"column:model_id;primaryKey;not null;index:idx_model_permissions_model_id" json:"modelId"`
	ModelType     string `gorm:"column:model_type;primaryKey;type:varchar(255);not null;index:idx_model_permissions_model_type" json:"modelType"`
	PermissionID  uint64 `gorm:"column:permission_id;primaryKey;not null;index:idx_model_permissions_permission_id" json:"permissionId"`
}

func (ModelPermission) TableName() string { return "model_permissions" }
