package model

type GroupBusinessUnit struct {
	BaseModel
	Status             bool
	GroupName          string
	BusinessUnits      []BusinessUnit `gorm:"foreignKey:GroupBusinessUnitId"`
	ScoreDistributions []ScoreDistribution
}
