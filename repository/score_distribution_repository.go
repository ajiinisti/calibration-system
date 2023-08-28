package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type ScoreDistributionRepo interface {
	BaseRepository[model.ScoreDistribution]
}

type scoreDistributionRepo struct {
	db *gorm.DB
}

func (r *scoreDistributionRepo) Save(payload *model.ScoreDistribution) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *scoreDistributionRepo) Get(id string) (*model.ScoreDistribution, error) {
	var scoreDistribution model.ScoreDistribution
	err := r.db.Preload("GroupBusinessUnit").Preload("Project").First(&scoreDistribution, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &scoreDistribution, nil
}

func (r *scoreDistributionRepo) List() ([]model.ScoreDistribution, error) {
	var scoreDistributions []model.ScoreDistribution
	err := r.db.Preload("GroupBusinessUnit").Preload("Project").Find(&scoreDistributions).Error
	if err != nil {
		return nil, err
	}
	return scoreDistributions, nil
}

func (r *scoreDistributionRepo) Delete(id string) error {
	result := r.db.Delete(&model.ScoreDistribution{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Score Distribution not found!")
	}
	return nil
}

func NewScoreDistributionRepo(db *gorm.DB) ScoreDistributionRepo {
	return &scoreDistributionRepo{
		db: db,
	}
}
