package entities

// Polymorphic model type values stored in model_roles.model_type.
const (
	ModelTypeUser  = "User"
	ModelTypeAdmin = "Admin"
)

// ModelRole is the pivot between a polymorphic model (users, admins, etc.) and roles.
type ModelRole struct {
	ModelID   uint64 `gorm:"column:model_id;primaryKey;not null;index:idx_model_roles_model_id" json:"modelId"`
	ModelType string `gorm:"column:model_type;primaryKey;type:varchar(255);not null;index:idx_model_roles_model_type" json:"modelType"`
	RoleID    uint64 `gorm:"column:role_id;primaryKey;not null;index:idx_model_roles_role_id" json:"roleId"`
}

func (ModelRole) TableName() string { return "model_roles" }
