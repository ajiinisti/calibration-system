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
	Update(payload *model.User) error
	Bulksave(payload *[]model.User) error
	PaginateList(pagination model.PaginationQuery) ([]model.User, response.Paging, error)
	PaginateByProjectId(pagination model.PaginationQuery, projectId string) ([]model.User, response.Paging, error)
	GetTotalRows() (int, error)
	GetTotalRowsByProjectID(projectId string) (int, error)
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
	batchSize := 100
	numFullBatches := len(*payload) / batchSize

	for i := 0; i < numFullBatches; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		currentBatch := (*payload)[start:end]
		return u.db.Save(&currentBatch).Error

	}
	remainingItems := (*payload)[numFullBatches*batchSize:]

	if len(remainingItems) > 0 {
		err := u.db.Save(&remainingItems)
		if err != nil {
			return u.db.Save(&remainingItems).Error
		}
	}
	return nil
}

func (u *userRepo) Get(id string) (*model.User, error) {
	var user model.User
	err := u.db.
		Preload("Roles").
		Preload("ActualScores").
		Preload("CalibrationScores").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) List() ([]model.User, error) {
	var users []model.User
	err := u.db.
		Preload("Roles").
		Preload("ActualScores").
		Preload("CalibrationScores").
		Preload("BusinessUnit").
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
	if err := u.db.Updates(&payload); err.Error != nil {
		return err.Error
	}
	return nil
}

func (u *userRepo) PaginateList(pagination model.PaginationQuery) ([]model.User, response.Paging, error) {
	var users []model.User
	err := u.db.
		Preload("Roles").
		Limit(pagination.Take).Offset(pagination.Skip).Find(&users).Error
	if err != nil {
		return nil, response.Paging{}, err
	}

	totalRows, err := u.GetTotalRows()
	if err != nil {
		return nil, response.Paging{}, err
	}

	return users, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (u *userRepo) PaginateByProjectId(pagination model.PaginationQuery, projectId string) ([]model.User, response.Paging, error) {
	var users []model.User
	err := u.db.
		Preload("Roles").
		Preload("ActualScores").
		Preload("CalibrationScores").
		Preload("CalibrationScores.Calibrator").
		Preload("CalibrationScores.Spmo").
		Preload("CalibrationScores.Hrbp").
		Preload("CalibrationScores.ProjectPhase").
		Preload("CalibrationScores.ProjectPhase.Phase").
		Joins("JOIN actual_scores ON users.id = actual_scores.employee_id").
		Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id").
		Where("actual_scores.project_id = ? OR calibrations.project_id = ?", projectId, projectId).
		Group("users.id").
		Limit(pagination.Take).Offset(pagination.Skip).
		Find(&users).Error
	if err != nil {
		return nil, response.Paging{}, err
	}

	totalRows, err := u.GetTotalRowsByProjectID(projectId)
	if err != nil {
		return nil, response.Paging{}, err
	}

	return users, utils.Paginate(pagination.Page, pagination.Take, totalRows), nil
}

func (u *userRepo) GetTotalRows() (int, error) {
	var count int64
	err := u.db.
		Model(&model.User{}).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (u *userRepo) GetTotalRowsByProjectID(projectId string) (int, error) {
	var count int64
	err := u.db.
		Model(&model.User{}).
		Joins("JOIN actual_scores ON users.id = actual_scores.employee_id").
		Joins("LEFT JOIN calibrations ON users.id = calibrations.employee_id").
		Where("actual_scores.project_id = ? OR calibrations.project_id = ?", projectId, projectId).
		Group("users.id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}
