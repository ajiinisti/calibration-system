package usecase

import (
	"fmt"
	"mime/multipart"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type BusinessUnitUsecase interface {
	BaseUsecase[model.BusinessUnit]
	BulkInsert(file *multipart.FileHeader) ([]string, error)
	FindPagination(param request.PaginationParam) ([]model.BusinessUnit, response.Paging, error)
}

type businessUnitUsecase struct {
	repo    repository.BusinessUnitRepo
	groupBu GroupBusinessUnitUsecase
}

func (r *businessUnitUsecase) FindAll() ([]model.BusinessUnit, error) {
	return r.repo.List()
}

func (r *businessUnitUsecase) FindPagination(param request.PaginationParam) ([]model.BusinessUnit, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)
	return r.repo.PaginateList(paginationQuery)
}

func (r *businessUnitUsecase) FindById(id string) (*model.BusinessUnit, error) {
	return r.repo.Get(id)
}

func (r *businessUnitUsecase) SaveData(payload *model.BusinessUnit) error {
	if payload.GroupBusinessUnitId != "" {
		_, err := r.groupBu.FindById(payload.GroupBusinessUnitId)
		if err != nil {
			return fmt.Errorf("Group Business Unit Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *businessUnitUsecase) DeleteData(id string) error {
	return r.repo.Delete(id)
}

func (r *businessUnitUsecase) BulkInsert(file *multipart.FileHeader) ([]string, error) {
	var logs []string
	var businessUnits []model.BusinessUnit
	var groupBu []model.GroupBusinessUnit

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

	sheetName := xlsFile.GetSheetName(1)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		buID := row[0]
		buName := row[1]
		gbuName := row[3]

		var found bool
		var gbuid string
		for _, gb := range groupBu {
			if gb.GroupName == gbuName {
				gbuid = gb.ID
				found = true
				break
			}
		}

		if !found {
			gbu, err := r.groupBu.FindByName(gbuName)
			if err != nil {
				logs = append(logs, fmt.Sprintf("Error Group Business Unit Name on Row %d ", i))
				passed = false
			}
			gbuid = gbu.ID
			groupBu = append(groupBu, *gbu)
		}

		if passed {
			bu := model.BusinessUnit{
				ID:                  buID,
				Status:              true,
				Name:                buName,
				GroupBusinessUnitId: gbuid,
			}

			businessUnits = append(businessUnits, bu)
		}
	}

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when insert data")
	}

	err = r.repo.Bulksave(&businessUnits)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func NewBusinessUnitUsecase(repo repository.BusinessUnitRepo, groupBu GroupBusinessUnitUsecase) BusinessUnitUsecase {
	return &businessUnitUsecase{
		repo:    repo,
		groupBu: groupBu,
	}
}
