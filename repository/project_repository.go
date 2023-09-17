package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type ProjectRepo interface {
	BaseRepository[model.Project]
	PaginateList(pagination model.PaginationQuery) ([]model.Project, response.Paging, error)
	GetTotalRows() (int, error)
	ActivateByID(id string) error
	DeactivateAllExceptID(id string) error
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
		Preload("ScoreDistributions.GroupBusinessUnit").
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
		Preload("ScoreDistributions.GroupBusinessUnit").
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

func (r *projectRepo) PaginateList(pagination model.PaginationQuery) ([]model.Project, response.Paging, error) {
	var projects []model.Project
	err := r.db.
		Preload("ActualScores").
		Preload("ProjectPhases").
		Preload("ProjectPhases.Phase").
		Preload("ScoreDistributions").
		Preload("ScoreDistributions.GroupBusinessUnit").
		Limit(pagination.Take).Offset(pagination.Skip).Find(&projects).Error
	if err != nil {
		return nil, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRows()
	if err != nil {
		return nil, response.Paging{}, err
	}

	return projects, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetTotalRows() (int, error) {
	var count int64
	err := r.db.Model(&model.Project{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *projectRepo) ActivateByID(id string) error {
	result := r.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]interface{}{"active": true})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *projectRepo) DeactivateAllExceptID(id string) error {
	// Update all rows where 'id' is not equal to the specified 'id'
	result := r.db.Model(&model.Project{}).Where("id <> ?", id).Update("active", false)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}
