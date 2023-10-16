package model

import (
	"time"

	"gorm.io/gorm"
)

type RatingQuota struct {
	CreatedAt      time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
	Project        Project
	ProjectID      string `gorm:"primaryKey"`
	BusinessUnit   BusinessUnit
	BusinessUnitID string `gorm:"primaryKey"`
	APlusQuota     float64
	AQuota         float64
	BPlusQuota     float64
	BQuota         float64
	CQuota         float64
	DQuota         float64
	Remaining      string
	Excess         string
	ScoringMethod  string
}
