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
	Get(id string) (*model.Calibration, error)
	List() ([]model.Calibration, error)
	GetActiveBySPMOID(spmoID string) ([]model.Calibration, error)
	GetAcceptedBySPMOID(spmoID string) ([]model.Calibration, error)
	GetRejectedBySPMOID(spmoID string) ([]model.Calibration, error)
	Delete(projectId, projectPhaseId, employeeId string) error
	Bulksave(payload *[]model.Calibration) error
	BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) error
	SaveChanges(payload *request.CalibrationRequest) error
	AcceptCalibration(payload *request.AcceptJustification, phaseOrder int) error
	AcceptMultipleCalibration(payload *request.AcceptMultipleJustification) error
	RejectCalibration(payload *request.RejectJustification) error
	GetSummaryBySPMOID(spmoID string) ([]response.SPMOSummaryResult, error)
	GetAllDetailCalibrationBySPMOID(spmoID, calibratorID, businessUnitID, department string, order int) ([]response.UserResponse, error)
	GetAllDetailCalibration2BySPMOID(spmoID, calibratorID, businessUnitID string, order int) ([]response.UserResponse, error)
}

type calibrationRepo struct {
	db *gorm.DB
}

func (r *calibrationRepo) Save(payload *model.Calibration) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *calibrationRepo) Bulksave(payload *[]model.Calibration) error {
	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		return r.db.Save(&currentBatch).Error

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := r.db.Save(&remainingItems).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *calibrationRepo) Get(id string) (*model.Calibration, error) {
	var calibration model.Calibration
	err := r.db.Preload("Project").Preload("Employee").Preload("ProjectPhase").Preload("Calibrator").First(&calibration, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &calibration, nil
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

func (r *calibrationRepo) Delete(projectId, projectPhaseId, employeeId string) error {
	result := r.db.Delete(&model.Calibration{
		ProjectID:      projectId,
		ProjectPhaseID: projectPhaseId,
		EmployeeID:     employeeId,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Calibration not found!")
	}
	return nil
}

func (r *calibrationRepo) BulkUpdate(payload *request.CalibrationRequest, projectPhase model.ProjectPhase) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	reviewSPMO := false
	var employeeCalibrationScore []model.Calibration
	for _, calibrations := range payload.RequestData {
		if projectPhase.ReviewSpmo == true {
			calibrations.SpmoStatus = "Waiting"
			reviewSPMO = true
		}
		calibrations.Status = "Complete"
		employeeCalibrationScore = append(employeeCalibrationScore, calibrations)

		result := tx.Updates(calibrations)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		} else if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("Calibrations not found!")
		}
	}
	for _, employeeCalibration := range employeeCalibrationScore {
		var calibrations []model.Calibration
		err := tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", projectPhase.Phase.Order, employeeCalibration.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		for _, c := range calibrations {
			c.CalibrationScore = employeeCalibration.CalibrationScore
			c.CalibrationRating = employeeCalibration.CalibrationRating
			c.Status = "Waiting"
			err := tx.Updates(&c).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		if !reviewSPMO && len(calibrations) > 0 {
			calibrations[0].Status = "Calibrate"
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
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
		err := tx.Updates(model.Calibration{
			ProjectID:      justification.ProjectID,
			ProjectPhaseID: justification.ProjectPhaseID,
			EmployeeID:     justification.EmployeeID,
			SpmoStatus:     "Accepted",
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		var calibration model.Calibration
		err = r.db.
			Preload("ProjectPhase").
			Preload("ProjectPhase.Phase").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Where("projects.active = ? AND calibrations.calibrator_id = ? ", true, justification.CalibratorID).
			First(&calibration).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		var calibrations []model.Calibration
		err = tx.Table("calibrations").
			Select("calibrations.*").
			Joins("JOIN projects ON projects.id = calibrations.project_id").
			Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
			Joins("JOIN phases ON phases.id = project_phases.phase_id").
			Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", calibration.ProjectPhase.Phase.Order, justification.EmployeeID).
			Order("phases.order ASC").
			Find(&calibrations).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		if len(calibrations) > 0 {
			calibrations[0].Status = "Calibrate"
			if err := tx.Updates(calibrations[0]).Error; err != nil {
				tx.Rollback()
				return err
			}
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

	err := tx.Updates(model.Calibration{
		ProjectID:      payload.ProjectID,
		ProjectPhaseID: payload.ProjectPhaseID,
		EmployeeID:     payload.EmployeeID,
		SpmoStatus:     "Accepted",
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var calibrations []model.Calibration
	err = tx.Table("calibrations").
		Select("calibrations.*").
		Joins("JOIN projects ON projects.id = calibrations.project_id").
		Joins("JOIN project_phases ON project_phases.id = calibrations.project_phase_id").
		Joins("JOIN phases ON phases.id = project_phases.phase_id").
		Where("projects.active = true AND phases.order > ? AND calibrations.employee_id = ?", phaseOrder, payload.EmployeeID).
		Order("phases.order ASC").
		Find(&calibrations).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	if len(calibrations) > 0 {
		calibrations[0].Status = "Calibrate"
		if err := tx.Updates(calibrations[0]).Error; err != nil {
			tx.Rollback()
			return err
		}
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
		Where("project_id = ? AND project_phase_id = ? AND employee_id = ?",
			payload.ProjectID,
			payload.ProjectPhaseID,
			payload.EmployeeID,
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
		Where("spmo_id = ? AND p.order NOT IN (SELECT MAX(\"order\") FROM phases)", spmoID).
		Group("u.business_unit_id, b.name, u2.name, c.calibrator_id, c.project_phase_id, p.order").
		Order("p.order ASC").
		Scan(&results).Error

	// err := tx.Table("calibrations c").
	// 	Select("COUNT(c.*) as count, u.business_unit_id, b.name as business_unit_name ,u.department, u2.name as calibrator_name, c.calibrator_id, c.project_phase_id, p.order").
	// 	Joins("JOIN project_phases pp on pp.id = c.project_phase_id").
	// 	Joins("JOIN phases p on pp.phase_id = p.id").
	// 	Joins("JOIN users u on c.employee_id = u.id").
	// 	Joins("JOIN business_units b on u.business_unit_id = b.id").
	// 	Joins("JOIN users u2 on c.calibrator_id = u2.id").
	// 	Where("spmo_id = ? AND p.order NOT IN (SELECT MAX(\"order\") FROM phases)", "23105b4e-e33c-4841-8e49-e23b8a745257").
	// 	Group("u.business_unit_id, b.name, u.department, u2.name, c.calibrator_id, c.project_phase_id, p.order").
	// 	Order("p.order ASC").
	// 	Scan(&results).Error

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
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.spmo_id = ? AND c1.calibrator_id = ?", spmoID, calibratorID).
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
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.spmo_id = ? AND c1.calibrator_id = ?", spmoID, calibratorID).
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("LEFT JOIN users u2 ON u.supervisor_nik = u2.nik").
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}

	return calibration, nil
}

func NewCalibrationRepo(db *gorm.DB) CalibrationRepo {
	return &calibrationRepo{
		db: db,
	}
}
