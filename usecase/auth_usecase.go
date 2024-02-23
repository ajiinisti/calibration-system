package usecase

import (
	"fmt"
	"strings"
	"time"

	"calibration-system.com/delivery/api/request"
	"calibration-system.com/model"
	"calibration-system.com/utils"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Login(payload request.Login) (*model.User, error)
	// ChangePassword(email string, requestData request.ChangePassword) error
	ForgetPassword(email string, resetToken string) error
	GetUserByEmail(email string) (*model.User, error)
	ResetPassword(email string, resetToken string, newPassword string, confirmPassword string) error
	CheckToken(token string) (*model.User, error)
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
func (a *authUsecase) ForgetPassword(email string, resetToken string) error {
	user, err := a.user.SearchEmail(email)
	if err != nil {
		return err
	}

	user.ResetPasswordToken = resetToken
	user.ExpiredPasswordToken = time.Now().Add(time.Minute * 10)
	return a.user.UpdateData(user)
}

func (a *authUsecase) ResetPassword(email string, resetToken string, newPassword string, confirmPassword string) error {
	user, err := a.user.SearchEmail(email)
	if err != nil {
		return err
	}

	if time.Now().After(user.ExpiredPasswordToken) {
		return fmt.Errorf("Your reset token has been expired")
	}

	if resetToken != user.ResetPasswordToken {
		return fmt.Errorf("Your reset token is invalid")
	}

	if newPassword != confirmPassword {
		return fmt.Errorf("Your new password and confirm new password isn't the same")
	}

	hashedPassword, err := utils.SaltPassword([]byte(newPassword))
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.LastPasswordChanged = time.Now()
	return a.user.UpdateData(user)
}

// update password from forgetForm

// verifyLogin implements AuthUsecase
func (a *authUsecase) Login(payload request.Login) (*model.User, error) {
	user, err := a.user.SearchEmail(strings.ToLower(payload.Email))
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

func (a *authUsecase) CheckToken(token string) (*model.User, error) {
	user, err := a.user.FindByGenerateToken(token)
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("Token invalid")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewAuthUsecase(user UserUsecase) AuthUsecase {
	return &authUsecase{
		user: user,
	}
}
