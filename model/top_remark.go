package model

import (
	"time"
)

type TopRemark struct {
	BaseModel
	Project        Project `gorm:"foreignKey:ProjectID" json:"-"`
	ProjectID      string
	Employee       User `gorm:"foreignKey:EmployeeID" json:"-"`
	EmployeeID     string
	ProjectPhase   ProjectPhase `gorm:"foreignKey:ProjectPhaseID" json:"-"`
	ProjectPhaseID string
	Initiative     string
	Description    string
	Result         string
	StartDate      *time.Time
	EndDate        *time.Time
	Comment        string
	EvidenceName   string
	Evidence       []byte `json:"-"`
	IsProject      bool
	IsInitiative   bool
	EvidenceLink   string
}
