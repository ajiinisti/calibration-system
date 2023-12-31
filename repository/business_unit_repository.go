package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type BusinessUnitRepo interface {
	BaseRepository[model.BusinessUnit]
	Bulksave(payload *[]model.BusinessUnit) error
	PaginateList(pagination model.PaginationQuery) ([]model.BusinessUnit, response.Paging, error)
	GetTotalRows(name string) (int, error)
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

func (r *businessUnitRepo) PaginateList(pagination model.PaginationQuery) ([]model.BusinessUnit, response.Paging, error) {
	var businessUnits []model.BusinessUnit
	var err error

	if pagination.Name == "" {
		err = r.db.
			Preload("GroupBusinessUnit").
			Limit(pagination.Take).Offset(pagination.Skip).
			Find(&businessUnits).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	} else {
		err = r.db.
			Preload("GroupBusinessUnit").
			Where("name ILIKE ?", "%"+pagination.Name+"%").
			Limit(pagination.Take).Offset(pagination.Skip).
			Find(&businessUnits).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	}

	totalRows, err := r.GetTotalRows(pagination.Name)
	if err != nil {
		return nil, response.Paging{}, err
	}

	return businessUnits, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *businessUnitRepo) GetTotalRows(name string) (int, error) {
	var count int64
	var err error
	if name == "" {
		err = r.db.
			Model(&model.BusinessUnit{}).
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	} else {
		err = r.db.
			Model(&model.BusinessUnit{}).
			Where("name ILIKE ?", "%"+name+"%").
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	}
	return int(count), nil
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
