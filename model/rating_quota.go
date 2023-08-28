package model

type RatingQuota struct {
	BaseModel
	Project        Project
	ProjectID      string
	BusinessUnit   BusinessUnit
	BusinessUnitID string
	APlusQuota     float64
	AQuota         float64
	BPlusQuota     float64
	BQuota         float64
	CQuota         float64
	DQuota         float64
}
