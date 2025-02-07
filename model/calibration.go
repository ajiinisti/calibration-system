package model

import (
	"time"

	"gorm.io/gorm"
)

type Calibration struct {
	CreatedAt      time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Project        Project        `json:"-"`
	ProjectID      string         `gorm:"primaryKey"`
	ProjectPhase   ProjectPhase
	ProjectPhaseID string `gorm:"primaryKey"`
	Employee       User   `json:"-"`
	EmployeeID     string `gorm:"primaryKey"`
	Calibrator     User   `json:"-"`
	CalibratorID   string
	Spmo           User `json:"-"`
	SpmoID         string
	Spmo2          User `json:"-"`
	Spmo2ID        *string
	Spmo3          User `json:"-"`
	Spmo3ID        *string
	// Hrbp                  User
	// HrbpID                string
	CalibrationScore          float64
	CalibrationRating         string
	Status                    string `gorm:"default:Waiting"`
	SpmoStatus                string `gorm:"default:'-'"`
	Comment                   string
	SpmoComment               string `gorm:"default:-"`
	JustificationType         string `gorm:"default:default"`
	JustificationReviewStatus bool   `gorm:"default:false"`
	SendBackDeadline          time.Time
	BottomRemark              BottomRemark `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	TopRemarks                []TopRemark  `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	FilledTopBottomMark       bool
}

type SeeCalibrationJustification struct {
	CreatedAt         time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"`
	Project           Project        `json:"-"`
	ProjectID         string         `gorm:"primaryKey" json:"-"`
	ProjectPhase      ProjectPhase
	ProjectPhaseID    string `gorm:"primaryKey" json:"-"`
	Employee          User   `json:"-"`
	EmployeeID        string `gorm:"primaryKey" json:"-"`
	Calibrator        User
	CalibratorID      string `json:"-"`
	Comment           string
	SpmoComment       string       `gorm:"default:-"`
	JustificationType string       `gorm:"default:default"`
	BottomRemark      BottomRemark `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	TopRemarks        []TopRemark  `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
}

type CalibrationForm struct {
	CreatedAt                 time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt                 time.Time      `json:"-"`
	DeletedAt                 gorm.DeletedAt `json:"-"`
	Project                   Project        `json:"-"`
	ProjectID                 string         `gorm:"primaryKey"`
	ProjectPhase              ProjectPhase
	ProjectPhaseID            string `gorm:"primaryKey"`
	Employee                  User
	EmployeeID                string `gorm:"primaryKey"`
	Calibrator                User
	CalibratorID              string
	Spmo                      User
	SpmoID                    string
	Spmo2                     User
	Spmo2ID                   *string
	Spmo3                     User
	Spmo3ID                   *string
	CalibrationScore          float64
	CalibrationRating         string
	Status                    string `gorm:"default:Waiting"`
	SpmoStatus                string `gorm:"default:'-'"`
	Comment                   string
	SpmoComment               string `gorm:"default:-"`
	JustificationType         string `gorm:"default:default"`
	JustificationReviewStatus bool   `gorm:"default:false"`
	SendBackDeadline          time.Time
	BottomRemark              BottomRemark `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	TopRemarks                []TopRemark  `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID;constraint:OnDelete:CASCADE"`
	FilledTopBottomMark       bool
}

func (CalibrationForm) TableName() string {
	return "calibrations"
}
