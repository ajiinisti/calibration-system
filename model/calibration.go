package model

type Calibration struct {
	BaseModel
	Project           Project
	ProjectID         string
	ProjectPhase      ProjectPhase
	ProjectPhaseID    string
	Employee          User
	EmployeeID        string
	Calibrator        User
	CalibratorID      string
	Spmo              User
	SpmoID            string
	CalibrationScore  float64
	CalibrationRating string
}
