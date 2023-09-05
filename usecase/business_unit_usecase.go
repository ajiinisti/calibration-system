package usecase

import (
	"fmt"
	"mime/multipart"

	"calibration-system.com/model"
	"calibration-system.com/repository"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type BusinessUnitUsecase interface {
	BaseUsecase[model.BusinessUnit]
	BulkInsert(file *multipart.FileHeader) ([]string, error)
}

type businessUnitUsecase struct {
	repo    repository.BusinessUnitRepo
	groupBu GroupBusinessUnitUsecase
}

func (r *businessUnitUsecase) FindAll() ([]model.BusinessUnit, error) {
	return r.repo.List()
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
		if i == 0 {
			// Skip the first row
			continue
		}

		buID := row[0]
		buName := row[1]
		gbuName := row[3]

		bu := model.BusinessUnit{
			ID:     buID,
			Status: true,
			Name:   buName,
		}

		var found bool
		for _, gb := range groupBu {
			if gb.GroupName == gbuName {
				bu.GroupBusinessUnitId = gb.ID
				found = true
				break
			}
		}

		if !found {
			gbu, err := r.groupBu.FindByName(gbuName)
			if err != nil {
				return nil, err
			}
			bu.GroupBusinessUnitId = gbu.ID
			groupBu = append(groupBu, *gbu)
		}

		businessUnits = append(businessUnits, bu)
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
