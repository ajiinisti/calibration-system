package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type BottomRemarkRepo interface {
	Save(payload *model.BottomRemark, projectPhases []model.ProjectPhase) error
	Get(projectID, employeeID, projectPhaseID string) (*model.BottomRemark, error)
	List() ([]model.BottomRemark, error)
	Delete(projectID, employeeID, projectPhaseID string) error
}

type bottomRemarkRepo struct {
	db *gorm.DB
}

func (r *bottomRemarkRepo) Save(payload *model.BottomRemark, projectPhases []model.ProjectPhase) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Save(&payload)
	if err.Error != nil {
		return err.Error
	}

	err = tx.Model(&model.Calibration{}).
		Where("project_id = ? AND project_phase_id = ? AND employee_id = ?",
			payload.ProjectID,
			payload.ProjectPhaseID,
			payload.EmployeeID,
		).Updates(map[string]interface{}{
		"filled_top_bottom_mark": true,
	})
	if err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	for _, projectPhase := range projectPhases {
		var calibrations []*model.Calibration
		err := r.db.
			Table("calibrations c").
			Where("c.employee_id = ? AND c.project_id = ? AND c.project_phase_id = ?", payload.EmployeeID, payload.ProjectID, projectPhase.ID).
			Find(&calibrations).
			Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if len(calibrations) > 0 {
			payload.ProjectPhaseID = projectPhase.ID
			err := r.db.Save(&payload)
			if err.Error != nil {
				tx.Rollback()
				return err.Error
			}
		}

		for _, calibration := range calibrations {
			calibration.JustificationType = "bottom"
			err := tx.Save(calibration)
			if err.Error != nil {
				tx.Rollback()
				return err.Error
			}
		}

	}

	tx.Commit()

	go func() {
		err := r.db.Exec("REFRESH MATERIALIZED VIEW materialized_user_view;").Error
		if err != nil {
			fmt.Printf("Failed to refresh materialized view: %v", err)
		}
	}()
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
