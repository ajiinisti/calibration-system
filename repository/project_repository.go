package repository

import (
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
	GetCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit, projectID string, phase int) (response.UserCalibration, error)
	GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorID, prevCalibrator, businessUnit string, phase int, exceptUsers []string) (response.UserCalibration, error)
	GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int, calibratorID, projectID string) (response.UserCalibration, error)
	GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, phase int) (response.UserCalibration, error)
	GetCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID string, phase int) (response.UserCalibration, error)
	GetCalibrationsByRating(calibratorID, rating, projectID string, phase int) (response.UserCalibration, error)
	GetAllBusinessUnitSummary(calibratorID, projectID string, phase int) ([]model.BusinessUnit, error)
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
			Preload("ActualScores").
			Preload("ProjectPhases").
			Preload("ProjectPhases.Phase").
			Preload("ScoreDistributions").
			Preload("ScoreDistributions.GroupBusinessUnit").
			Preload("RemarkSettings").
			Limit(pagination.Take).Offset(pagination.Skip).Find(&projects).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	} else {
		err = r.db.
			Preload("ActualScores").
			Preload("ProjectPhases").
			Preload("ProjectPhases.Phase").
			Preload("ScoreDistributions").
			Preload("ScoreDistributions.GroupBusinessUnit").
			Preload("RemarkSettings").
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
		Table("users u").
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.id = ? AND p.order <= ?", projectID, calibratorPhase).
				Order("p.order ASC")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
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
		Where("p2.order <= ? AND c1.calibrator_id = ?", calibratorPhase, calibratorID).
		Group("u.id").
		Order("calibration_count DESC").
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
	var users []model.User
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("users u").
		Select("u.id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("LEFT JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
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
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Where("p.order <= ? AND u.id IN (?)", phase, subqueryResults).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
		}

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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}
	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetCalibrationsByBusinessUnit(calibratorID, businessUnit, projectID string, phase int) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("INNER JOIN business_units b ON u.business_unit_id = b.id").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.project_id = ?", projectID).
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Joins("INNER JOIN calibrations c2 ON c2.employee_id = u.id AND c2.deleted_at is NULL").
		Joins("INNER JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("INNER JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("INNER JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order <= ?", phase).
		Where("p.order = ? AND c1.calibrator_id = ? AND b.id = ? and c1.project_id = ?", phase, calibratorID, businessUnit, projectID).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}
	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
		}

		fmt.Println(user.Name, user.Nik, "==========================================DATA USER ACTUAL SCORE==========================", user.ActualScores)
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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}
	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
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

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
		}

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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}
	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int, calibratorID, projectID string) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	// prev calibrator
	queryPrevCalibrator := r.db.
		Table("users u2").
		Select("u2.id").
		Joins("JOIN calibrations c2 ON c2.calibrator_id = u2.id AND c2.deleted_at IS NULL").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
		Where("u2.business_unit_id = ?", businessUnit)

	// Subquery
	subquery := r.db.
		Table("users u2").
		Select("u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u2.id AND c2.deleted_at IS NULL").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order < ?", phase).
		Where("u2.business_unit_id = ?", businessUnit)

	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.calibrator_id = ? AND c1.project_id = ?", calibratorID, projectID).
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id AND p.order = ?", phase).
		Where("u.business_unit_id = ? AND u.id NOT IN (?) AND u.id NOT IN (?)", businessUnit, subquery, queryPrevCalibrator).
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
		}
		fmt.Println(user.Name, user.Nik, "==========================================DATA USER ACTUAL SCORE==========================", user.ActualScores)
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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}
	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: false,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetCalibrationsByPrevCalibratorBusinessUnitAndRating(calibratorID, prevCalibrator, businessUnit, rating, projectID string, phase int) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("users u").
		Select("u.id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("LEFT JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
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
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.project_id = ?", projectID).
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Joins("INNER JOIN business_units b ON u.business_unit_id = b.id").
		Where("p.order <= ? AND u.id IN (?) AND c1.calibration_rating = ?", phase, subqueryResults, rating).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	NPlusOneManagerFlag := false
	SendToManagerFlag := false
	SendBackFlag := false

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
		}

		// if len(user.CalibrationScores) > 1 {
		// 	if user.CalibrationScores[len(user.CalibrationScores)-2].ProjectPhase.Phase.Order == 1 {
		// 		NPlusOneManagerFlag = NPlusOneManagerFlag || true

		// 		if user.CalibrationScores[len(user.CalibrationScores)-2].Status != "Waiting" || user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Complete" {
		// 			SendToManagerFlag = SendToManagerFlag || true
		// 		}
		// 	}

		// 	if user.CalibrationScores[len(user.CalibrationScores)-1].Status == "Calibrate" {
		// 		SendBackFlag = SendBackFlag || true
		// 	}
		// }

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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}

	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     NPlusOneManagerFlag,
		SendToManager:       SendToManagerFlag,
		SendBackCalibration: SendBackFlag,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetCalibrationsByBusinessUnitAndRating(calibratorID, businessUnit, rating, projectID string, phase int) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("INNER JOIN business_units b ON u.business_unit_id = b.id").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.project_id = ?", projectID).
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Joins("INNER JOIN calibrations c2 ON c2.employee_id = u.id AND c2.deleted_at is NULL").
		Joins("INNER JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("INNER JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("INNER JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order <= ?", phase).
		Where("p.order = ? AND c1.calibrator_id = ? AND b.id = ? AND c1.calibration_rating = ?", phase, calibratorID, businessUnit, rating).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}

	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     false,
		SendToManager:       false,
		SendBackCalibration: false,
		UserData:            resultUsers,
	}, nil
}

func (r *projectRepo) GetCalibrationsByRating(calibratorID, rating, projectID string, phase int) (response.UserCalibration, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	err := r.db.
		Table("users u").
		Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
			return db.Where("project_id = ?", projectID)
		}).
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("proj2.id = ? AND p.order <= ?", projectID, phase).
				Order("p.order")
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("INNER JOIN business_units b ON u.business_unit_id = b.id").
		Joins("INNER JOIN calibrations c1 ON c1.employee_id = u.id AND c1.deleted_at IS NULL AND c1.project_id = ?", projectID).
		Joins("INNER JOIN projects pr ON pr.id = c1.project_id AND pr.id = ?", projectID).
		Joins("INNER JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("INNER JOIN phases p ON p.id = pp.phase_id").
		Joins("INNER JOIN calibrations c2 ON c2.employee_id = u.id AND c2.deleted_at is NULL").
		Joins("INNER JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.id = ?", projectID).
		Joins("INNER JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("INNER JOIN phases p2 ON p2.id = pp2.phase_id AND p2.order <= ?", phase).
		Where("p.order = ? AND c1.calibrator_id = ? AND c1.calibration_rating = ?", phase, calibratorID, rating).
		Group("u.id").
		Order("calibration_count ASC").
		Find(&users).Error
	if err != nil {
		return response.UserCalibration{}, err
	}

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return response.UserCalibration{}, err
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
					IsProject:      topRemark.IsProject,
					IsInitiative:   topRemark.IsInitiative,
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

		resultUsers = append(resultUsers, *dataOneResponse)
	}
	if err != nil {
		return response.UserCalibration{}, err
	}

	return response.UserCalibration{
		NPlusOneManager:     false,
		SendToManager:       false,
		SendBackCalibration: false,
		UserData:            resultUsers,
	}, nil
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

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}
