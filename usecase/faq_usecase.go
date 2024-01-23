package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type FaqUsecase interface {
	FindAll() ([]*model.Faq, error)
	FindAllActive() ([]*model.Faq, error)
	FindById(id string) (*model.Faq, error)
	SaveData(payload *model.Faq) error
	DeleteData(id string) error
	FindByName(name string) (*model.Faq, error)
}

type faqUsecase struct {
	repo repository.FaqRepo
}

func (r *faqUsecase) FindByName(name string) (*model.Faq, error) {
	return r.repo.SearchByName(name)
}

func (r *faqUsecase) FindAll() ([]*model.Faq, error) {
	return r.repo.List()
}

func (r *faqUsecase) FindAllActive() ([]*model.Faq, error) {
	return r.repo.ListActive()
}

func (r *faqUsecase) FindById(id string) (*model.Faq, error) {
	return r.repo.Get(id)
}

func (r *faqUsecase) SaveData(payload *model.Faq) error {
	return r.repo.Save(payload)
}

func (r *faqUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewFaqUsecase(repo repository.FaqRepo) FaqUsecase {
	return &faqUsecase{
		repo: repo,
	}
}
