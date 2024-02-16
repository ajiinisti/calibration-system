package model

import (
	"time"
)

type User struct {
	BaseModel
	CreatedBy              string `gorm:"default:admin" json:"-"`
	UpdatedBy              string `gorm:"default:admin" json:"-"`
	Email                  string `gorm:"unique"`
	Name                   string
	Nik                    string
	DateOfBirth            time.Time `gorm:"type:timestamp without time zone"`
	SupervisorNik          string
	BusinessUnit           BusinessUnit
	BusinessUnitId         *string
	OrganizationUnit       string
	Division               string
	Directorate            string
	Department             string
	JoinDate               time.Time `gorm:"type:timestamp without time zone"`
	Grade                  string
	HRBP                   string
	Position               string
	PhoneNumber            string
	GeneratePassword       bool   `gorm:"default:false"`
	Password               string `json:"-"`
	Roles                  []Role `gorm:"many2many:user_roles"`
	ResetPasswordToken     string
	ScoringMethod          string        `gorm:"default:Score"`
	LastLogin              time.Time     `gorm:"type:timestamp without time zone"`
	ExpiredPasswordToken   time.Time     `gorm:"type:timestamp without time zone"`
	LastPasswordChanged    time.Time     `gorm:"type:timestamp without time zone"`
	ActualScores           []ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores      []Calibration `gorm:"foreignKey:EmployeeID"`
	SpmoCalibrations       []Calibration `gorm:"foreignKey:SpmoID"`
	CalibratorCalibrations []Calibration `gorm:"foreignKey:CalibratorID"`
	AccessTokenGenerate    string        `gorm:"unique;type:uuid;default:gen_random_uuid()"`
}
