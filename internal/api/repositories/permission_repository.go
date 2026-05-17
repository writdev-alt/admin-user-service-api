package repositories

import (
	"context"

	"github.com/google/uuid"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/turahe/pkg/types"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IPermissionRepository interface {
	repositories.IBaseRepository
	FindByID(ctx context.Context, id uint64) (*entities.Permission, error)
	FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Permission, error)
}

type PermissionRepository struct {
	repositories.IBaseRepository
}

func NewPermissionRepository() IPermissionRepository {
	return &PermissionRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *PermissionRepository) FindByID(ctx context.Context, id uint64) (*entities.Permission, error) {
	var row entities.Permission
	notFound, err := r.First(ctx, &row, types.Conditions{"id": id})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &row, nil
}

func (r *PermissionRepository) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Permission, error) {
	var row entities.Permission
	notFound, err := r.First(ctx, &row, types.Conditions{"uuid": id.String()})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &row, nil
}
