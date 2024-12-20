package usecase

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type TopRemarkUsecase interface {
	FindAll() ([]model.TopRemark, error)
	FindById(id string) (*model.TopRemark, error)
	FindByForeignKeyID(projectID, employeeID, projectPhaseID string) ([]*response.TopRemarkResponse, error)
	SaveData(payload *model.TopRemark) error
	SaveDataByProject(payload []*model.TopRemark) error
	DeleteData(projectID, employeeID, projectPhaseID string) error
	BulkDeleteData(payload request.DeleteTopRemarks) error
}

type topRemarkUsecase struct {
	repo         repository.TopRemarkRepo
	project      ProjectUsecase
	employee     UserUsecase
	projectPhase ProjectPhaseUsecase
}

func (r *topRemarkUsecase) FindAll() ([]model.TopRemark, error) {
	return r.repo.List()
}

func (r *topRemarkUsecase) FindById(id string) (*model.TopRemark, error) {
	data, err := r.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("ERROR DI DB %s", err.Error()))
	}
	return data, nil
}

func (r *topRemarkUsecase) FindByForeignKeyID(projectID, employeeID, projectPhaseID string) ([]*response.TopRemarkResponse, error) {
	topRemarks, err := r.repo.Get(projectID, employeeID, projectPhaseID)
	if err != nil {
		return nil, err
	}

	var topRemarksResponse []*response.TopRemarkResponse
	for _, data := range topRemarks {
		topRemarksResponse = append(topRemarksResponse, &response.TopRemarkResponse{
			BaseModel:      data.BaseModel,
			ProjectID:      data.ProjectID,
			EmployeeID:     data.EmployeeID,
			ProjectPhaseID: data.ProjectPhaseID,
			Initiative:     data.Initiative,
			Description:    data.Description,
			Result:         data.Result,
			StartDate:      data.StartDate,
			EndDate:        data.EndDate,
			Comment:        data.Comment,
			EvidenceName:   data.EvidenceName,
			IsProject:      data.IsProject,
			IsInitiative:   data.IsInitiative,
		})
	}

	return topRemarksResponse, nil
}

func (r *topRemarkUsecase) SaveData(payload *model.TopRemark) error {
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
	return r.repo.Save(payload)
}

func (r *topRemarkUsecase) SaveDataByProject(payload []*model.TopRemark) error {
	projectPhases, err := r.projectPhase.FindAllActiveHigherThanID(payload[0].ProjectPhaseID)
	if err != nil {
		return err
	}
	return r.repo.BulkSave(payload, projectPhases)
}

func (r *topRemarkUsecase) DeleteData(projectID, employeeID, projectPhaseID string) error {
	return r.repo.Delete(projectID, employeeID, projectPhaseID)
}

func (r *topRemarkUsecase) BulkDeleteData(payload request.DeleteTopRemarks) error {
	return r.repo.BulkDelete(payload)
}

func NewTopRemarkUsecase(
	repo repository.TopRemarkRepo,
	project ProjectUsecase,
	employee UserUsecase,
	projectPhase ProjectPhaseUsecase,
) TopRemarkUsecase {
	return &topRemarkUsecase{
		repo:         repo,
		project:      project,
		employee:     employee,
		projectPhase: projectPhase,
	}
}
