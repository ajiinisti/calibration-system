package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type RoleUsecase interface {
	BaseUsecase[model.Role]
	FindByName(name string) (*model.Role, error)
}

type roleUsecase struct {
	repo repository.RoleRepo
}

func (r *roleUsecase) FindByName(name string) (*model.Role, error) {
	return r.repo.SearchByName(name)
}

func (r *roleUsecase) FindAll() ([]model.Role, error) {
	return r.repo.List()
}

func (r *roleUsecase) FindById(id string) (*model.Role, error) {
	return r.repo.Get(id)
}

func (r *roleUsecase) SaveData(payload *model.Role) error {
	return r.repo.Save(payload)
}

func (r *roleUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewRoleUsecase(repo repository.RoleRepo) RoleUsecase {
	return &roleUsecase{
		repo: repo,
	}
}
