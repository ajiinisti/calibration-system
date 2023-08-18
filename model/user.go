package model

import "time"

type User struct {
	BaseModel
	CreatedBy string `gorm:"default:admin" json:"-"`
	UpdatedBy string `gorm:"default:admin" json:"-"`
	Email     string `gorm:"unique" `
	Password  string `json:"-"`
	// Username  string `gorm:"unique" `
	// EmployeeId string `gorm:"unique"`
	LastLogin            time.Time
	Role                 Role
	RoleID               string
	ResetPasswordToken   string `gorm:"type:uuid"`
	ExpiredPasswordToken time.Time
	LastPasswordChanged  time.Time
}
