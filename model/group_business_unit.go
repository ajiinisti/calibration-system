package model

type GroupBusinessUnit struct {
	BaseModel
	Status             bool
	GroupName          string
	BusinessUnits      []BusinessUnit
	ScoreDistributions []ScoreDistribution
}
