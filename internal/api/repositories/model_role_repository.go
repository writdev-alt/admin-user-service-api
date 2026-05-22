package repositories

import (
	"context"
	"errors"

	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"gorm.io/gorm"
)

type IModelRoleRepository interface {
	repositories.IBaseRepository
	AssignRole(ctx context.Context, modelID uint64, roleID uint64) error
	ReplaceRole(ctx context.Context, modelID uint64, roleID uint64) error
	RoleByUserID(ctx context.Context, userID uint64) (*entities.Role, error)
	RolesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]*entities.Role, error)
}

type ModelRoleRepository struct {
	repositories.IBaseRepository
}

func NewModelRoleRepository() IModelRoleRepository {
	return &ModelRoleRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *ModelRoleRepository) AssignRole(ctx context.Context, modelID uint64, roleID uint64) error {
	log := logger.WithContext(ctx)
	log.Infof("ModelRoleRepository.AssignRole: request_received")
	modelRole := entities.ModelRole{
		ModelID:   modelID,
		ModelType: entities.ModelTypeUser,
		RoleID:    roleID,
	}
	if err := r.Save(ctx, &modelRole); err != nil {
		log.Warnf("ModelRoleRepository.AssignRole: error saving model role err=%v", err)
		return err
	}
	log.Infof("ModelRoleRepository.AssignRole: ok model_id=%d role_id=%d", modelID, roleID)
	return nil
}

// ReplaceRole sets the single role assignment for a model (removes prior rows first).
func (r *ModelRoleRepository) ReplaceRole(ctx context.Context, modelID uint64, roleID uint64) error {
	log := logger.WithContext(ctx)
	log.Infof("ModelRoleRepository.ReplaceRole: request_received")
	db := database.GetDB().WithContext(ctx)
	if err := db.Where("model_id = ? AND model_type = ?", modelID, entities.ModelTypeUser).Delete(&entities.ModelRole{}).Error; err != nil {
		log.Warnf("ModelRoleRepository.ReplaceRole: error deleting model role err=%v", err)
		return err
	}
	modelRole := entities.ModelRole{
		ModelID:   modelID,
		ModelType: entities.ModelTypeUser,
		RoleID:    roleID,
	}
	if err := r.Save(ctx, &modelRole); err != nil {
		log.Warnf("ModelRoleRepository.ReplaceRole: error saving model role err=%v", err)
		return err
	}
	log.Infof("ModelRoleRepository.ReplaceRole: model role saved successfully model_id=%d role_id=%d", modelID, roleID)
	return nil
}

func (r *ModelRoleRepository) RoleByUserID(ctx context.Context, userID uint64) (*entities.Role, error) {
	var role entities.Role
	err := database.GetDB().WithContext(ctx).
		Table("roles").
		Joins("INNER JOIN model_roles ON model_roles.role_id = roles.id").
		Where("model_roles.model_id = ? AND model_roles.model_type = ?", userID, entities.ModelTypeUser).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *ModelRoleRepository) RolesByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]*entities.Role, error) {
	out := make(map[uint64]*entities.Role)
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
	userToRoleID := make(map[uint64]uint64, len(links))
	for _, link := range links {
		userToRoleID[link.ModelID] = link.RoleID
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
	roleByID := make(map[uint64]*entities.Role, len(roles))
	for i := range roles {
		roleByID[roles[i].ID] = &roles[i]
	}
	for userID, roleID := range userToRoleID {
		if role, ok := roleByID[roleID]; ok {
			out[userID] = role
		}
	}
	return out, nil
}
