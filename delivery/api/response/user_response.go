package response

import (
	"time"

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
	ActualScores      []model.ActualScore `gorm:"foreignKey:EmployeeID"`
	CalibrationScores []model.Calibration `gorm:"foreignKey:EmployeeID"`
}

type ActualScoreResponse struct {
	ProjectID    string `gorm:"primaryKey"`
	EmployeeID   string `gorm:"primaryKey"`
	ActualScore  float64
	ActualRating string
	Y1Rating     string
	Y2Rating     string
	PTTScore     float64
	PATScore     float64
	Score360     float64
}

type BusinessUnitResponse struct {
	ID                  string `gorm:"unique"`
	Status              bool
	Name                string
	GroupBusinessUnitId string
}

type CalibrationResponse struct {
	ProjectID                 string               `gorm:"primaryKey"`
	ProjectPhase              ProjectPhaseResponse `gorm:"foreignKey:ProjectPhaseID;references:ID"`
	ProjectPhaseID            string               `gorm:"primaryKey"`
	EmployeeID                string               `gorm:"primaryKey"`
	Calibrator                CalibratorResponse   `gorm:"foreignKey:CalibratorID;references:ID"`
	CalibratorID              string
	CalibrationScore          float64
	CalibrationRating         string
	Status                    string `gorm:"default:Waiting"`
	SpmoStatus                string `gorm:"default:'-'"`
	Comment                   string
	SpmoComment               string `gorm:"default:-"`
	JustificationType         string `gorm:"default:default"`
	JustificationReviewStatus bool   `gorm:"default:false"`
	SendBackDeadline          time.Time
	BottomRemark              BottomRemarkResponse `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	TopRemarks                []TopRemarkResponse  `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
}

type ProjectPhaseResponse struct {
	model.BaseModel
	Phase     PhaseResponse `gorm:"foreignKey:PhaseID"`
	PhaseID   string
	StartDate time.Time
	EndDate   time.Time
}

type PhaseResponse struct {
	model.BaseModel
	Order int
}

type CalibratorResponse struct {
	model.BaseModel
	Name string
}

type BottomRemarkResponse struct {
	ProjectID      string `gorm:"primaryKey"`
	EmployeeID     string `gorm:"primaryKey"`
	ProjectPhaseID string `gorm:"primaryKey"`
	LowPerformance string
	Indisipliner   string
	Attitude       string
	WarningLetter  string
}

type UserCalibration struct {
	NPlusOneManager     bool
	SendToManager       bool
	SendBackCalibration bool
	UserData            []UserResponse
}

type UserCalibrationNew struct {
	NPlusOneManager     bool
	SendToManager       bool
	SendBackCalibration bool
	UserData            []model.UserCalibration
}
