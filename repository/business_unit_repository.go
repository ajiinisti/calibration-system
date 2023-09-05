package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type BusinessUnitRepo interface {
	BaseRepository[model.BusinessUnit]
	Bulksave(payload *[]model.BusinessUnit) error
}

type businessUnitRepo struct {
	db *gorm.DB
}

func (r *businessUnitRepo) Save(payload *model.BusinessUnit) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *businessUnitRepo) Bulksave(payload *[]model.BusinessUnit) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *businessUnitRepo) Get(id string) (*model.BusinessUnit, error) {
	var businessUnit model.BusinessUnit
	err := r.db.First(&businessUnit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &businessUnit, nil
}

func (r *businessUnitRepo) List() ([]model.BusinessUnit, error) {
	var businessUnits []model.BusinessUnit
	err := r.db.Preload("GroupBusinessUnit").Find(&businessUnits).Error
	if err != nil {
		return nil, err
	}
	return businessUnits, nil
}

func (r *businessUnitRepo) Delete(id string) error {
	result := r.db.Delete(&model.BusinessUnit{
		ID: id,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Business Unit not found!")
	}
	return nil
}

func NewBusinessUnitRepo(db *gorm.DB) BusinessUnitRepo {
	return &businessUnitRepo{
		db: db,
	}
}
