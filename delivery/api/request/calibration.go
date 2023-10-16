package request

import "calibration-system.com/model"

type CalibrationRequest struct {
	RequestData map[string]model.Calibration
}
