package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/turahe/pkg/logger"
	repositories "github.com/turahe/pkg/repositories"
	"github.com/turahe/pkg/types"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
)

type IUserRepository interface {
	repositories.IBaseRepository
	FindByID(ctx context.Context, id uint64) (*entities.User, error)
	FindByUUID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
}

type UserRepository struct {
	repositories.IBaseRepository
}

func NewUserRepository() IUserRepository {
	return &UserRepository{IBaseRepository: repositories.NewBaseRepository()}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*entities.User, error) {
	return r.findBy(ctx, "id", id)
}

func (r *UserRepository) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return r.findBy(ctx, "uuid", id.String())
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return r.findBy(ctx, "email", email)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	return r.findBy(ctx, "username", username)
}

func (r *UserRepository) findBy(ctx context.Context, key string, value interface{}) (*entities.User, error) {
	var user entities.User
	notFound, err := r.First(ctx, &user, types.Conditions{key: value})
	if err != nil {
		logger.WithContext(ctx).Errorf("UserRepository.findBy: %v", err)
		return nil, err
	}
	if notFound {
		return nil, nil
	}
	return &user, nil
}
