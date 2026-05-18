package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/turahe/pkg/database"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo repositories.IUserRepository
	base repositories.IBaseRepository
}

var User = &UserService{
	repo: repositories.Repo.User,
	base: repositories.Repo.Base,
}

func normalizePage(pageNumber, pageSize int) (int, int) {
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	return pageNumber, pageSize
}

func (s *UserService) List(ctx context.Context, req request.UserListRequest) ([]entities.User, int64, error) {
	pageNumber, pageSize := normalizePage(req.PageNumber, req.PageSize)
	db := database.GetDB().WithContext(ctx).Model(&entities.User{})
	if req.Email != nil && *req.Email != "" {
		db = db.Where("email LIKE ?", "%"+*req.Email+"%")
	}
	if req.Username != nil && *req.Username != "" {
		db = db.Where("username LIKE ?", "%"+*req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var users []entities.User
	offset := (pageNumber - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return s.repo.FindByUUID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, req request.UserCreateRequest, actorID uint64) (*entities.User, error) {
	if existing, _ := s.repo.FindByEmail(ctx, req.Email); existing != nil {
		return nil, errors.New("email already registered")
	}
	if existing, _ := s.repo.FindByUsername(ctx, req.Username); existing != nil {
		return nil, errors.New("username already taken")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := true
	if req.Status != nil {
		status = *req.Status
	}
	twoFA := false
	if req.TwoFactorEnabled != nil {
		twoFA = *req.TwoFactorEnabled
	}
	user := entities.User{
		Username:         req.Username,
		Email:            req.Email,
		Password:         string(hashed),
		Phone:            req.Phone,
		Country:          req.Country,
		Status:           status,
		TwoFactorEnabled: twoFA,
		CreatedBy:        actorID,
		UpdatedBy:        actorID,
	}
	if err := s.base.Save(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req request.UserUpdateRequest, actorID uint64) (*entities.User, error) {
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}
	if req.Username != nil && *req.Username != user.Username {
		if existing, _ := s.repo.FindByUsername(ctx, *req.Username); existing != nil && existing.ID != user.ID {
			return nil, errors.New("username already taken")
		}
		user.Username = *req.Username
	}
	if req.Email != nil && *req.Email != user.Email {
		if existing, _ := s.repo.FindByEmail(ctx, *req.Email); existing != nil && existing.ID != user.ID {
			return nil, errors.New("email already registered")
		}
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Country != nil {
		user.Country = req.Country
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.TwoFactorEnabled != nil {
		user.TwoFactorEnabled = *req.TwoFactorEnabled
	}
	user.UpdatedBy = actorID
	if err := s.base.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID, actorID uint64) error {
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.DeletedBy = actorID
	if err := s.base.Save(ctx, user); err != nil {
		return err
	}
	return database.GetDB().WithContext(ctx).Delete(user).Error
}

func (s *UserService) ToggleStatus(ctx context.Context, id uuid.UUID, actorID uint64) (*entities.User, error) {
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.Status = !user.Status
	user.UpdatedBy = actorID
	if err := s.base.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
