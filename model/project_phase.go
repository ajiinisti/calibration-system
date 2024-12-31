package model

import "time"

type ProjectPhase struct {
	BaseModel
	Phase      Phase
	PhaseID    string
	Project    Project `json:"-"`
	ProjectID  string
	ReviewSpmo bool
	StartDate  time.Time
	EndDate    time.Time
	Guideline  bool
	ShowChart  bool
}
