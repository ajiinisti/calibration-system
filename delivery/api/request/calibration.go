package request

import "calibration-system.com/model"

type CalibrationRequest struct {
	RequestData map[string]*model.Calibration
}

type CalibrationForm struct {
	CalibrationDataForms []CalibrationDataForm
	ActualScore          float64
	ActualRating         string
	Y1Rating             string
	Y2Rating             string
	PTTScore             float64
	PATScore             float64
	Score360             float64
}

type CalibrationDataForm struct {
	ProjectID      string
	ProjectPhaseID string
	EmployeeID     string
	CalibratorID   string
	SpmoID         string
	Spmo2ID        string
	Spmo3ID        string
}
