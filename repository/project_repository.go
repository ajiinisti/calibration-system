package repository

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type ProjectRepo interface {
	BaseRepository[model.Project]
	PaginateList(pagination model.PaginationQuery) ([]model.Project, response.Paging, error)
	GetTotalRows(name string) (int, error)
	ActivateByID(id string) error
	NonactivateByID(id string) error
	DeactivateAllExceptID(id string) error
	GetProjectPhaseOrder(calibratorID, projectID string) (int, error)
	GetProjectPhase(calibratorID, projectID string) (*model.ProjectPhase, error)
	GetActiveProject() ([]model.Project, error)
	GetActiveProjectPhase(projectID string) ([]model.ProjectPhase, error)
	GetActiveManagerPhase() (model.ProjectPhase, error)
	GetScoreDistributionByCalibratorID(businessUnitID, projectID string) (*model.Project, error)
	GetRatingQuotaByCalibratorID(businessUnitID, projectID string) (*model.Project, error)
	GetNumberOneUserWhoCalibrator(calibratorID, businessUnit, projectID string, calibratorPhase int) ([]string, error)
	GetAllUserCalibrationByCalibratorID(calibratorID, projectID string, calibratorPhase int) ([]model.User, error)
	GetCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string, phase int) (response.UserCalibration, error)
	GetCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string, phase int) (response.UserCalibration, error)
	GetCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit string, phase int, exceptUsers []string) (response.UserCalibration, error)
	GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int, calibratorID, projectID string) (response.UserCalibration, error)
	GetNMinusOneCalibrationsByBusinessUnitPaginate(businessUnit string, phase int, calibratorID, projectID string, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetCalibrationsByRating(calibratorID, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error)
	GetAllBusinessUnitSummary(calibratorID, projectID string, phase int) ([]model.BusinessUnit, error)
	GetCalibrationsForSummaryHelper(types, calibratorID, prevCalibrator, businessUnit, projectID string, phase int) (int, error)
	FindIfCalibratorOnPhaseBefore(calibratorID, projectID string, phase int) (bool, error)
	GetAllActiveProjectByCalibratorID(calibratorID string) ([]model.Project, error)
	GetAllActiveProjectBySpmoID(spmoID string) ([]model.Project, error)
	GetAllEmployeeName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	GetAllSupervisorName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	GetAllGrade(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error)
	GetTotalRowsCalibration(calibratorID, prevCalibrator, businessUnitName, types, projectID, rating string, pagination model.PaginationQuery) (int, error)
	GetCalibratedRating(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalCalibratedRating, error)
	GetAverageScore(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (float32, error)
}

type projectRepo struct {
	db *gorm.DB
}

func (r *projectRepo) Save(payload *model.Project) error {
	// err := r.db.Save(&payload)
	// if err.Error != nil {
	// 	return err.Error
	// }

	if payload.ID == "" {
		var project model.Project
		_ = r.db.
			Preload("ScoreDistributions").
			Preload("ScoreDistributions.GroupBusinessUnit").
			Preload("RemarkSettings").
			Order("created_at DESC").
			First(&project).Error
		// if err != nil {
		// 	return err
		// }

		for _, scoreD := range project.ScoreDistributions {
			payload.ScoreDistributions = append(payload.ScoreDistributions, model.ScoreDistribution{
				ProjectID:           payload.ID,
				GroupBusinessUnitID: scoreD.GroupBusinessUnitID,
				APlusUpperLimit:     scoreD.APlusUpperLimit,
				APlusLowerLimit:     scoreD.APlusLowerLimit,
				AUpperLimit:         scoreD.AUpperLimit,
				ALowerLimit:         scoreD.ALowerLimit,
				BPlusUpperLimit:     scoreD.BPlusUpperLimit,
				BPlusLowerLimit:     scoreD.BLowerLimit,
				BUpperLimit:         scoreD.BUpperLimit,
				BLowerLimit:         scoreD.BLowerLimit,
				CUpperLimit:         scoreD.CUpperLimit,
				CLowerLimit:         scoreD.CLowerLimit,
				DUpperLimit:         scoreD.DUpperLimit,
				DLowerLimit:         scoreD.DLowerLimit,
			})
		}

		for _, remarks := range project.RemarkSettings {
			payload.RemarkSettings = append(payload.RemarkSettings, model.RemarkSetting{
				ProjectID:         payload.ID,
				JustificationType: remarks.JustificationType,
				ScoringType:       remarks.ScoringType,
				Level:             remarks.Level,
				From:              remarks.From,
				To:                remarks.To,
			})
		}
	}

	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}

	return nil
}

func (r *projectRepo) Get(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("ProjectPhases", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN phases p ON project_phases.phase_id = p.id").
				Order("p.order ASC")
		}).
		Preload("ProjectPhases.Phase").
		Preload("ScoreDistributions").
		Preload("ScoreDistributions.GroupBusinessUnit").
		Preload("RemarkSettings").
		First(&project, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) List() ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Preload("ActualScores").
		Preload("ProjectPhases").
		Preload("ProjectPhases.Phase").
		Preload("ScoreDistributions").
		Preload("ScoreDistributions.GroupBusinessUnit").
		Preload("RemarkSettings").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepo) Delete(id string) error {
	result := r.db.Delete(&model.Project{
		BaseModel: model.BaseModel{
			ID: id,
		},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Project not found!")
	}
	return nil
}

func (r *projectRepo) PaginateList(pagination model.PaginationQuery) ([]model.Project, response.Paging, error) {
	var projects []model.Project
	var err error

	if pagination.Name == "" {
		err = r.db.
			// Preload("ActualScores").
			// Preload("ProjectPhases").
			// Preload("ProjectPhases.Phase").
			// Preload("ScoreDistributions").
			// Preload("ScoreDistributions.GroupBusinessUnit").
			// Preload("RemarkSettings").
			Limit(pagination.Take).Offset(pagination.Skip).Find(&projects).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	} else {
		err = r.db.
			// Preload("ActualScores").
			// Preload("ProjectPhases").
			// Preload("ProjectPhases.Phase").
			// Preload("ScoreDistributions").
			// Preload("ScoreDistributions.GroupBusinessUnit").
			// Preload("RemarkSettings").
			Where("name ILIKE ?", "%"+pagination.Name+"%").
			Limit(pagination.Take).Offset(pagination.Skip).Find(&projects).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	}

	totalRows, err := r.GetTotalRows(pagination.Name)
	if err != nil {
		return nil, response.Paging{}, err
	}

	return projects, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetTotalRows(name string) (int, error) {
	var count int64
	var err error
	if name == "" {
		err = r.db.Model(&model.Project{}).
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	} else {
		err = r.db.Model(&model.Project{}).
			Where("name ILIKE ?", "%"+name+"%").
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	}
	return int(count), nil
}

func (r *projectRepo) ActivateByID(id string) error {
	result := r.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]interface{}{"active": true})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *projectRepo) NonactivateByID(id string) error {
	result := r.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]interface{}{"active": false})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *projectRepo) DeactivateAllExceptID(id string) error {
	// Update all rows where 'id' is not equal to the specified 'id'
	result := r.db.Model(&model.Project{}).Where("id <> ?", id).Update("active", false)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *projectRepo) GetActiveProject() ([]model.Project, error) {
	var project []model.Project
	err := r.db.
		Preload("RemarkSettings").
		Where("active = ?", true).
		Find(&project).
		Error
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *projectRepo) GetActiveProjectPhase(projectID string) ([]model.ProjectPhase, error) {
	var project model.Project
	err := r.db.
		Preload("ProjectPhases", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN phases p ON project_phases.phase_id = p.id").
				Order("p.order ASC")
		}).
		Preload("ProjectPhases.Phase").
		First(&project, "id = ?", projectID).
		Error
	if err != nil {
		return nil, err
	}
	return project.ProjectPhases, nil
}

func (r *projectRepo) GetActiveManagerPhase() (model.ProjectPhase, error) {
	var project model.Project
	err := r.db.
		Preload("ProjectPhases", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN phases p ON project_phases.phase_id = p.id").
				Order("p.order ASC")
		}).
		Preload("ProjectPhases.Phase").
		First(&project, "active = ?", true).
		Error
	if err != nil {
		return model.ProjectPhase{}, err
	}
	return project.ProjectPhases[0], nil
}

func (r *projectRepo) GetProjectPhase(calibratorID, projectID string) (*model.ProjectPhase, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		// Joins("JOIN projects ON projects.id = calibrations.project_id").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	return &calibration.ProjectPhase, nil
}

func (r *projectRepo) GetProjectPhaseOrder(calibratorID, projectID string) (int, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return -1, err
	}

	return calibration.ProjectPhase.Phase.Order, nil
}

func (r *projectRepo) GetScoreDistributionByCalibratorID(businessUnitID, projectID string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("ScoreDistributions", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN group_business_units AS gbu ON gbu.id = score_distributions.group_business_unit_id").
				Joins("JOIN business_units as bu ON bu.group_business_unit_id = gbu.id AND bu.id = ?", businessUnitID)
		}).
		Where("projects.id = ?", projectID).
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) GetRatingQuotaByCalibratorID(businessUnitID, projectID string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("RatingQuotas", func(db *gorm.DB) *gorm.DB {
			return db.Where("business_unit_id = ?", businessUnitID)
		}).
		First(&project, "projects.id = ?", projectID).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) GetAllUserCalibrationByCalibratorID(calibratorID, projectID string, calibratorPhase int) ([]model.User, error) {
	var calibration []model.User
	err := r.db.
		Table("materialized_user_view").
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, calibratorPhase).
				Order("p.order ASC")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("*").
		Where("calibrator_id = ? AND project_id = ? AND phase_order = ?", calibratorID, projectID, calibratorPhase).
		Find(&calibration).Error
	if err != nil {
		return nil, err
	}

	return calibration, nil
}

func (r *projectRepo) GetNumberOneUserWhoCalibrator(calibratorID, businessUnit, projectID string, calibratorPhase int) ([]string, error) {
	var users []model.User
	err := r.db.
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.id = ? AND p.order <= ?", projectID, calibratorPhase).
				Order("p.order DESC")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Table("users u").
		Select("u.*").
		Distinct().
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("p.order = ? AND p2.order < ? AND b.id = ? AND c1.calibrator_id = ?",
			calibratorPhase, calibratorPhase, businessUnit, calibratorID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	filteredUsers := make([]string, 0)
	filteredUsersName := make([]string, 0)
	// fmt.Println("SORTING :=", len(users))
	for _, user := range users {
		// fmt.Println("SORTING :=", user.Name)
		for _, calibration := range user.CalibrationScores {
			if calibration.ProjectPhase.Phase.Order < calibratorPhase {
				// fmt.Println("DATA := ", user.Name, calibration.ProjectPhase.Phase.Order, calibration.Calibrator.Name)
				filteredUsers = append(filteredUsers, calibration.Calibrator.ID)
				filteredUsersName = append(filteredUsersName, calibration.Calibrator.Name+strconv.Itoa(calibration.ProjectPhase.Phase.Order))
				break
			}
			filteredUsers = append(filteredUsers, user.ID)
			filteredUsersName = append(filteredUsersName, user.Name)
		}
	}

	// fmt.Println("TEST := ", filteredUsersName)
	return filteredUsers, nil
}

func (r *projectRepo) GetCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string, phase int) (response.UserCalibration, error) {
	// var users []model.UserCalibration
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("materialized_user_view mv1").
		Select("mv1.id").
		Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnit, projectID).
		Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnit, phase, projectID)

	var subqueryResults []string
	if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
		return response.UserCalibration{}, err
	}

	subQueryCount := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	// First get the base users
	err := r.db.Table("materialized_user_view m").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Select("m.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = m.id", subQueryCount).
		Where("m.phase_order = ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
		Order("calibration_count ASC").
		Find(&resultUsers).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	for _, user := range resultUsers {
		if len(user.CalibrationScores) > 1 {
			if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.Phase.Order == 1 {
				NPlusOneManagerFlag = NPlusOneManagerFlag || true

				if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.EndDate.After(time.Now()) && user.CalibrationScores[len(user.CalibrationScores)-2].Status == "Waiting" && user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
					SendToManagerFlag = SendToManagerFlag || false
				} else {
					SendToManagerFlag = SendToManagerFlag || true
				}
			}

			if user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
				SendBackFlag = SendBackFlag || true
			}
		}
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetCalibrationsByPrevCalibratorBusinessUnitPaginate(calibratorID, prevCalibrator, businessUnit, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	var users []model.UserCalibration

	// subquery := r.db.
	// 	Table("materialized_user_view m1").
	// 	Select("m1.id").
	// 	Joins("JOIN materialized_user_view m2 on m1.employee_id = m2.id").
	// 	Where("m2.phase_order < ? AND m2.calibrator_id = ? AND m1.business_unit_id = ? AND m1.calibrator_id = ? AND m1.project_id = ?",
	// 		phase, prevCalibrator, businessUnit, calibratorID, projectID).
	// 	Or("m1.id = ? AND m1.business_unit_id = ? AND m1.phase_order = ? AND m2.phase_order = ? AND m2.project_id = ?",
	// 		prevCalibrator, businessUnit, phase, phase, projectID)

	subquery := r.db.
		Table("materialized_user_view mv1").
		Select("mv1.id").
		Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnit, projectID).
		Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnit, phase, projectID)

	var subqueryResults []string
	if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	subQueryCount := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	order := getOrder(pagination)
	fmt.Println("DATA SUPERVISOR NAME================================================", pagination.EmployeeName)
	err := r.db.Table("materialized_user_view m").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.ProjectPhase.Phase").
		Select("m.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = m.id", subQueryCount).
		Where("m.phase_order = ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(pagination.SupervisorName) > 0 {
				db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
			}
			if len(pagination.EmployeeName) > 0 {
				db = db.Where("m.name IN ?", pagination.EmployeeName)
			}
			if len(pagination.Grade) > 0 {
				db = db.Where("m.grade IN ?", pagination.Grade)
			}
			if pagination.FilterCalibrationRating != "" {
				db = db.Where("m.calibration_rating = ?", pagination.FilterCalibrationRating)
			}
			return db
		}).
		Order(order).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&users).Error
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	for _, user := range users {
		if len(user.CalibrationScores) > 1 {
			if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.Phase.Order == 1 {
				NPlusOneManagerFlag = NPlusOneManagerFlag || true

				if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.EndDate.After(time.Now()) && user.CalibrationScores[len(user.CalibrationScores)-2].Status == "Waiting" && user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
					SendToManagerFlag = SendToManagerFlag || false
				} else {
					SendToManagerFlag = SendToManagerFlag || true
				}
			}

			if user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
				SendBackFlag = SendBackFlag || true
			}
		}
	}

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, prevCalibrator, businessUnit, "default", projectID, "", pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            users,
	}, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string, phase int) (response.UserCalibration, error) {
	var users []response.UserResponse

	err := r.db.
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Table("materialized_user_view m").
		Joins("JOIN users u2 on u2.id = m.id").
		Select("u2.*, COUNT(m.id) AS calibration_count").
		Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnit).
		Group("u2.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false
	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            users,
	}, nil
}

func (r *projectRepo) GetCalibrationsByBusinessUnitPaginate(calibratorID, businessUnit, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	var users []model.UserCalibration
	// var resultUsers []response.UserResponse

	subQuery := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	order := getOrder(pagination)
	err := r.db.
		Table("materialized_user_view m").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.ProjectPhase.Phase").
		Select("m.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = m.id", subQuery).
		Where("(m.phase_order = ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnit).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(pagination.SupervisorName) > 0 {
				db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
			}
			if len(pagination.EmployeeName) > 0 {
				db = db.Where("m.name IN ?", pagination.EmployeeName)
			}
			if len(pagination.Grade) > 0 {
				db = db.Where("m.grade IN ?", pagination.Grade)
			}
			if pagination.FilterCalibrationRating != "" {
				db = db.Where("m.calibration_rating = ?", pagination.FilterCalibrationRating)
			}
			return db
		}).
		Order(order).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&users).Error

	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, "", businessUnit, "all", projectID, "", pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            users,
	}, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit string, phase int, exceptUsers []string) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("users u").
		Select("u.id").
		Distinct().
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("p2.order < ? AND u3.id = ? AND b.id = ? AND c1.calibrator_id = ?",
			phase, prevCalibrator, businessUnit, calibratorID).
		Or("u.id = ? AND b.id = ? AND p.order = ? AND p2.order = ?",
			prevCalibrator, businessUnit, phase, phase)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return response.UserCalibration{}, err
	}

	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj1 ON actual_scores.project_id = proj1.id").
				Where("proj1.active = ?", true)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.active = ? AND p.order <= ?", true, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Joins("INNER JOIN business_units b ON u.business_unit_id = b.id").
		Joins("INNER JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("INNER JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("INNER JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("INNER JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("INNER JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("INNER JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("(p.order <= ? AND u.id IN (?))", phase, subqueryResults).
		Or("(p.order = ? AND c1.calibrator_id = ? AND b.id = ? AND u.id NOT IN (?))", phase, calibratorID, businessUnit, exceptUsers).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	// for _, user := range users {
	// 	var supervisorName string
	// 	err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
	// 	if err != nil {
	// 		return response.UserCalibration{}, err
	// 	}

	// 	if len(user.CalibrationScores) > 1 {
	// 		if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.Phase.Order == 1 {
	// 			NPlusOneManagerFlag = NPlusOneManagerFlag || true

	// 			if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.EndDate.After(time.Now()) && user.CalibrationScores[len(user.CalibrationScores)-2].Status == "Waiting" && user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
	// 				SendToManagerFlag = SendToManagerFlag || false
	// 			} else {
	// 				SendToManagerFlag = SendToManagerFlag || true
	// 			}
	// 		}

	// 		if user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
	// 			SendBackFlag = SendBackFlag || true
	// 		}
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

	// 	resultUsers = append(resultUsers, *dataOneResponse)
	// }
	// if err != nil {
	// 	return response.UserCalibration{}, err
	// }

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int, calibratorID, projectID string) (response.UserCalibration, error) {
	var users []response.UserResponse

	// prev calibrator
	queryPrevCalibrator := r.db.
		Table("users u2").
		Select("u2.id").
		Distinct().
		Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
		Joins("JOIN users u3 on c2.employee_id = u3.id").
		Where("u3.business_unit_id = ?", businessUnit)

	var queryPrevCalibratorResults []string
	if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
		return response.UserCalibration{}, err
	}

	// Subquery
	subquery := r.db.
		Table("materialized_user_view m1").
		Select("m1.id").
		Distinct().
		Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnit)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return response.UserCalibration{}, err
	}

	if len(queryPrevCalibratorResults) == 0 {
		queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
	}
	if len(subqueryResults) == 0 {
		subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
	}

	err := r.db.
		Table("materialized_user_view m1").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Select("m1.*").
		Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
			calibratorID, projectID, phase, businessUnit, queryPrevCalibrator, subqueryResults).
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: false,
		UserData:            users,
	}, nil
}

func (r *projectRepo) GetNMinusOneCalibrationsByBusinessUnitPaginate(businessUnit string, phase int, calibratorID, projectID string, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	// var users []model.User
	var resultUsers []model.UserCalibration

	// prev calibrator
	queryPrevCalibrator := r.db.
		Table("users u2").
		Select("u2.id").
		Distinct().
		Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
		Joins("JOIN users u3 on c2.employee_id = u3.id").
		Where("u3.business_unit_id = ?", businessUnit)

	var queryPrevCalibratorResults []string
	if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}
	// Subquery
	subquery := r.db.
		Table("materialized_user_view m1").
		Select("m1.id").
		Distinct().
		Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnit)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	if len(queryPrevCalibratorResults) == 0 {
		queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
	}
	if len(subqueryResults) == 0 {
		subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
	}

	order := getOrder(pagination)
	err := r.db.
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Table("materialized_user_view m1").
		Select("m1.*").
		Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
			calibratorID, projectID, phase, businessUnit, queryPrevCalibrator, subqueryResults).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(pagination.SupervisorName) > 0 {
				db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
			}
			if len(pagination.EmployeeName) > 0 {
				db = db.Where("m.name IN ?", pagination.EmployeeName)
			}
			if len(pagination.Grade) > 0 {
				db = db.Where("m.grade IN ?", pagination.Grade)
			}
			if pagination.FilterCalibrationRating != "" {
				db = db.Where("m.calibration_rating = ?", pagination.FilterCalibrationRating)
			}
			return db
		}).
		Order(order).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&resultUsers).Error
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, "", businessUnit, "n-1", projectID, "", pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     false,
		SendToManager:       false,
		SendBackCalibration: false,
		UserData:            resultUsers,
	}, utils.Paginate(pagination.Take, pagination.Skip, totalRows), nil
}

func (r *projectRepo) GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	var resultUsers []model.UserCalibration

	subquery := r.db.
		Table("materialized_user_view mv1").
		Select("mv1.id").
		Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnit, projectID).
		Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnit, phase, projectID)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	countQuery := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	err := r.db.
		Table("materialized_user_view m").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.ProjectPhase.Phase").
		Select("m.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = m.id", countQuery).
		Where("m.phase_order = ? AND m.id IN (?) AND m.project_id = ? AND m.calibration_rating = ?", phase, subqueryResults, projectID, rating).
		Order(`calibration_count ASC, 
			CASE calibration_rating 
				WHEN 'A+' THEN 1 
				WHEN 'A' THEN 2 
				WHEN 'B+' THEN 3 
				WHEN 'B' THEN 4 
				WHEN 'C' THEN 5 
				WHEN 'D' THEN 6 
				ELSE 7 
			END ASC, 
			m.calibration_score DESC,
			m.grade DESC,
			m.name ASC
		`).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&resultUsers).Error
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, prevCalibrator, businessUnit, "rating-prev", projectID, rating, pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	var users []model.UserCalibration
	subQuery := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	err := r.db.
		Table("materialized_user_view mv").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.ProjectPhase.Phase").
		Select("mv.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = mv.id", subQuery).
		Where("mv.phase_order = ? AND mv.calibrator_id = ? AND mv.business_unit_id = ? AND mv.calibration_rating = ? AND mv.project_id = ?",
			phase, calibratorID, businessUnit, rating, projectID).
		Order(`calibration_count ASC, 
			CASE mv.calibration_rating 
				WHEN 'A+' THEN 1 
				WHEN 'A' THEN 2 
				WHEN 'B+' THEN 3 
				WHEN 'B' THEN 4 
				WHEN 'C' THEN 5 
				WHEN 'D' THEN 6 
				ELSE 7 
			END ASC, 
			mv.calibration_score DESC,
			mv.grade DESC,
			mv.name ASC
		`).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&users).Error
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, "", businessUnit, "rating-bu", projectID, rating, pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     false,
		SendToManager:       false,
		SendBackCalibration: false,
		UserData:            users,
	}, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetCalibrationsByRating(calibratorID, rating, projectID string, phase int, pagination model.PaginationQuery) (response.UserCalibrationNew, response.Paging, error) {
	var users []model.UserCalibration

	subQuery := r.db.Table("materialized_user_view mv2").
		Select("mv2.employee_id, COUNT(mv2.id) as calibration_count").
		Where("mv2.project_id = ? AND mv2.phase_order <= ?", projectID, phase).
		Group("mv2.employee_id")

	err := r.db.
		Table("materialized_user_view mv").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("calibrations.project_id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.ProjectPhase.Phase").
		Select("mv.*, sq.calibration_count").
		Joins("LEFT JOIN (?) as sq ON sq.employee_id = mv.id", subQuery).
		Where("mv.project_id = ? AND mv.phase_order = ? AND mv.calibrator_id = ? AND mv.calibration_rating = ?", projectID, phase, calibratorID, rating).
		Order(`calibration_count ASC, 
			CASE mv.calibration_rating 
				WHEN 'A+' THEN 1 
				WHEN 'A' THEN 2 
				WHEN 'B+' THEN 3 
				WHEN 'B' THEN 4 
				WHEN 'C' THEN 5 
				WHEN 'D' THEN 6 
				ELSE 7 
			END ASC, 
			mv.calibration_score DESC,
			mv.grade DESC,
			mv.name ASC
		`).
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&users).Error
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	totalRows, err := r.GetTotalRowsCalibration(calibratorID, "", "", "rating-all", projectID, rating, pagination)
	if err != nil {
		return response.UserCalibrationNew{}, response.Paging{}, err
	}

	return response.UserCalibrationNew{
		NPlusOneManager:     false,
		SendToManager:       false,
		SendBackCalibration: false,
		UserData:            users,
	}, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetAllBusinessUnitSummary(calibratorID, projectID string, phase int) ([]model.BusinessUnit, error) {
	var results []model.BusinessUnit
	err := r.db.
		Table("users u").
		Select("b.*").
		Distinct().
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Where("c1.calibrator_id = ? AND p.order = ?", calibratorID, phase).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *projectRepo) GetCalibrationsForSummaryHelper(types, calibratorID, prevCalibrator, businessUnit, projectID string, phase int) (int, error) {
	var count int64

	if types == "numberOne" {
		return -1, nil
	}
	if types == "n-1" {
		// Query for previous calibrators
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnit)

		// Subquery for previous phases
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?",
				projectID, phase, businessUnit)

		// Main count query
		err := r.db.
			Table("users u").
			Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.calibrator_id = ? AND c1.project_id = ?", calibratorID, projectID).
			Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
			Joins("INNER JOIN phases p ON p.id = pp.phase_id AND p.order = ?", phase).
			Where("u.business_unit_id = ?", businessUnit).
			Where("u.id NOT IN (?)", subquery).
			Where("u.id NOT IN (?)", queryPrevCalibrator).
			Count(&count).Error
		if err != nil {
			return -1, err
		}

	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?",
				phase, prevCalibrator, businessUnit, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?",
				prevCalibrator, businessUnit, phase, projectID)

		err := r.db.Table("materialized_user_view m").
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ?", phase, subquery, projectID).
			Count(&count).Error

		if err != nil {
			return -1, err
		}

	} else {
		err := r.db.Table("materialized_user_view m").
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?",
				phase, calibratorID, projectID, businessUnit).
			Count(&count).Error

		if err != nil {
			return -1, err
		}
	}

	return int(count), nil
}

func (r *projectRepo) FindIfCalibratorOnPhaseBefore(calibratorID, projectID string, phase int) (bool, error) {
	var cal *model.Calibration

	err := r.db.Table("calibrations c").
		Select("c.*").
		Joins("INNER JOIN project_phases pp ON pp.id = c.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Where("p.order < ? AND c.calibrator_id = ? and c.project_id = ?", phase, calibratorID, projectID).
		First(&cal).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		// Other errors should be returned
		return false, err
	}

	// If we found the record, return true
	return true, nil
}

func (r *projectRepo) GetAllActiveProjectByCalibratorID(calibratorID string) ([]model.Project, error) {
	var projectList []model.Project
	err := r.db.Table("projects p").
		Select("p.*").
		Distinct().
		Joins("INNER JOIN calibrations c ON p.id = c.project_id AND p.active = ? AND c.calibrator_id = ?", true, calibratorID).
		Find(&projectList).Error
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (r *projectRepo) GetAllActiveProjectBySpmoID(spmoID string) ([]model.Project, error) {
	var projectList []model.Project
	err := r.db.Table("projects p").
		Select("p.*").
		Distinct().
		Joins("INNER JOIN calibrations c ON p.id = c.project_id AND p.active = ? AND (spmo_id = ? OR spmo2_id = ? OR spmo3_id = ?)", true, spmoID, spmoID, spmoID).
		Find(&projectList).Error
	if err != nil {
		return nil, err
	}

	return projectList, nil
}

func (r *projectRepo) GetAllEmployeeName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	var employeeName []string
	var err error

	var calibration model.Calibration
	err = r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	phase := calibration.ProjectPhase.Phase.Order

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return nil, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("m1.name").
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Order("m1.name ASC").
			Find(&employeeName).Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("m.name").
			Distinct().
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
			Order("m.name ASC").
			Find(&employeeName).Error
	} else {
		err = r.db.
			Table("materialized_user_view m").
			Select("m.name").
			Distinct().
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnitName).
			Order("m.name ASC").
			Find(&employeeName).Error
	}

	if err != nil {
		return nil, err
	}
	return employeeName, nil
}

func (r *projectRepo) GetAllSupervisorName(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	var employeeName []string
	var err error

	var calibration model.Calibration
	err = r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	phase := calibration.ProjectPhase.Phase.Order

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return nil, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("m1.supervisor_names").
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Order("m1.supervisor_names ASC").
			Find(&employeeName).Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("m.supervisor_names").
			Distinct().
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ? AND m.supervisor_names IS NOT NULL", phase, subqueryResults, projectID).
			Order("m.supervisor_names ASC").
			Find(&employeeName).Error
	} else {
		err = r.db.
			Table("materialized_user_view m").
			Select("m.supervisor_names").
			Distinct().
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ? AND m.supervisor_names IS NOT NULL", phase, calibratorID, projectID, businessUnitName).
			Order("m.supervisor_names ASC").
			Find(&employeeName).Error
	}

	if err != nil {
		return nil, err
	}
	return employeeName, nil
}

func (r *projectRepo) GetAllGrade(calibratorID, prevCalibrator, businessUnitName, types, projectID string) ([]string, error) {
	var employeeName []string
	var err error

	var calibration model.Calibration
	err = r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	phase := calibration.ProjectPhase.Phase.Order

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return nil, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("m1.grade").
			Distinct().
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Order("m1.grade ASC").
			Find(&employeeName).Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("m.grade").
			Distinct().
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
			Order("m.grade ASC").
			Find(&employeeName).Error
	} else {
		err = r.db.
			Table("materialized_user_view m").
			Select("m.grade").
			Distinct().
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnitName).
			Order("m.grade ASC").
			Find(&employeeName).Error
	}

	if err != nil {
		return nil, err
	}
	return employeeName, nil
}

func (r *projectRepo) GetTotalRowsCalibration(calibratorID, prevCalibrator, businessUnitName, types, projectID, rating string, pagination model.PaginationQuery) (int, error) {
	var count int64
	var err error

	var calibration model.Calibration
	err = r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return -1, err
	}

	phase := calibration.ProjectPhase.Phase.Order

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return -1, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return -1, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("m1.*").
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Scopes(func(db *gorm.DB) *gorm.DB {
				if len(pagination.SupervisorName) > 0 {
					db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
				}
				if len(pagination.EmployeeName) > 0 {
					db = db.Where("m.name IN ?", pagination.EmployeeName)
				}
				if len(pagination.Grade) > 0 {
					db = db.Where("m.grade IN ?", pagination.Grade)
				}
				return db
			}).
			Count(&count).Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return -1, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("m.id").
			Distinct().
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
			Scopes(func(db *gorm.DB) *gorm.DB {
				if len(pagination.SupervisorName) > 0 {
					db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
				}
				if len(pagination.EmployeeName) > 0 {
					db = db.Where("m.name IN ?", pagination.EmployeeName)
				}
				if len(pagination.Grade) > 0 {
					db = db.Where("m.grade IN ?", pagination.Grade)
				}
				return db
			}).
			Count(&count).Error
	} else if types == "all" {
		err = r.db.
			Table("materialized_user_view m").
			Select("m.id").
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnitName).
			Scopes(func(db *gorm.DB) *gorm.DB {
				if len(pagination.SupervisorName) > 0 {
					db = db.Where("m.supervisor_names IN ?", pagination.SupervisorName)
				}
				if len(pagination.EmployeeName) > 0 {
					db = db.Where("m.name IN ?", pagination.EmployeeName)
				}
				if len(pagination.Grade) > 0 {
					db = db.Where("m.grade IN ?", pagination.Grade)
				}
				return db
			}).
			Count(&count).Error
	} else if types == "rating-prev" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return -1, err
		}

		err = r.db.
			Table("materialized_user_view m").
			Select("m.id").
			Where("m.phase_order = ? AND m.id IN (?) AND m.project_id = ? AND m.calibration_rating = ?", phase, subqueryResults, projectID, rating).
			Count(&count).Error
	} else if types == "rating-bu" {
		err = r.db.
			Table("materialized_user_view mv").
			Select("mv.id").
			Where("mv.phase_order = ? AND mv.calibrator_id = ? AND mv.business_unit_id = ? AND mv.calibration_rating = ? AND mv.project_id = ?",
				phase, calibratorID, businessUnitName, rating, projectID).
			Count(&count).Error
	} else {
		err = r.db.
			Table("materialized_user_view mv").
			Select("mv.id").
			Where("mv.project_id = ? AND mv.phase_order = ? AND mv.calibrator_id = ? AND mv.calibration_rating = ?", projectID, phase, calibratorID, rating).
			Count(&count).Error
	}

	if err != nil {
		return -1, err
	}

	return int(count), nil
}

func (r *projectRepo) GetCalibratedRating(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (*response.TotalCalibratedRating, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return nil, err
	}
	phase := calibration.ProjectPhase.Phase.Order

	var groupedResults []struct {
		CalibrationRating string
		Count             int
	}

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return nil, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("m.calibration_rating AS calibration_rating, COUNT(*) as count").
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?)",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Group("m1.calibration_rating").
			Scan(&groupedResults).
			Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return nil, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("m.calibration_rating AS calibration_rating, COUNT(*) as count").
			Distinct().
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ?", phase, subqueryResults, projectID).
			Group("m.calibration_rating").
			Scan(&groupedResults).
			Error
	} else if types == "all" {
		err = r.db.Table("materialized_user_view m").
			Select("m.calibration_rating AS calibration_rating, COUNT(*) as count").
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ?", phase, calibratorID, projectID, businessUnitName).
			Group("m.calibration_rating").
			Scan(&groupedResults).
			Error
	}

	if err != nil {
		return nil, err
	}

	totalCalibratedRating := &response.TotalCalibratedRating{
		APlus: 0,
		A:     0,
		BPlus: 0,
		B:     0,
		C:     0,
		D:     0,
		Total: 0,
	}

	for _, result := range groupedResults {
		switch result.CalibrationRating {
		case "A+":
			totalCalibratedRating.APlus += result.Count
		case "A":
			totalCalibratedRating.A += result.Count
		case "B+":
			totalCalibratedRating.BPlus += result.Count
		case "B":
			totalCalibratedRating.B += result.Count
		case "C":
			totalCalibratedRating.C += result.Count
		case "D":
			totalCalibratedRating.D += result.Count
		}
		totalCalibratedRating.Total += result.Count
	}
	return totalCalibratedRating, nil
}

func (r *projectRepo) GetAverageScore(calibratorID, prevCalibrator, businessUnitName, types, projectID string) (float32, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Where("project_id = ? AND calibrator_id = ? ", projectID, calibratorID).
		First(&calibration).Error
	if err != nil {
		return 0, err
	}
	phase := calibration.ProjectPhase.Phase.Order

	var averageScore *float32

	if types == "n-1" {
		// prev calibrator
		queryPrevCalibrator := r.db.
			Table("users u2").
			Select("u2.id").
			Distinct().
			Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL AND c2.project_id = ?", projectID).
			Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
			Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
			Joins("JOIN users u3 on c2.employee_id = u3.id").
			Where("u3.business_unit_id = ?", businessUnitName)

		var queryPrevCalibratorResults []string
		if err := queryPrevCalibrator.Pluck("u.id", &queryPrevCalibratorResults).Error; err != nil {
			return 0, err
		}
		// Subquery
		subquery := r.db.
			Table("materialized_user_view m1").
			Select("m1.id").
			Distinct().
			Where("m1.project_id = ? AND m1.phase_order < ? AND m1.business_unit_id = ?", projectID, phase, businessUnitName)

		var subqueryResults []string
		if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
			return 0, err
		}

		if len(queryPrevCalibratorResults) == 0 {
			queryPrevCalibratorResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}
		if len(subqueryResults) == 0 {
			subqueryResults = []string{"00000000-0000-0000-0000-000000000000"} // Placeholder UUID
		}

		err = r.db.
			Table("materialized_user_view m1").
			Select("AVG(m1.calibration_score)").
			Where("m1.calibrator_id = ? AND m1.project_id = ? and m1.phase_order = ? AND m1.business_unit_id = ? AND m1.id NOT IN (?) AND m1.id NOT IN (?) AND m1.scoring_method='Score'",
				calibratorID, projectID, phase, businessUnitName, queryPrevCalibrator, subqueryResults).
			Scan(&averageScore).
			Error
	} else if types == "default" {
		subquery := r.db.
			Table("materialized_user_view mv1").
			Select("mv1.id").
			Where("mv1.phase_order < ? AND mv1.calibrator_id = ? AND mv1.business_unit_id = ? AND mv1.project_id = ?", phase, prevCalibrator, businessUnitName, projectID).
			Or("mv1.id = ? AND mv1.business_unit_id = ? AND mv1.phase_order = ? AND mv1.project_id = ?", prevCalibrator, businessUnitName, phase, projectID)

		var subqueryResults []string
		if err := subquery.Pluck("id", &subqueryResults).Error; err != nil {
			return 0, err
		}

		// First get the base users
		err = r.db.Table("materialized_user_view m").
			Select("AVG(m.calibration_score)").
			Where("m.phase_order <= ? AND m.id IN (?) AND m.project_id = ? AND m1.scoring_method='Score'", phase, subqueryResults, projectID).
			Scan(&averageScore).
			Error
	} else if types == "all" {
		err = r.db.Table("materialized_user_view m").
			Select("AVG(m.calibration_score)").
			Where("(m.phase_order <= ? AND m.calibrator_id = ?) AND m.project_id = ? AND m.business_unit_id = ? AND m.scoring_method='Score'", phase, calibratorID, projectID, businessUnitName).
			Scan(&averageScore).
			Error
	}

	if err != nil {
		return 0, err
	}
	return *averageScore, nil
}

func getOrder(pagination model.PaginationQuery) string {
	orderBy := ""
	if pagination.OrderEmployeeName != "default" {
		if pagination.OrderEmployeeName == "ascending" {
			orderBy += "m.name ASC"
		} else if pagination.OrderEmployeeName == "descending" {
			orderBy += "m.name DESC"
		}
	}

	if pagination.OrderGrade != "default" {
		if orderBy != "" {
			orderBy += ", "
		}
		if pagination.OrderGrade == "ascending" {
			orderBy += `CASE 
				WHEN m.grade ~ '^[0-9]+$' THEN m.grade::INTEGER 
				ELSE 0 
			END ASC`
		} else if pagination.OrderGrade == "descending" {
			orderBy += `CASE 
				WHEN m.grade ~ '^[0-9]+$' THEN m.grade::INTEGER 
				ELSE 0 
			END DESC`
		}
	}

	if pagination.OrderCalibrationRating != "default" {
		if orderBy != "" {
			orderBy += ", "
		}
		if pagination.OrderCalibrationRating == "ascending" {
			orderBy += `CASE m.calibration_rating 
				WHEN 'A+' THEN 1 
				WHEN 'A' THEN 2 
				WHEN 'B+' THEN 3 
				WHEN 'B' THEN 4 
				WHEN 'C' THEN 5 
				WHEN 'D' THEN 6 
				ELSE 7 
			END DESC`
		} else if pagination.OrderCalibrationRating == "descending" {
			orderBy += `CASE m.calibration_rating 
				WHEN 'A+' THEN 1 
				WHEN 'A' THEN 2 
				WHEN 'B+' THEN 3 
				WHEN 'B' THEN 4 
				WHEN 'C' THEN 5 
				WHEN 'D' THEN 6 
				ELSE 7 
			END ASC`
		}
	}

	if pagination.OrderCalibrationScore != "default" {
		if orderBy != "" {
			orderBy += ", "
		}
		if pagination.OrderCalibrationScore == "ascending" {
			orderBy += `m.calibration_score ASC`
		} else if pagination.OrderCalibrationScore == "descending" {
			orderBy += `m.calibration_score DESC`
		}
	}

	if orderBy == "" {
		orderBy = `
        calibration_count ASC, 
        CASE m.calibration_rating 
            WHEN 'A+' THEN 1 
            WHEN 'A' THEN 2 
            WHEN 'B+' THEN 3 
            WHEN 'B' THEN 4 
            WHEN 'C' THEN 5 
            WHEN 'D' THEN 6 
            ELSE 7 
        END ASC, 
        m.calibration_score DESC, 
		CASE 
			WHEN m.grade ~ '^[0-9]+$' THEN m.grade::INTEGER 
			ELSE 0 
		END DESC,
        m.name ASC
    `
	}

	return orderBy
}
func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}
