package manager

import "calibration-system.com/repository"

type RepoManager interface {
	RoleRepo() repository.RoleRepo
	UserRepo() repository.UserRepo
	EmployeeRepo() repository.EmployeeRepo
	BusinessUnitRepo() repository.BusinessUnitRepo
	GroupBusinessUnitRepo() repository.GroupBusinessUnitRepo
}

type repoManager struct {
	infra InfraManager
}

func (r *repoManager) RoleRepo() repository.RoleRepo {
	return repository.NewRoleRepo(r.infra.Conn())
}

func (r *repoManager) UserRepo() repository.UserRepo {
	return repository.NewUserRepo(r.infra.Conn())
}

func (r *repoManager) EmployeeRepo() repository.EmployeeRepo {
	return repository.NewEmployeeRepo(r.infra.Conn())
}

func (r *repoManager) BusinessUnitRepo() repository.BusinessUnitRepo {
	return repository.NewBusinessUnitRepo(r.infra.Conn())
}

func (r *repoManager) GroupBusinessUnitRepo() repository.GroupBusinessUnitRepo {
	return repository.NewGroupBusinessUnitRepo(r.infra.Conn())
}

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{
		infra: infra,
	}
}
