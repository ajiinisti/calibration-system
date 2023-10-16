package response

import (
	"time"

	"calibration-system.com/model"
)

type UserResponse struct {
	model.BaseModel
	CreatedBy              string `gorm:"default:admin" json:"-"`
	UpdatedBy              string `gorm:"default:admin" json:"-"`
	Email                  string `gorm:"unique" `
	Name                   string
	Nik                    string
	DateOfBirth            time.Time `gorm:"type:timestamp without time zone"`
	SupervisorNames        string
	BusinessUnit           model.BusinessUnit
	BusinessUnitId         *string
	OrganizationUnit       string
	Division               string
	Department             string
	JoinDate               time.Time `gorm:"type:timestamp without time zone"`
	Grade                  string
	HRBP                   string
	Position               string
	GeneratePassword       bool         `gorm:"default:false"`
	Password               string       `json:"-"`
	Roles                  []model.Role `gorm:"many2many:user_roles"`
	ResetPasswordToken     string
	LastLogin              time.Time           `gorm:"type:timestamp without time zone"`
	ExpiredPasswordToken   time.Time           `gorm:"type:timestamp without time zone"`
	LastPasswordChanged    time.Time           `gorm:"type:timestamp without time zone"`
	ActualScores           []model.ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores      []model.Calibration `gorm:"foreignKey:EmployeeID"`
	SpmoCalibrations       []model.Calibration `gorm:"foreignKey:SpmoID"`
	CalibratorCalibrations []model.Calibration `gorm:"foreignKey:CalibratorID"`
}
