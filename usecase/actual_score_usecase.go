package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type ActualScoreUsecase interface {
	BaseUsecase[model.ActualScore]
}

type actualScoreUsecase struct {
	repo     repository.ActualScoreRepo
	employee UserUsecase
	project  ProjectUsecase
}

func (r *actualScoreUsecase) FindAll() ([]model.ActualScore, error) {
	return r.repo.List()
}

func (r *actualScoreUsecase) FindById(id string) (*model.ActualScore, error) {
	return r.repo.Get(id)
}

func (r *actualScoreUsecase) SaveData(payload *model.ActualScore) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.EmployeeID != "" {
		_, err := r.employee.FindById(payload.EmployeeID)
		if err != nil {
			return fmt.Errorf("Employee Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *actualScoreUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewActualScoreUsecase(repo repository.ActualScoreRepo, employee UserUsecase, project ProjectUsecase) ActualScoreUsecase {
	return &actualScoreUsecase{
		repo:     repo,
		employee: employee,
		project:  project,
	}
}
