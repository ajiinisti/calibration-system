package model

type ActualScore struct {
	BaseModel
	Project      Project
	ProjectID    string
	Employee     User
	EmployeeID   string
	ActualScore  float64
	ActualRating string
	Y1Rating     string
	Y2Rating     string
}
