package repository

import (
	"fmt"

	"calibration-system.com/delivery/api/response"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type UserRepo interface {
	BaseRepository[model.User]
	SearchByEmail(email string) (*model.User, error)
	SearchByNik(nik string) (*model.User, error)
	SearchByGenerateToken(generateToken string) (*model.User, error)
	Update(payload *model.User) error
	Bulksave(payload *[]model.User) error
	PaginateList(pagination model.PaginationQuery) ([]model.User, response.Paging, error)
	PaginateByProjectId(pagination model.PaginationQuery, projectId string) ([]model.User, response.Paging, error)
	GetTotalRows(name string) (int, error)
	GetTotalRowsByProjectID(projectId, name string) (int, error)
	ListUserAdmin() ([]model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func (u *userRepo) SearchByNik(nik string) (*model.User, error) {
	var user model.User
	err := u.db.Preload("Roles").First(&user, "nik = ?", nik).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) SearchByGenerateToken(generateToken string) (*model.User, error) {
	var user model.User
	err := u.db.Preload("Roles").First(&user, "access_token_generate = ?", generateToken).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) SearchByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.db.Preload("Roles").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Save(payload *model.User) error {
	err := u.db.Save(&payload)
	if err.Error != nil {
		return fmt.Errorf(err.Error.Error() + payload.ID)
	}
	return nil
}

func (u *userRepo) Bulksave(payload *[]model.User) error {
	tx := u.db.Begin()
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

func (u *userRepo) Get(id string) (*model.User, error) {
	var user model.User
	if id != "" {
		err := u.db.
			Preload("Roles").
			// Preload("ActualScores").
			// Preload("CalibrationScores").
			Preload("BusinessUnit").
			First(&user, "id = ?", id).Error
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepo) List() ([]model.User, error) {
	var users []model.User
	err := u.db.
		Preload("Roles").
		Preload("BusinessUnit").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepo) ListUserAdmin() ([]model.User, error) {
	subquery := u.db.
		Table("users u").
		Select("u.id").
		Joins("JOIN user_roles ur on u.id = ur.user_id").
		Joins("JOIN roles r on ur.role_id = r.id").
		Where("r.name = 'exclude'")

	var subqueryResults []string
	if err := subquery.Pluck("u.id", &subqueryResults).Error; err != nil {
		return nil, err
	}

	var users []model.User
	err := u.db.
		Table("users u").
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Where("name <> 'exclude'")
		}).
		Preload("BusinessUnit").
		Select("u.*").
		Where("u.id NOT IN (?)", subqueryResults).
		Distinct().
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepo) Delete(id string) error {
	result := u.db.Delete(&model.User{
		BaseModel: model.BaseModel{ID: id},
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("Uer not found!")
	}
	return nil
}

func (u *userRepo) Update(payload *model.User) error {
	payloadRole := payload.Roles
	err := u.db.Model(&payload).Association("Roles").Clear()
	if err != nil {
		return err
	}

	payload.Roles = payloadRole

	if err := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&payload); err.Error != nil {
		return err.Error
	}
	return nil
}

func (u *userRepo) PaginateList(pagination model.PaginationQuery) ([]model.User, response.Paging, error) {
	var users []model.User
	var err error
	if pagination.Name == "" {
		err := u.db.
			Preload("Roles").
			Preload("BusinessUnit").
			Limit(pagination.Take).Offset(pagination.Skip).Find(&users).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	} else {
		fmt.Println(len(pagination.Name), pagination.Name)
		err := u.db.
			Preload("Roles").
			Preload("BusinessUnit").
			Where("name ILIKE ?", "%"+pagination.Name+"%").
			Or("nik ILIKE ?", "%"+pagination.Name+"%").
			Limit(pagination.Take).Offset(pagination.Skip).Find(&users).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	}

	totalRows, err := u.GetTotalRows(pagination.Name)
	if err != nil {
		return nil, response.Paging{}, err
	}

	return users, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (u *userRepo) PaginateByProjectId(pagination model.PaginationQuery, projectId string) ([]model.User, response.Paging, error) {
	var users []model.User
	var err error

	if pagination.Name == "" {
		err = u.db.
			Preload("Roles").
			Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
				return db.
					Joins("JOIN projects proj2 ON actual_scores.project_id = proj2.id").
					Where("proj2.id = ?", projectId)
			}).
			Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
				return db.
					Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
					Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
					Joins("JOIN phases p ON p.id = pp.phase_id ").
					Where("proj2.id = ?", projectId).
					Order("p.order ASC")
			}).
			Preload("CalibrationScores.Calibrator").
			Preload("CalibrationScores.Spmo").
			Preload("CalibrationScores.ProjectPhase").
			Preload("CalibrationScores.ProjectPhase.Phase").
			Joins("LEFT JOIN actual_scores ON users.id = actual_scores.employee_id AND actual_scores.deleted_at IS NULL").
			Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id AND calibrations.deleted_at IS NULL").
			Where("(actual_scores.project_id = ? OR calibrations.project_id = ?)", projectId, projectId).
			Group("users.id").
			Limit(pagination.Take).Offset(pagination.Skip).
			Find(&users).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	} else {
		err = u.db.
			Preload("Roles").
			Preload("ActualScores", func(db *gorm.DB) *gorm.DB {
				return db.
					Joins("JOIN projects proj2 ON actual_scores.project_id = proj2.id").
					Where("proj2.active = ?", true)
			}).
			Preload("CalibrationScores", func(db *gorm.DB) *gorm.DB {
				return db.
					Joins("JOIN projects proj2 ON calibrations.project_id = proj2.id").
					Joins("JOIN project_phases pp ON pp.id = calibrations.project_phase_id").
					Joins("JOIN phases p ON p.id = pp.phase_id ").
					Where("proj2.active = ?", true).
					Order("p.order ASC")
			}).
			Preload("CalibrationScores.Calibrator").
			Preload("CalibrationScores.Spmo").
			Preload("CalibrationScores.ProjectPhase").
			Preload("CalibrationScores.ProjectPhase.Phase").
			Joins("LEFT JOIN actual_scores ON users.id = actual_scores.employee_id AND actual_scores.deleted_at IS NULL").
			Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id AND calibrations.deleted_at IS NULL").
			Where("(actual_scores.project_id = ? AND calibrations.project_id = ?) AND (name ILIKE ? OR nik ILIKE ?)", projectId, projectId, "%"+pagination.Name+"%", "%"+pagination.Name+"%").
			Group("users.id").
			Limit(pagination.Take).Offset(pagination.Skip).
			Find(&users).Error
		if err != nil {
			return nil, response.Paging{}, err
		}
	}

	totalRows, err := u.GetTotalRowsByProjectID(projectId, pagination.Name)
	if err != nil {
		return nil, response.Paging{}, err
	}

	return users, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (u *userRepo) GetTotalRows(name string) (int, error) {
	var count int64
	var err error
	if name == "" {
		err = u.db.
			Model(&model.User{}).
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	} else {
		err = u.db.
			Model(&model.User{}).
			Where("name ILIKE ?", "%"+name+"%").
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	}
	return int(count), nil
}

func (u *userRepo) GetTotalRowsByProjectID(projectId, name string) (int, error) {
	var count int64
	var err error

	if name == "" {
		err = u.db.
			Model(&model.User{}).
			Joins("JOIN actual_scores ON users.id = actual_scores.employee_id").
			Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id").
			Where("actual_scores.project_id = ? OR calibrations.project_id = ?", projectId, projectId).
			Group("users.id").
			Count(&count).Error
		if err != nil {
			return 0, err
		}

	} else {
		err = u.db.
			Model(&model.User{}).
			Joins("JOIN actual_scores ON users.id = actual_scores.employee_id").
			Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id").
			Where("(actual_scores.project_id = ? OR calibrations.project_id = ?) AND name ILIKE ?", projectId, projectId, "%"+name+"%").
			Group("users.id").
			Count(&count).Error
		if err != nil {
			return 0, err
		}
	}
	return int(count), nil
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}
