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

func (s *RoleService) List(ctx context.Context, req request.RoleListRequest) ([]entities.Role, int64, error) {
	pageNumber, pageSize := normalizePage(req.PageNumber, req.PageSize)
	db := database.GetDB().WithContext(ctx).Model(&entities.Role{})
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.GuardName != nil && *req.GuardName != "" {
		db = db.Where("guard_name = ?", *req.GuardName)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var roles []entities.Role
	offset := (pageNumber - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (s *RoleService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	return s.repo.FindByUUID(ctx, id)
}
