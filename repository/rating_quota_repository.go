package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type RatingQuotaRepo interface {
	Save(payload *model.RatingQuota) error
	Get(id string) (*model.RatingQuota, error)
	List() ([]model.RatingQuota, error)
	Delete(projectId, businessunitId string) error
	Bulksave(payload *[]model.RatingQuota) error
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
	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		return r.db.Save(&currentBatch).Error

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := r.db.Save(&remainingItems)
		if err != nil {
			return r.db.Save(&remainingItems).Error
		}
	}
	return nil
}

func (r *ratingQuotaRepo) Get(id string) (*model.RatingQuota, error) {
	var ratingQuota model.RatingQuota
	err := r.db.Preload("Project").Preload("BusinessUnit").First(&ratingQuota, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ratingQuota, nil
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

func NewRatingQuotaRepo(db *gorm.DB) RatingQuotaRepo {
	return &ratingQuotaRepo{
		db: db,
	}
}
