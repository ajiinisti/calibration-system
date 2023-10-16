package usecase

import (
	"fmt"
	"mime/multipart"
	"sort"
	"time"

	"calibration-system.com/config"
	"calibration-system.com/delivery/api/request"
	"calibration-system.com/delivery/api/response"
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
	FindByNik(nik string) (*model.User, error)
	FindPagination(param request.PaginationParam) ([]model.User, response.Paging, error)
	FindByProjectIdPagination(param request.PaginationParam, projectId string) ([]model.User, response.Paging, error)
	GeneratePasswordById(id string) error
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

func (u *userUsecase) FindByNik(nik string) (*model.User, error) {
	return u.repo.SearchByNik(nik)
}

func (u *userUsecase) FindAll() ([]model.User, error) {
	return u.repo.List()
}

func (u *userUsecase) FindById(id string) (*model.User, error) {
	return u.repo.Get(id)
}

func (u *userUsecase) FindPagination(param request.PaginationParam) ([]model.User, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)
	return u.repo.PaginateList(paginationQuery)
}

func (u *userUsecase) FindByProjectIdPagination(param request.PaginationParam, projectId string) ([]model.User, response.Paging, error) {
	paginationQuery := utils.GetPaginationParams(param)

	users, paging, err := u.repo.PaginateByProjectId(paginationQuery, projectId)
	if err != nil {
		return []model.User{}, response.Paging{}, err
	}

	for _, user := range users {
		sort.Slice(user.CalibrationScores, func(i, j int) bool {
			return user.CalibrationScores[i].ProjectPhase.Phase.Order < user.CalibrationScores[j].ProjectPhase.Phase.Order
		})
	}

	return users, paging, err
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
	var dateLogs []string
	var buLogs []string
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

	sheetName := xlsFile.GetSheetName(2)
	rows := xlsFile.GetRows(sheetName)

	for i, row := range rows {
		passed := true
		if i == 0 {
			continue
		}

		layout := "01-02-06"
		parsedTime, err := time.Parse(layout, row[2])
		if err != nil {
			dateLogs = append(dateLogs, fmt.Sprintf("Cannot parse date on Employee NIK %s", row[1]))
			passed = false
		}

		var found bool
		buId := row[4]
		for _, bu := range businessunits {
			if bu.ID == buId {
				found = true
				break
			}
		}

		if !found {
			bu, err := u.bu.FindById(buId)
			if err != nil {
				buLogs = append(buLogs, fmt.Sprintf("Error Business Unit Id on Row %d", i+1))
				passed = false
			} else {
				businessunits = append(businessunits, *bu)
			}
		}

		user := model.User{
			Email:            row[11],
			Name:             row[1],
			Nik:              row[0],
			DateOfBirth:      time.Time{},
			SupervisorNik:    row[3],
			BusinessUnitId:   nil,
			OrganizationUnit: row[5],
			Division:         row[6],
			Department:       row[7],
			JoinDate:         parsedTime,
			Grade:            row[9],
			HRBP:             row[10],
			Position:         row[8],
			GeneratePassword: false,
		}

		if passed {
			user.BusinessUnitId = &buId
		}
		users = append(users, user)
	}

	logs = append(logs, dateLogs...)
	logs = append(logs, buLogs...)

	err = u.repo.Bulksave(&users)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (u *userUsecase) GeneratePasswordById(id string) error {
	password, err := utils.SaltPassword([]byte("password"))
	if err != nil {
		return err
	}

	user := model.User{
		BaseModel: model.BaseModel{
			ID: id,
		},
		Password:         password,
		GeneratePassword: true,
	}

	if err := u.repo.Update(&user); err != nil {
		return err
	}

	return nil
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
