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
	GetAllPreviousEmployeeCalibrationByActiveProject(employeeID, projectID string, phaseOrder int) ([]model.Calibration, error)
	GetByProjectEmployeeID(projectID, employeeID string) ([]model.CalibrationForm, error)
	List() ([]model.Calibration, error)
	GetActiveUserBySPMOID(spmoID string) ([]model.UserChange, error)
	GetAcceptedBySPMOID(spmoID string) ([]model.Calibration, error)
	GetRejectedBySPMOID(spmoID string) ([]model.Calibration, error)
	Delete(projectId, employeeId string) error
	DeleteCalibrationPhase(projectId, projectPhaseId, employeeId string) error
	Bulksave(payload *[]model.Calibration) error
	BulkUpdate(payload []response.UserResponse, projectPhase model.ProjectPhase, projectID string) ([]string, []*response.NotificationModel, error)
	UpdateManagerCalibrations(payload []response.UserResponse, projectPhase model.ProjectPhase) ([]string, string, error)
	UpdateCalibrationsOnePhaseBefore(payload []response.UserResponse, projectPhase model.ProjectPhase) ([]response.NotificationModel, error)
	SaveChanges(payload *request.CalibrationRequest) error
	SaveCommentCalibration(payload *model.Calibration) error
	SaveScoreAndRating(payload *model.Calibration) error
	AcceptCalibration(payload *request.AcceptJustification, phaseOrder int) error
	AcceptMultipleCalibration(payload *request.AcceptMultipleJustification) error
	RejectCalibration(payload *request.RejectJustification) error
	SubmitReview(payload *request.AcceptMultipleJustification) ([]response.NotificationModel, error)
	GetSummaryBySPMOID(spmoID, projectID string) ([]response.SPMOSummaryResult, error)
	GetAllDetailCalibrationBySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID, projectID string, order int) ([]response.UserResponse, error)
	GetCalibrateCalibration() ([]model.Calibration, error)
	GetAllCalibrationByCalibratorID(calibratorId string) ([]model.Calibration, error)
	GetLatestJustification(projectID, calibratorID, employeeID string) ([]model.SeeCalibrationJustification, error)
	CheckConditionBeforeSubmitCalibration(projectID string, payload []response.UserResponse,
		projectPhase model.ProjectPhase, countCalibrated response.TotalCalibratedRating, countRatingQuota response.RatingQuota,
	) (bool, error)
}

type calibrationRepo struct {
	db *gorm.DB
}

func (r *calibrationRepo) GetLatestJustification(projectID, calibratorID, employeeID string) ([]model.SeeCalibrationJustification, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	var calibrations []model.SeeCalibrationJustification
	err = r.db.
		Table("calibrations c").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Preload("Calibrator").
		Preload("TopRemarks").
		Preload("BottomRemark").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.employee_id = ? AND c.project_id = ? AND p.order <= ?", employeeID, projectID, calibration.ProjectPhase.Phase.Order).
		Order("").
		Find(&calibrations).
		Error

	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) GetAllCalibrationByCalibratorID(calibratorId string) ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.
		Table("calibrations c").
		Joins("JOIN projects pr ON pr.id = c.project_id AND pr.active = true").
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

	projectPhaseStartIndex := 0
	var projectProjectPhase model.ProjectPhase
	err := r.db.Preload("Phase").First(&projectProjectPhase, "id = ?", payload.CalibrationDataForms[0].ProjectPhaseID).Error
	if err != nil {
		return err
	}

	if projectProjectPhase.Phase.Order == 1 {
		projectPhaseStartIndex = 1
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

		if index == projectPhaseStartIndex {
			data.Status = "Calibrate"
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
			data.FilledTopBottomMark = getCalibration.FilledTopBottomMark
			data.Comment = getCalibration.Comment
		}

		if calibrationData.Spmo2ID != "" {
			data.Spmo2ID = &calibrationData.Spmo2ID
		}

		if calibrationData.Spmo3ID != "" {
			data.Spmo3ID = &calibrationData.Spmo3ID
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

	go func() {
		err := r.db.Exec("REFRESH MATERIALIZED VIEW materialized_user_view;").Error
		if err != nil {
			fmt.Printf("Failed to refresh materialized view: %v", err)
		}
	}()
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

func (r *calibrationRepo) GetAllPreviousEmployeeCalibrationByActiveProject(employeeID, projectID string, phaseOrder int) ([]model.Calibration, error) {
	var calibrations []model.Calibration
	err := r.db.
		Table("calibrations c").
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c.employee_id = ? AND p.order <= ? AND c.project_id = ?", employeeID, phaseOrder, projectID).
		Order("p.order ASC").
		Find(&calibrations).
		Error
	if err != nil {
		return nil, err
	}
	return calibrations, nil
}

func (r *calibrationRepo) GetByProjectEmployeeID(projectID, employeeID string) ([]model.CalibrationForm, error) {
	var calibration []model.CalibrationForm
	err := r.db.
		Table("calibrations").
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

func (r *calibrationRepo) GetActiveUserBySPMOID(spmoID string) ([]model.UserChange, error) {
	var calibration []model.UserChange
	err := r.db.
		Table("users u").
		Select("u.id as id, u.email as email, u.name as name, u.division as division, u.nik as nik, b.name as business_unit_name").
		Joins("JOIN business_units b on u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON (c1.employee_id = u.id OR c1.calibrator_id = u.id) AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?) AND c1.deleted_at IS NULL", spmoID, spmoID, spmoID).
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("LEFT JOIN users u2 ON u.supervisor_nik = u2.nik").
		Distinct().
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

func (r *calibrationRepo) BulkUpdate(payload []response.UserResponse, projectPhase model.ProjectPhase, projectID string) ([]string, []*response.NotificationModel, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reviewSPMO := false
	var spmoID []string

	var employeeCalibrationScore []*model.Calibration
	var nextCalibrator []*response.NotificationModel
	for _, dataPayload := range payload {
		var getCalibration *model.Calibration
		err := tx.Table("calibrations c").
			Select("c.*").
			Where("c.project_id = ? AND c.project_phase_id = ? AND c.employee_id = ?", projectID, projectPhase.ID, dataPayload.ID).
			First(&getCalibration).Error

		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}

		spmoID = append(spmoID, getCalibration.SpmoID)
		if getCalibration.Spmo2ID != nil {
			spmoID = append(spmoID, *getCalibration.Spmo2ID)
		}
		if getCalibration.Spmo3ID != nil {
			spmoID = append(spmoID, *getCalibration.Spmo3ID)
		}

		if projectPhase.ReviewSpmo {
			getCalibration.SpmoStatus = "Waiting"
			reviewSPMO = true
		}
		getCalibration.Status = "Complete"
		employeeCalibrationScore = append(employeeCalibrationScore, getCalibration)

		result := tx.Updates(getCalibration)
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
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("calibrations.project_id = ? AND phases.order > ? AND calibrations.employee_id = ?", employeeCalibration.ProjectID, projectPhase.Phase.Order, employeeCalibration.EmployeeID).
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
				ProjectID:    calibrations[0].ProjectID,
				Deadline:     calibrations[0].ProjectPhase.EndDate,
			})
			err := tx.Updates(&model.Calibration{
				ProjectID:      calibrations[0].ProjectID,
				ProjectPhaseID: calibrations[0].ProjectPhaseID,
				EmployeeID:     calibrations[0].EmployeeID,
				Comment:        employeeCalibration.Comment,
				Status:         "Calibrate",
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}
			if err != nil {
				tx.Rollback()
				return nil, nil, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := r.db.Exec("REFRESH MATERIALIZED VIEW materialized_user_view;").Error
		if err != nil {
			fmt.Printf("Failed to refresh materialized view: %v", err)
		}
	}()
	return spmoID, nextCalibrator, nil
}

func (r *calibrationRepo) UpdateManagerCalibrations(payload []response.UserResponse, projectPhase model.ProjectPhase) ([]string, string, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employeeCalibrationScore []*model.Calibration
	for _, user := range payload {
		employeeCalibrations, err := r.GetAllPreviousEmployeeCalibrationByActiveProject(user.ID, projectPhase.ProjectID, projectPhase.Phase.Order)
		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		var calibration *model.Calibration
		err = tx.Table("calibrations").
			Select("calibrations.*").
			Where("project_id = ? AND project_phase_id = ? AND employee_id = ?", projectPhase.ProjectID, projectPhase.ID, user.ID).
			First(&calibration).Error
		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		if len(employeeCalibrations) == 2 {
			calibration.Status = "Waiting"
			calibration.SpmoStatus = "-"
			calibration.JustificationReviewStatus = false
			employeeCalibrationScore = append(employeeCalibrationScore, calibration)
		}

		result := tx.Updates(calibration)
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
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Preload("Calibrator").
			Where("calibrations.project_id = ? AND phases.order < ? AND calibrations.employee_id = ?", projectPhase.ProjectID, projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, "", err
		}

		if len(calibrations) > 0 {
			ppId = calibrations[0].ProjectPhaseID
			managerCalibratorIDs = append(managerCalibratorIDs, calibrations[0].CalibratorID)
			// err := tx.Model(&model.Calibration{}).
			// 	Where("project_id = ? AND project_phase_id = ? AND employee_id = ? AND calibrator_id = ?",
			// 		calibrations[0].ProjectID,
			// 		calibrations[0].ProjectPhaseID,
			// 		calibrations[0].EmployeeID,
			// 		calibrations[0].CalibratorID,
			// 	).Updates(map[string]interface{}{
			// 	"spmo_status":                 "-",
			// 	"status":                      "Calibrate",
			// 	"justification_review_status": false,
			// }).Error
			err := tx.Updates(&model.Calibration{
				ProjectID:                 calibrations[0].ProjectID,
				ProjectPhaseID:            calibrations[0].ProjectPhaseID,
				EmployeeID:                calibrations[0].EmployeeID,
				Status:                    "Calibrate",
				JustificationReviewStatus: false,
				SpmoStatus:                "-",
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, "", err
			}
		}
	}

	tx.Commit()
	return managerCalibratorIDs, ppId, nil
}

func (r *calibrationRepo) UpdateCalibrationsOnePhaseBefore(payload []response.UserResponse, projectPhase model.ProjectPhase) ([]response.NotificationModel, error) {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var employeeCalibrationScore []*model.Calibration
	for _, user := range payload {
		employeeCalibrations, err := r.GetAllPreviousEmployeeCalibrationByActiveProject(user.ID, projectPhase.ProjectID, projectPhase.Phase.Order)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var calibration *model.Calibration
		err = tx.Table("calibrations").
			Select("calibrations.*").
			Where("project_id = ? AND project_phase_id = ? AND employee_id = ?", projectPhase.ProjectID, projectPhase.ID, user.ID).
			First(&calibration).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(employeeCalibrations) > 1 {
			calibration.Status = "Waiting"
			calibration.JustificationReviewStatus = false
			calibration.SpmoStatus = "-"
			employeeCalibrationScore = append(employeeCalibrationScore, calibration)
		}

		result := tx.Updates(calibration)
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return nil, fmt.Errorf("Calibrations not found!")
		}
	}

	mapResult := make(map[string]response.NotificationModel)
	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []*model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("calibrations.project_id = ? AND phases.order < ? AND calibrations.employee_id = ?",
				projectPhase.ProjectID, projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(calibrations) > 0 {
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
				"spmo_status":                 "-",
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			if len(calibrations) > 2 {
				if _, ok := mapResult[calibrations[len(calibrations)-1].CalibratorID]; !ok {
					var projectPhaseNextCalibrator *model.ProjectPhase
					err := tx.Table("project_phases").
						Select("project_phases.*").
						Preload("Phase").
						Where("project_phases.id= ?", calibrations[len(calibrations)-1].ProjectPhaseID).
						First(&projectPhaseNextCalibrator).Error
					if err != nil {
						tx.Rollback()
						return nil, err
					}

					var nextCal *model.User
					err = tx.Where("id = ?", employeeCalibration.CalibratorID).First(&nextCal).Error
					if err != nil {
						return nil, err
					}

					var prevCal *model.User
					err = tx.Where("id = ?", calibrations[len(calibrations)-2].CalibratorID).First(&prevCal).Error
					if err != nil {
						return nil, err
					}
					mapResult[calibrations[len(calibrations)-1].CalibratorID] = response.NotificationModel{
						CalibratorID:           calibrations[len(calibrations)-1].CalibratorID,
						ProjectID:              calibrations[len(calibrations)-1].ProjectID,
						ProjectPhase:           projectPhaseNextCalibrator.Phase.Order,
						Deadline:               projectPhaseNextCalibrator.EndDate,
						NextCalibrator:         nextCal.Name,
						PreviousCalibrator:     prevCal.Name,
						PreviousCalibratorID:   prevCal.ID,
						PreviousBusinessUnitID: *prevCal.BusinessUnitId,
					}
				}
			} else {
				if _, ok := mapResult[calibrations[len(calibrations)-1].CalibratorID]; !ok {
					var projectPhaseNextCalibrator *model.ProjectPhase
					err := tx.Table("project_phases").
						Select("project_phases.*").
						Preload("Phase").
						Where("project_phases.id= ?", calibrations[len(calibrations)-1].ProjectPhaseID).
						First(&projectPhaseNextCalibrator).Error
					if err != nil {
						tx.Rollback()
						return nil, err
					}

					var nextCal *model.User
					err = tx.Where("id = ?", employeeCalibration.CalibratorID).First(&nextCal).Error
					if err != nil {
						return nil, err
					}

					mapResult[calibrations[len(calibrations)-1].CalibratorID] = response.NotificationModel{
						CalibratorID:   calibrations[len(calibrations)-1].CalibratorID,
						ProjectID:      calibrations[len(calibrations)-1].ProjectID,
						ProjectPhase:   projectPhaseNextCalibrator.Phase.Order,
						Deadline:       projectPhaseNextCalibrator.EndDate,
						NextCalibrator: nextCal.Name,
					}
				}
			}

		}
	}

	var nextCalibrator []response.NotificationModel
	for _, data := range mapResult {
		nextCalibrator = append(nextCalibrator, data)
	}

	tx.Commit()
	return nextCalibrator, nil
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

func (r *calibrationRepo) SaveCommentCalibration(payload *model.Calibration) error {
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
		"comment": payload.Comment,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var phases *model.ProjectPhase
	err = tx.Table("project_phases").
		Preload("Phase").
		Select("project_phases.*").
		Where("id = ?", payload.ProjectPhaseID).
		Find(&phases).Error

	if err != nil {
		return err
	}

	var calibrations []*model.Calibration
	err = tx.Table("calibrations").
		Select("calibrations.*").
		Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
		Joins("JOIN phases ON phases.id = project_phases.phase_id").
		Where("calibrations.project_id = ? AND phases.order > ? AND calibrations.employee_id = ?", payload.ProjectID, phases.Phase.Order, payload.EmployeeID).
		Order("phases.order ASC").
		Find(&calibrations).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, c := range calibrations {
		c.Comment = payload.Comment
		err := tx.Updates(c).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (r *calibrationRepo) SaveScoreAndRating(payload *model.Calibration) error {
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
		"calibration_score":      payload.CalibrationScore,
		"calibration_rating":     payload.CalibrationRating,
		"justification_type":     payload.JustificationType,
		"filled_top_bottom_mark": payload.FilledTopBottomMark,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var phases *model.ProjectPhase
	err = tx.Table("project_phases").
		Preload("Phase").
		Select("project_phases.*").
		Where("id = ?", payload.ProjectPhaseID).
		Find(&phases).Error

	if err != nil {
		return err
	}

	var calibrations []*model.Calibration
	err = tx.Table("calibrations").
		Select("calibrations.*").
		Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
		Joins("JOIN phases ON phases.id = project_phases.phase_id").
		Where("calibrations.project_id = ? AND phases.order > ? AND calibrations.employee_id = ?", payload.ProjectID, phases.Phase.Order, payload.EmployeeID).
		Order("phases.order ASC").
		Find(&calibrations).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, c := range calibrations {
		c.CalibrationRating = payload.CalibrationRating
		c.CalibrationScore = payload.CalibrationScore
		err := tx.Updates(c).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	go func() {
		err := r.db.Exec("REFRESH MATERIALIZED VIEW materialized_user_view;").Error
		if err != nil {
			fmt.Printf("Failed to refresh materialized view: %v", err)
		}
	}()
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
		err := tx.Updates(&model.Calibration{
			ProjectID:                 justification.ProjectID,
			ProjectPhaseID:            justification.ProjectPhaseID,
			EmployeeID:                justification.EmployeeID,
			JustificationReviewStatus: true,
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

		// fmt.Println('1', calibrations[0])

		var calibrationBefore *model.Calibration
		err = tx.Table("calibrations").
			Select("calibrations.*").
			Where("employee_id = ? AND calibrator_id = ? and project_id = ?", justification.EmployeeID, justification.CalibratorID, justification.ProjectID).
			First(&calibrationBefore).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// fmt.Println("2", calibrationBefore.JustificationReviewStatus, len(calibrations))

		if len(calibrations) > 0 {
			err := tx.Updates(&model.Calibration{
				ProjectID:      calibrations[0].ProjectID,
				ProjectPhaseID: calibrations[0].ProjectPhaseID,
				EmployeeID:     calibrations[0].EmployeeID,
				Comment:        calibrationBefore.Comment,
				Status:         "Calibrate",
			}).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			if _, ok := mapResult[calibrations[0].CalibratorID]; !ok {
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
					tx.Rollback()
					return nil, err
				}

				if prevCal.BusinessUnitId != nil {
					mapResult[calibrations[0].CalibratorID] = response.NotificationModel{
						CalibratorID:           calibrations[0].CalibratorID,
						ProjectID:              calibrations[0].ProjectID,
						ProjectPhase:           projectPhaseNextCalibrator.Phase.Order,
						Deadline:               projectPhaseNextCalibrator.EndDate,
						PreviousCalibrator:     prevCal.Name,
						PreviousCalibratorID:   prevCal.ID,
						PreviousBusinessUnitID: *prevCal.BusinessUnitId,
					}
				} else {
					mapResult[calibrations[0].CalibratorID] = response.NotificationModel{
						CalibratorID:         calibrations[0].CalibratorID,
						ProjectID:            calibrations[0].ProjectID,
						ProjectPhase:         projectPhaseNextCalibrator.Phase.Order,
						Deadline:             projectPhaseNextCalibrator.EndDate,
						PreviousCalibrator:   prevCal.Name,
						PreviousCalibratorID: prevCal.ID,
					}
				}
			}
		}
	}
	tx.Commit()

	var nextCalibrator []response.NotificationModel
	for _, data := range mapResult {
		nextCalibrator = append(nextCalibrator, data)
	}

	return nextCalibrator, nil
}

func (r *calibrationRepo) GetSummaryBySPMOID(spmoID, projectID string) ([]response.SPMOSummaryResult, error) {
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
		Joins("JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("JOIN phases p ON pp.phase_id = p.id").
		Joins("JOIN users u ON c.employee_id = u.id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN users u2 ON c.calibrator_id = u2.id").
		Joins("JOIN projects pr ON pr.id = c.project_id").
		Where("pr.id = ? AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?)", projectID, spmoID, spmoID, spmoID).
		Where("p.order < 6").
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

func (r *calibrationRepo) GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID, projectID string, order int) ([]response.UserResponse, error) {
	var calibrations []response.UserResponse
	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj1 ON actual_scores.project_id = proj1.id AND actual_scores.deleted_at IS NULL").
				Where("proj1.id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, order).
				Order("p.order ASC")
		}).
		// Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Select("u.*, u2.name as supervisor_names").
		Distinct().
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?) AND c1.calibrator_id = ? AND c1.deleted_at IS NULL AND c1.project_id = ?", spmoID, spmoID, spmoID, calibratorID, projectID).
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id AND p.order = ?", order).
		Joins("LEFT JOIN users u2 ON u.supervisor_nik = u2.nik").
		Where("u.business_unit_id = ?", businessUnitID).
		Find(&calibrations).Error
	if err != nil {
		return nil, err
	}

	// var calibrations []response.UserResponse
	// for _, user := range calibration {
	// 	var supervisorName string
	// 	err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	dataOneResponse := &response.UserResponse{
	// 		BaseModel: model.BaseModel{
	// 			ID:        user.ID,
	// 			CreatedAt: user.CreatedAt,
	// 			UpdatedAt: user.UpdatedAt,
	// 			DeletedAt: user.DeletedAt,
	// 		},
	// 		CreatedBy:       user.CreatedBy,
	// 		UpdatedBy:       user.UpdatedBy,
	// 		Email:           user.Email,
	// 		Name:            user.Name,
	// 		Nik:             user.Nik,
	// 		SupervisorNames: supervisorName,
	// 		BusinessUnit: response.BusinessUnitResponse{
	// 			ID:                  user.BusinessUnit.ID,
	// 			Status:              user.BusinessUnit.Status,
	// 			Name:                user.BusinessUnit.Name,
	// 			GroupBusinessUnitId: user.BusinessUnit.GroupBusinessUnitId,
	// 		},
	// 		BusinessUnitId:   user.BusinessUnitId,
	// 		OrganizationUnit: user.OrganizationUnit,
	// 		Division:         user.Division,
	// 		Department:       user.Department,
	// 		Grade:            user.Grade,
	// 		Position:         user.Position,
	// 		Roles:            user.Roles,
	// 		ActualScores: []response.ActualScoreResponse{{
	// 			ProjectID:    user.ActualScores[0].ProjectID,
	// 			EmployeeID:   user.ActualScores[0].EmployeeID,
	// 			ActualScore:  user.ActualScores[0].ActualScore,
	// 			ActualRating: user.ActualScores[0].ActualRating,
	// 			Y1Rating:     user.ActualScores[0].Y1Rating,
	// 			Y2Rating:     user.ActualScores[0].Y2Rating,
	// 			PTTScore:     user.ActualScores[0].PTTScore,
	// 			PATScore:     user.ActualScores[0].PATScore,
	// 			Score360:     user.ActualScores[0].Score360,
	// 		}},
	// 		CalibrationScores: []response.CalibrationResponse{},
	// 		ScoringMethod:     user.ScoringMethod,
	// 		Directorate:       user.Directorate,
	// 	}

	// 	for _, data := range user.CalibrationScores {
	// 		topRemarks := []response.TopRemarkResponse{}
	// 		for _, topRemark := range data.TopRemarks {
	// 			topRemarks = append(topRemarks, response.TopRemarkResponse{
	// 				BaseModel:      topRemark.BaseModel,
	// 				ProjectID:      topRemark.ProjectID,
	// 				EmployeeID:     topRemark.EmployeeID,
	// 				ProjectPhaseID: topRemark.ProjectPhaseID,
	// 				Initiative:     topRemark.Initiative,
	// 				Description:    topRemark.Description,
	// 				Result:         topRemark.Result,
	// 				StartDate:      topRemark.StartDate,
	// 				EndDate:        topRemark.EndDate,
	// 				Comment:        topRemark.Comment,
	// 				EvidenceName:   topRemark.EvidenceName,
	// 				IsProject:      topRemark.IsProject,
	// 				IsInitiative:   topRemark.IsInitiative,
	// 			})
	// 		}
	// 		dataOneResponse.CalibrationScores = append(dataOneResponse.CalibrationScores, response.CalibrationResponse{
	// 			ProjectID: data.ProjectID,
	// 			ProjectPhase: response.ProjectPhaseResponse{
	// 				Phase: response.PhaseResponse{
	// 					Order: data.ProjectPhase.Phase.Order,
	// 				},
	// 				StartDate: data.ProjectPhase.StartDate,
	// 				EndDate:   data.ProjectPhase.EndDate,
	// 			},
	// 			ProjectPhaseID: data.ProjectPhaseID,
	// 			EmployeeID:     data.EmployeeID,
	// 			Calibrator: response.CalibratorResponse{
	// 				Name: data.Calibrator.Name,
	// 			},
	// 			CalibratorID:              data.CalibratorID,
	// 			CalibrationScore:          data.CalibrationScore,
	// 			CalibrationRating:         data.CalibrationRating,
	// 			Status:                    data.Status,
	// 			SpmoStatus:                data.SpmoStatus,
	// 			Comment:                   data.Comment,
	// 			SpmoComment:               data.SpmoComment,
	// 			JustificationType:         data.JustificationType,
	// 			JustificationReviewStatus: data.JustificationReviewStatus,
	// 			SendBackDeadline:          data.SendBackDeadline,
	// 			BottomRemark: response.BottomRemarkResponse{
	// 				ProjectID:      data.BottomRemark.ProjectID,
	// 				EmployeeID:     data.BottomRemark.EmployeeID,
	// 				ProjectPhaseID: data.BottomRemark.ProjectPhaseID,
	// 				LowPerformance: data.BottomRemark.LowPerformance,
	// 				Indisipliner:   data.BottomRemark.Indisipliner,
	// 				Attitude:       data.BottomRemark.Attitude,
	// 				WarningLetter:  data.BottomRemark.WarningLetter,
	// 			},
	// 			TopRemarks: topRemarks,
	// 		})
	// 	}

	// 	calibrations = append(calibrations, *dataOneResponse)
	// }

	return calibrations, nil
}

func (r *calibrationRepo) CheckConditionBeforeSubmitCalibration(projectID string, payload []response.UserResponse,
	projectPhase model.ProjectPhase, countCalibrated response.TotalCalibratedRating, countRatingQuota response.RatingQuota,
) (bool, error) {

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	finalValue := true
	for _, calibrationData := range payload {
		var calibration *model.Calibration
		err := tx.Table("calibrations c").
			Where("c.project_id = ? AND c.employee_id = ? AND c.project_phase_id = ?", projectID, calibrationData.ID, projectPhase.ID).
			First(&calibration).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		if !calibration.FilledTopBottomMark {
			return false, fmt.Errorf("You should fill justification remark on top/bottom remark.")
		}
	}

	if projectPhase.Guideline {
		var project *model.Project
		err := tx.Table("projects").
			Where("id = ?", projectID).Find(&project).Error
		if err != nil {
			tx.Rollback()
			return false, err
		}

		if !project.APlusExcess && countCalibrated.APlus > countRatingQuota.APlus {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
		if !project.AExcess && countCalibrated.A > countRatingQuota.A {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
		if !project.BPlusExcess && countCalibrated.BPlus > countRatingQuota.BPlus {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
		if !project.BExcess && countCalibrated.B > countRatingQuota.B {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
		if !project.CExcess && countCalibrated.C > countRatingQuota.C {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
		if !project.DExcess && countCalibrated.D > countRatingQuota.D {
			return false, fmt.Errorf("Calibration Score Exceeds Rating Quota.")
		}
	}

	tx.Commit()
	return finalValue, nil
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
