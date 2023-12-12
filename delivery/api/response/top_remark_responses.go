package response

import (
	"time"

	"calibration-system.com/model"
)

type TopRemarkResponse struct {
	model.BaseModel
	Project        model.Project
	ProjectID      string
	Employee       model.User
	EmployeeID     string
	ProjectPhase   model.ProjectPhase
	ProjectPhaseID string
	Initiative     string
	Description    string
	Result         string
	StartDate      time.Time
	EndDate        time.Time
	Comment        string
	EvidenceLink   string
	EvidenceName   string
}
