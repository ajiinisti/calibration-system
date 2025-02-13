package model

import (
	"time"
)

type User struct {
	BaseModel
	CreatedBy              string `gorm:"default:admin" json:"-"`
	UpdatedBy              string `gorm:"default:admin" json:"-"`
	Email                  string
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
	SupervisorNames        string
}

type UserCalibration struct {
	BaseModel
	CreatedBy              string `gorm:"default:admin" json:"-"`
	UpdatedBy              string `gorm:"default:admin" json:"-"`
	Email                  string `json:"-"`
	Name                   string
	Nik                    string
	DateOfBirth            time.Time    `gorm:"type:timestamp without time zone" json: "-"`
	SupervisorNik          string       `json:"-"`
	BusinessUnit           BusinessUnit `json:"-"`
	BusinessUnitId         *string      `json:"-"`
	OrganizationUnit       string
	Division               string
	Directorate            string
	Department             string
	JoinDate               time.Time `gorm:"type:timestamp without time zone" json:"-"`
	Grade                  string
	HRBP                   string `json:"-"`
	Position               string
	PhoneNumber            string        `json:"-"`
	GeneratePassword       bool          `gorm:"default:false" json:"-"`
	Password               string        `json:"-"`
	Roles                  []Role        `gorm:"many2many:user_roles" json:"-"`
	ResetPasswordToken     string        `json:"-"`
	ScoringMethod          string        `gorm:"default:Score"`
	LastLogin              time.Time     `gorm:"type:timestamp without time zone" json:"-"`
	ExpiredPasswordToken   time.Time     `gorm:"type:timestamp without time zone" json:"-"`
	LastPasswordChanged    time.Time     `gorm:"type:timestamp without time zone" json:"-"`
	ActualScores           []ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores      []Calibration `gorm:"foreignKey:EmployeeID"`
	SpmoCalibrations       []Calibration `gorm:"foreignKey:SpmoID" json:"-"`
	CalibratorCalibrations []Calibration `gorm:"foreignKey:CalibratorID" json:"-"`
	AccessTokenGenerate    string        `gorm:"unique;type:uuid;default:gen_random_uuid()" json:"-"`
	SupervisorNames        string
}

type UserChange struct {
	ID               string
	Email            string
	Division         string
	Name             string
	Nik              string
	BusinessUnitName string
	Roles            []string `gorm:"type:text[]"`
}

type UserShow struct {
	BaseModel
	CreatedBy         string `gorm:"default:admin" json:"-"`
	UpdatedBy         string `gorm:"default:admin" json:"-"`
	Email             string
	Name              string
	Nik               string
	DateOfBirth       time.Time `gorm:"type:timestamp without time zone"`
	SupervisorNik     string
	BusinessUnit      BusinessUnit
	BusinessUnitId    *string
	OrganizationUnit  string
	Division          string
	Directorate       string
	Department        string
	JoinDate          time.Time `gorm:"type:timestamp without time zone"`
	Grade             string
	HRBP              string
	Position          string
	PhoneNumber       string
	ScoringMethod     string            `gorm:"default:Score"`
	CalibrationScores []CalibrationForm `gorm:"foreignKey:EmployeeID"`
	ActualScores      []ActualScore     `gorm:"foreignKey:EmployeeID"`
	SupervisorNames   string
}
