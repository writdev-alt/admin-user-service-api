package request

type UserListRequest struct {
	PageInfo
	Email    *string `form:"email" json:"email"`
	Username *string `form:"username" json:"username"`
	Status   *bool   `form:"status" json:"status"`
}

type UserCreateRequest struct {
	Username        string  `json:"username" binding:"required,min=3"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=6"`
	Phone           *string `json:"phone"`
	Country         *string `json:"country"`
	Status          *bool   `json:"status"`
	TwoFactorEnabled *bool  `json:"twoFactorEnabled"`
}

type UserUpdateRequest struct {
	Username         *string `json:"username"`
	Email            *string `json:"email"`
	Phone            *string `json:"phone"`
	Country          *string `json:"country"`
	Status           *bool   `json:"status"`
	TwoFactorEnabled *bool   `json:"twoFactorEnabled"`
}

type RoleListRequest struct {
	Name      *string `form:"name" json:"name"`
	GuardName *string `form:"guardName" json:"guardName"`
}

type PermissionListRequest struct {
	Name      *string `form:"name" json:"name"`
	GuardName *string `form:"guardName" json:"guardName"`
}
