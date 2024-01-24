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
	BusinessUnit      BusinessUnitResponse
	BusinessUnitId    *string
	OrganizationUnit  string
	Division          string
	Department        string
	Grade             string
	Position          string
	Directorate       string
	ScoringMethod     string
	Roles             []model.Role          `gorm:"many2many:user_roles"`
	ActualScores      []ActualScoreResponse `gorm:"foreignKey:EmployeeID"`
	CalibrationScores []CalibrationResponse `gorm:"foreignKey:EmployeeID"`
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
	ProjectID                 string `gorm:"primaryKey"`
	ProjectPhase              ProjectPhaseResponse
	ProjectPhaseID            string `gorm:"primaryKey"`
	EmployeeID                string `gorm:"primaryKey"`
	Calibrator                CalibratorResponse
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
	Phase     PhaseResponse
	StartDate time.Time
	EndDate   time.Time
}

type PhaseResponse struct {
	Order int
}

type CalibratorResponse struct {
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
