package response

import (
	"time"

	"calibration-system.com/model"
)

type TopRemarkResponse struct {
	model.BaseModel
	ProjectID      string
	EmployeeID     string
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
