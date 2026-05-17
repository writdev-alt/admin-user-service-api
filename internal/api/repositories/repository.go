package repositories

import pkgrepos "github.com/turahe/pkg/repositories"

type IBaseRepository = pkgrepos.IBaseRepository

func baseRepo() IBaseRepository {
	return pkgrepos.NewBaseRepository()
}

type Repository struct {
	Base       IBaseRepository
	User       IUserRepository
	Role       IRoleRepository
	Permission IPermissionRepository
}

func NewRepository() *Repository {
	return &Repository{
		Base:       baseRepo(),
		User:       NewUserRepository(),
		Role:       NewRoleRepository(),
		Permission: NewPermissionRepository(),
	}
}

var Repo = NewRepository()
