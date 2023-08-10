package manager

import "calibration-system.com/repository"

type RepoManager interface {
	RoleRepo() repository.RoleRepo
	UserRepo() repository.UserRepo
	EmployeeRepo() repository.EmployeeRepo
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

func NewRepoManager(infra InfraManager) RepoManager {
	return &repoManager{
		infra: infra,
	}
}
