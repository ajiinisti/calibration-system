package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type AnnouncementRepo interface {
	Save(payload *model.Announcement) error
	Get(id string) (*model.Announcement, error)
	List() ([]*model.Announcement, error)
	ListActive() ([]*model.Announcement, error)
	Delete(id string) error
	SearchByName(name string) (*model.Announcement, error)
}

type announcementRepo struct {
	db *gorm.DB
}

func (r *announcementRepo) Save(payload *model.Announcement) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *announcementRepo) Get(id string) (*model.Announcement, error) {
	var Announcement model.Announcement
	err := r.db.First(&Announcement, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &Announcement, nil
}

func (r *announcementRepo) SearchByName(name string) (*model.Announcement, error) {
	var Announcement model.Announcement
	err := r.db.First(&Announcement, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &Announcement, nil
}

func (r *announcementRepo) List() ([]*model.Announcement, error) {
	var announcements []*model.Announcement
	err := r.db.Find(&announcements).Error
	if err != nil {
		return nil, err
	}
	return announcements, nil
}

func (r *announcementRepo) ListActive() ([]*model.Announcement, error) {
	var announcements []*model.Announcement
	err := r.db.
		Table("announcements a").
		Where("a.active = ?", true).
		Order("a.order ASC").
		Order("a.created_at DESC").
		Find(&announcements).Error
	if err != nil {
		return nil, err
	}
	return announcements, nil
}

func (r *announcementRepo) Delete(id string) error {
	result := r.db.Delete(&model.Announcement{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Announcement not found!")
	}
	return nil
}

func NewAnnouncementRepo(db *gorm.DB) AnnouncementRepo {
	return &announcementRepo{
		db: db,
	}
}
