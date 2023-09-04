package model

type Project struct {
	BaseModel
	Name               string
	Year               int
	Status             bool
	ActualScores       []ActualScore  `gorm:"foreignKey:ProjectID"`
	ProjectPhases      []ProjectPhase `gorm:"foreignKey:ProjectID"`
	ScoreDistributions []ScoreDistribution
}
