package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type CalibrationUsecase interface {
	BaseUsecase[model.Calibration]
}

type calibrationUsecase struct {
	repo         repository.CalibrationRepo
	user         UserUsecase
	project      ProjectUsecase
	projectPhase ProjectPhaseUsecase
}

func (r *calibrationUsecase) FindAll() ([]model.Calibration, error) {
	return r.repo.List()
}

func (r *calibrationUsecase) FindById(id string) (*model.Calibration, error) {
	return r.repo.Get(id)
}

func (r *calibrationUsecase) SaveData(payload *model.Calibration) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.ProjectPhaseID != "" {
		_, err := r.projectPhase.FindById(payload.ProjectPhaseID)
		if err != nil {
			return fmt.Errorf("Project Phase Not Found")
		}
	}

	if payload.EmployeeID != "" {
		_, err := r.user.FindById(payload.EmployeeID)
		if err != nil {
			return fmt.Errorf("Employee Not Found")
		}
	}

	if payload.CalibratorID != "" {
		_, err := r.user.FindById(payload.CalibratorID)
		if err != nil {
			return fmt.Errorf("Calibrator Not Found")
		}
	}

	if payload.SpmoID != "" {
		_, err := r.user.FindById(payload.SpmoID)
		if err != nil {
			return fmt.Errorf("SPMO Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *calibrationUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewCalibrationUsecase(repo repository.CalibrationRepo, user UserUsecase, project ProjectUsecase, projectPhase ProjectPhaseUsecase) CalibrationUsecase {
	return &calibrationUsecase{
		repo:         repo,
		user:         user,
		project:      project,
		projectPhase: projectPhase,
	}
}
