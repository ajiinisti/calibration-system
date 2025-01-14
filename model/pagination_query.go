package model

type PaginationQuery struct {
	Page                          int
	Take                          int
	Skip                          int
	Name                          string
	SupervisorName                []string
	Grade                         []string
	EmployeeName                  []string
	OrderGrade                    string
	OrderEmployeeName             string
	OrderCalibrationScore         string
	OrderCalibrationRating        string
	FilterCalibrationRating       string
	CalibrationPhaseBefore        int
	OrderCalibrationScoreBefore   string
	OrderCalibrationRatingBefore  string
	FilterCalibrationRatingBefore string
	RatingChangedStatus           string
	RatingChanged                 int
}
