package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/pkg/database"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

type RoleService struct {
	repo repositories.IRoleRepository
	base repositories.IBaseRepository
}

var Role = &RoleService{
	repo: repositories.Repo.Role,
	base: repositories.Repo.Base,
}

func (s *RoleService) List(ctx context.Context, req request.RoleListRequest) ([]entities.Role, error) {
	db := database.GetDB().WithContext(ctx).Model(&entities.Role{})
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	db = db.Where("guard_name = ?", req.GuardName)
	var roles []entities.Role
	if err := db.Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	return s.repo.FindByUUID(ctx, id)
}
