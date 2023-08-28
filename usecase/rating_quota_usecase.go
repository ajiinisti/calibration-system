package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type RatingQuotaUsecase interface {
	BaseUsecase[model.RatingQuota]
}

type ratingQuotaUsecase struct {
	repo         repository.RatingQuotaRepo
	businessUnit BusinessUnitUsecase
	project      ProjectUsecase
}

func (r *ratingQuotaUsecase) FindAll() ([]model.RatingQuota, error) {
	return r.repo.List()
}

func (r *ratingQuotaUsecase) FindById(id string) (*model.RatingQuota, error) {
	return r.repo.Get(id)
}

func (r *ratingQuotaUsecase) SaveData(payload *model.RatingQuota) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.BusinessUnitID != "" {
		_, err := r.businessUnit.FindById(payload.BusinessUnitID)
		if err != nil {
			return fmt.Errorf("Business Unit Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *ratingQuotaUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewRatingQuotaUsecase(repo repository.RatingQuotaRepo, businessUnit BusinessUnitUsecase, project ProjectUsecase) RatingQuotaUsecase {
	return &ratingQuotaUsecase{
		repo:         repo,
		businessUnit: businessUnit,
		project:      project,
	}
}
