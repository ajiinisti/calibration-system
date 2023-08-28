package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type ProjectPhaseRepo interface {
	BaseRepository[model.ProjectPhase]
}

type projectProjectPhaseRepo struct {
	db *gorm.DB
}

func (r *projectProjectPhaseRepo) Save(payload *model.ProjectPhase) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *projectProjectPhaseRepo) Get(id string) (*model.ProjectPhase, error) {
	var projectProjectPhase model.ProjectPhase
	err := r.db.Preload("Phase").Preload("Project").First(&projectProjectPhase, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &projectProjectPhase, nil
}

func (r *projectProjectPhaseRepo) List() ([]model.ProjectPhase, error) {
	var projectProjectPhases []model.ProjectPhase
	err := r.db.Preload("Phase").Preload("Project").Find(&projectProjectPhases).Error
	if err != nil {
		return nil, err
	}
	return projectProjectPhases, nil
}

func (r *projectProjectPhaseRepo) Delete(id string) error {
	result := r.db.Delete(&model.ProjectPhase{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Project Phase not found!")
	}
	return nil
}

func NewProjectPhaseRepo(db *gorm.DB) ProjectPhaseRepo {
	return &projectProjectPhaseRepo{
		db: db,
	}
}
