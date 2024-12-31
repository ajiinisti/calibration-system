package model

import (
	"time"

	"gorm.io/gorm"
)

type BottomRemark struct {
	CreatedAt      time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Project        Project        `gorm:"foreignKey:ProjectID" json:"-"`
	ProjectID      string         `gorm:"primaryKey"`
	Employee       User           `gorm:"foreignKey:EmployeeID" json:"-"`
	EmployeeID     string         `gorm:"primaryKey"`
	ProjectPhase   ProjectPhase   `gorm:"foreignKey:ProjectPhaseID" json:"-"`
	ProjectPhaseID string         `gorm:"primaryKey"`
	LowPerformance string
	Indisipliner   string
	Attitude       string
	WarningLetter  string
}
