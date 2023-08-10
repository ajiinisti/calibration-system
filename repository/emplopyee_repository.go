package repository

import (
	"fmt"

	"calibration-system.com/model"
	"gorm.io/gorm"
)

type EmployeeRepo interface {
	BaseRepository[model.Employee]
	GetByEmail(email string) (*model.Employee, error)
}

type employeeRepo struct {
	db *gorm.DB
}

func (r *employeeRepo) Save(payload *model.Employee) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *employeeRepo) GetByEmail(email string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.First(&employee, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepo) Get(id string) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.First(&employee, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepo) List() ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepo) Delete(id string) error {
	result := r.db.Delete(&model.Employee{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Employee not found!")
	}
	return nil
}

func NewEmployeeRepo(db *gorm.DB) EmployeeRepo {
	return &employeeRepo{
		db: db,
	}
}
