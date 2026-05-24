package repositories

import (
	"context"

	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IRolePermissionRepository interface {
	repositories.IBaseRepository
	SetPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error
	PermissionsByRoleIDs(ctx context.Context, roleIDs []uint64) (map[uint64][]entities.Permission, error)
}

type RolePermissionRepository struct {
	repositories.IBaseRepository
}

func NewRolePermissionRepository() IRolePermissionRepository {
	return &RolePermissionRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *RolePermissionRepository) SetPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	log := logger.WithContext(ctx)
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("role_id = ?", roleID).Delete(&entities.RolePermission{}).Error; err != nil {
		log.Warnf("RolePermissionRepository.SetPermissions: delete failed err=%v", err)
		return err
	}
	for _, permissionID := range permissionIDs {
		row := entities.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		if err := r.Save(ctx, &row); err != nil {
			log.Warnf("RolePermissionRepository.SetPermissions: save failed err=%v", err)
			return err
		}
	}
	return nil
}

func (r *RolePermissionRepository) PermissionsByRoleIDs(ctx context.Context, roleIDs []uint64) (map[uint64][]entities.Permission, error) {
	out := make(map[uint64][]entities.Permission)
	if len(roleIDs) == 0 {
		return out, nil
	}
	var links []entities.RolePermission
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("role_id IN ?", roleIDs).Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return out, nil
	}
	permissionIDSet := make(map[uint64]struct{})
	roleToPermissionIDs := make(map[uint64][]uint64)
	for _, link := range links {
		roleToPermissionIDs[link.RoleID] = append(roleToPermissionIDs[link.RoleID], link.PermissionID)
		permissionIDSet[link.PermissionID] = struct{}{}
	}
	permissionIDs := make([]uint64, 0, len(permissionIDSet))
	for id := range permissionIDSet {
		permissionIDs = append(permissionIDs, id)
	}
	var permissions []entities.Permission
	if err := db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return nil, err
	}
	permissionByID := make(map[uint64]entities.Permission, len(permissions))
	for _, p := range permissions {
		permissionByID[p.ID] = p
	}
	for roleID, ids := range roleToPermissionIDs {
		for _, pid := range ids {
			if p, ok := permissionByID[pid]; ok {
				out[roleID] = append(out[roleID], p)
			}
		}
	}
	return out, nil
}
