package usecase

import (
	"calibration-system.com/model"
	"calibration-system.com/repository"
)

type ProjectUsecase interface {
	BaseUsecase[model.Project]
}

type projectUsecase struct {
	repo repository.ProjectRepo
}

func (r *projectUsecase) FindAll() ([]model.Project, error) {
	return r.repo.List()
}

func (r *projectUsecase) FindById(id string) (*model.Project, error) {
	return r.repo.Get(id)
}

func (r *projectUsecase) SaveData(payload *model.Project) error {
	return r.repo.Save(payload)
}

func (r *projectUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func NewProjectUsecase(repo repository.ProjectRepo) ProjectUsecase {
	return &projectUsecase{
		repo: repo,
	}
}
