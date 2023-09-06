package model

import (
	"time"
)

type User struct {
	BaseModel
	CreatedBy              string `gorm:"default:admin" json:"-"`
	UpdatedBy              string `gorm:"default:admin" json:"-"`
	Email                  string `gorm:"unique" `
	Name                   string
	Nik                    string
	DateOfBirth            time.Time `gorm:"type:timestamp without time zone"`
	SupervisorName         string
	BusinessUnit           BusinessUnit
	BusinessUnitId         *string
	OrganizationUnit       string
	Division               string
	Department             string
	JoinDate               time.Time `gorm:"type:timestamp without time zone"`
	Grade                  string
	HRBP                   string
	Position               string
	Password               string `json:"-"`
	Roles                  []Role `gorm:"many2many:user_roles"`
	ResetPasswordToken     string
	LastLogin              time.Time     `gorm:"type:timestamp without time zone"`
	ExpiredPasswordToken   time.Time     `gorm:"type:timestamp without time zone"`
	LastPasswordChanged    time.Time     `gorm:"type:timestamp without time zone"`
	ActualScores           []ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores      []Calibration `gorm:"foreignKey:EmployeeID"`
	SpmoCalibrations       []Calibration `gorm:"foreignKey:SpmoID"`
	CalibratorCalibrations []Calibration `gorm:"foreignKey:CalibratorID"`
	// Username  string `gorm:"unique" `
	// EmployeeId string `gorm:"unique"`
}
