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
	Delete(projectId, projectPhaseId, employeeId string) error
	Bulksave(payload *[]model.Calibration) error
	BulkUpdate(payload *request.CalibrationRequest, phaseOrder int) error
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

func (r *calibrationRepo) BulkUpdate(payload *request.CalibrationRequest, phaseOrder int) error {
	var employeeId []string
	var employeeCalibrationScore []model.Calibration
	for _, calibrations := range payload.RequestData {
		employeeId = append(employeeId, calibrations.EmployeeID)
		employeeCalibrationScore = append(employeeCalibrationScore, calibrations)
	}

	err := r.db.Save(employeeCalibrationScore).Error
	if err != nil {
		return err
	}

	for _, employeeCalibration := range employeeCalibrationScore {
		err := r.db.Model(&model.Calibration{}).
			Joins("ProjectPhase").
			Joins("Phase").
			Joins("Project").
			Where("projects.active = ? AND phases.order > ?", true, phaseOrder).
			Where("employee_id = ?", employeeCalibration.EmployeeID).
			Update("calibration_score", employeeCalibration.CalibrationScore).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *calibrationRepo) SaveChanges(payload *request.CalibrationRequest) error {
	var employeeId []string
	var employeeCalibrationScore []model.Calibration
	for _, calibrations := range payload.RequestData {
		employeeId = append(employeeId, calibrations.EmployeeID)
		employeeCalibrationScore = append(employeeCalibrationScore, calibrations)
	}

	return r.db.Save(employeeCalibrationScore).Error
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
