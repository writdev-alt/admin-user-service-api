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

var ErrRoleNotFound = errors.New("role not found")

type RoleService struct {
	repo         repositories.IRoleRepository
	permRepo     repositories.IPermissionRepository
	rolePermRepo repositories.IRolePermissionRepository
	base         repositories.IBaseRepository
}

var Role = &RoleService{
	repo:         repositories.Repo.Role,
	permRepo:     repositories.Repo.Permission,
	rolePermRepo: repositories.Repo.RolePermission,
	base:         repositories.Repo.Base,
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
	if err := s.attachPermissions(ctx, roles); err != nil {
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
	roles := []entities.Role{*role}
	if err := s.attachPermissions(ctx, roles); err != nil {
		return nil, err
	}
	role.Permissions = roles[0].Permissions
	return role, nil
}

func (s *RoleService) GetPermissions(ctx context.Context, roleUUID uuid.UUID) ([]entities.Permission, error) {
	log := logger.WithContext(ctx)
	log.Infof("RoleService.GetPermissions: request_received")
	role, err := s.repo.FindByUUID(ctx, roleUUID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	permsByRole, err := s.rolePermRepo.PermissionsByRoleIDs(ctx, []uint64{role.ID})
	if err != nil {
		return nil, err
	}
	return permsByRole[role.ID], nil
}

func (s *RoleService) SetPermissions(ctx context.Context, roleUUID uuid.UUID, permissionIDs []uint64) (*entities.Role, error) {
	log := logger.WithContext(ctx)
	log.Infof("RoleService.SetPermissions: request_received")
	role, err := s.repo.FindByUUID(ctx, roleUUID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrRoleNotFound
	}
	permissionIDs = dedupeUint64(permissionIDs)
	for _, id := range permissionIDs {
		perm, err := s.permRepo.FindByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if perm == nil {
			return nil, ErrPermissionNotFound
		}
	}
	if err := s.rolePermRepo.SetPermissions(ctx, role.ID, permissionIDs); err != nil {
		log.Warnf("RoleService.SetPermissions: error saving role permissions err=%v", err)
		return nil, err
	}
	roles := []entities.Role{*role}
	if err := s.attachPermissions(ctx, roles); err != nil {
		return nil, err
	}
	role.Permissions = roles[0].Permissions
	return role, nil
}

func (s *RoleService) attachPermissions(ctx context.Context, roles []entities.Role) error {
	if len(roles) == 0 {
		return nil
	}
	ids := make([]uint64, len(roles))
	for i := range roles {
		ids[i] = roles[i].ID
	}
	permsByRole, err := s.rolePermRepo.PermissionsByRoleIDs(ctx, ids)
	if err != nil {
		return err
	}
	for i := range roles {
		roles[i].Permissions = permsByRole[roles[i].ID]
	}
	return nil
}
