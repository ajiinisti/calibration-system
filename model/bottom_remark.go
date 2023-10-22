package model

import (
	"time"

	"gorm.io/gorm"
)

type BottomRemark struct {
	CreatedAt      time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Project        Project        `gorm:"foreignKey:ProjectID"`
	ProjectID      string         `gorm:"primaryKey"`
	Employee       User           `gorm:"foreignKey:EmployeeID"`
	EmployeeID     string         `gorm:"primaryKey"`
	ProjectPhase   ProjectPhase   `gorm:"foreignKey:ProjectPhaseID"`
	ProjectPhaseID string         `gorm:"primaryKey"`
	LowPerformance string
	Indisipliner   string
	Attitude       string
	WarningLetter  string
}
