package usecase

import (
	"fmt"

	"calibration-system.com/config"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
)

type UserUsecase interface {
	BaseUsecase[model.User]
	SearchEmail(email string) (*model.User, error)
	CreateUser(payload model.User, role []string) error
	SaveUser(payload model.User, role []string) error
	UpdateData(payload *model.User) error
}

type userUsecase struct {
	repo repository.UserRepo
	role RoleUsecase
	cfg  *config.Config
}

func (u *userUsecase) SearchEmail(email string) (*model.User, error) {
	return u.repo.SearchByEmail(email)
}
func (u *userUsecase) FindAll() ([]model.User, error) {
	return u.repo.List()
}

func (u *userUsecase) FindById(id string) (*model.User, error) {
	return u.repo.Get(id)
}

func (u *userUsecase) CreateUser(payload model.User, role []string) error {
	var password string
	if len(role) > 0 {
		var err error
		password, err = utils.SaltPassword([]byte("password"))
		if err != nil {
			return err
		}

	}

	//Find Role
	var roles []model.Role
	for _, v := range role {
		getRole, err := u.role.FindByName(v)
		if err != nil {
			return err
		}
		roles = append(roles, *getRole)
	}

	payload.Password = password
	payload.Roles = roles

	if err := u.repo.Save(&payload); err != nil {
		return err
	}

	// body := fmt.Sprintf("Hi %s, You are registered to TalentConnect Platform\n\nYour Password is <b>%s</b>", payload.FirstName, password)
	// log.Println(body)
	// if err := utils.SendMail([]string{payload.Email}, "TalentConnect Registration", body, u.cfg.SMTPConfig); err != nil {
	// 	return err
	// }
	return nil
}

func (u *userUsecase) SaveUser(payload model.User, role []string) error {
	//Find Role
	var roles []model.Role
	for _, v := range role {
		getRole, err := u.role.FindByName(v)
		if err != nil {
			return err
		}
		roles = append(roles, *getRole)
	}
	payload.Roles = roles
	fmt.Println("DATA", payload)

	if err := u.repo.Update(&payload); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) SaveData(payload *model.User) error {
	return u.repo.Save(payload)
}

func (u *userUsecase) DeleteData(id string) error {
	return u.repo.Delete(id)
}

func (u *userUsecase) UpdateData(payload *model.User) error {
	return u.repo.Update(payload)
}

func NewUserUseCase(repo repository.UserRepo, role RoleUsecase, cfg *config.Config) UserUsecase {
	return &userUsecase{
		repo: repo,
		role: role,
		cfg:  cfg,
	}
}
