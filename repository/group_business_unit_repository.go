package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type GroupBusinessUnitRepo interface {
	BaseRepository[model.GroupBusinessUnit]
}

type groupBusinessUnitRepo struct {
	db *gorm.DB
}

func (r *groupBusinessUnitRepo) Save(payload *model.GroupBusinessUnit) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *groupBusinessUnitRepo) Get(id string) (*model.GroupBusinessUnit, error) {
	var groupBusinessUnit model.GroupBusinessUnit
	err := r.db.First(&groupBusinessUnit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &groupBusinessUnit, nil
}

func (r *groupBusinessUnitRepo) List() ([]model.GroupBusinessUnit, error) {
	var groupBusinessUnits []model.GroupBusinessUnit
	err := r.db.Find(&groupBusinessUnits).Error
	if err != nil {
		return nil, err
	}
	return groupBusinessUnits, nil
}

func (r *groupBusinessUnitRepo) Delete(id string) error {
	result := r.db.Delete(&model.Employee{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Business Unit not found!")
	}
	return nil
}

func NewGroupBusinessUnitRepo(db *gorm.DB) GroupBusinessUnitRepo {
	return &groupBusinessUnitRepo{
		db: db,
	}
}
