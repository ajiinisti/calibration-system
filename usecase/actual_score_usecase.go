package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"calibration-system.com/model"
	"calibration-system.com/repository"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type ActualScoreUsecase interface {
	FindAll() ([]model.ActualScore, error)
	FindById(projectId, employeeId string) (*model.ActualScore, error)
	SaveData(payload *model.ActualScore) error
	DeleteData(projectId, employeeId string) error
	BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error)
}

type actualScoreUsecase struct {
	repo     repository.ActualScoreRepo
	employee UserUsecase
	project  ProjectUsecase
}

func (r *actualScoreUsecase) FindAll() ([]model.ActualScore, error) {
	return r.repo.List()
}

func (r *actualScoreUsecase) FindById(projectId, employeeId string) (*model.ActualScore, error) {
	return r.repo.Get(projectId, employeeId)
}

func (r *actualScoreUsecase) SaveData(payload *model.ActualScore) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.EmployeeID != "" {
		_, err := r.employee.FindById(payload.EmployeeID)
		if err != nil {
			return fmt.Errorf("Employee Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *actualScoreUsecase) DeleteData(projectId, employeeId string) error {
	return r.repo.Delete(projectId, employeeId)
}

func (r *actualScoreUsecase) BulkInsert(file *multipart.FileHeader, projectId string) ([]string, error) {
	logs := map[string]string{}
	var actualScores []model.ActualScore

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

	sheetName := xlsFile.GetSheetName(3)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		nik := row[0]
		y1rating := row[1]
		y2rating := row[2]
		actualRating := row[7]

		employee, err := r.employee.FindByNik(nik)
		if err != nil {
			if _, ok := logs[nik]; !ok {
				logs[nik] = nik
			}
			passed = false
		}

		ptt, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			if _, ok := logs[nik]; !ok {
				logs[nik] = nik
			}
			passed = false
		}

		pat, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			if _, ok := logs[nik]; !ok {
				logs[nik] = nik
			}
			passed = false
		}

		score360, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			if _, ok := logs[nik]; !ok {
				logs[nik] = nik
			}
			passed = false
		}

		actualScore, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			if _, ok := logs[nik]; !ok {
				logs[nik] = nik
			}
			passed = false
		}

		if passed {
			actualScore := model.ActualScore{
				ProjectID:    projectId,
				EmployeeID:   employee.ID,
				ActualScore:  actualScore,
				ActualRating: actualRating,
				Y1Rating:     y1rating,
				Y2Rating:     y2rating,
				PTTScore:     ptt,
				PATScore:     pat,
				Score360:     score360,
			}
			actualScores = append(actualScores, actualScore)
		}
	}

	var dataError []string
	for _, key := range logs {
		dataError = append(dataError, key)
	}

	if len(dataError) > 0 {
		return dataError, fmt.Errorf("Error when insert data")
	}

	err = r.repo.Bulksave(&actualScores)
	if err != nil {
		return nil, err
	}

	return dataError, nil
}

func NewActualScoreUsecase(repo repository.ActualScoreRepo, employee UserUsecase, project ProjectUsecase) ActualScoreUsecase {
	return &actualScoreUsecase{
		repo:     repo,
		employee: employee,
		project:  project,
	}
}
