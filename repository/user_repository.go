package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	BaseRepository[model.User]
	SearchByEmail(email string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func (u *userRepo) SearchByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.db.Preload("Role").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Save(payload *model.User) error {
	err := u.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (u *userRepo) Get(id string) (*model.User, error) {
	var user model.User
	err := u.db.Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) List() ([]model.User, error) {
	var users []model.User
	err := u.db.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepo) Delete(id string) error {
	result := u.db.Delete(&model.User{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Uer not found!")
	}
	return nil
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}
