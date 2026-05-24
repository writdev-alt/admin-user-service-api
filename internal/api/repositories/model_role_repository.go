package repositories

import (
	"context"

	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IModelRoleRepository interface {
	repositories.IBaseRepository
	SetRoles(ctx context.Context, modelID uint64, roleIDs []uint64) error
	RolesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]entities.Role, error)
}

type ModelRoleRepository struct {
	repositories.IBaseRepository
}

func NewModelRoleRepository() IModelRoleRepository {
	return &ModelRoleRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *ModelRoleRepository) SetRoles(ctx context.Context, modelID uint64, roleIDs []uint64) error {
	log := logger.WithContext(ctx)
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("model_id = ? AND model_type = ?", modelID, entities.ModelTypeUser).Delete(&entities.ModelRole{}).Error; err != nil {
		log.Warnf("ModelRoleRepository.SetRoles: delete failed err=%v", err)
		return err
	}
	for _, roleID := range roleIDs {
		row := entities.ModelRole{
			ModelID:   modelID,
			ModelType: entities.ModelTypeUser,
			RoleID:    roleID,
		}
		if err := r.Save(ctx, &row); err != nil {
			log.Warnf("ModelRoleRepository.SetRoles: save failed err=%v", err)
			return err
		}
	}
	return nil
}

func (r *ModelRoleRepository) RolesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64][]entities.Role, error) {
	out := make(map[uint64][]entities.Role)
	if len(userIDs) == 0 {
		return out, nil
	}
	var links []entities.ModelRole
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("model_type = ? AND model_id IN ?", entities.ModelTypeUser, userIDs).Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return out, nil
	}
	roleIDSet := make(map[uint64]struct{})
	userToRoleIDs := make(map[uint64][]uint64)
	for _, link := range links {
		userToRoleIDs[link.ModelID] = append(userToRoleIDs[link.ModelID], link.RoleID)
		roleIDSet[link.RoleID] = struct{}{}
	}
	roleIDs := make([]uint64, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}
	var roles []entities.Role
	if err := db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	roleByID := make(map[uint64]entities.Role, len(roles))
	for _, role := range roles {
		roleByID[role.ID] = role
	}
	for userID, roleIDs := range userToRoleIDs {
		for _, roleID := range roleIDs {
			if role, ok := roleByID[roleID]; ok {
				out[userID] = append(out[userID], role)
			}
		}
	}
	return out, nil
}
