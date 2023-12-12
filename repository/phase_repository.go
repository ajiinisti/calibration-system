package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PhaseRepo interface {
	BaseRepository[model.Phase]
}

type phaseRepo struct {
	db *gorm.DB
}

func (r *phaseRepo) Save(payload *model.Phase) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *phaseRepo) Get(id string) (*model.Phase, error) {
	var phase model.Phase
	err := r.db.First(&phase, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &phase, nil
}

func (r *phaseRepo) List() ([]model.Phase, error) {
	var phases []model.Phase
	err := r.db.Order(clause.OrderByColumn{Column: clause.Column{Name: "order"}, Desc: false}).Find(&phases).Error
	if err != nil {
		return nil, err
	}
	return phases, nil
}

func (r *phaseRepo) Delete(id string) error {
	result := r.db.Delete(&model.Phase{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Phase not found!")
	}
	return nil
}

func NewPhaseRepo(db *gorm.DB) PhaseRepo {
	return &phaseRepo{
		db: db,
	}
}
