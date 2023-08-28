package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type PhaseUsecase interface {
	BaseUsecase[model.Phase]
}

type phaseUsecase struct {
	repo repository.PhaseRepo
}

func (r *phaseUsecase) FindAll() ([]model.Phase, error) {
	return r.repo.List()
}

func (r *phaseUsecase) FindById(id string) (*model.Phase, error) {
	return r.repo.Get(id)
}

func (r *phaseUsecase) SaveData(payload *model.Phase) error {
	return r.repo.Save(payload)
}

func (r *phaseUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewPhaseUsecase(repo repository.PhaseRepo) PhaseUsecase {
	return &phaseUsecase{
		repo: repo,
	}
}
