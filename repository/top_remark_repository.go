package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"gorm.io/gorm"
)

type TopRemarkRepo interface {
	Save(payload *model.TopRemark) error
	BulkSave(payload []*model.TopRemark) error
	Get(projectID, employeeID, projectPhaseID string) ([]*model.TopRemark, error)
	GetByID(id string) (*model.TopRemark, error)
	List() ([]model.TopRemark, error)
	Delete(projectID, employeeID, projectPhaseID string) error
	BulkDelete(payload request.DeleteTopRemarks) error
}

type topRemarkRepo struct {
	db *gorm.DB
}

func (r *topRemarkRepo) Save(payload *model.TopRemark) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *topRemarkRepo) BulkSave(payload []*model.TopRemark) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, remarks := range payload {
		if remarks.ID != "" {
			topRemarks, err := r.GetByID(remarks.ID)
			if err != nil {
				return err
			}

			if topRemarks != nil {
				if topRemarks.EvidenceName != "" && remarks.EvidenceName == "" {
					remarks.EvidenceName = topRemarks.EvidenceName
					remarks.Evidence = topRemarks.Evidence
				}

			}
		}
		err := r.db.Save(&remarks)
		if err.Error != nil {
			tx.Rollback()
			return err.Error
		}

	}

	tx.Commit()
	return nil
}

func (r *topRemarkRepo) GetByID(id string) (*model.TopRemark, error) {
	var topRemark *model.TopRemark
	err := r.db.First(&topRemark, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return topRemark, nil
}

func (r *topRemarkRepo) Get(projectID, employeeID, projectPhaseID string) ([]*model.TopRemark, error) {
	var topRemark []*model.TopRemark
	err := r.db.Find(&topRemark, "project_id = ? AND employee_id = ? AND project_phase_id = ? ", projectID, employeeID, projectPhaseID).Error
	if err != nil {
		return nil, err
	}
	return topRemark, nil
}

func (r *topRemarkRepo) List() ([]model.TopRemark, error) {
	var topRemarks []model.TopRemark
	err := r.db.Find(&topRemarks).Error
	if err != nil {
		return nil, err
	}
	return topRemarks, nil
}

func (r *topRemarkRepo) Delete(projectID, employeeID, projectPhaseID string) error {
	result := r.db.Delete(&model.TopRemark{
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

func (r *topRemarkRepo) BulkDelete(payload request.DeleteTopRemarks) error {
	// fmt.Println("ALL ID:=", payload.IDs)
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, deleted := range payload.IDs {
		result := tx.Where("project_id = ? AND employee_id = ? AND project_phase_id = ?",
			deleted.ProjectID,
			deleted.EmployeeID,
			deleted.ProjectPhaseID).
			Delete(&model.TopRemark{})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Remark Settings not found!")
		}
	}

	tx.Commit()
	return nil
}

func NewTopRemarkRepo(db *gorm.DB) TopRemarkRepo {
	return &topRemarkRepo{
		db: db,
	}
}
