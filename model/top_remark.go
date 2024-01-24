package model

import (
	"time"
)

type TopRemark struct {
	BaseModel
	Project        Project `gorm:"foreignKey:ProjectID"`
	ProjectID      string
	Employee       User `gorm:"foreignKey:EmployeeID"`
	EmployeeID     string
	ProjectPhase   ProjectPhase `gorm:"foreignKey:ProjectPhaseID"`
	ProjectPhaseID string
	Initiative     string
	Description    string
	Result         string
	StartDate      time.Time
	EndDate        time.Time
	Comment        string
	EvidenceName   string
	Evidence       []byte `json:"-"`
}
