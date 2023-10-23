package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"gorm.io/gorm"
)

type CalibrationRepo interface {
	Save(payload *model.Calibration) error
	Get(id string) (*model.Calibration, error)
	List() ([]model.Calibration, error)
	GetActiveBySPMOID(spmoID string) ([]model.Calibration, error)
	Delete(projectId, projectPhaseId, employeeId string) error
	Bulksave(payload *[]model.Calibration) error
	BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) error
	SaveChanges(payload *request.CalibrationRequest) error
}

type calibrationRepo struct {
	db *gorm.DB
}

func (r *calibrationRepo) Save(payload *model.Calibration) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *calibrationRepo) Bulksave(payload *[]model.Calibration) error {
	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		return r.db.Save(&currentBatch).Error

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := r.db.Save(&remainingItems)
		if err != nil {
			return r.db.Save(&remainingItems).Error
		}
	}
	return nil
}

func (r *calibrationRepo) Get(id string) (*model.Calibration, error) {
	var calibration model.Calibration
	err := r.db.Preload("Project").Preload("Employee").Preload("ProjectPhase").Preload("Calibrator").First(&calibration, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &calibration, nil
}

func (r *calibrationRepo) GetActiveBySPMOID(id string) ([]model.Calibration, error) {
	var calibration []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("Employee").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Preload("BottomRemark").
		Preload("TopRemarks").
		Select("c.*").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id AND pp.review_spmo = true").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("spmo_id", id).
		// Where("c.spmo_id AND c.status = 'Waiting' ", id).
		Order("p.order ASC").
		// Limit(10).Offset(0).
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *calibrationRepo) List() ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.Preload("Project").Preload("Employee").Preload("ProjectPhase").Preload("Calibrator").Find(&calibrations).Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) Delete(projectId, projectPhaseId, employeeId string) error {
	result := r.db.Delete(&model.Calibration{
		ProjectID:      projectId,
		ProjectPhaseID: projectPhaseId,
		EmployeeID:     employeeId,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Calibration not found!")
	}
	return nil
}

func (r *calibrationRepo) BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reviewSPMO := false
	var employeeCalibrationScore []model.Calibration
	for _, calibrations := range payload.RequestData {
		if projectPhase.ReviewSpmo == true {
			calibrations.SpmoStatus = "Waiting"
			reviewSPMO = true
		}
		calibrations.Status = "Complete"
		employeeCalibrationScore = append(employeeCalibrationScore, calibrations)

		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Calibrations not found!")
		}
	}
	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		for _, c := range calibrations {
			c.CalibrationScore = employeeCalibration.CalibrationScore
			c.CalibrationRating = employeeCalibration.CalibrationRating
			c.Status = "Waiting"
			err := tx.Updates(&c).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		if !reviewSPMO {
			calibrations[0].Status = "Scoring"
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) SaveChanges(payload *request.CalibrationRequest) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, calibrations := range payload.RequestData {
		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Calibrations not found!")
		}
	}

	tx.Commit()
	return nil
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
