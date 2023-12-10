package usecase

import (
	"fmt"
	"math"
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
	FindById(projectID, projectPhaseID, employeeID string) (*model.Calibration, error)
	FindByProjectEmployeeId(projectID, employeeID string) ([]model.Calibration, error)
	SaveData(payload *model.Calibration) error
	SaveDataByUser(payload *request.CalibrationForm) error
	DeleteData(projectId, employeeId string) error
	CheckEmployee(file *multipart.FileHeader, projectId string) ([]string, error)
	CheckCalibrator(file *multipart.FileHeader, projectId string) ([]string, error)
	BulkInsert(file *multipart.FileHeader, projectId string) error
	SubmitCalibrations(payload *request.CalibrationRequest, calibratorID string) error
	SendCalibrationsToManager(payload *request.CalibrationRequest, calibratorID string) error
	SaveCalibrations(payload *request.CalibrationRequest) error
	SpmoAcceptApproval(payload *request.AcceptJustification) error
	SpmoAcceptMultipleApproval(payload *request.AcceptMultipleJustification) error
	SpmoRejectApproval(payload *request.RejectJustification) error
	SpmoSubmit(payload *request.AcceptMultipleJustification) error
	FindSummaryCalibrationBySPMOID(spmoID string) (response.SummarySPMO, error)
	FindAllDetailCalibrationbySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error)
	SendNotificationToCurrentCalibrator() error
	FindRatingQuotaSPMOByCalibratorID(spmoID, calibratorID, businessUnitID string, order int) (*response.RatingQuota, error)
}

type calibrationUsecase struct {
	repo         repository.CalibrationRepo
	user         UserUsecase
	project      ProjectUsecase
	projectPhase ProjectPhaseUsecase
	notification NotificationUsecase
	actualScore  ActualScoreUsecase
}

func (r *calibrationUsecase) SendNotificationToCurrentCalibrator() error {
	calibrations, err := r.repo.GetCalibrateCalibration()
	if err != nil {
		return err
	}

	uniqueCalibratorIDs := make(map[string]response.NotificationModel)
	for _, data := range calibrations {
		if _, ok := uniqueCalibratorIDs[data.CalibratorID]; !ok {
			uniqueCalibratorIDs[data.CalibratorID] = response.NotificationModel{
				CalibratorID: data.CalibratorID,
				ProjectPhase: data.ProjectPhase.Phase.Order,
				Deadline:     data.ProjectPhase.EndDate,
			}
		}
	}

	var uniqueCalibratorIDsSlice []response.NotificationModel
	for _, value := range uniqueCalibratorIDs {
		uniqueCalibratorIDsSlice = append(uniqueCalibratorIDsSlice, value)
	}

	var currentCalibrators []response.NotificationModel
	for _, calibratorData := range uniqueCalibratorIDsSlice {
		calibrations, err := r.repo.GetAllCalibrationByCalibratorID(calibratorData.CalibratorID)
		if err != nil {
			return err
		}

		flag := true
		for _, calibrationData := range calibrations {
			if calibrationData.Status != "Calibrate" {
				flag = flag && false
				break
			}
		}

		if flag {
			currentCalibrators = append(currentCalibrators, calibratorData)
		}

	}

	err = r.notification.NotifyThisCurrentCalibrators(currentCalibrators)
	if err != nil {
		return err
	}
	return nil
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

func (r *calibrationUsecase) FindById(projectID, projectPhaseID, employeeID string) (*model.Calibration, error) {
	return r.repo.Get(projectID, projectPhaseID, employeeID)
}

func (r *calibrationUsecase) FindByProjectEmployeeId(projectID, employeeID string) ([]model.Calibration, error) {
	return r.repo.GetByProjectEmployeeID(projectID, employeeID)
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

func (r *calibrationUsecase) SaveDataByUser(payload *request.CalibrationForm) error {
	project, err := r.project.FindById(payload.CalibrationDataForms[0].ProjectID)
	if err != nil {
		return err
	}

	err = r.repo.SaveByUser(payload, project)
	if err != nil {
		return err
	}

	err = r.actualScore.SaveData(&model.ActualScore{
		ProjectID:    payload.CalibrationDataForms[0].ProjectID,
		EmployeeID:   payload.CalibrationDataForms[0].EmployeeID,
		ActualScore:  payload.ActualScore,
		ActualRating: payload.ActualRating,
		Y1Rating:     payload.Y1Rating,
		Y2Rating:     payload.Y2Rating,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *calibrationUsecase) DeleteData(projectId, employeeId string) error {
	err := r.repo.Delete(projectId, employeeId)
	if err != nil {
		return err
	}

	err = r.actualScore.DeleteData(projectId, employeeId)
	if err != nil {
		return err
	}

	return nil
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
	logs := map[string]string{}
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
					if _, ok := logs[calibratorNik]; !ok {
						logs[calibratorNik] = calibratorNik
					}
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
	// fmt.Println(calibrators)
	var dataError []string
	for _, key := range logs {
		dataError = append(dataError, key)
	}

	if len(dataError) > 0 {
		return dataError, fmt.Errorf("Error when checking nik")
	}

	return dataError, nil
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

	var phaseOne model.ProjectPhase
	//Asc by phase order
	for _, v := range project.ProjectPhases {
		phases = append(phases, v.ID)
		if v.Phase.Order == 1 {
			phaseOne = v
		}
	}

	// fmt.Println("phase: ", project.ProjectPhases)

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

				// hrbp, err := r.user.FindByNik(row[lenProjectPhase+2])
				// if err != nil {
				// 	return fmt.Errorf("HRBP ID not available in database %s", row[lenProjectPhase+2])
				// }
				cali := model.Calibration{
					ProjectID:      projectId,
					ProjectPhaseID: phases[j-1],
					EmployeeID:     employee.ID,
					CalibratorID:   calibratorId,
					SpmoID:         spmo.ID,
					// HrbpID:         hrbp.ID,
				}
				calibrations = append(calibrations, cali)
			} else {
				cal, _ := r.repo.Get(projectId, phases[j-1], employee.ID)
				if cal != nil {
					err := r.repo.DeleteCalibrationPhase(projectId, phases[j-1], employee.ID)
					if err != nil {
						return err
					}
				}
			}
			supervisorNIK = calibratorNik
		}

		if calibrations[len(calibrations)-1].ProjectPhaseID == phaseOne.ID {
			calibrations[len(calibrations)-2].Status = "Calibrate"
		} else {
			calibrations[len(calibrations)-1].Status = "Calibrate"
		}
	}

	return r.repo.Bulksave(&calibrations)
}

func (r *calibrationUsecase) SubmitCalibrations(payload *request.CalibrationRequest, calibratorID string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID)
	if err != nil {
		return err
	}

	spmoIDs, err := r.repo.BulkUpdate(payload, *projectPhase)
	if err != nil {
		return err
	}

	calibrator, err := r.user.FindById(calibratorID)
	if err != nil {
		return err
	}

	mapSpmo := make(map[string]*model.User)
	for _, spmoID := range spmoIDs {
		if spmoID != nil {
			spmo, err := r.user.FindById(*spmoID)
			if err != nil {
				return err
			}

			if _, ok := mapSpmo[*spmoID]; !ok {
				mapSpmo[*spmoID] = spmo
			}
		}
	}

	var listSpmo []*model.User
	for _, data := range mapSpmo {
		listSpmo = append(listSpmo, data)
	}

	err = r.notification.NotifyCalibrationToSpmo(calibrator, listSpmo)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) SendCalibrationsToManager(payload *request.CalibrationRequest, calibratorID string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID)
	if err != nil {
		return err
	}

	managerCalibratorIDs, projectPhaseId, err := r.repo.UpdateManagerCalibrations(payload, *projectPhase)
	if err != nil {
		return err
	}

	projectPhaseNew, err := r.projectPhase.FindById(projectPhaseId)
	if err != nil {
		return err
	}

	uniqueCalibrator := removeDuplicates(managerCalibratorIDs)

	err = r.notification.NotifyCalibrators(uniqueCalibrator, projectPhaseNew.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) SaveCalibrations(payload *request.CalibrationRequest) error {
	return r.repo.SaveChanges(payload)
}

func (r *calibrationUsecase) SpmoAcceptApproval(payload *request.AcceptJustification) error {
	projectPhase, err := r.project.FindCalibratorPhase(payload.CalibratorID)
	if err != nil {
		return err
	}

	err = r.repo.AcceptCalibration(payload, projectPhase.Phase.Order)
	if err != nil {
		return err
	}

	// err = r.notification.NotifyApprovedCalibrationToCalibrator([]string{payload.CalibratorID})
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *calibrationUsecase) SpmoAcceptMultipleApproval(payload *request.AcceptMultipleJustification) error {
	err := r.repo.AcceptMultipleCalibration(payload)
	if err != nil {
		return err
	}

	ids := []string{}
	for _, acceptJustification := range payload.ArrayOfAcceptsJustification {
		ids = append(ids, acceptJustification.CalibratorID)
	}

	// results := removeDuplicates(ids)
	// err = r.notification.NotifyApprovedCalibrationToCalibrator(results)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

func (r *calibrationUsecase) SpmoRejectApproval(payload *request.RejectJustification) error {
	err := r.repo.RejectCalibration(payload)
	if err != nil {
		return err
	}

	employee, err := r.user.FindById(payload.EmployeeID)
	if err != nil {
		return err
	}

	employeeName := fmt.Sprintf("%s(%s) - %s - %s", employee.Name, employee.Nik, employee.BusinessUnit.Name, employee.OrganizationUnit)

	err = r.notification.NotifyRejectedCalibrationToCalibrator(payload.CalibratorID, employeeName, payload.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) SpmoSubmit(payload *request.AcceptMultipleJustification) error {
	nextCalibrator, err := r.repo.SubmitReview(payload)
	if err != nil {
		return err
	}

	prevCalibrator := make(map[string]response.NotificationModel)
	for _, requestData := range payload.ArrayOfAcceptsJustification {
		if _, ok := prevCalibrator[requestData.CalibratorID]; !ok {
			prevCalibrator[requestData.CalibratorID] = response.NotificationModel{
				CalibratorID: requestData.CalibratorID,
			}
		}
	}

	var uniquePrev []response.NotificationModel
	for _, data := range prevCalibrator {
		uniquePrev = append(uniquePrev, data)
	}

	err = r.notification.NotifyApprovedCalibrationToCalibrators(uniquePrev)
	if err != nil {
		return err
	}
	err = r.notification.NotifyThisCalibrators(nextCalibrator)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) FindSummaryCalibrationBySPMOID(spmoID string) (response.SummarySPMO, error) {
	results, err := r.repo.GetSummaryBySPMOID(spmoID)
	if err != nil {
		return response.SummarySPMO{}, err
	}

	groupedData := make(map[string]*response.BUPerformanceSummarySPMO)
	for _, d := range results {
		key := d.BusinessUnitID

		calibratorSummary := response.CalibratorSummary{
			CalibratorName: d.CalibratorName,
			CalibratorID:   d.CalibratorID,
			Count:          d.Count,
			Status:         "Pending",
		}

		phaseSummary := response.ProjectPhaseSummarySPMO{
			ProjectPhaseID:     d.ProjectPhaseID,
			Order:              d.Order,
			CalibratorSummarys: []*response.CalibratorSummary{&calibratorSummary},
		}

		if _, ok := groupedData[key]; !ok {
			groupedData[key] = &response.BUPerformanceSummarySPMO{
				BusinessUnitName:    d.BusinessUnitName,
				BusinessUnitID:      d.BusinessUnitID,
				ProjectPhaseSummary: []*response.ProjectPhaseSummarySPMO{&phaseSummary},
			}
		} else {
			data := groupedData[key]
			found := false
			for _, existingPhase := range data.ProjectPhaseSummary {
				if existingPhase.ProjectPhaseID == phaseSummary.ProjectPhaseID {
					existingPhase.CalibratorSummarys = append(existingPhase.CalibratorSummarys, &calibratorSummary)
					found = true
					break
				}
			}
			if !found {
				data.ProjectPhaseSummary = append(data.ProjectPhaseSummary, &phaseSummary)
			}
			groupedData[key] = data
		}
		// fmt.Println("Grouped DATA: =", groupedData)
	}
	// fmt.Println("Grouped DATA: =", groupedData)

	// Transforming map to slice
	var finalResult []*response.BUPerformanceSummarySPMO
	for _, value := range groupedData {
		finalResult = append(finalResult, value)
	}

	// Assigning the final result to the SummarySPMO struct
	summary := response.SummarySPMO{SummaryData: finalResult}
	sort.Slice(summary.SummaryData, func(i, j int) bool {
		return summary.SummaryData[i].BusinessUnitName < summary.SummaryData[j].BusinessUnitName
	})

	for _, smry := range summary.SummaryData {
		for _, projectPhase := range smry.ProjectPhaseSummary {
			countMaximum := 0
			for _, calibrationSummary := range projectPhase.CalibratorSummarys {
				countMaximum += 1
				status := "-"
				data, err := r.repo.GetAllDetailCalibration2BySPMOID(spmoID, calibrationSummary.CalibratorID, smry.BusinessUnitID, projectPhase.Order)
				if err != nil {
					return response.SummarySPMO{}, err
				}

				allSubmitted := true
				for _, user := range data {
					lastCalibrationStatus := user.CalibrationScores[len(user.CalibrationScores)-1].SpmoStatus
					if lastCalibrationStatus == "Waiting" || (lastCalibrationStatus == "Accepted" && user.CalibrationScores[len(user.CalibrationScores)-1].JustificationReviewStatus == false) {
						status = "Pending"
						allSubmitted = allSubmitted && false
						break
					} else if (lastCalibrationStatus == "Accepted" && user.CalibrationScores[len(user.CalibrationScores)-1].JustificationReviewStatus == true) || lastCalibrationStatus == "Rejected" {
						allSubmitted = allSubmitted && true
					} else {
						allSubmitted = allSubmitted && false
					}
				}

				if allSubmitted {
					status = "Completed"
				}
				calibrationSummary.Status = status

			}

			if countMaximum > smry.MaximumTotalData {
				smry.MaximumTotalData = countMaximum
			}
			projectPhase.DataCount = countMaximum
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

func (r *calibrationUsecase) FindRatingQuotaSPMOByCalibratorID(spmoID, calibratorID, businessUnitID string, order int) (*response.RatingQuota, error) {
	users, err := r.FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, order)
	if err != nil {
		return nil, err
	}

	projects, err := r.project.FindProjectRatingQuotaByBusinessUnit(businessUnitID)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := len(users)
	responses := response.RatingQuota{
		APlus: int(math.Floor(((ratingQuota.APlusQuota) / float64(100)) * float64(totalCalibrations))),
		A:     int(math.Floor(((ratingQuota.AQuota) / float64(100)) * float64(totalCalibrations))),
		BPlus: int(math.Floor(((ratingQuota.BPlusQuota) / float64(100)) * float64(totalCalibrations))),
		B:     int(math.Floor(((ratingQuota.BQuota) / float64(100)) * float64(totalCalibrations))),
		C:     int(math.Floor(((ratingQuota.CQuota) / float64(100)) * float64(totalCalibrations))),
		D:     int(math.Floor(((ratingQuota.DQuota) / float64(100)) * float64(totalCalibrations))),
	}

	var total = responses.APlus + responses.A +
		responses.BPlus + responses.B +
		responses.C + responses.D

	if total < totalCalibrations {
		if ratingQuota.Remaining == "A+" {
			responses.APlus += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "A" {
			responses.A += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "B+" {
			responses.BPlus += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "B" {
			responses.B += (totalCalibrations - total)
		} else if ratingQuota.Remaining == "C" {
			responses.C += (totalCalibrations - total)
		} else {
			responses.D += (totalCalibrations - total)
		}
		total += (totalCalibrations - total)
	}

	if total > totalCalibrations {
		if ratingQuota.Excess == "A+" {
			responses.APlus -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "A" {
			responses.A -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "B+" {
			responses.BPlus -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "B" {
			responses.B -= (total - totalCalibrations)
		} else if ratingQuota.Excess == "C" {
			responses.C -= (total - totalCalibrations)
		} else {
			responses.D -= (total - totalCalibrations)
		}
		total -= (total - totalCalibrations)
	}
	responses.Total = total

	return &responses, nil
}

func NewCalibrationUsecase(repo repository.CalibrationRepo, user UserUsecase, project ProjectUsecase, projectPhase ProjectPhaseUsecase, notification NotificationUsecase, actualScore ActualScoreUsecase) CalibrationUsecase {
	return &calibrationUsecase{
		repo:         repo,
		user:         user,
		project:      project,
		projectPhase: projectPhase,
		notification: notification,
		actualScore:  actualScore,
	}
}
