package request

type RejectJustification struct {
	ProjectID      string
	EmployeeID     string
	ProjectPhaseID string
	Comment        string
}

type AcceptJustification struct {
	ProjectID      string
	EmployeeID     string
	ProjectPhaseID string
	CalibratorID   string
}

type AcceptMultipleJustification struct {
	ArrayOfAcceptsJustification []AcceptJustification
}
