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
	FindById(id string) (*model.ActualScore, error)
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

func (r *actualScoreUsecase) FindById(id string) (*model.ActualScore, error) {
	return r.repo.Get(id)
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
	var logs []string
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

	sheetName := xlsFile.GetSheetName(4)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		nik := row[0]
		y1rating := row[1]
		y2rating := row[2]
		actualRating := row[4]

		employee, err := r.employee.FindByNik(nik)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot get employee nik on row %d ", i))
			passed = false
		}

		as, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Error cannot convert actual score on row %d ", i))
			passed = false
		}

		if passed {
			actualScore := model.ActualScore{
				ProjectID:    projectId,
				EmployeeID:   employee.ID,
				ActualScore:  as,
				ActualRating: actualRating,
				Y1Rating:     y1rating,
				Y2Rating:     y2rating,
			}
			actualScores = append(actualScores, actualScore)
		}
	}

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when insert data")
	}

	err = r.repo.Bulksave(&actualScores)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func NewActualScoreUsecase(repo repository.ActualScoreRepo, employee UserUsecase, project ProjectUsecase) ActualScoreUsecase {
	return &actualScoreUsecase{
		repo:     repo,
		employee: employee,
		project:  project,
	}
}
