package request

import "calibration-system.com/model"

type CalibrationRequest struct {
	RequestData map[string]model.Calibration
}

type CalibrationForm struct {
	CalibrationDataForms []CalibrationDataForm
}

type CalibrationDataForm struct {
	ProjectID      string
	ProjectPhaseID string
	EmployeeID     string
	CalibratorID   string
	SpmoID         string
	Spmo2ID        string
	Spmo3ID        string
	HrbpID         string
}
