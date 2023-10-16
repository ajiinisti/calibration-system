package model

type Project struct {
	BaseModel
	Name               string
	Year               int
	Calibrations       []Calibration  `gorm:"foreignKey:ProjectID;references:ID"`
	ActualScores       []ActualScore  `gorm:"foreignKey:ProjectID"`
	ProjectPhases      []ProjectPhase `gorm:"foreignKey:ProjectID"`
	ScoreDistributions []ScoreDistribution
	RatingQuotas       []RatingQuota   `gorm:"foreignKey:ProjectID;references:ID"`
	RemarkSettings     []RemarkSetting `gorm:"foreignKey:ProjectID;references:ID"`
	Active             bool            `gorm:"default:false"`
	APlusExcess        bool            `gorm:"default:false"`
	AExcess            bool            `gorm:"default:false"`
	BPlusExcess        bool            `gorm:"default:false"`
	BExcess            bool            `gorm:"default:false"`
	CExcess            bool            `gorm:"default:false"`
	DExcess            bool            `gorm:"default:false"`
}
