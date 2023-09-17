package usecase

import (
	"fmt"
	"mime/multipart"

	"calibration-system.com/model"
	"calibration-system.com/repository"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type CalibrationUsecase interface {
	FindAll() ([]model.Calibration, error)
	FindById(id string) (*model.Calibration, error)
	SaveData(payload *model.Calibration) error
	DeleteData(projectId, projectPhaseId, employeeId string) error
	CheckEmployee(file *multipart.FileHeader, projectId string) ([]string, error)
	CheckCalibrator(file *multipart.FileHeader, projectId string) ([]string, error)
	BulkInsert(file *multipart.FileHeader, projectId string) error
}

type calibrationUsecase struct {
	repo         repository.CalibrationRepo
	user         UserUsecase
	project      ProjectUsecase
	projectPhase ProjectPhaseUsecase
}

func (r *calibrationUsecase) FindAll() ([]model.Calibration, error) {
	return r.repo.List()
}

func (r *calibrationUsecase) FindById(id string) (*model.Calibration, error) {
	return r.repo.Get(id)
}

func (r *calibrationUsecase) SaveData(payload *model.Calibration) error {
	if payload.ProjectID != "" {
		_, err := r.project.FindById(payload.ProjectID)
		if err != nil {
			return fmt.Errorf("Project Not Found")
		}
	}

	if payload.ProjectPhaseID != "" {
		_, err := r.projectPhase.FindById(payload.ProjectPhaseID)
		if err != nil {
			return fmt.Errorf("Project Phase Not Found")
		}
	}

	if payload.EmployeeID != "" {
		_, err := r.user.FindById(payload.EmployeeID)
		if err != nil {
			return fmt.Errorf("Employee Not Found")
		}
	}

	if payload.CalibratorID != "" {
		_, err := r.user.FindById(payload.CalibratorID)
		if err != nil {
			return fmt.Errorf("Calibrator Not Found")
		}
	}

	if payload.SpmoID != "" {
		_, err := r.user.FindById(payload.SpmoID)
		if err != nil {
			return fmt.Errorf("SPMO Not Found")
		}
	}
	return r.repo.Save(payload)
}

func (r *calibrationUsecase) DeleteData(projectId, projectPhaseId, employeeId string) error {
	return r.repo.Delete(projectId, projectPhaseId, employeeId)
}

func (r *calibrationUsecase) CheckEmployee(file *multipart.FileHeader, projectId string) ([]string, error) {
	var logs []string

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
		if i == 0 {
			continue
		}

		nik := row[0]

		_, err = r.user.FindByNik(nik)
		if err != nil {
			logs = append(logs, fmt.Sprintf("Employee not available in database %s", nik))
		}
	}

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when checking nik")
	}

	return logs, nil
}

func (r *calibrationUsecase) CheckCalibrator(file *multipart.FileHeader, projectId string) ([]string, error) {
	var logs []string
	type EmployeeSupervisor struct {
		EmployeeID    string
		SupervisorNIK string
	}
	calibrators := map[string]EmployeeSupervisor{}

	project, err := r.project.FindById(projectId)
	if err != nil {
		return nil, err
	}

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
	fmt.Println("rows: ", rows)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		lenProjectPhase := len(project.ProjectPhases)
		var supervisorNIK string
		for j := lenProjectPhase; j > 0; j-- {
			calibratorNik := row[j]
			_, exist := calibrators[calibratorNik]
			// fmt.Println(fmt.Sprintln("NIK: ", calibratorNik))
			if !exist && calibratorNik != "None" {
				// fmt.Println("Not exist")
				calibrator, err := r.user.FindByNik(calibratorNik)
				if err != nil {
					logs = append(logs, fmt.Sprintf("Calibrator not available in database %s", calibratorNik))
				} else {
					calibrators[calibratorNik] = EmployeeSupervisor{
						EmployeeID:    calibrator.ID,
						SupervisorNIK: supervisorNIK,
					}
				}
			}
			supervisorNIK = calibratorNik
		}
	}
	fmt.Println(calibrators)

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when checking nik")
	}

	return logs, nil
}

func (r *calibrationUsecase) BulkInsert(file *multipart.FileHeader, projectId string) error {
	type EmployeeSupervisor struct {
		EmployeeID    string
		SupervisorNIK string
	}
	calibrators := map[string]EmployeeSupervisor{}
	calibrations := []model.Calibration{}
	phases := []string{}

	project, err := r.project.FindById(projectId)
	if err != nil {
		return err
	}

	//Asc by phase order
	for _, v := range project.ProjectPhases {
		phases = append(phases, v.ID)
	}

	fmt.Println("phase: ", project.ProjectPhases)

	excelFile, err := file.Open()
	if err != nil {
		return err
	}
	defer excelFile.Close()

	xlsFile, err := excelize.OpenReader(excelFile)
	if err != nil {
		return err
	}

	sheetName := xlsFile.GetSheetName(5)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		if i == 0 {
			continue
		}

		var supervisorNIK string
		lenProjectPhase := len(project.ProjectPhases)
		employee, err := r.user.FindByNik(row[0])
		if err != nil {
			return fmt.Errorf("Employee NIK not available in database %s", row[0])
		}
		for j := lenProjectPhase; j > 0; j-- {
			calibratorNik := row[j]
			calibratorEs, exist := calibrators[calibratorNik]
			fmt.Sprintln(calibratorNik != "None", " NIK: ", calibratorNik)
			if calibratorNik != "None" {
				var calibratorId string
				if !exist {
					// fmt.Println("Not exist")
					calibratorUser, err := r.user.FindByNik(calibratorNik)
					if err != nil {
						return fmt.Errorf("Calibrator not available in database %s", calibratorNik)
					} else {
						calibrators[calibratorNik] = EmployeeSupervisor{
							EmployeeID:    calibratorUser.ID,
							SupervisorNIK: supervisorNIK,
						}
						calibratorId = calibratorUser.ID
					}
				} else {
					calibratorId = calibratorEs.EmployeeID
				}
				cali := model.Calibration{
					ProjectID:      projectId,
					ProjectPhaseID: phases[j-1],
					EmployeeID:     employee.ID,
					CalibratorID:   calibratorId,
					SpmoID:         calibratorId,
					// CalibrationScore:  5.4,
					// CalibrationRating: "A",
				}
				// fmt.Println("KALIBRASI XX:", cali)
				calibrations = append(calibrations, cali)
			}
			supervisorNIK = calibratorNik
		}
	}
	// fmt.Println("KALIBRASI X:", calibrations)

	return r.repo.Bulksave(&calibrations)
}

func NewCalibrationUsecase(repo repository.CalibrationRepo, user UserUsecase, project ProjectUsecase, projectPhase ProjectPhaseUsecase) CalibrationUsecase {
	return &calibrationUsecase{
		repo:         repo,
		user:         user,
		project:      project,
		projectPhase: projectPhase,
	}
}
