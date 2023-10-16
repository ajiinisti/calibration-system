package usecase

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type RemarkSettingUsecase interface {
	FindAll() ([]model.RemarkSetting, error)
	FindById(id string) ([]*model.RemarkSetting, error)
	SaveData(payload *model.RemarkSetting) error
	SaveDataByProject(payload []*model.RemarkSetting) error
	DeleteData(id string) error
	BulkDeleteData(payload request.DeleteRemark) error
	FindPagination(param request.PaginationParam, id string) ([]model.RemarkSetting, response.Paging, error)
}

type remarkSettingUsecase struct {
	repo    repository.RemarkSettingRepo
	project ProjectUsecase
}

func (r *remarkSettingUsecase) FindAll() ([]model.RemarkSetting, error) {
	return r.repo.List()
}

func (r *remarkSettingUsecase) FindById(id string) ([]*model.RemarkSetting, error) {
	return r.repo.Get(id)
}

func (r *remarkSettingUsecase) FindPagination(param request.PaginationParam, id string) ([]model.RemarkSetting, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.PaginateList(paginationQuery, id)
}

func (r *remarkSettingUsecase) SaveData(payload *model.RemarkSetting) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *remarkSettingUsecase) SaveDataByProject(payload []*model.RemarkSetting) error {
	return r.repo.BulkSave(payload)
}

func (r *remarkSettingUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func (r *remarkSettingUsecase) BulkDeleteData(payload request.DeleteRemark) error {
	return r.repo.BulkDelete(payload)
}

func NewRemarkSettingUsecase(repo repository.RemarkSettingRepo, project ProjectUsecase) RemarkSettingUsecase {
	return &remarkSettingUsecase{
		repo:    repo,
		project: project,
	}
}
