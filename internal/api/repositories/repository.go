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
	ModelRole        IModelRoleRepository
	ModelPermission  IModelPermissionRepository
	RolePermission   IRolePermissionRepository
	Permission       IPermissionRepository
}

func NewRepository() *Repository {
	return &Repository{
		Base:            baseRepo(),
		Admin:           NewAdminRepository(),
		User:            NewUserRepository(),
		Role:            NewRoleRepository(),
		ModelRole:        NewModelRoleRepository(),
		ModelPermission:  NewModelPermissionRepository(),
		RolePermission:   NewRolePermissionRepository(),
		Permission:       NewPermissionRepository(),
	}
}

var Repo = NewRepository()
