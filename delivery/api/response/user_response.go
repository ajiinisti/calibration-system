package response

import (
	"calibration-system.com/model"
)

type UserResponse struct {
	model.BaseModel
	CreatedBy         string `gorm:"default:admin" json:"-"`
	UpdatedBy         string `gorm:"default:admin" json:"-"`
	Email             string `gorm:"unique" `
	Name              string
	Nik               string
	SupervisorNames   string
	BusinessUnit      model.BusinessUnit
	BusinessUnitId    *string
	OrganizationUnit  string
	Division          string
	Department        string
	Grade             string
	Position          string
	Directorate       string
	ScoringMethod     string
	Roles             []model.Role        `gorm:"many2many:user_roles"`
	ActualScores      []model.ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores []model.Calibration `gorm:"foreignKey:EmployeeID"`
}

type UserCalibration struct {
	NPlusOneManager     bool
	SendToManager       bool
	SendBackCalibration bool
	UserData            []UserResponse
}
