package request

type DeleteTopRemarks struct {
	IDs []string
}

type DeleteTopRemark struct {
	ProjectID      string
	EmployeeID     string
	ProjectPhaseID string
}
