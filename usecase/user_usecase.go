package usecase

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/model"
	"calibration-system.com/repository"
	"calibration-system.com/utils"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type UserUsecase interface {
	BaseUsecase[model.User]
	SearchEmail(email string) (*model.User, error)
	CreateUser(payload model.User, role []string) error
	SaveUser(payload model.User, role []string) error
	UpdateData(payload *model.User) error
	BulkInsert(file *multipart.FileHeader) ([]string, error)
}

type userUsecase struct {
	repo repository.UserRepo
	role RoleUsecase
	bu   BusinessUnitUsecase
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

func (u *userUsecase) BulkInsert(file *multipart.FileHeader) ([]string, error) {
	var logs []string
	var users []model.User
	var businessunits []model.BusinessUnit

	// Membuka file Excel yang diunggah
	excelFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()

	xlsFile, err := excelize.OpenReader(excelFile)
	if err != nil {
		return nil, err
	}

	sheetName := xlsFile.GetSheetName(2) // Ganti dengan nama sheet yang sesuai
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		if i == 0 {
			// Skip the first row
			continue
		}

		num, err := strconv.Atoi(row[2])
		if err != nil {
			return logs, err
		}
		dateValue := time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, num-1)

		nik := row[0]
		name := row[1]
		supervisor := row[3]
		buId := row[4]
		orgUnit := row[5]
		division := row[6]
		department := row[7]
		position := row[8]
		grade := row[9]
		hrbp := row[10]
		email := row[11]

		user := model.User{
			Email:            email,
			Name:             name,
			Nik:              nik,
			DateOfBirth:      time.Time{},
			SupervisorName:   supervisor,
			BusinessUnitId:   buId,
			OrganizationUnit: orgUnit,
			Division:         division,
			Department:       department,
			JoinDate:         dateValue,
			Grade:            grade,
			HRBP:             hrbp,
			Position:         position,
		}

		var found bool
		for _, bu := range businessunits {
			if bu.ID == buId {
				user.BusinessUnitId = bu.ID
				found = true
				break
			}
		}

		if !found {
			bu, err := u.bu.FindById(buId)
			if err != nil {
				logs = append(logs, fmt.Sprintf("Error Business Unit on Row %d Employee %s", i, name))
				break
			}
			user.BusinessUnitId = bu.ID
			businessunits = append(businessunits, *bu)
		}

		users = append(users, user)
	}

	if len(logs) > 0 {
		return logs, fmt.Errorf("Error when insert data")
	}

	err = u.repo.Bulksave(&users)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func NewUserUseCase(
	repo repository.UserRepo,
	role RoleUsecase,
	bu BusinessUnitUsecase,
	cfg *config.Config,
) UserUsecase {
	return &userUsecase{
		repo: repo,
		role: role,
		bu:   bu,
		cfg:  cfg,
	}
}
