package repositories

import (
	"context"

	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IModelPermissionRepository interface {
	repositories.IBaseRepository
	SetPermissions(ctx context.Context, modelID uint64, permissionIDs []uint64) error
	PermissionsByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]entities.Permission, error)
}

type ModelPermissionRepository struct {
	repositories.IBaseRepository
}

func NewModelPermissionRepository() IModelPermissionRepository {
	return &ModelPermissionRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *ModelPermissionRepository) SetPermissions(ctx context.Context, modelID uint64, permissionIDs []uint64) error {
	log := logger.WithContext(ctx)
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("model_id = ? AND model_type = ?", modelID, entities.ModelTypeUser).
		Delete(&entities.ModelPermission{}).Error; err != nil {
		log.Warnf("ModelPermissionRepository.SetPermissions: delete failed err=%v", err)
		return err
	}
	for _, permissionID := range permissionIDs {
		row := entities.ModelPermission{
			ModelID:      modelID,
			ModelType:    entities.ModelTypeUser,
			PermissionID: permissionID,
		}
		if err := r.Save(ctx, &row); err != nil {
			log.Warnf("ModelPermissionRepository.SetPermissions: save failed err=%v", err)
			return err
		}
	}
	return nil
}

func (r *ModelPermissionRepository) PermissionsByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]entities.Permission, error) {
	out := make(map[uint64][]entities.Permission)
	if len(userIDs) == 0 {
		return out, nil
	}
	var links []entities.ModelPermission
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("model_type = ? AND model_id IN ?", entities.ModelTypeUser, userIDs).Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return out, nil
	}
	permissionIDSet := make(map[uint64]struct{})
	userToPermissionIDs := make(map[uint64][]uint64)
	for _, link := range links {
		userToPermissionIDs[link.ModelID] = append(userToPermissionIDs[link.ModelID], link.PermissionID)
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
	for userID, ids := range userToPermissionIDs {
		for _, pid := range ids {
			if p, ok := permissionByID[pid]; ok {
				out[userID] = append(out[userID], p)
			}
		}
	}
	return out, nil
}
