package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type BusinessUnitUsecase interface {
	BaseUsecase[model.BusinessUnit]
}

type businessUnitUsecase struct {
	repo    repository.BusinessUnitRepo
	groupBu GroupBusinessUnitUsecase
	// groupBu repository.
}

func (r *businessUnitUsecase) FindAll() ([]model.BusinessUnit, error) {
	return r.repo.List()
}

func (r *businessUnitUsecase) FindById(id string) (*model.BusinessUnit, error) {
	return r.repo.Get(id)
}

func (r *businessUnitUsecase) SaveData(payload *model.BusinessUnit) error {
	if payload.GroupBusinessUnitId != "" {
		_, err := r.groupBu.FindById(payload.GroupBusinessUnitId)
		if err != nil {
			return fmt.Errorf("Group Business Unit Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *businessUnitUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewBusinessUnitUsecase(repo repository.BusinessUnitRepo, groupBu GroupBusinessUnitUsecase) BusinessUnitUsecase {
	return &businessUnitUsecase{
		repo:    repo,
		groupBu: groupBu,
	}
}
