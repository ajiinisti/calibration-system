package repository

import (
	"fmt"
	"strconv"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type ProjectRepo interface {
	BaseRepository[model.Project]
	PaginateList(pagination model.PaginationQuery) ([]model.Project, response.Paging, error)
	GetTotalRows() (int, error)
	ActivateByID(id string) error
	DeactivateAllExceptID(id string) error
	GetProjectPhaseOrder(calibratorID string) (int, error)
	GetProjectPhase(calibratorID string) (*model.ProjectPhase, error)
	GetActiveProject() (*model.Project, error)
	GetScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error)
	GetRatingQuotaByCalibratorID(businessUnitName string) (*model.Project, error)
	GetNumberOneUserWhoCalibrator(calibratorID string, businessUnit string, calibratorPhase int) ([]string, error)
	GetAllCalibrationByCalibratorID(calibratorID string, calibratorPhase int) ([]model.User, error)
	GetCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string, phase int) ([]response.UserResponse, error)
	GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string, phase int, exceptPrevCalibrator []string) ([]response.UserResponse, error)
	GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int) ([]response.UserResponse, error)
}

type projectRepo struct {
	db *gorm.DB
}

func (r *projectRepo) Save(payload *model.Project) error {
	err := r.db.Save(&payload)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *projectRepo) Get(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("ActualScores").
		Preload("ProjectPhases").
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
	err := r.db.
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

	totalRows, err := r.GetTotalRows()
	if err != nil {
		return nil, response.Paging{}, err
	}

	return projects, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (r *projectRepo) GetTotalRows() (int, error) {
	var count int64
	err := r.db.Model(&model.Project{}).Count(&count).Error
	if err != nil {
		return 0, err
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

func (r *projectRepo) DeactivateAllExceptID(id string) error {
	// Update all rows where 'id' is not equal to the specified 'id'
	result := r.db.Model(&model.Project{}).Where("id <> ?", id).Update("active", false)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *projectRepo) GetActiveProject() (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("RemarkSettings").
		First(&project, "active = ?", true).
		Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) GetProjectPhase(calibratorID string) (*model.ProjectPhase, error) {
	var calibration model.Calibration
	err := r.db.
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		Joins("JOIN projects ON projects.id = calibrations.project_id").
		Where("projects.active = ? AND calibrations.calibrator_id = ? ", true, calibratorID).
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		First(&calibration).Error
	if err != nil {
		return nil, err
	}

	return &calibration.ProjectPhase, nil
}

func (r *projectRepo) GetProjectPhaseOrder(calibratorID string) (int, error) {
	var calibration model.Calibration
	err := r.db.
		Joins("JOIN projects ON projects.id = calibrations.project_id").
		Where("projects.active = ? AND calibrations.calibrator_id = ? ", true, calibratorID).
		Preload("ProjectPhase").
		Preload("ProjectPhase.Phase").
		First(&calibration).Error
	if err != nil {
		return -1, err
	}

	return calibration.ProjectPhase.Phase.Order, nil
}

func (r *projectRepo) GetScoreDistributionByCalibratorID(businessUnitName string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("ScoreDistributions", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN group_business_units AS gbu ON gbu.id = score_distributions.group_business_unit_id").
				Joins("JOIN business_units as bu ON bu.group_business_unit_id = gbu.id").
				Where("bu.name = ?", businessUnitName)
		}).
		First(&project, "projects.active = ?", true).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) GetRatingQuotaByCalibratorID(businessUnitName string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("RatingQuotas", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN business_units AS bu ON rating_quota.business_unit_id = bu.id").
				Where("bu.name= ?", businessUnitName)
		}).
		First(&project, "projects.active = ?", true).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *projectRepo) GetAllCalibrationByCalibratorID(calibratorID string, calibratorPhase int) ([]model.User, error) {
	var calibration []model.User
	err := r.db.
		Table("users u").
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.active = ? AND p.order <= ?", true, calibratorPhase)
		}).
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*, COUNT(u.id) AS calibration_count").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
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

func (r *projectRepo) GetNumberOneUserWhoCalibrator(calibratorID string, businessUnit string, calibratorPhase int) ([]string, error) {
	var users []model.User
	err := r.db.
		Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN projects ON calibrations.project_id = projects.id").
				Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
				Joins("JOIN phases p ON p.id = pp.phase_id ").
				Where("projects.active = ? AND p.order <= ?", true, calibratorPhase).
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
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("p.order = ? AND p2.order < ? AND b.name = ? AND c1.calibrator_id = ?",
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

func (r *projectRepo) GetCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string, phase int) ([]response.UserResponse, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("users u").
		Select("u.id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("p2.order < ? AND u3.name = ? AND b.name = ? AND c1.calibrator_id = ?",
			phase, prevCalibrator, businessUnit, calibratorId).
		Or("u.name = ? AND b.name = ? AND p.order = ? AND p2.order = ?",
			prevCalibrator, businessUnit, phase, phase)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return nil, err
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
		Limit(10).Offset(0).
		Find(&users).Error

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return nil, err
		}

		resultUsers = append(resultUsers, response.UserResponse{
			BaseModel: model.BaseModel{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			CreatedBy:              user.CreatedBy,
			UpdatedBy:              user.UpdatedBy,
			Email:                  user.Email,
			Name:                   user.Name,
			Nik:                    user.Nik,
			DateOfBirth:            user.DateOfBirth,
			SupervisorNames:        supervisorName,
			BusinessUnit:           user.BusinessUnit,
			BusinessUnitId:         user.BusinessUnitId,
			OrganizationUnit:       user.OrganizationUnit,
			Division:               user.Division,
			Department:             user.Department,
			JoinDate:               user.JoinDate,
			Grade:                  user.Grade,
			HRBP:                   user.HRBP,
			Position:               user.Position,
			Roles:                  user.Roles,
			ResetPasswordToken:     user.ResetPasswordToken,
			LastLogin:              user.LastLogin,
			ExpiredPasswordToken:   user.ExpiredPasswordToken,
			LastPasswordChanged:    user.LastPasswordChanged,
			ActualScores:           user.ActualScores,
			CalibrationScores:      user.CalibrationScores,
			SpmoCalibrations:       user.SpmoCalibrations,
			CalibratorCalibrations: user.CalibratorCalibrations,
			ScoringMethod:          user.ScoringMethod,
		})
	}
	if err != nil {
		return nil, err
	}

	return resultUsers, nil
}

func (r *projectRepo) GetNumberOneCalibrationsByPrevCalibratorBusinessUnit(calibratorId, prevCalibrator, businessUnit string, phase int, exceptUsers []string) ([]response.UserResponse, error) {
	var users []model.User
	var resultUsers []response.UserResponse

	subquery := r.db.
		Table("users u").
		Select("u.id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("p2.order < ? AND u3.name = ? AND b.name = ? AND c1.calibrator_id = ?",
			phase, prevCalibrator, businessUnit, calibratorId).
		Or("u.name = ? AND b.name = ? AND p.order = ? AND p2.order = ?",
			prevCalibrator, businessUnit, phase, phase)

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return nil, err
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
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Joins("JOIN calibrations c2 ON c2.employee_id = u.id").
		Joins("JOIN projects pr2 ON pr2.id = c2.project_id AND pr2.active = true").
		Joins("JOIN project_phases pp2 ON pp2.id = c2.project_phase_id").
		Joins("JOIN phases p2 ON p2.id = pp2.phase_id").
		Joins("JOIN users u3 ON c2.calibrator_id = u3.id").
		Where("(p.order <= ? AND u.id IN (?))", phase, subqueryResults).
		Or("(p.order = ? AND c1.calibrator_id = ? AND b.name = ? AND u.id NOT IN (?))", phase, calibratorId, businessUnit, exceptUsers).
		Group("u.id").
		Order("calibration_count ASC").
		Limit(10).Offset(0).
		Find(&users).Error

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return nil, err
		}

		resultUsers = append(resultUsers, response.UserResponse{
			BaseModel: model.BaseModel{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			CreatedBy:              user.CreatedBy,
			UpdatedBy:              user.UpdatedBy,
			Email:                  user.Email,
			Name:                   user.Name,
			Nik:                    user.Nik,
			DateOfBirth:            user.DateOfBirth,
			SupervisorNames:        supervisorName,
			BusinessUnit:           user.BusinessUnit,
			BusinessUnitId:         user.BusinessUnitId,
			OrganizationUnit:       user.OrganizationUnit,
			Division:               user.Division,
			Department:             user.Department,
			JoinDate:               user.JoinDate,
			Grade:                  user.Grade,
			HRBP:                   user.HRBP,
			Position:               user.Position,
			Roles:                  user.Roles,
			ResetPasswordToken:     user.ResetPasswordToken,
			LastLogin:              user.LastLogin,
			ExpiredPasswordToken:   user.ExpiredPasswordToken,
			LastPasswordChanged:    user.LastPasswordChanged,
			ActualScores:           user.ActualScores,
			CalibrationScores:      user.CalibrationScores,
			SpmoCalibrations:       user.SpmoCalibrations,
			CalibratorCalibrations: user.CalibratorCalibrations,
			ScoringMethod:          user.ScoringMethod,
		})
	}
	if err != nil {
		return nil, err
	}

	return resultUsers, nil
}

func (r *projectRepo) GetNMinusOneCalibrationsByBusinessUnit(businessUnit string, phase int) ([]response.UserResponse, error) {
	var users []model.User
	var resultUsers []response.UserResponse

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
		Preload("CalibrationScores.TopRemarks").
		Preload("CalibrationScores.BottomRemark").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Preload("BusinessUnit").
		Select("u.*").
		Joins("JOIN calibrations c1 ON c1.employee_id = u.id").
		Joins("JOIN projects pr ON pr.id = c1.project_id AND pr.active = true").
		Joins("JOIN project_phases pp ON pp.id = c1.project_phase_id").
		Joins("JOIN phases p ON p.id = pp.phase_id").
		Joins("JOIN business_units b ON u.business_unit_id = b.id").
		Joins("JOIN users u2 ON c1.calibrator_id = u2.id").
		Where("p.order = ? AND b.name = ? ", phase, businessUnit).
		Limit(10).Offset(0).
		Find(&users).Error

	for _, user := range users {
		var supervisorName string
		err = r.db.Raw("SELECT name FROM users WHERE nik = ?", user.SupervisorNik).Scan(&supervisorName).Error
		if err != nil {
			return nil, err
		}

		resultUsers = append(resultUsers, response.UserResponse{
			BaseModel: model.BaseModel{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			CreatedBy:              user.CreatedBy,
			UpdatedBy:              user.UpdatedBy,
			Email:                  user.Email,
			Name:                   user.Name,
			Nik:                    user.Nik,
			DateOfBirth:            user.DateOfBirth,
			SupervisorNames:        supervisorName,
			BusinessUnit:           user.BusinessUnit,
			BusinessUnitId:         user.BusinessUnitId,
			OrganizationUnit:       user.OrganizationUnit,
			Division:               user.Division,
			Department:             user.Department,
			JoinDate:               user.JoinDate,
			Grade:                  user.Grade,
			HRBP:                   user.HRBP,
			Position:               user.Position,
			Roles:                  user.Roles,
			ResetPasswordToken:     user.ResetPasswordToken,
			LastLogin:              user.LastLogin,
			ExpiredPasswordToken:   user.ExpiredPasswordToken,
			LastPasswordChanged:    user.LastPasswordChanged,
			ActualScores:           user.ActualScores,
			CalibrationScores:      user.CalibrationScores,
			SpmoCalibrations:       user.SpmoCalibrations,
			CalibratorCalibrations: user.CalibratorCalibrations,
			ScoringMethod:          user.ScoringMethod,
		})
	}
	if err != nil {
		return nil, err
	}

	for _, data := range resultUsers {
		fmt.Println(data.Name)
		fmt.Println(data.CalibrationScores)
	}

	return resultUsers, nil
}

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &projectRepo{
		db: db,
	}
}
