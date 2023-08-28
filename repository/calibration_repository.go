package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type CalibrationRepo interface {
	BaseRepository[model.Calibration]
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

func (r *calibrationRepo) Delete(id string) error {
	result := r.db.Delete(&model.Calibration{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Calibration not found!")
	}
	return nil
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
