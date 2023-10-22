package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type BottomRemarkRepo interface {
	Save(payload *model.BottomRemark) error
	Get(projectID, employeeID, projectPhaseID string) (*model.BottomRemark, error)
	List() ([]model.BottomRemark, error)
	Delete(projectID, employeeID, projectPhaseID string) error
}

type bottomRemarkRepo struct {
	db *gorm.DB
}

func (r *bottomRemarkRepo) Save(payload *model.BottomRemark) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *bottomRemarkRepo) Get(projectID, employeeID, projectPhaseID string) (*model.BottomRemark, error) {
	var bottomRemark *model.BottomRemark
	err := r.db.Find(&bottomRemark, "project_id = ? AND employee_id = ? AND project_phase_id = ? ", projectID, employeeID, projectPhaseID).Error
	if err != nil {
		return nil, err
	}
	return bottomRemark, nil
}

func (r *bottomRemarkRepo) List() ([]model.BottomRemark, error) {
	var bottomRemarks []model.BottomRemark
	err := r.db.Find(&bottomRemarks).Error
	if err != nil {
		return nil, err
	}
	return bottomRemarks, nil
}

func (r *bottomRemarkRepo) Delete(projectID, employeeID, projectPhaseID string) error {
	result := r.db.Delete(&model.BottomRemark{
		ProjectID:      projectID,
		EmployeeID:     employeeID,
		ProjectPhaseID: projectPhaseID,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Remark Settings not found!")
	}
	return nil
}

func NewBottomRemarkRepo(db *gorm.DB) BottomRemarkRepo {
	return &bottomRemarkRepo{
		db: db,
	}
}
