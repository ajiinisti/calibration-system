package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type GroupBusinessUnitUsecase interface {
	BaseUsecase[model.GroupBusinessUnit]
	FindByName(name string) (*model.GroupBusinessUnit, error)
}

type groupBusinessUnitUsecase struct {
	repo repository.GroupBusinessUnitRepo
}

func (r *groupBusinessUnitUsecase) FindAll() ([]model.GroupBusinessUnit, error) {
	return r.repo.List()
}

func (r *groupBusinessUnitUsecase) FindById(id string) (*model.GroupBusinessUnit, error) {
	return r.repo.Get(id)
}

func (r *groupBusinessUnitUsecase) FindByName(name string) (*model.GroupBusinessUnit, error) {
	return r.repo.GetByName(name)
}

func (r *groupBusinessUnitUsecase) SaveData(payload *model.GroupBusinessUnit) error {
	return r.repo.Save(payload)
}

func (r *groupBusinessUnitUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewGroupBusinessUnitUsecase(repo repository.GroupBusinessUnitRepo) GroupBusinessUnitUsecase {
	return &groupBusinessUnitUsecase{
		repo: repo,
	}
}
