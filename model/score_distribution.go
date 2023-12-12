package model

import (
	"time"

	"gorm.io/gorm"
)

type ScoreDistribution struct {
	CreatedAt           time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `json:"-"`
	Project             Project
	ProjectID           string `gorm:"primaryKey"`
	GroupBusinessUnit   GroupBusinessUnit
	GroupBusinessUnitID string `gorm:"primaryKey"`
	APlusUpperLimit     float64
	APlusLowerLimit     float64
	AUpperLimit         float64
	ALowerLimit         float64
	BPlusUpperLimit     float64
	BPlusLowerLimit     float64
	BUpperLimit         float64
	BLowerLimit         float64
	CUpperLimit         float64
	CLowerLimit         float64
	DUpperLimit         float64
	DLowerLimit         float64
}
