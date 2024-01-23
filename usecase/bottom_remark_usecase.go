package usecase

import (
	"fmt"

	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type BottomRemarkUsecase interface {
	FindAll() ([]model.BottomRemark, error)
	FindByForeignKeyID(projectID, employeeID, projectPhaseID string) (*model.BottomRemark, error)
	SaveData(payload *model.BottomRemark) error
	DeleteData(projectID, employeeID, projectPhaseID string) error
}

type bottomRemarkUsecase struct {
	repo         repository.BottomRemarkRepo
	project      ProjectUsecase
	employee     UserUsecase
	projectPhase ProjectPhaseUsecase
}

func (r *bottomRemarkUsecase) FindAll() ([]model.BottomRemark, error) {
	return r.repo.List()
}

func (r *bottomRemarkUsecase) FindByForeignKeyID(projectID, employeeID, projectPhaseID string) (*model.BottomRemark, error) {
	return r.repo.Get(projectID, employeeID, projectPhaseID)
}

func (r *bottomRemarkUsecase) SaveData(payload *model.BottomRemark) error {
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

	if payload.ProjectPhaseID != "" {
		_, err := r.projectPhase.FindById(payload.ProjectPhaseID)
		if err != nil {
			return fmt.Errorf("Project Phase Not Found")
		}
	}

	projectPhases, err := r.projectPhase.FindAllActiveHigherThanID(payload.ProjectPhaseID)
	if err != nil {
		return err
	}

	return r.repo.Save(payload, projectPhases)
}

func (r *bottomRemarkUsecase) DeleteData(projectID, employeeID, projectPhaseID string) error {
	return r.repo.Delete(projectID, employeeID, projectPhaseID)
}

func NewBottomRemarkUsecase(
	repo repository.BottomRemarkRepo,
	project ProjectUsecase,
	employee UserUsecase,
	projectPhase ProjectPhaseUsecase,
) BottomRemarkUsecase {
	return &bottomRemarkUsecase{
		repo:         repo,
		project:      project,
		employee:     employee,
		projectPhase: projectPhase,
	}
}
