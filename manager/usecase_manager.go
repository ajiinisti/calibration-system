package manager

import (
	"calibration-system.com/config"
	"calibration-system.com/usecase"
)

type UsecaseManager interface {
	UserUc() usecase.UserUsecase
	RoleUc() usecase.RoleUsecase
	AuthUc() usecase.AuthUsecase
	EmployeeUc() usecase.EmployeeUsecase
	BusinessUnitUc() usecase.BusinessUnitUsecase
	GroupBusinessUnitUc() usecase.GroupBusinessUnitUsecase
}

type usecaseManager struct {
	repo RepoManager
	cfg  *config.Config
}

func (u *usecaseManager) RoleUc() usecase.RoleUsecase {
	return usecase.NewRoleUsecase(u.repo.RoleRepo())
}

func (u *usecaseManager) UserUc() usecase.UserUsecase {
	return usecase.NewUserUseCase(u.repo.UserRepo(), u.RoleUc(), u.EmployeeUc(), u.cfg)
}

func (u *usecaseManager) AuthUc() usecase.AuthUsecase {
	return usecase.NewAuthUsecase(u.UserUc())
}

func (u *usecaseManager) EmployeeUc() usecase.EmployeeUsecase {
	return usecase.NewEmployeeUsecase(u.repo.EmployeeRepo())
}

func (u *usecaseManager) GroupBusinessUnitUc() usecase.GroupBusinessUnitUsecase {
	return usecase.NewGroupBusinessUnitUsecase(u.repo.GroupBusinessUnitRepo())
}

func (u *usecaseManager) BusinessUnitUc() usecase.BusinessUnitUsecase {
	return usecase.NewBusinessUnitUsecase(u.repo.BusinessUnitRepo(), u.GroupBusinessUnitUc())
}

func NewUsecaseManager(repo RepoManager, cfg *config.Config) UsecaseManager {
	return &usecaseManager{
		repo: repo,
		cfg:  cfg,
	}
}
