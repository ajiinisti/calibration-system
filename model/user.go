package model

import "time"

type User struct {
	BaseModel
	CreatedBy            string `gorm:"default:admin" json:"-"`
	UpdatedBy            string `gorm:"default:admin" json:"-"`
	Email                string `gorm:"unique" `
	Name                 string
	Password             string `json:"-"`
	LastLogin            time.Time
	Roles                []Role `gorm:"many2many:user_roles"`
	ResetPasswordToken   string
	ExpiredPasswordToken time.Time
	LastPasswordChanged  time.Time
	// Username  string `gorm:"unique" `
	// EmployeeId string `gorm:"unique"`
}
