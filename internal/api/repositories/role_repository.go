package repositories

import (
	"context"

	"github.com/google/uuid"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/turahe/pkg/types"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IRoleRepository interface {
	repositories.IBaseRepository
	FindByID(ctx context.Context, id uint64) (*entities.Role, error)
	FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Role, error)
	FindByName(ctx context.Context, name string) (*entities.Role, error)
}

type RoleRepository struct {
	repositories.IBaseRepository
}

func NewRoleRepository() IRoleRepository {
	return &RoleRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *RoleRepository) FindByID(ctx context.Context, id uint64) (*entities.Role, error) {
	var role entities.Role
	notFound, err := r.First(ctx, &role, types.Conditions{"id": id})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &role, nil
}

func (r *RoleRepository) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Role, error) {
	var role entities.Role
	notFound, err := r.First(ctx, &role, types.Conditions{"uuid": id.String()})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &role, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	notFound, err := r.First(ctx, &role, types.Conditions{"name": name})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &role, nil
}
