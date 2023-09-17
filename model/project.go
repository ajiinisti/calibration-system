package model

type Project struct {
	BaseModel
	Name               string
	Year               int
	ActualScores       []ActualScore  `gorm:"foreignKey:ProjectID"`
	ProjectPhases      []ProjectPhase `gorm:"foreignKey:ProjectID"`
	ScoreDistributions []ScoreDistribution
	Active             bool `gorm:"default:false"`
}
