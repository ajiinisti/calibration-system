package model

import (
	"time"

	"gorm.io/gorm"
)

type BusinessUnit struct {
	ID                  string         `gorm:"primaryKey;unique"`
	CreatedAt           time.Time      `gorm:"<-:create" json:"-"`
	UpdatedAt           time.Time      `json:"-"`
	DeletedAt           gorm.DeletedAt `json:"-"`
	Status              bool
	Name                string
	GroupBusinessUnit   GroupBusinessUnit
	GroupBusinessUnitId string
	RatingQuotas        []RatingQuota `gorm:"constraint:OnDelete:CASCADE"`
}
