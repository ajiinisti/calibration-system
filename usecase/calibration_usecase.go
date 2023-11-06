package usecase

import (
	"fmt"
	"mime/multipart"
	"sort"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type CalibrationUsecase interface {
	FindAll() ([]model.Calibration, error)
	FindActiveBySPMOID(spmoID string) ([]model.Calibration, error)
	FindAcceptedBySPMOID(spmoID string) ([]model.Calibration, error)
	FindRejectedBySPMOID(spmoID string) ([]model.Calibration, error)
	FindById(id string) (*model.Calibration, error)
	SaveData(payload *model.Calibration) error
	DeleteData(projectId, projectPhaseId, employeeId string) error
	CheckEmployee(file *multipart.FileHeader, projectId string) ([]string, error)
	CheckCalibrator(file *multipart.FileHeader, projectId string) ([]string, error)
	BulkInsert(file *multipart.FileHeader, projectId string) error
	SubmitCalibrations(payload *request.CalibrationRequest, calibratorID string) error
	SaveCalibrations(payload *request.CalibrationRequest) error
	SpmoAcceptApproval(payload *request.AcceptJustification) error
	SpmoAcceptMultipleApproval(payload *request.AcceptMultipleJustification) error
	SpmoRejectApproval(payload *request.RejectJustification) error
	FindSummaryCalibrationBySPMOID(spmoID string) (response.SummarySPMO, error)
	FindAllDetailCalibrationbySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error)
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

func (r *calibrationUsecase) FindActiveBySPMOID(spmoID string) ([]model.Calibration, error) {
	return r.repo.GetActiveBySPMOID(spmoID)
}

func (r *calibrationUsecase) FindAcceptedBySPMOID(spmoID string) ([]model.Calibration, error) {
	return r.repo.GetAcceptedBySPMOID(spmoID)
}

func (r *calibrationUsecase) FindRejectedBySPMOID(spmoID string) ([]model.Calibration, error) {
	return r.repo.GetRejectedBySPMOID(spmoID)
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
			if !exist && calibratorNik != "None" {
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

				spmo, err := r.user.FindByNik(row[lenProjectPhase+1])
				if err != nil {
					return fmt.Errorf("SPMO ID not available in database %s", row[lenProjectPhase+1])
				}

				hrbp, err := r.user.FindByNik(row[lenProjectPhase+2])
				if err != nil {
					return fmt.Errorf("HRBP ID not available in database %s", row[lenProjectPhase+2])
				}
				cali := model.Calibration{
					ProjectID:      projectId,
					ProjectPhaseID: phases[j-1],
					EmployeeID:     employee.ID,
					CalibratorID:   calibratorId,
					SpmoID:         spmo.ID,
					HrbpID:         hrbp.ID,
				}
				calibrations = append(calibrations, cali)
			}
			supervisorNIK = calibratorNik
		}
		calibrations[len(calibrations)-1].Status = "Calibrate"
	}

	return r.repo.Bulksave(&calibrations)
}

func (r *calibrationUsecase) SubmitCalibrations(payload *request.CalibrationRequest, calibratorID string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID)
	if err != nil {
		return err
	}
	return r.repo.BulkUpdate(payload, *projectPhase)
}

func (r *calibrationUsecase) SaveCalibrations(payload *request.CalibrationRequest) error {
	return r.repo.SaveChanges(payload)
}

func (r *calibrationUsecase) SpmoAcceptApproval(payload *request.AcceptJustification) error {
	projectPhase, err := r.project.FindCalibratorPhase(payload.CalibratorID)
	if err != nil {
		return err
	}
	return r.repo.AcceptCalibration(payload, projectPhase.Phase.Order)
}

func (r *calibrationUsecase) SpmoAcceptMultipleApproval(payload *request.AcceptMultipleJustification) error {
	return r.repo.AcceptMultipleCalibration(payload)
}

func (r *calibrationUsecase) SpmoRejectApproval(payload *request.RejectJustification) error {
	return r.repo.RejectCalibration(payload)
}

type ByBusinessUnitName []response.BUPerformanceSummarySPMO

func (a ByBusinessUnitName) Len() int           { return len(a) }
func (a ByBusinessUnitName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBusinessUnitName) Less(i, j int) bool { return a[i].BusinessUnitName < a[j].BusinessUnitName }

func (r *calibrationUsecase) FindSummaryCalibrationBySPMOID(spmoID string) (response.SummarySPMO, error) {
	results, err := r.repo.GetSummaryBySPMOID(spmoID)
	if err != nil {
		return response.SummarySPMO{}, err
	}

	// Grouping the query output
	groupedList := make(map[string]response.BUPerformanceSummarySPMO)

	for _, res := range results {
		departmentCount := response.ProjectPhaseSummarySPMO{
			CalibratorName: res.CalibratorName,
			CalibratorID:   res.CalibratorID,
			ProjectPhaseID: res.ProjectPhaseID,
			Order:          res.Order,
			Count:          res.Count,
			Status:         "Pending",
		}

		departmentData := response.DepartmentCountSummarySPMO{
			DepartmentName:   res.Department,
			ProjectPhaseData: []*response.ProjectPhaseSummarySPMO{&departmentCount},
		}

		businessUnitData, ok := groupedList[res.BusinessUnitID]

		if ok {
			departmentExists := false
			for i, department := range businessUnitData.DepartmentData {
				if department.DepartmentName == res.Department {
					groupedList[res.BusinessUnitID].DepartmentData[i].ProjectPhaseData = append(department.ProjectPhaseData, &departmentCount)
					departmentExists = true
					break
				}
			}

			if !departmentExists {
				groupedList[res.BusinessUnitID] = response.BUPerformanceSummarySPMO{
					BusinessUnitName: res.BusinessUnitName,
					BusinessUnitID:   res.BusinessUnitID,
					DepartmentData:   append(businessUnitData.DepartmentData, departmentData),
				}
			}
		} else {
			businessUnitData = response.BUPerformanceSummarySPMO{
				BusinessUnitName: res.BusinessUnitName,
				BusinessUnitID:   res.BusinessUnitID,
				DepartmentData:   []response.DepartmentCountSummarySPMO{departmentData},
			}
			groupedList[res.BusinessUnitID] = businessUnitData
		}
	}

	// Transforming map to slice
	var finalResult []response.BUPerformanceSummarySPMO
	for _, value := range groupedList {
		finalResult = append(finalResult, value)
	}

	// Assigning the final result to the SummarySPMO struct
	summary := response.SummarySPMO{SummaryData: finalResult}
	sort.Sort(ByBusinessUnitName(summary.SummaryData))

	for _, smry := range summary.SummaryData {
		for _, department := range smry.DepartmentData {
			for _, projectPhase := range department.ProjectPhaseData {
				status := "-"
				data, err := r.repo.GetAllDetailCalibration2BySPMOID(spmoID, projectPhase.CalibratorID, smry.BusinessUnitID, projectPhase.Order)
				if err != nil {
					return response.SummarySPMO{}, err
				}

				allSubmitted := true
				for _, user := range data {
					lastCalibrationStatus := user.CalibrationScores[len(user.CalibrationScores)-1].SpmoStatus
					if lastCalibrationStatus == "Waiting" {
						status = "Pending"
						allSubmitted = allSubmitted && false
						break
					} else if lastCalibrationStatus == "Accepted" || lastCalibrationStatus == "Rejected" {
						allSubmitted = allSubmitted && true
					} else {
						allSubmitted = allSubmitted && false
					}
				}

				if allSubmitted {
					status = "Completed"
				}
				projectPhase.Status = status
			}
		}
	}

	return summary, nil
}

func (r *calibrationUsecase) FindAllDetailCalibrationbySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error) {
	return r.repo.GetAllDetailCalibrationBySPMOID(spmoID, calibratorID, businessUnitID, department, order)
}

func (r *calibrationUsecase) FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error) {
	return r.repo.GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID, order)
}

func NewCalibrationUsecase(repo repository.CalibrationRepo, user UserUsecase, project ProjectUsecase, projectPhase ProjectPhaseUsecase) CalibrationUsecase {
	return &calibrationUsecase{
		repo:         repo,
		user:         user,
		project:      project,
		projectPhase: projectPhase,
	}
}
