package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type ActualScoreRepo interface {
	BaseRepository[model.ActualScore]
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

func (r *actualScoreRepo) Get(id string) (*model.ActualScore, error) {
	var actualScore model.ActualScore
	err := r.db.Preload("Project").Preload("Employee").First(&actualScore, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &actualScore, nil
}

func (r *actualScoreRepo) List() ([]model.ActualScore, error) {
	var actualScores []model.ActualScore
	err := r.db.Preload("Project").Preload("Employee").Find(&actualScores).Error
	if err != nil {
		return nil, err
	}
	return actualScores, nil
}

func (r *actualScoreRepo) Delete(id string) error {
	result := r.db.Delete(&model.ActualScore{
		BaseModel: model.BaseModel{
			ID: id,
		},
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
