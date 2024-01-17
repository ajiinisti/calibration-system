package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type FaqRepo interface {
	Save(payload *model.Faq) error
	Get(id string) (*model.Faq, error)
	List() ([]*model.Faq, error)
	ListActive() ([]*model.Faq, error)
	Delete(id string) error
	SearchByName(name string) (*model.Faq, error)
}

type faqRepo struct {
	db *gorm.DB
}

func (r *faqRepo) Save(payload *model.Faq) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *faqRepo) Get(id string) (*model.Faq, error) {
	var Faq model.Faq
	err := r.db.First(&Faq, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &Faq, nil
}

func (r *faqRepo) SearchByName(name string) (*model.Faq, error) {
	var Faq model.Faq
	err := r.db.First(&Faq, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &Faq, nil
}

func (r *faqRepo) List() ([]*model.Faq, error) {
	var faqs []*model.Faq
	err := r.db.Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}

func (r *faqRepo) ListActive() ([]*model.Faq, error) {
	var faqs []*model.Faq
	err := r.db.
		Table("faqs a").
		Where("a.active = ?", true).
		Order("a.order ASC").
		Order("a.created_at DESC").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}

func (r *faqRepo) Delete(id string) error {
	result := r.db.Delete(&model.Faq{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Faq not found!")
	}
	return nil
}

func NewFaqRepo(db *gorm.DB) FaqRepo {
	return &faqRepo{
		db: db,
	}
}
