package request

type UserListRequest struct {
	PageInfo
	Email    *string `form:"email" json:"email"`
	Username *string `form:"username" json:"username"`
	Search   *string `form:"search" json:"search"`
	Status   *bool   `form:"status" json:"status"`
}

type UserCreateRequest struct {
	Username         string   `json:"username" binding:"required,min=3"`
	Email            string   `json:"email" binding:"required,email"`
	Password         string   `json:"password" binding:"required,min=6"`
	Phone            *string  `json:"phone"`
	Country          *string  `json:"country"`
	Status           *bool    `json:"status"`
	TwoFactorEnabled *bool    `json:"twoFactorEnabled"`
	RoleIDs          []uint64 `json:"roleIds"`
	PermissionIDs    []uint64 `json:"permissionIds"`
}

type UserUpdateRequest struct {
	Username         *string   `json:"username"`
	Email            *string   `json:"email"`
	Phone            *string   `json:"phone"`
	Country          *string   `json:"country"`
	Status           *bool     `json:"status"`
	RoleIDs          *[]uint64 `json:"roleIds"`
	PermissionIDs    *[]uint64 `json:"permissionIds"`
	TwoFactorEnabled *bool     `json:"twoFactorEnabled"`
}

// ChangeUserPasswordRequest sets a user's password (admin action; no current password).
type ChangeUserPasswordRequest struct {
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=NewPassword"`
}

type RoleListRequest struct {
	Name      *string `form:"name" json:"name"`
	GuardName *string `form:"guardName" json:"guardName"`
}

type PermissionListRequest struct {
	Name      *string `form:"name" json:"name"`
	GuardName *string `form:"guardName" json:"guardName"`
}
