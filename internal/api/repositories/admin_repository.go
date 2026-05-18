package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/pkg/types"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IAdminRepository interface {
	IBaseRepository
	FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Admin, error)
}

type adminRepository struct {
	IBaseRepository
}

func NewAdminRepository() IAdminRepository {
	return &adminRepository{IBaseRepository: baseRepo()}
}

func (r *adminRepository) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.Admin, error) {
	var admin entities.Admin
	notFound, err := r.First(ctx, &admin, types.Conditions{"uuid": id.String()})
	if err != nil {
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &admin, nil
}
