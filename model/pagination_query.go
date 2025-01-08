package model

type PaginationQuery struct {
	Page                    int
	Take                    int
	Skip                    int
	Name                    string
	SupervisorName          []string
	Grade                   []string
	EmployeeName            []string
	OrderGrade              string
	OrderEmployeeName       string
	OrderCalibrationScore   string
	OrderCalibrationRating  string
	FilterCalibrationRating string
	RatingChanged           int
}
