package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type RoleRepo interface {
	BaseRepository[model.Role]
	SearchByName(name string) (*model.Role, error)
}

type roleRepo struct {
	db *gorm.DB
}

func (r *roleRepo) Save(payload *model.Role) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *roleRepo) Get(id string) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) SearchByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) List() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepo) Delete(id string) error {
	result := r.db.Delete(&model.Role{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Role not found!")
	}
	return nil
}

func NewRoleRepo(db *gorm.DB) RoleRepo {
	return &roleRepo{
		db: db,
	}
}
