package repositories

import pkgrepos "github.com/turahe/pkg/repositories"

type IBaseRepository = pkgrepos.IBaseRepository

func baseRepo() IBaseRepository {
	return pkgrepos.NewBaseRepository()
}

type Repository struct {
	Base       IBaseRepository
	Admin      IAdminRepository
	User       IUserRepository
	Role       IRoleRepository
	Permission IPermissionRepository
}

func NewRepository() *Repository {
	return &Repository{
		Base:       baseRepo(),
		Admin:      NewAdminRepository(),
		User:       NewUserRepository(),
		Role:       NewRoleRepository(),
		Permission: NewPermissionRepository(),
	}
}

var Repo = NewRepository()
