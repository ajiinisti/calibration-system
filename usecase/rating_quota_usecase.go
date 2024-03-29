package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type RatingQuotaUsecase interface {
	FindAll() ([]model.RatingQuota, error)
	FindByProject(id string) ([]*model.RatingQuota, error)
	FindById(projectID, businessUnitID string) (*model.RatingQuota, error)
	SaveData(payload *model.RatingQuota) error
	DeleteData(projectId, businessUnitId string) error
	BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error)
	FindPagination(param request.PaginationParam, id string) ([]model.RatingQuota, response.Paging, error)
}

type ratingQuotaUsecase struct {
	repo         repository.RatingQuotaRepo
	businessUnit BusinessUnitUsecase
	project      ProjectUsecase
}

func (r *ratingQuotaUsecase) FindAll() ([]model.RatingQuota, error) {
	return r.repo.List()
}

func (r *ratingQuotaUsecase) FindByProject(id string) ([]*model.RatingQuota, error) {
	return r.repo.GetByProject(id)
}

func (r *ratingQuotaUsecase) FindById(projectID, businessUnitID string) (*model.RatingQuota, error) {
	return r.repo.Get(projectID, businessUnitID)
}

func (r *ratingQuotaUsecase) FindPagination(param request.PaginationParam, id string) ([]model.RatingQuota, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.PaginateList(paginationQuery, id)
}

func (r *ratingQuotaUsecase) SaveData(payload *model.RatingQuota) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.BusinessUnitID != "" {
		_, err := r.businessUnit.FindById(payload.BusinessUnitID)
		if err != nil {
			return fmt.Errorf("Business Unit Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *ratingQuotaUsecase) DeleteData(projectId, businessUnitId string) error {
	return r.repo.Delete(projectId, businessUnitId)
}

func (r *ratingQuotaUsecase) BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error) {
	logs := map[string]string{}
	var ratingQuotas []model.RatingQuota

	_, err := r.project.FindById(projectId)
	if err != nil {
		return nil, err
	}

	// Membuka file Excel yang diunggah
	excelFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()

	xlsFile, err := excelize.OpenReader(excelFile)
	if err != nil {
		return nil, err
	}

	sheetName := xlsFile.GetSheetName(5)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		buId := row[0]
		_, err := r.businessUnit.FindById(buId)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		aPlusQuota, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		aQuota, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		bPlusQuota, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		bQuota, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		cQuota, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		dQuota, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			if _, ok := logs[buId]; !ok {
				logs[buId] = buId
			}
			passed = false
		}

		if passed {
			ratingQuota := model.RatingQuota{
				ProjectID:      projectId,
				BusinessUnitID: buId,
				APlusQuota:     aPlusQuota,
				AQuota:         aQuota,
				BPlusQuota:     bPlusQuota,
				BQuota:         bQuota,
				CQuota:         cQuota,
				DQuota:         dQuota,
				Remaining:      row[7],
				Excess:         row[8],
			}
			ratingQuotas = append(ratingQuotas, ratingQuota)
		}
	}

	var dataError []string
	for _, key := range logs {
		dataError = append(dataError, key)
	}

	if len(dataError) > 0 {
		return dataError, fmt.Errorf("Error when insert data")
	}

	err = r.repo.Bulksave(&ratingQuotas)
	if err != nil {
		return nil, err
	}

	return dataError, nil
}

func NewRatingQuotaUsecase(repo repository.RatingQuotaRepo, businessUnit BusinessUnitUsecase, project ProjectUsecase) RatingQuotaUsecase {
	return &ratingQuotaUsecase{
		repo:         repo,
		businessUnit: businessUnit,
		project:      project,
	}
}
