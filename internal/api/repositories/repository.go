package repositories

import pkgrepos "github.com/turahe/pkg/repositories"

type IBaseRepository = pkgrepos.IBaseRepository

func baseRepo() IBaseRepository {
	return pkgrepos.NewBaseRepository()
}

type Repository struct {
	Base            IBaseRepository
	Admin           IAdminRepository
	User            IUserRepository
	Role            IRoleRepository
	ModelRole       IModelRoleRepository
	ModelPermission IModelPermissionRepository
	Permission      IPermissionRepository
}

func NewRepository() *Repository {
	return &Repository{
		Base:            baseRepo(),
		Admin:           NewAdminRepository(),
		User:            NewUserRepository(),
		Role:            NewRoleRepository(),
		ModelRole:       NewModelRoleRepository(),
		ModelPermission: NewModelPermissionRepository(),
		Permission:      NewPermissionRepository(),
	}
}

var Repo = NewRepository()
