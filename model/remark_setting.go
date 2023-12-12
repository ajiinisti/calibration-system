package model

type RemarkSetting struct {
	BaseModel
	Project           Project
	ProjectID         string
	JustificationType string //	level/rating
	ScoringType       string //	bottom/top remark
	Level             int
	From              string
	To                string
}

// Define
// A+ = 6
// A = 5
// B+ = 4
// B = 3
// C = 2
// D = 1
