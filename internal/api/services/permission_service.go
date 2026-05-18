package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/pkg/database"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

type PermissionService struct {
	repo repositories.IPermissionRepository
	base repositories.IBaseRepository
}

var Permission = &PermissionService{
	repo: repositories.Repo.Permission,
	base: repositories.Repo.Base,
}

func (s *PermissionService) List(ctx context.Context, req request.PermissionListRequest) ([]entities.Permission, error) {
	db := database.GetDB().WithContext(ctx).Model(&entities.Permission{})
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.GuardName != nil && *req.GuardName != "" {
		db = db.Where("guard_name = ?", *req.GuardName)
	}
	var rows []entities.Permission
	if err := db.Order("created_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *PermissionService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	return s.repo.FindByUUID(ctx, id)
}
