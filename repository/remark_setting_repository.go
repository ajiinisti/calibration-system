package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type RemarkSettingRepo interface {
	Save(payload *model.RemarkSetting) error
	BulkSave(payload []*model.RemarkSetting) error
	Get(id string) ([]*model.RemarkSetting, error)
	List() ([]model.RemarkSetting, error)
	Delete(id string) error
	BulkDelete(payload request.DeleteRemark) error
	PaginateList(pagination model.PaginationQuery, id string) ([]model.RemarkSetting, response.Paging, error)
	GetTotalRows() (int, error)
}

type remarkSettingRepo struct {
	db *gorm.DB
}

func (r *remarkSettingRepo) Save(payload *model.RemarkSetting) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *remarkSettingRepo) BulkSave(payload []*model.RemarkSetting) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *remarkSettingRepo) Get(id string) ([]*model.RemarkSetting, error) {
	var remarkSetting []*model.RemarkSetting
	err := r.db.Preload("Project").Find(&remarkSetting, "project_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return remarkSetting, nil
}

func (r *remarkSettingRepo) List() ([]model.RemarkSetting, error) {
	var remarkSettings []model.RemarkSetting
	err := r.db.Preload("Project").Find(&remarkSettings).Error
	if err != nil {
		return nil, err
	}
	return remarkSettings, nil
}

func (r *remarkSettingRepo) Delete(id string) error {
	result := r.db.Delete(&model.RemarkSetting{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Remark Settings not found!")
	}
	return nil
}

func (r *remarkSettingRepo) BulkDelete(payload request.DeleteRemark) error {
	// fmt.Println("ALL ID:=", payload.IDs)
	result := r.db.Where("id IN (?)", payload.IDs).Delete(&model.RemarkSetting{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Remark Settings not found!")
	}
	return nil
}

func (r *remarkSettingRepo) PaginateList(pagination model.PaginationQuery, id string) ([]model.RemarkSetting, response.Paging, error) {
	var remarkSetting []model.RemarkSetting
	err := r.db.Preload("Project").Limit(pagination.Take).Offset(pagination.Skip).Find(&remarkSetting, "project_id = ?", id).Error
	if err != nil {
		return nil, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRows()
	if err != nil {
		return nil, response.Paging{}, err
	}

	return remarkSetting, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *remarkSettingRepo) GetTotalRows() (int, error) {
	var count int64
	err := r.db.Model(&model.RemarkSetting{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func NewRemarkSettingRepo(db *gorm.DB) RemarkSettingRepo {
	return &remarkSettingRepo{
		db: db,
	}
}
