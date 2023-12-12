package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type ActualScoreRepo interface {
	Save(payload *model.ActualScore) error
	Get(projectId, employeeId string) (*model.ActualScore, error)
	List() ([]model.ActualScore, error)
	Delete(projectId, employeeId string) error
	Bulksave(payload *[]model.ActualScore) error
}

type actualScoreRepo struct {
	db *gorm.DB
}

func (r *actualScoreRepo) Save(payload *model.ActualScore) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *actualScoreRepo) Bulksave(payload *[]model.ActualScore) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		err := tx.Save(&currentBatch).Error
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := tx.Save(&remainingItems).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *actualScoreRepo) Get(projectId, employeeId string) (*model.ActualScore, error) {
	var actualScore model.ActualScore
	err := r.db.Preload("Project").Preload("Employee").Preload("Employee.BusinessUnit").First(&actualScore, "project_id = ? AND employee_id = ?", projectId, employeeId).Error
	if err != nil {
		return nil, err
	}
	return &actualScore, nil
}

func (r *actualScoreRepo) List() ([]model.ActualScore, error) {
	var actualScores []model.ActualScore
	err := r.db.Preload("Project").Preload("Employee").Preload("Employee.BusinessUnit").Find(&actualScores).Error
	if err != nil {
		return nil, err
	}
	return actualScores, nil
}

func (r *actualScoreRepo) Delete(projectId, employeeId string) error {
	result := r.db.Delete(&model.ActualScore{
		ProjectID:  projectId,
		EmployeeID: employeeId,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Actual Score not found!")
	}
	return nil
}

func NewActualScoreRepo(db *gorm.DB) ActualScoreRepo {
	return &actualScoreRepo{
		db: db,
	}
}
