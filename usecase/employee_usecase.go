package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type EmployeeUsecase interface {
	BaseUsecase[model.Employee]
	FindByEmail(email string) (*model.Employee, error)
}

type employeeUsecase struct {
	repo repository.EmployeeRepo
}

func (r *employeeUsecase) FindByEmail(email string) (*model.Employee, error) {
	return r.repo.GetByEmail(email)
}

func (r *employeeUsecase) FindAll() ([]model.Employee, error) {
	return r.repo.List()
}

func (r *employeeUsecase) FindById(id string) (*model.Employee, error) {
	return r.repo.Get(id)
}

func (r *employeeUsecase) SaveData(payload *model.Employee) error {
	return r.repo.Save(payload)
}

func (r *employeeUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewEmployeeUsecase(repo repository.EmployeeRepo) EmployeeUsecase {
	return &employeeUsecase{
		repo: repo,
	}
}
