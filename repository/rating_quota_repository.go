package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type RatingQuotaRepo interface {
	Save(payload *model.RatingQuota) error
	Get(projectID, businessUnitID string) (*model.RatingQuota, error)
	GetByProject(id string) ([]*model.RatingQuota, error)
	List() ([]model.RatingQuota, error)
	Delete(projectId, businessunitId string) error
	Bulksave(payload *[]model.RatingQuota) error
	PaginateList(pagination model.PaginationQuery, id string) ([]model.RatingQuota, response.Paging, error)
	GetTotalRows() (int, error)
}

type ratingQuotaRepo struct {
	db *gorm.DB
}

func (r *ratingQuotaRepo) Save(payload *model.RatingQuota) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *ratingQuotaRepo) Bulksave(payload *[]model.RatingQuota) error {
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

func (r *ratingQuotaRepo) GetByProject(id string) ([]*model.RatingQuota, error) {
	var ratingQuota []*model.RatingQuota
	err := r.db.Preload("Project").Preload("BusinessUnit").Find(&ratingQuota, "project_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return ratingQuota, nil
}

func (r *ratingQuotaRepo) Get(projectID, businessUnitID string) (*model.RatingQuota, error) {
	var ratingQuota *model.RatingQuota
	err := r.db.Preload("Project").Preload("BusinessUnit").Find(&ratingQuota, "project_id = ? AND business_unit_id = ?", projectID, businessUnitID).Error
	if err != nil {
		return nil, err
	}
	return ratingQuota, nil
}

func (r *ratingQuotaRepo) List() ([]model.RatingQuota, error) {
	var ratingQuotas []model.RatingQuota
	err := r.db.Preload("Project").Preload("BusinessUnit").Find(&ratingQuotas).Error
	if err != nil {
		return nil, err
	}
	return ratingQuotas, nil
}

func (r *ratingQuotaRepo) Delete(projectId, businessunitId string) error {
	result := r.db.Delete(&model.RatingQuota{
		ProjectID:      projectId,
		BusinessUnit:   model.BusinessUnit{},
		BusinessUnitID: businessunitId,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Rating Quota not found!")
	}
	return nil
}

func (r *ratingQuotaRepo) PaginateList(pagination model.PaginationQuery, id string) ([]model.RatingQuota, response.Paging, error) {
	var ratingQuota []model.RatingQuota
	err := r.db.Preload("Project").Preload("BusinessUnit").Limit(pagination.Take).Offset(pagination.Skip).Find(&ratingQuota, "project_id = ?", id).Error
	if err != nil {
		return nil, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRows()
	if err != nil {
		return nil, response.Paging{}, err
	}

	return ratingQuota, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *ratingQuotaRepo) GetTotalRows() (int, error) {
	var count int64
	err := r.db.Model(&model.RatingQuota{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func NewRatingQuotaRepo(db *gorm.DB) RatingQuotaRepo {
	return &ratingQuotaRepo{
		db: db,
	}
}
