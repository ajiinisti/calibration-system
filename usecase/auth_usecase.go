package usecase

import (
	"fmt"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(payload request.Login) (*model.User, error)
	// ChangePassword(email string, requestData request.ChangePassword) error
	ForgetPassword() error
	GetUserByEmail(email string) (*model.User, error)
}

type authUsecase struct {
	user UserUsecase
}

// changePassword implements AuthUsecase
// func (a *authUsecase) ChangePassword(email string, requestData request.ChangePassword) error {
// 	user, err := a.user.SearchEmail(email)
// 	if err != nil {
// 		return err
// 	}
// 	if !utils.ComparePassword(user.Password, []byte(requestData.CurrentPassword)) {
// 		return fmt.Errorf("Password not valid")
// 	}
// 	newPassword, err := utils.SaltPassword([]byte(requestData.NewPassword))
// 	if err != nil {
// 		return err
// 	}
// 	user.Password = newPassword
// 	return a.user.UpdateData(user)
// 	panic("unimplemented")
// }

// forgetPassword implements AuthUsecase
func (*authUsecase) ForgetPassword() error {
	// check if email exists
	// add key to redis (different DB)
	// send email
	panic("unimplemented")
}

// update password from forgetForm

// verifyLogin implements AuthUsecase
func (a *authUsecase) Login(payload request.Login) (*model.User, error) {
	user, err := a.user.SearchEmail(payload.Email)
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("Email/Password invalid")
	}
	if err != nil {
		return nil, err
	}
	if !utils.ComparePassword(user.Password, []byte(payload.Password)) {
		return nil, fmt.Errorf("Email/Password invalid")
	}
	return user, nil
}

func (a *authUsecase) GetUserByEmail(email string) (*model.User, error) {
	return a.user.SearchEmail(email)
}

func NewAuthUsecase(user UserUsecase) AuthUsecase {
	return &authUsecase{
		user: user,
	}
}
