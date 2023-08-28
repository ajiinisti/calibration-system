package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type ProjectPhaseUsecase interface {
	BaseUsecase[model.ProjectPhase]
}

type projectPhaseUsecase struct {
	repo    repository.ProjectPhaseRepo
	phase   PhaseUsecase
	project ProjectUsecase
}

func (r *projectPhaseUsecase) FindAll() ([]model.ProjectPhase, error) {
	return r.repo.List()
}

func (r *projectPhaseUsecase) FindById(id string) (*model.ProjectPhase, error) {
	return r.repo.Get(id)
}

func (r *projectPhaseUsecase) SaveData(payload *model.ProjectPhase) error {
	if payload.PhaseID != "" {
		_, err := r.phase.FindById(payload.PhaseID)
		if err != nil {
			return fmt.Errorf("Phase Not Found")
		}
	}

	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *projectPhaseUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewProjectPhaseUsecase(repo repository.ProjectPhaseRepo, phase PhaseUsecase) ProjectPhaseUsecase {
	return &projectPhaseUsecase{
		repo:  repo,
		phase: phase,
	}
}
