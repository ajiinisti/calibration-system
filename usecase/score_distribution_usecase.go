package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type ScoreDistributionUsecase interface {
	BaseUsecase[model.ScoreDistribution]
}

type scoreDistributionUsecase struct {
	repo    repository.ScoreDistributionRepo
	gbu     GroupBusinessUnitUsecase
	project ProjectUsecase
}

func (r *scoreDistributionUsecase) FindAll() ([]model.ScoreDistribution, error) {
	return r.repo.List()
}

func (r *scoreDistributionUsecase) FindById(id string) (*model.ScoreDistribution, error) {
	return r.repo.Get(id)
}

func (r *scoreDistributionUsecase) SaveData(payload *model.ScoreDistribution) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.GroupBusinessUnitID != "" {
		_, err := r.gbu.FindById(payload.GroupBusinessUnitID)
		if err != nil {
			return fmt.Errorf("Business Unit Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *scoreDistributionUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewScoreDistributionUsecase(repo repository.ScoreDistributionRepo, gbu GroupBusinessUnitUsecase, project ProjectUsecase) ScoreDistributionUsecase {
	return &scoreDistributionUsecase{
		repo:    repo,
		gbu:     gbu,
		project: project,
	}
}
