package response

import "time"

type NotificationModel struct {
	CalibratorID           string
	ProjectPhase           int
	Deadline               time.Time
	NextCalibrator         string
	PreviousCalibrator     string
	PreviousCalibratorID   string
	PreviousBusinessUnitID string
}
