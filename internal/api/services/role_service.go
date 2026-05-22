package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
)

var (
	ErrRoleNotFound = errors.New("role not found")
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
	log := logger.WithContext(ctx)
	log.Infof("RoleService.List: request_received")
	db := database.GetDB().WithContext(ctx).Model(&entities.Role{})
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.GuardName != nil && *req.GuardName != "" {
		db = db.Where("guard_name = ?", *req.GuardName)
	}
	var roles []entities.Role
	if err := db.Order("created_at DESC").Find(&roles).Error; err != nil {
		log.Warnf("RoleService.List: error finding roles err=%v", err)
		return nil, err
	}
	log.Infof("RoleService.List: roles found successfully count=%d", len(roles))
	return roles, nil
}

func (s *RoleService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	log := logger.WithContext(ctx)
	log.Infof("RoleService.FindByUUID: request_received")
	role, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("RoleService.FindByUUID: error finding role err=%v", err)
		return nil, err
	}
	if role == nil {
		log.Warnf("RoleService.FindByUUID: role not found")
		return nil, ErrRoleNotFound
	}
	return role, nil
}
