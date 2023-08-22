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
	CreateUser(email string, role []string) error
	UpdateData(payload *model.User) error
}

type userUsecase struct {
	repo     repository.UserRepo
	role     RoleUsecase
	employee EmployeeUsecase
	cfg      *config.Config
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

func (u *userUsecase) CreateUser(email string, role []string) error {
	//Find user in Employee
	// _, err := u.employee.FindByEmail(email)
	// if err != nil {
	// 	return err
	// }

	// password, err := utils.GeneratePassword()
	// if err != nil {
	// 	return err
	// }

	fmt.Println("MASUK")
	password, err := utils.SaltPassword([]byte("password"))
	if err != nil {
		return err
	}

	fmt.Println("MASUK2")
	//Find Role
	var roles []model.Role
	for _, v := range role {
		getRole, err := u.role.FindByName(v)
		if err != nil {
			return err
		}
		roles = append(roles, *getRole)
	}

	fmt.Println("MASUK3")
	user := model.User{
		Email:    email,
		Password: password,
		Roles:    roles,
	}

	fmt.Println("New user", user)
	if err := u.repo.Save(&user); err != nil {
		return err
	}

	// body := fmt.Sprintf("Hi %s, You are registered to TalentConnect Platform\n\nYour Password is <b>%s</b>", payload.FirstName, password)
	// log.Println(body)
	// if err := utils.SendMail([]string{payload.Email}, "TalentConnect Registration", body, u.cfg.SMTPConfig); err != nil {
	// 	return err
	// }
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

func NewUserUseCase(repo repository.UserRepo, role RoleUsecase, employee EmployeeUsecase, cfg *config.Config) UserUsecase {
	return &userUsecase{
		repo:     repo,
		role:     role,
		employee: employee,
		cfg:      cfg,
	}
}
