package model

type GroupBusinessUnit struct {
	BaseModel
	Status             bool
	GroupName          string
	BusinessUnits      []BusinessUnit      `gorm:"foreignKey:GroupBusinessUnitId;constraint:OnDelete:CASCADE"`
	ScoreDistributions []ScoreDistribution `gorm:"constraint:OnDelete:CASCADE"`
}
