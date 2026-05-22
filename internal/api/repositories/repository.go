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
	ModelRole  IModelRoleRepository
	Permission IPermissionRepository
}

func NewRepository() *Repository {
	return &Repository{
		Base:       baseRepo(),
		Admin:      NewAdminRepository(),
		User:       NewUserRepository(),
		Role:       NewRoleRepository(),
		ModelRole:  NewModelRoleRepository(),
		Permission: NewPermissionRepository(),
	}
}

var Repo = NewRepository()
