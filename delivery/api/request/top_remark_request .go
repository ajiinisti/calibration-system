package request

type DeleteTopRemarks struct {
	IDs []DeleteTopRemark
}

type DeleteTopRemark struct {
	ProjectID      string
	EmployeeID     string
	ProjectPhaseID string
}
