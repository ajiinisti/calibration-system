package model

import (
	"time"

	"gorm.io/gorm"
)

type Calibration struct {
	CreatedAt         time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-"`
	Project           Project
	ProjectID         string `gorm:"primaryKey"`
	ProjectPhase      ProjectPhase
	ProjectPhaseID    string `gorm:"primaryKey"`
	Employee          User
	EmployeeID        string `gorm:"primaryKey"`
	Calibrator        User
	CalibratorID      string
	Spmo              User
	SpmoID            string
	CalibrationScore  float64
	CalibrationRating string
	Status            string
	SpmoStatus        string
}
