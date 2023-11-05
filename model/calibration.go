package model

import (
	"time"

	"gorm.io/gorm"
)

type Calibration struct {
	CreatedAt             time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt             time.Time      `json:"-"`
	DeletedAt             gorm.DeletedAt `json:"-"`
	Project               Project
	ProjectID             string `gorm:"primaryKey"`
	ProjectPhase          ProjectPhase
	ProjectPhaseID        string `gorm:"primaryKey"`
	Employee              User
	EmployeeID            string `gorm:"primaryKey"`
	Calibrator            User
	CalibratorID          string
	Spmo                  User
	SpmoID                string
	Spmo2                 User
	Spmo2ID               *string
	Spmo3                 User
	Spmo3ID               *string
	Hrbp                  User
	HrbpID                string
	CalibrationScore      float64
	CalibrationRating     string
	Status                string `gorm:"default:'Waiting'"`
	SpmoStatus            string `gorm:"default:'-'"`
	Comment               string
	SpmoComment           string `gorm:"default:'-'"`
	CalibrationSpmoStatus string
	JustificationType     string       `gorm:"default:'default'"`
	BottomRemark          BottomRemark `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID"`
	TopRemarks            []TopRemark  `gorm:"foreignKey:ProjectID,EmployeeID,ProjectPhaseID;references:ProjectID,EmployeeID,ProjectPhaseID"`
}
