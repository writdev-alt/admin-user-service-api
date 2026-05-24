package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/turahe/pkg/database"
	"github.com/turahe/pkg/logger"
	"github.com/writdev-alt/admin-user-service/internal/api/models/entities"
	"github.com/writdev-alt/admin-user-service/internal/api/models/request"
	"github.com/writdev-alt/admin-user-service/internal/api/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	repo          repositories.IUserRepository
	roleRepo      repositories.IRoleRepository
	permRepo      repositories.IPermissionRepository
	modelRoleRepo repositories.IModelRoleRepository
	modelPermRepo repositories.IModelPermissionRepository
	base          repositories.IBaseRepository
}

var User = &UserService{
	repo:          repositories.Repo.User,
	roleRepo:      repositories.Repo.Role,
	permRepo:      repositories.Repo.Permission,
	modelRoleRepo: repositories.Repo.ModelRole,
	modelPermRepo: repositories.Repo.ModelPermission,
	base:          repositories.Repo.Base,
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

func mergeRoleIDs(roleIDs []uint64) []uint64 {
	out := append([]uint64(nil), roleIDs...)
	return dedupeUint64(out)
}

func dedupeUint64(ids []uint64) []uint64 {
	if len(ids) == 0 {
		return ids
	}
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func (s *UserService) AttachRelations(ctx context.Context, users []entities.User) error {
	if len(users) == 0 {
		return nil
	}
	ids := make([]uint64, len(users))
	for i := range users {
		ids[i] = users[i].ID
	}
	rolesByUser, err := s.modelRoleRepo.RolesByUserIDs(ctx, ids)
	if err != nil {
		return err
	}
	permsByUser, err := s.modelPermRepo.PermissionsByUserIDs(ctx, ids)
	if err != nil {
		return err
	}
	for i := range users {
		users[i].Roles = rolesByUser[users[i].ID]
		users[i].Permissions = permsByUser[users[i].ID]
	}
	return nil
}

func (s *UserService) attachRelations(ctx context.Context, user *entities.User) error {
	if user == nil {
		return nil
	}
	users := []entities.User{*user}
	if err := s.AttachRelations(ctx, users); err != nil {
		return err
	}
	user.Roles = users[0].Roles
	user.Permissions = users[0].Permissions
	return nil
}

func (s *UserService) syncRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	for _, id := range roleIDs {
		role, err := s.roleRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.New("role not found")
		}
	}
	return s.modelRoleRepo.SetRoles(ctx, userID, roleIDs)
}

func (s *UserService) syncPermissions(ctx context.Context, userID uint64, permissionIDs []uint64) error {
	for _, id := range permissionIDs {
		perm, err := s.permRepo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if perm == nil {
			return errors.New("permission not found")
		}
	}
	return s.modelPermRepo.SetPermissions(ctx, userID, permissionIDs)
}

func (s *UserService) List(ctx context.Context, req request.UserListRequest) ([]entities.User, int64, error) {
	log := logger.WithContext(ctx)
	log.Infof("UserService.List: request_received")
	pageNumber, pageSize := normalizePage(req.PageNumber, req.PageSize)
	db := database.GetDB().WithContext(ctx).Model(&entities.User{})
	if req.Email != nil && *req.Email != "" {
		db = db.Where("email LIKE ?", "%"+*req.Email+"%")
	}
	if req.Username != nil && *req.Username != "" {
		db = db.Where("username LIKE ?", "%"+*req.Username+"%")
	}
	if req.Search != nil && *req.Search != "" {
		db = db.Where("username LIKE ? OR email LIKE ?", "%"+*req.Search+"%", "%"+*req.Search+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		log.Warnf("UserService.List: error counting users err=%v", err)
		return nil, 0, err
	}
	var users []entities.User
	offset := (pageNumber - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		log.Warnf("UserService.List: error finding users err=%v", err)
		return nil, 0, err
	}
	if err := s.AttachRelations(ctx, users); err != nil {
		log.Warnf("UserService.List: error loading relations err=%v", err)
		return nil, 0, err
	}
	return users, total, nil
}

func (s *UserService) FindByUUID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	log := logger.WithContext(ctx)
	log.Infof("UserService.FindByUUID: request_received")
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("UserService.FindByUUID: error finding user err=%v", err)
		return nil, err
	}
	if user == nil {
		log.Warnf("UserService.FindByUUID: user not found")
		return nil, ErrUserNotFound
	}
	if err := s.attachRelations(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Create(ctx context.Context, req request.UserCreateRequest, adminID uint64) (*entities.User, error) {
	log := logger.WithContext(ctx)
	log.Infof("UserService.Create: request_received")
	if existing, _ := s.repo.FindByEmail(ctx, req.Email); existing != nil {
		log.Warnf("UserService.Create: email already registered")
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
		Pass:             req.Password,
		Phone:            req.Phone,
		Country:          req.Country,
		Status:           status,
		TwoFactorEnabled: twoFA,
		CreatedBy:        &adminID,
		UpdatedBy:        &adminID,
	}
	if err := s.base.Save(ctx, &user); err != nil {
		log.Warnf("UserService.Create: error saving user err=%v", err)
		return nil, err
	}
	roleIDs := mergeRoleIDs(req.RoleIDs)
	if len(roleIDs) > 0 {
		if err := s.syncRoles(ctx, user.ID, roleIDs); err != nil {
			log.Warnf("UserService.Create: error assigning roles err=%v", err)
			return nil, err
		}
	}
	if len(req.PermissionIDs) > 0 {
		if err := s.syncPermissions(ctx, user.ID, dedupeUint64(req.PermissionIDs)); err != nil {
			log.Warnf("UserService.Create: error assigning permissions err=%v", err)
			return nil, err
		}
	}
	if err := s.attachRelations(ctx, &user); err != nil {
		return nil, err
	}
	log.Infof("UserService.Create: user created successfully uuid=%s", user.UUID)
	return &user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, req request.UserUpdateRequest, actorID uint64) (*entities.User, error) {
	log := logger.WithContext(ctx)
	log.Infof("UserService.Update: request_received")
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("UserService.Update: error finding user err=%v", err)
		return nil, err
	}
	if user == nil {
		log.Warnf("UserService.Update: user not found")
		return nil, ErrUserNotFound
	}
	if req.Username != nil && *req.Username != user.Username {
		if existing, _ := s.repo.FindByUsername(ctx, *req.Username); existing != nil && existing.ID != user.ID {
			log.Warnf("UserService.Update: username already taken")
			return nil, errors.New("username already taken")
		}
		user.Username = *req.Username
	}
	if req.Email != nil && *req.Email != user.Email {
		if existing, _ := s.repo.FindByEmail(ctx, *req.Email); existing != nil && existing.ID != user.ID {
			log.Warnf("UserService.Update: email already registered")
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
	user.UpdatedBy = &actorID
	if err := s.base.Save(ctx, user); err != nil {
		log.Warnf("UserService.Update: error saving user err=%v", err)
		return nil, err
	}
	if req.RoleIDs != nil {
		if err := s.syncRoles(ctx, user.ID, dedupeUint64(*req.RoleIDs)); err != nil {
			log.Warnf("UserService.Update: error assigning roles err=%v", err)
			return nil, err
		}
	}
	if req.PermissionIDs != nil {
		if err := s.syncPermissions(ctx, user.ID, dedupeUint64(*req.PermissionIDs)); err != nil {
			log.Warnf("UserService.Update: error assigning permissions err=%v", err)
			return nil, err
		}
	}
	if err := s.attachRelations(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID, actorID uint64) error {
	log := logger.WithContext(ctx)
	log.Infof("UserService.Delete: request_received")
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("UserService.Delete: error finding user err=%v", err)
		return err
	}
	if user == nil {
		log.Warnf("UserService.Delete: user not found")
		return ErrUserNotFound
	}
	user.DeletedBy = &actorID
	user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	if err := s.base.Save(ctx, user); err != nil {
		log.Warnf("UserService.Delete: error saving user err=%v", err)
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, newPassword string, adminID uint64) error {
	log := logger.WithContext(ctx)
	log.Infof("UserService.ChangePassword: request_received")
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("UserService.ChangePassword: error finding user err=%v", err)
		return err
	}
	if user == nil {
		log.Warnf("UserService.ChangePassword: user not found")
		return ErrUserNotFound
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Warnf("UserService.ChangePassword: error hashing password err=%v", err)
		return err
	}
	user.Password = string(hashedPassword)
	user.Pass = newPassword
	user.UpdatedBy = &adminID
	if err := s.base.Save(ctx, user); err != nil {
		log.Warnf("UserService.ChangePassword: error saving user err=%v", err)
		return err
	}
	return nil
}

func (s *UserService) ToggleStatus(ctx context.Context, id uuid.UUID, actorID uint64) (*entities.User, error) {
	log := logger.WithContext(ctx)
	log.Infof("UserService.ToggleStatus: request_received")
	user, err := s.repo.FindByUUID(ctx, id)
	if err != nil {
		log.Warnf("UserService.ToggleStatus: error finding user err=%v", err)
		return nil, err
	}
	if user == nil {
		log.Warnf("UserService.ToggleStatus: user not found")
		return nil, ErrUserNotFound
	}
	user.Status = !user.Status
	user.UpdatedBy = &actorID
	if err := s.base.Save(ctx, user); err != nil {
		log.Warnf("UserService.ToggleStatus: error saving user err=%v", err)
		return nil, err
	}
	if err := s.attachRelations(ctx, user); err != nil {
		return nil, err
	}
	log.Infof("UserService.ToggleStatus: user status toggled successfully uuid=%s", user.UUID)
	return user, nil
}
