package model

import (
	"time"

	"gorm.io/gorm"
)

type ActualScore struct {
	CreatedAt    time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `json:"-"`
	Project      Project
	ProjectID    string `gorm:"primaryKey"`
	Employee     User
	EmployeeID   string `gorm:"primaryKey"`
	ActualScore  float64
	ActualRating string
	Y1Rating     string
	Y2Rating     string
	PTTScore     float64
	PATScore     float64
	Score360     float64
}
