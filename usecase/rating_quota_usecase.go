package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"calibration-system.com/model"
	"calibration-system.com/repository"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type RatingQuotaUsecase interface {
	BaseUsecase[model.RatingQuota]
	BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error)
}

type ratingQuotaUsecase struct {
	repo         repository.RatingQuotaRepo
	businessUnit BusinessUnitUsecase
	project      ProjectUsecase
}

func (r *ratingQuotaUsecase) FindAll() ([]model.RatingQuota, error) {
	return r.repo.List()
}

func (r *ratingQuotaUsecase) FindById(id string) (*model.RatingQuota, error) {
	return r.repo.Get(id)
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

func (r *ratingQuotaUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func (r *ratingQuotaUsecase) BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error) {
	var logs []string
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

	sheetName := xlsFile.GetSheetName(6)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		buId := row[0]
		_, err := r.businessUnit.FindById(buId)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot get business unit id on database %s", buId))
			passed = false
		}

		aPlusQuota, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert A plus quota on business unit %s ", buId))
			passed = false
		}

		aQuota, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert A quota on business unit %s ", buId))
			passed = false
		}

		bPlusQuota, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert B plus quota on business unit %s ", buId))
			passed = false
		}

		bQuota, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert B quota on business unit %s ", buId))
			passed = false
		}

		cQuota, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert C quota on business unit %s ", buId))
			passed = false
		}

		dQuota, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert D quota on business unit %s ", buId))
			passed = false
		}

		if passed {
			actualScore := model.RatingQuota{
				ProjectID:      projectId,
				BusinessUnitID: buId,
				APlusQuota:     aPlusQuota,
				AQuota:         aQuota,
				BPlusQuota:     bPlusQuota,
				BQuota:         bQuota,
				CQuota:         cQuota,
				DQuota:         dQuota,
			}
			ratingQuotas = append(ratingQuotas, actualScore)
		}
	}

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when insert data")
	}

	err = r.repo.Bulksave(&ratingQuotas)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func NewRatingQuotaUsecase(repo repository.RatingQuotaRepo, businessUnit BusinessUnitUsecase, project ProjectUsecase) RatingQuotaUsecase {
	return &ratingQuotaUsecase{
		repo:         repo,
		businessUnit: businessUnit,
		project:      project,
	}
}
