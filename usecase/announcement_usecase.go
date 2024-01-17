package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type AnnouncementUsecase interface {
	FindAll() ([]*model.Announcement, error)
	FindAllActive() ([]*model.Announcement, error)
	FindById(id string) (*model.Announcement, error)
	SaveData(payload *model.Announcement) error
	DeleteData(id string) error
	FindByName(name string) (*model.Announcement, error)
}

type announcementUsecase struct {
	repo repository.AnnouncementRepo
}

func (r *announcementUsecase) FindByName(name string) (*model.Announcement, error) {
	return r.repo.SearchByName(name)
}

func (r *announcementUsecase) FindAll() ([]*model.Announcement, error) {
	return r.repo.List()
}

func (r *announcementUsecase) FindAllActive() ([]*model.Announcement, error) {
	return r.repo.ListActive()
}

func (r *announcementUsecase) FindById(id string) (*model.Announcement, error) {
	return r.repo.Get(id)
}

func (r *announcementUsecase) SaveData(payload *model.Announcement) error {
	return r.repo.Save(payload)
}

func (r *announcementUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewAnnouncementUsecase(repo repository.AnnouncementRepo) AnnouncementUsecase {
	return &announcementUsecase{
		repo: repo,
	}
}
