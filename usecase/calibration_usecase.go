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
	FindActiveUserBySPMOID(spmoID string) ([]model.UserChange, error)
	FindAcceptedBySPMOID(spmoID string) ([]model.Calibration, error)
	FindRejectedBySPMOID(spmoID string) ([]model.Calibration, error)
	FindById(projectID, projectPhaseID, employeeID string) (*model.Calibration, error)
	FindByProjectEmployeeId(projectID, employeeID string) ([]model.CalibrationForm, error)
	SaveData(payload *model.Calibration) error
	SaveDataByUser(payload *request.CalibrationForm) error
	DeleteData(projectId, employeeId string) error
	CheckEmployee(file *multipart.FileHeader, projectId string) ([]string, error)
	CheckCalibrator(file *multipart.FileHeader, projectId string) ([]string, error)
	BulkInsert(file *multipart.FileHeader, projectId string) error
	SubmitCalibrations(calibratorID, projectID, businessUnit string) error
	SaveCalibrations(payload *request.CalibrationRequest) error
	SaveCommentCalibration(payload *model.Calibration) error
	SaveScoreAndRating(payload *model.Calibration) error
	SendCalibrationsToManager(calibratorID, projectID, prevCalibrator, businessUnit string) error
	SendBackCalibrationsToOnePhaseBefore(calibratorID, projectID, prevCalibrator, businessUnit string) error
	SpmoAcceptApproval(payload *request.AcceptJustification) error
	SpmoAcceptMultipleApproval(payload *request.AcceptMultipleJustification) error
	SpmoRejectApproval(payload *request.RejectJustification) error
	SpmoSubmit(payload *request.AcceptMultipleJustification) error
	FindSummaryCalibrationBySPMOID(spmoID, projectID string) (response.SummarySPMO, error)
	FindAllDetailCalibrationbySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, projectID string, order int) ([]response.UserResponse, error)
	SendNotificationToCurrentCalibrator(projectID string) ([]response.NotificationModel, error)
	FindRatingQuotaSPMOByCalibratorID(spmoID, calibratorID, businessUnitID, projectID string, order int) (*response.RatingQuota, error)
	FindLatestJustification(projectID, calibratorID, employeeID string) ([]model.SeeCalibrationJustification, error)
}

type calibrationUsecase struct {
	repo         repository.CalibrationRepo
	user         UserUsecase
	project      ProjectUsecase
	projectPhase ProjectPhaseUsecase
	notification NotificationUsecase
	actualScore  ActualScoreUsecase
}

func (r *calibrationUsecase) FindLatestJustification(projectID, calibratorID, employeeID string) ([]model.SeeCalibrationJustification, error) {
	return r.repo.GetLatestJustification(projectID, calibratorID, employeeID)
}

func (r *calibrationUsecase) SendNotificationToCurrentCalibrator(projectID string) ([]response.NotificationModel, error) {
	calibrations, err := r.repo.GetCalibrateCalibrationByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	uniqueCalibratorIDs := make(map[string]response.NotificationModel)
	for _, data := range calibrations {
		if _, ok := uniqueCalibratorIDs[data.CalibratorID]; !ok {
			user, _ := r.user.FindById(data.CalibratorID)
			uniqueCalibratorIDs[data.CalibratorID] = response.NotificationModel{
				CalibratorID:   data.CalibratorID,
				ProjectPhase:   data.ProjectPhase.Phase.Order,
				Deadline:       data.ProjectPhase.EndDate,
				NextCalibrator: user.Name,
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
			return nil, err
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

	err = r.notification.NotifyFirstCurrentCalibrators(currentCalibrators)
	if err != nil {
		return nil, err
	}
	return currentCalibrators, nil
}

func (r *calibrationUsecase) FindAll() ([]model.Calibration, error) {
	return r.repo.List()
}

func (r *calibrationUsecase) FindActiveUserBySPMOID(spmoID string) ([]model.UserChange, error) {
	return r.repo.GetActiveUserBySPMOID(spmoID)
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

func (r *calibrationUsecase) FindByProjectEmployeeId(projectID, employeeID string) ([]model.CalibrationForm, error) {
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

	err = r.repo.SaveByUser(payload, project, payload.ActualScore, payload.ActualRating)
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
		PTTScore:     payload.PTTScore,
		PATScore:     payload.PATScore,
		Score360:     payload.Score360,
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

	sheetName := xlsFile.GetSheetName(4)
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

	sheetName := xlsFile.GetSheetName(4)
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

	sheetName := xlsFile.GetSheetName(4)
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
				score, err := r.actualScore.FindById(projectId, employee.ID)
				if err != nil {
					return fmt.Errorf("Employee %s doesn't have actual score inputted", row[0])
				}

				cali := model.Calibration{
					ProjectID:           projectId,
					ProjectPhaseID:      phases[j-1],
					EmployeeID:          employee.ID,
					CalibratorID:        calibratorId,
					SpmoID:              spmo.ID,
					CalibrationRating:   score.ActualRating,
					CalibrationScore:    score.ActualScore,
					FilledTopBottomMark: true,
					// HrbpID:         hrbp.ID,
				}

				if row[lenProjectPhase+2] != "None" {
					spmo2, err := r.user.FindByNik(row[lenProjectPhase+2])
					if err != nil {
						return fmt.Errorf("SPMO ID not available in database %s", row[lenProjectPhase+2])
					}
					if spmo2 != nil {
						cali.Spmo2ID = &spmo2.ID
					}
				}

				if row[lenProjectPhase+3] != "None" {
					spmo3, err := r.user.FindByNik(row[lenProjectPhase+3])
					if err != nil {
						return fmt.Errorf("SPMO ID not available in database %s", row[lenProjectPhase+3])
					}

					if spmo3 != nil {
						cali.Spmo3ID = &spmo3.ID
					}
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
			justificationType := returnRemarkType(calibrations[len(calibrations)-2].CalibrationRating, project.RemarkSettings)
			calibrations[len(calibrations)-2].Status = "Calibrate"
			calibrations[len(calibrations)-2].JustificationType = justificationType
			if justificationType != "default" {
				calibrations[len(calibrations)-2].FilledTopBottomMark = false
				calibrations[len(calibrations)-1].FilledTopBottomMark = false
			}
			fmt.Println("=========================ISI TOP BOTTOMNYA=====================", justificationType, calibrations[len(calibrations)-2].FilledTopBottomMark)
		} else {
			justificationType := returnRemarkType(calibrations[len(calibrations)-1].CalibrationRating, project.RemarkSettings)
			calibrations[len(calibrations)-1].Status = "Calibrate"
			calibrations[len(calibrations)-1].JustificationType = justificationType
			if justificationType != "default" {
				calibrations[len(calibrations)-1].FilledTopBottomMark = false
			}
			fmt.Println("=========================ISI TOP BOTTOMNYA=====================", justificationType, calibrations[len(calibrations)-1].FilledTopBottomMark)
		}
	}

	return r.repo.Bulksave(&calibrations)
}

func returnRemarkType(score string, projectRemark []model.RemarkSetting) string {
	value := "default"
	for _, data := range projectRemark {
		if data.JustificationType == "top" {
			if data.ScoringType == "rating" {
				if score == data.From && score == data.To {
					value = "top"
					break
				}
			}
		} else if data.JustificationType == "bottom" {
			if data.ScoringType == "rating" {
				if score == data.From && score == data.To {
					value = "bottom"
					break
				}
			}
		}
	}
	return value
}

func (r *calibrationUsecase) SubmitCalibrations(calibratorID, projectID, businessUnit string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID, projectID)
	if err != nil {
		return err
	}

	payload, err := r.project.FindCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID)
	if err != nil {
		return err
	}

	totalCalibrated, err := r.project.FindTotalCalibratedByCalibratorID(calibratorID, "", businessUnit, "all", projectID)
	if err != nil {
		return err
	}

	totalRatingQuota, err := r.project.FindRatingQuotaByCalibratorID(calibratorID, "", businessUnit, "all", projectID, 0)
	if err != nil {
		return err
	}

	checkCondition, err := r.repo.CheckConditionBeforeSubmitCalibration(projectID, payload.UserData, *projectPhase, *totalCalibrated, *totalRatingQuota)
	if err != nil {
		return err
	}

	if checkCondition {
		spmoIDs, nextCalibrator, err := r.repo.BulkUpdate(payload.UserData, *projectPhase, projectID)
		if err != nil {
			return err
		}

		calibrator, err := r.user.FindById(calibratorID)
		if err != nil {
			return err
		}

		if projectPhase.ReviewSpmo {
			mapSpmo := make(map[string]*model.User)
			for _, spmoID := range spmoIDs {
				spmo, err := r.user.FindById(spmoID)
				if err != nil {
					return err
				}

				if _, ok := mapSpmo[spmo.ID]; !ok {
					mapSpmo[spmo.ID] = spmo
				}
			}

			var listSpmo []*model.User
			for _, data := range mapSpmo {
				listSpmo = append(listSpmo, data)
			}

			err = r.notification.NotifySubmittedCalibrationToSpmo(calibrator, listSpmo, projectPhase.Phase.Order, projectID)
			if err != nil {
				return err
			}
			return nil
		}

		nCalibrator := make(map[string]response.NotificationModel)
		for _, requestData := range nextCalibrator {
			if _, ok := nCalibrator[requestData.CalibratorID]; !ok {
				nCalibrator[requestData.CalibratorID] = response.NotificationModel{
					CalibratorID:           requestData.CalibratorID,
					ProjectPhase:           requestData.ProjectPhase,
					Deadline:               requestData.Deadline,
					PreviousCalibrator:     calibrator.Name,
					PreviousCalibratorID:   calibratorID,
					PreviousBusinessUnitID: *calibrator.BusinessUnitId,
				}
			}
		}

		var uniqueNextCalibrator []response.NotificationModel
		for _, data := range nCalibrator {
			uniqueNextCalibrator = append(uniqueNextCalibrator, data)
		}

		err = r.notification.NotifySubmittedCalibrationToNextCalibratorsWithoutReview(response.NotificationModel{
			CalibratorID: calibratorID,
			// PreviousCalibratorID: ,
		})
		if err != nil {
			return err
		}

		err = r.notification.NotifyNextCalibrators(uniqueNextCalibrator)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *calibrationUsecase) SendCalibrationsToManager(calibratorID, projectID, prevCalibrator, businessUnit string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID, projectID)
	if err != nil {
		return err
	}

	projectData, err := r.project.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID)
	if err != nil {
		return err
	}

	managerCalibratorIDs, projectPhaseId, err := r.repo.UpdateManagerCalibrations(projectData.UserData, *projectPhase)
	if err != nil {
		return err
	}

	projectPhaseNew, err := r.projectPhase.FindById(projectPhaseId)
	if err != nil {
		return err
	}

	uniqueCalibrator := removeDuplicates(managerCalibratorIDs)
	err = r.notification.NotifyManager(uniqueCalibrator, projectPhaseNew.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) SendBackCalibrationsToOnePhaseBefore(calibratorID, projectID, prevCalibrator, businessUnit string) error {
	projectPhase, err := r.project.FindCalibratorPhase(calibratorID, projectID)
	if err != nil {
		return err
	}

	projectData, err := r.project.FindCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID)
	if err != nil {
		return err
	}

	managerCalibratorIDs, err := r.repo.UpdateCalibrationsOnePhaseBefore(projectData.UserData, *projectPhase)
	if err != nil {
		return err
	}

	err = r.notification.NotifySendBackCalibrators(managerCalibratorIDs)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) SaveCalibrations(payload *request.CalibrationRequest) error {
	return r.repo.SaveChanges(payload)
}

func (r *calibrationUsecase) SaveCommentCalibration(payload *model.Calibration) error {
	return r.repo.SaveCommentCalibration(payload)
}

func (r *calibrationUsecase) SaveScoreAndRating(payload *model.Calibration) error {
	return r.repo.SaveScoreAndRating(payload)
}

func (r *calibrationUsecase) SpmoAcceptApproval(payload *request.AcceptJustification) error {
	projectPhase, err := r.project.FindCalibratorPhase(payload.CalibratorID, payload.ProjectID)
	if err != nil {
		return err
	}

	err = r.repo.AcceptCalibration(payload, projectPhase.Phase.Order)
	if err != nil {
		return err
	}

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

	err = r.notification.NotifyRejectedCalibrationToCalibrator(payload.CalibratorID, employeeName, payload.Comment, payload.ProjectID)
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

	err = r.notification.NotifyNextCalibrators(nextCalibrator)
	if err != nil {
		return err
	}
	return nil
}

func (r *calibrationUsecase) FindSummaryCalibrationBySPMOID(spmoID, projectID string) (response.SummarySPMO, error) {
	results, err := r.repo.GetSummaryBySPMOID(spmoID, projectID)
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
			LastLogin:      d.LastLogin,
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
				data, err := r.repo.GetAllDetailCalibration2BySPMOID(spmoID, calibrationSummary.CalibratorID, smry.BusinessUnitID, projectID, projectPhase.Order)
				if err != nil {
					return response.SummarySPMO{}, err
				}

				allSubmitted := true
				for _, user := range data {
					// fmt.Println("====================================u", user.CalibrationScores, len(user.CalibrationScores)-1, user.Name, user.Nik, projectPhase.Order)
					lastCalibrationStatus := user.CalibrationScores[len(user.CalibrationScores)-1].SpmoStatus
					if lastCalibrationStatus == "Waiting" || (lastCalibrationStatus == "Accepted" && user.CalibrationScores[len(user.CalibrationScores)-1].JustificationReviewStatus == false) {
						status = "Pending"
						allSubmitted = allSubmitted && false
					} else if lastCalibrationStatus == "Accepted" && user.CalibrationScores[len(user.CalibrationScores)-1].JustificationReviewStatus == true {
						allSubmitted = allSubmitted && true
					} else if lastCalibrationStatus == "Rejected" {
						if status == "-" || status == "Pending" {
							status = "Pending"
						} else {
							status = "-"
						}
						allSubmitted = allSubmitted && false
						break
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

func (r *calibrationUsecase) FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, projectID string, order int) ([]response.UserResponse, error) {
	return r.repo.GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID, projectID, order)
}

func (r *calibrationUsecase) FindRatingQuotaSPMOByCalibratorID(spmoID, calibratorID, businessUnitID, projectID string, order int) (*response.RatingQuota, error) {
	users, err := r.FindAllDetailCalibration2bySPMOID(spmoID, calibratorID, businessUnitID, projectID, order)
	if err != nil {
		return nil, err
	}

	projects, err := r.project.FindProjectRatingQuotaByBusinessUnit(businessUnitID, projectID)
	if err != nil {
		return nil, err
	}

	ratingQuota := projects.RatingQuotas[0]
	totalCalibrations := len(users)
	responses := response.RatingQuota{
		APlus: int(math.Round(((ratingQuota.APlusQuota) / float64(100)) * float64(totalCalibrations))),
		A:     int(math.Round(((ratingQuota.AQuota) / float64(100)) * float64(totalCalibrations))),
		BPlus: int(math.Round(((ratingQuota.BPlusQuota) / float64(100)) * float64(totalCalibrations))),
		B:     int(math.Round(((ratingQuota.BQuota) / float64(100)) * float64(totalCalibrations))),
		C:     int(math.Round(((ratingQuota.CQuota) / float64(100)) * float64(totalCalibrations))),
		D:     int(math.Round(((ratingQuota.DQuota) / float64(100)) * float64(totalCalibrations))),
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
			if responses.APlus-(total-totalCalibrations) > 0 {
				responses.APlus -= (total - totalCalibrations)
			} else {
				responses.BPlus -= (total - totalCalibrations)
			}
		} else if ratingQuota.Excess == "A" {
			if responses.A-(total-totalCalibrations) > 0 {
				responses.A -= (total - totalCalibrations)
			} else {
				responses.BPlus -= (total - totalCalibrations)
			}
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
