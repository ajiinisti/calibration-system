package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type ProjectRepo interface {
	BaseRepository[model.Project]
}

type projectRepo struct {
	db *gorm.DB
}

func (r *projectRepo) Save(payload *model.Project) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *projectRepo) Get(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("ActualScores").
		Preload("ProjectPhases").
		Preload("ProjectPhases.Phase").
		Preload("ScoreDistributions").
		First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) List() ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Preload("ActualScores").
		Preload("ProjectPhases").
		Preload("ProjectPhases.Phase").
		Preload("ScoreDistributions").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) Delete(id string) error {
	result := r.db.Delete(&model.Project{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Project not found!")
	}
	return nil
}

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}
