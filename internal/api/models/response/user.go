package response

import (
	"time"

	"github.com/google/uuid"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type UserResponse struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uint64     `json:"userId"`
	Username         string     `json:"username"`
	Email            string     `json:"email"`
	Phone            *string    `json:"phone,omitempty"`
	Country          *string    `json:"country,omitempty"`
	EmailVerifiedAt  *time.Time `json:"emailVerifiedAt,omitempty"`
	TwoFactorEnabled bool       `json:"twoFactorEnabled"`
	Status           bool       `json:"status"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
}

func ToUserResponse(u entities.User) UserResponse {
	return UserResponse{
		ID:               u.UUID,
		UserID:           u.ID,
		Username:         u.Username,
		Email:            u.Email,
		Phone:            u.Phone,
		Country:          u.Country,
		EmailVerifiedAt:  u.EmailVerifiedAt,
		TwoFactorEnabled: u.TwoFactorEnabled,
		Status:           u.Status,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

type RoleResponse struct {
	ID          uuid.UUID  `json:"id"`
	RoleID      uint64     `json:"roleId"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	GuardName   string     `json:"guardName"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

func ToRoleResponse(r entities.Role) RoleResponse {
	return RoleResponse{
		ID:          r.UUID,
		RoleID:      r.ID,
		Name:        r.Name,
		Description: r.Description,
		GuardName:   r.GuardName,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

type PermissionResponse struct {
	ID          uuid.UUID  `json:"id"`
	PermissionID uint64    `json:"permissionId"`
	Name        string     `json:"name"`
	GuardName   string     `json:"guardName"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

func ToPermissionResponse(p entities.Permission) PermissionResponse {
	return PermissionResponse{
		ID:           p.UUID,
		PermissionID: p.ID,
		Name:         p.Name,
		GuardName:    p.GuardName,
		Description:  p.Description,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
