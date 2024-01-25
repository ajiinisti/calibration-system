package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"gorm.io/gorm"
)

type CalibrationRepo interface {
	Save(payload *model.Calibration) error
	SaveByUser(payload *request.CalibrationForm, project *model.Project) error
	Get(projectID, projectPhaseID, employeeID string) (*model.Calibration, error)
	GetAllPreviousEmployeeCalibrationByActiveProject(employeeID string, phaseOrder int) ([]model.Calibration, error)
	GetByProjectEmployeeID(projectID, employeeID string) ([]model.Calibration, error)
	List() ([]model.Calibration, error)
	GetActiveBySPMOID(spmoID string) ([]model.Calibration, error)
	GetAcceptedBySPMOID(spmoID string) ([]model.Calibration, error)
	GetRejectedBySPMOID(spmoID string) ([]model.Calibration, error)
	Delete(projectId, employeeId string) error
	DeleteCalibrationPhase(projectId, projectPhaseId, employeeId string) error
	Bulksave(payload *[]model.Calibration) error
	BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]*string, []*response.NotificationModel, error)
	UpdateManagerCalibrations(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]string, string, error)
	UpdateCalibrationsOnePhaseBefore(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]string, error)
	SaveChanges(payload *request.CalibrationRequest) error
	AcceptCalibration(payload *request.AcceptJustification, phaseOrder int) error
	AcceptMultipleCalibration(payload *request.AcceptMultipleJustification) error
	RejectCalibration(payload *request.RejectJustification) error
	SubmitReview(payload *request.AcceptMultipleJustification) ([]response.NotificationModel, error)
	GetSummaryBySPMOID(spmoID string) ([]response.SPMOSummaryResult, error)
	GetAllDetailCalibrationBySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error)
	GetCalibrateCalibration() ([]model.Calibration, error)
	GetAllCalibrationByCalibratorID(calibratorId string) ([]model.Calibration, error)
}

type calibrationRepo struct {
	db *gorm.DB
}

func (r *calibrationRepo) GetAllCalibrationByCalibratorID(calibratorId string) ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.
		Table("calibrations c").
		Where("c.calibrator_id = ?", calibratorId).
		Find(&calibrations).
		Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) GetCalibrateCalibration() ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Where("c.status = ? AND c.spmo_status = ? ", "Calibrate", "-").
		Find(&calibrations).
		Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) Save(payload *model.Calibration) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *calibrationRepo) SaveByUser(payload *request.CalibrationForm, project *model.Project) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	excludeProjectPhase := map[string]string{}
	for _, projectPhase := range project.ProjectPhases {
		for _, calibrationData := range payload.CalibrationDataForms {
			if projectPhase.ID != calibrationData.ProjectPhaseID {
				if _, ok := excludeProjectPhase[projectPhase.ID]; !ok {
					excludeProjectPhase[projectPhase.ID] = projectPhase.ID
				}
			} else {
				break
			}
		}
	}

	for _, projectPhaseID := range excludeProjectPhase {
		cal, _ := r.Get(payload.CalibrationDataForms[0].ProjectID, projectPhaseID, payload.CalibrationDataForms[0].EmployeeID)
		if cal != nil {
			err := r.DeleteCalibrationPhase(payload.CalibrationDataForms[0].ProjectID, projectPhaseID, payload.CalibrationDataForms[0].EmployeeID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for index, calibrationData := range payload.CalibrationDataForms {
		data := model.Calibration{
			ProjectID:         calibrationData.ProjectID,
			ProjectPhaseID:    calibrationData.ProjectPhaseID,
			EmployeeID:        calibrationData.EmployeeID,
			CalibratorID:      calibrationData.CalibratorID,
			SpmoID:            calibrationData.SpmoID,
			Spmo2ID:           nil,
			Spmo3ID:           nil,
			Status:            "Waiting",
			SpmoStatus:        "-",
			SpmoComment:       "-",
			JustificationType: "default",
		}

		getCalibration, _ := r.Get(calibrationData.ProjectID, calibrationData.ProjectPhaseID, calibrationData.EmployeeID)
		if getCalibration != nil {
			data.BottomRemark = getCalibration.BottomRemark
			data.TopRemarks = getCalibration.TopRemarks
			data.CalibrationScore = getCalibration.CalibrationScore
			data.CalibrationRating = getCalibration.CalibrationRating
			data.Status = getCalibration.Status
			data.SpmoStatus = getCalibration.SpmoStatus
			data.SpmoComment = getCalibration.SpmoComment
			data.JustificationType = getCalibration.JustificationType
			data.JustificationReviewStatus = getCalibration.JustificationReviewStatus
		}

		if calibrationData.Spmo2ID != "" {
			data.Spmo2ID = &calibrationData.Spmo2ID
		}

		if calibrationData.Spmo3ID != "" {
			data.Spmo3ID = &calibrationData.Spmo3ID
		}

		if index == 0 {
			data.Status = "Calibrate"
		}

		err := tx.Save(&data).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) Bulksave(payload *[]model.Calibration) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		err := tx.Save(&currentBatch).Error
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := tx.Save(&remainingItems).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) Get(projectID, projectPhaseID, employeeID string) (*model.Calibration, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("Project").
		Preload("Employee").
		Preload("ProjectPhase").
		Preload("Calibrator").
		First(&calibration, "project_id = ? AND project_phase_id = ? AND employee_id = ?", projectID, projectPhaseID, employeeID).Error
	if err != nil {
		return nil, err
	}
	return &calibration, nil
}

func (r *calibrationRepo) GetAllPreviousEmployeeCalibrationByActiveProject(employeeID string, phaseOrder int) ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.employee_id = ? AND p.order <= ?", employeeID, phaseOrder).
		Order("p.order ASC").
		Find(&calibrations).
		Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) GetByProjectEmployeeID(projectID, employeeID string) ([]model.Calibration, error) {
	var calibration []model.Calibration
	err := r.db.
		Preload("Employee").
		Preload("Employee.BusinessUnit").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Preload("Calibrator").
		Preload("Calibrator.BusinessUnit").
		Preload("Spmo").
		Preload("Spmo.BusinessUnit").
		Preload("Spmo2").
		Preload("Spmo2.BusinessUnit").
		Preload("Spmo3").
		Preload("Spmo3.BusinessUnit").
		Find(&calibration, "project_id = ? AND employee_id = ?", projectID, employeeID).Error
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *calibrationRepo) GetActiveBySPMOID(id string) ([]model.Calibration, error) {
	var calibration []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("Calibrator").
		Preload("Employee").
		Preload("ProjectPhase", "review_spmo = ?", true).
		Preload("ProjectPhase.Phase").
		Preload("BottomRemark").
		Preload("TopRemarks").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.spmo_id = ? AND c.spmo_status = 'Waiting'", id).
		Order("p.order ASC").
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *calibrationRepo) GetAcceptedBySPMOID(id string) ([]model.Calibration, error) {
	var calibration []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("Calibrator").
		Preload("Employee").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Preload("BottomRemark").
		Preload("TopRemarks").
		Select("c.*").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id AND pp.review_spmo = true").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.spmo_id = ? AND c.spmo_status = 'Accepted' ", id).
		Order("p.order ASC").
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *calibrationRepo) GetRejectedBySPMOID(id string) ([]model.Calibration, error) {
	var calibration []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("Calibrator").
		Preload("Employee").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Preload("BottomRemark").
		Preload("TopRemarks").
		Select("c.*").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id AND pp.review_spmo = true").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.spmo_id = ? AND c.spmo_status = 'Rejected' ", id).
		Order("p.order ASC").
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}
	return calibration, nil
}

func (r *calibrationRepo) List() ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.Preload("Project").Preload("Employee").Preload("ProjectPhase").Preload("Calibrator").Find(&calibrations).Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) Delete(projectId, employeeId string) error {
	result := r.db.Unscoped().Where("project_id = ? AND employee_id = ?", projectId, employeeId).Delete(&model.Calibration{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Calibration not found!")
	}
	return nil
}

func (r *calibrationRepo) DeleteCalibrationPhase(projectId, projectPhaseId, employeeId string) error {
	result := r.db.Unscoped().Where("project_id = ? AND employee_id = ? AND project_phase_id = ?", projectId, employeeId, projectPhaseId).Delete(&model.Calibration{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Calibration not found!")
	}
	return nil
}

func (r *calibrationRepo) BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]*string, []*response.NotificationModel, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reviewSPMO := false
	var spmoID []*string

	var employeeCalibrationScore []*model.Calibration
	var nextCalibrator []*response.NotificationModel
	for _, calibrations := range payload.RequestData {
		spmoID = append(spmoID, &calibrations.SpmoID)
		spmoID = append(spmoID, calibrations.Spmo2ID)
		spmoID = append(spmoID, calibrations.Spmo3ID)

		if projectPhase.ReviewSpmo == true {
			calibrations.SpmoStatus = "Waiting"
			reviewSPMO = true
		}
		calibrations.Status = "Complete"
		employeeCalibrationScore = append(employeeCalibrationScore, calibrations)

		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return nil, nil, result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, nil, fmt.Errorf("Calibrations not found!")
		}
	}

	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []*model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Preload("ProjectPhase").
			Preload("ProjectPhase.Phase").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		for _, c := range calibrations {
			c.CalibrationScore = employeeCalibration.CalibrationScore
			c.CalibrationRating = employeeCalibration.CalibrationRating
			c.Status = "Waiting"
			err := tx.Updates(c).Error
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}
		}

		if !reviewSPMO && len(calibrations) > 0 {
			nextCalibrator = append(nextCalibrator, &response.NotificationModel{
				CalibratorID: calibrations[0].CalibratorID,
				ProjectPhase: calibrations[0].ProjectPhase.Phase.Order,
				Deadline:     calibrations[0].ProjectPhase.EndDate,
			})
			calibrations[0].Status = "Calibrate"
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return nil, nil, err
			}
		}
	}

	tx.Commit()
	return spmoID, nextCalibrator, nil
}

func (r *calibrationRepo) UpdateManagerCalibrations(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]string, string, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employeeCalibrationScore []*model.Calibration
	for _, calibrations := range payload.RequestData {
		employeeCalibrations, err := r.GetAllPreviousEmployeeCalibrationByActiveProject(calibrations.EmployeeID, projectPhase.Phase.Order)
		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		if len(employeeCalibrations) == 2 {
			calibrations.Status = "Waiting"
			employeeCalibrationScore = append(employeeCalibrationScore, calibrations)
		}

		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return nil, "", result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, "", fmt.Errorf("Calibrations not found!")
		}
	}

	var ppId string
	var managerCalibratorIDs []string
	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []*model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Preload("Calibrator").
			Where("projects.active = true AND phases.order < ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		if len(calibrations) > 0 {
			ppId = calibrations[0].ProjectPhaseID
			calibrations[0].Status = "Calibrate"
			managerCalibratorIDs = append(managerCalibratorIDs, calibrations[0].CalibratorID)
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return nil, "", err
			}
		}
	}

	tx.Commit()
	return managerCalibratorIDs, ppId, nil
}

func (r *calibrationRepo) UpdateCalibrationsOnePhaseBefore(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) ([]string, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employeeCalibrationScore []*model.Calibration
	for _, calibrations := range payload.RequestData {
		employeeCalibrations, err := r.GetAllPreviousEmployeeCalibrationByActiveProject(calibrations.EmployeeID, projectPhase.Phase.Order)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(employeeCalibrations) > 1 {
			calibrations.Status = "Waiting"
			calibrations.JustificationReviewStatus = false
			employeeCalibrationScore = append(employeeCalibrationScore, calibrations)
		}

		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("Calibrations not found!")
		}
	}

	var managerCalibratorIDs []string
	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []*model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Preload("Calibrator").
			Where("projects.active = true AND phases.order < ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(calibrations) > 0 {
			managerCalibratorIDs = append(managerCalibratorIDs, calibrations[len(calibrations)-1].CalibratorID)
			err := tx.Model(&model.Calibration{}).
				Where("project_id = ? AND project_phase_id = ? AND employee_id = ? AND calibrator_id = ?",
					calibrations[len(calibrations)-1].ProjectID,
					calibrations[len(calibrations)-1].ProjectPhaseID,
					calibrations[len(calibrations)-1].EmployeeID,
					calibrations[len(calibrations)-1].CalibratorID,
				).Updates(map[string]interface{}{
				"status":                      "Calibrate",
				"send_back_deadline":          projectPhase.EndDate,
				"justification_review_status": false,
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	tx.Commit()
	return managerCalibratorIDs, nil
}

func (r *calibrationRepo) SaveChanges(payload *request.CalibrationRequest) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, calibrations := range payload.RequestData {
		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Calibrations not found!")
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) AcceptMultipleCalibration(payload *request.AcceptMultipleJustification) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, justification := range payload.ArrayOfAcceptsJustification {
		err := tx.Updates(&model.Calibration{
			ProjectID:      justification.ProjectID,
			ProjectPhaseID: justification.ProjectPhaseID,
			EmployeeID:     justification.EmployeeID,
			SpmoStatus:     "Accepted",
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) AcceptCalibration(payload *request.AcceptJustification, phaseOrder int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Updates(&model.Calibration{
		ProjectID:      payload.ProjectID,
		ProjectPhaseID: payload.ProjectPhaseID,
		EmployeeID:     payload.EmployeeID,
		SpmoStatus:     "Accepted",
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *calibrationRepo) RejectCalibration(payload *request.RejectJustification) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&model.Calibration{}).
		Where("project_id = ? AND project_phase_id = ? AND employee_id = ? AND calibrator_id = ?",
			payload.ProjectID,
			payload.ProjectPhaseID,
			payload.EmployeeID,
			payload.CalibratorID,
		).Updates(map[string]interface{}{
		"spmo_status":  "Rejected",
		"spmo_comment": payload.Comment,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) SubmitReview(payload *request.AcceptMultipleJustification) ([]response.NotificationModel, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	mapResult := make(map[string]response.NotificationModel)
	for _, justification := range payload.ArrayOfAcceptsJustification {
		err := tx.Model(&model.Calibration{}).
			Where("project_id = ? AND project_phase_id = ? AND employee_id = ? AND calibrator_id = ?",
				justification.ProjectID,
				justification.ProjectPhaseID,
				justification.EmployeeID,
				justification.CalibratorID,
			).Updates(map[string]interface{}{
			"justification_review_status": true,
		}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var projectPhase *model.ProjectPhase
		err = tx.Table("project_phases").
			Select("project_phases.*").
			Preload("Phase").
			Where("project_phases.id= ?", justification.ProjectPhaseID).
			First(&projectPhase).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var calibrations []*model.Calibration
		err = tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, justification.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(calibrations) > 0 {
			calibrations[0].Status = "Calibrate"
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return nil, err
			}

			if _, ok := mapResult[calibrations[0].CalibratorID]; !ok {
				fmt.Println("DATA PROJECT", calibrations[0].ProjectPhaseID)
				var projectPhaseNextCalibrator *model.ProjectPhase
				err := tx.Table("project_phases").
					Select("project_phases.*").
					Preload("Phase").
					Where("project_phases.id= ?", calibrations[0].ProjectPhaseID).
					First(&projectPhaseNextCalibrator).Error
				if err != nil {
					tx.Rollback()
					return nil, err
				}

				var prevCal *model.User
				err = tx.Where("id = ?", justification.CalibratorID).First(&prevCal).Error
				if err != nil {
					return nil, err
				}
				mapResult[calibrations[0].CalibratorID] = response.NotificationModel{
					CalibratorID:           calibrations[0].CalibratorID,
					ProjectPhase:           projectPhaseNextCalibrator.Phase.Order,
					Deadline:               projectPhaseNextCalibrator.EndDate,
					PreviousCalibrator:     prevCal.Name,
					PreviousCalibratorID:   prevCal.ID,
					PreviousBusinessUnitID: *prevCal.BusinessUnitId,
				}
			}
		}
	}

	fmt.Println("DATA NEXT", mapResult)
	var nextCalibrator []response.NotificationModel
	for _, data := range mapResult {
		nextCalibrator = append(nextCalibrator, data)
	}

	tx.Commit()
	return nextCalibrator, nil
}

func (r *calibrationRepo) GetSummaryBySPMOID(spmoID string) ([]response.SPMOSummaryResult, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Assuming db is your GORM database instance
	var results []response.SPMOSummaryResult
	err := tx.Table("calibrations c").
		Select("COUNT(c.*) as count, u.business_unit_id, b.name as business_unit_name, u2.name as calibrator_name, c.calibrator_id, c.project_phase_id, p.order").
		Joins("JOIN project_phases pp on pp.id = c.project_phase_id").
		Joins("JOIN phases p on pp.phase_id = p.id").
		Joins("JOIN users u on c.employee_id = u.id").
		Joins("JOIN business_units b on u.business_unit_id = b.id").
		Joins("JOIN users u2 on c.calibrator_id = u2.id").
		Joins("JOIN projects pr on pr.id = c.project_id AND pr.active = true").
		Where("(spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?) AND p.order NOT IN (SELECT MAX(\"order\") FROM phases)", spmoID, spmoID, spmoID).
		Group("u.business_unit_id, b.name, u2.name, c.calibrator_id, c.project_phase_id, p.order").
		Order("p.order ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return results, nil
}

func (r *calibrationRepo) GetAllDetailCalibrationBySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error) {
	var calibration []response.UserResponse
	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj1 ON actual_scores.project_id = proj1.id").
				Where("proj1.active = ?", true)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.active = true AND p.order <= ?", order).
				Order("p.order ASC")
		}).
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Select("u.*, u2.name as supervisor_names").
		Joins("JOIN business_units b ON u.business_unit_id = b.id AND b.id = ?", businessUnitID).
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?) AND c1.calibrator_id = ?", spmoID, spmoID, spmoID, calibratorID).
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("LEFT JOIN users u2 ON u.supervisor_nik = u2.nik").
		Where("u.department = ?", department).
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}

	return calibration, nil
}

func (r *calibrationRepo) GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error) {
	var calibration []model.User
	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj1 ON actual_scores.project_id = proj1.id AND actual_scores.deleted_at IS NULL").
				Where("proj1.active = ?", true)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.active = true AND p.order <= ?", order).
				Order("p.order ASC")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Select("u.*, u2.name as supervisor_names").
		Joins("JOIN business_units b ON u.business_unit_id = b.id AND b.id = ?", businessUnitID).
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?) AND c1.calibrator_id = ? AND c1.deleted_at IS NULL", spmoID, spmoID, spmoID, calibratorID).
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("LEFT JOIN users u2 ON u.supervisor_nik = u2.nik").
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}

	var calibrations []response.UserResponse
	for _, user := range calibration {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return nil, err
		}

		dataOneResponse := &response.UserResponse{
			BaseModel: model.BaseModel{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			CreatedBy:       user.CreatedBy,
			UpdatedBy:       user.UpdatedBy,
			Email:           user.Email,
			Name:            user.Name,
			Nik:             user.Nik,
			SupervisorNames: supervisorName,
			BusinessUnit: response.BusinessUnitResponse{
				ID:                  user.BusinessUnit.ID,
				Status:              user.BusinessUnit.Status,
				Name:                user.BusinessUnit.Name,
				GroupBusinessUnitId: user.BusinessUnit.GroupBusinessUnitId,
			},
			BusinessUnitId:   user.BusinessUnitId,
			OrganizationUnit: user.OrganizationUnit,
			Division:         user.Division,
			Department:       user.Department,
			Grade:            user.Grade,
			Position:         user.Position,
			Roles:            user.Roles,
			ActualScores: []response.ActualScoreResponse{{
				ProjectID:    user.ActualScores[0].ProjectID,
				EmployeeID:   user.ActualScores[0].EmployeeID,
				ActualScore:  user.ActualScores[0].ActualScore,
				ActualRating: user.ActualScores[0].ActualRating,
				Y1Rating:     user.ActualScores[0].Y1Rating,
				Y2Rating:     user.ActualScores[0].Y2Rating,
				PTTScore:     user.ActualScores[0].PTTScore,
				PATScore:     user.ActualScores[0].PATScore,
				Score360:     user.ActualScores[0].Score360,
			}},
			CalibrationScores: []response.CalibrationResponse{},
			ScoringMethod:     user.ScoringMethod,
			Directorate:       user.Directorate,
		}

		for _, data := range user.CalibrationScores {
			topRemarks := []response.TopRemarkResponse{}
			for _, topRemark := range data.TopRemarks {
				topRemarks = append(topRemarks, response.TopRemarkResponse{
					BaseModel:      topRemark.BaseModel,
					ProjectID:      topRemark.ProjectID,
					EmployeeID:     topRemark.EmployeeID,
					ProjectPhaseID: topRemark.ProjectPhaseID,
					Initiative:     topRemark.Initiative,
					Description:    topRemark.Description,
					Result:         topRemark.Result,
					StartDate:      topRemark.StartDate,
					EndDate:        topRemark.EndDate,
					Comment:        topRemark.Comment,
					EvidenceName:   topRemark.EvidenceName,
				})
			}
			dataOneResponse.CalibrationScores = append(dataOneResponse.CalibrationScores, response.CalibrationResponse{
				ProjectID: data.ProjectID,
				ProjectPhase: response.ProjectPhaseResponse{
					Phase: response.PhaseResponse{
						Order: data.ProjectPhase.Phase.Order,
					},
					StartDate: data.ProjectPhase.StartDate,
					EndDate:   data.ProjectPhase.EndDate,
				},
				ProjectPhaseID: data.ProjectPhaseID,
				EmployeeID:     data.EmployeeID,
				Calibrator: response.CalibratorResponse{
					Name: data.Calibrator.Name,
				},
				CalibratorID:              data.CalibratorID,
				CalibrationScore:          data.CalibrationScore,
				CalibrationRating:         data.CalibrationRating,
				Status:                    data.Status,
				SpmoStatus:                data.SpmoStatus,
				Comment:                   data.Comment,
				SpmoComment:               data.SpmoComment,
				JustificationType:         data.JustificationType,
				JustificationReviewStatus: data.JustificationReviewStatus,
				SendBackDeadline:          data.SendBackDeadline,
				BottomRemark: response.BottomRemarkResponse{
					ProjectID:      data.BottomRemark.ProjectID,
					EmployeeID:     data.BottomRemark.EmployeeID,
					ProjectPhaseID: data.BottomRemark.ProjectPhaseID,
					LowPerformance: data.BottomRemark.LowPerformance,
					Indisipliner:   data.BottomRemark.Indisipliner,
					Attitude:       data.BottomRemark.Attitude,
					WarningLetter:  data.BottomRemark.WarningLetter,
				},
				TopRemarks: topRemarks,
			})
		}

		calibrations = append(calibrations, *dataOneResponse)
	}

	return calibrations, nil
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
