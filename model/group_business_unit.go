package model

type GroupBusinessUnit struct {
	BaseModel
	Status          bool
	GroupName       string
	APlusUpperLimit float64
	APlusLowerLimit float64
	AUpperLimit     float64
	ALowerLimit     float64
	BPlusUpperLimit float64
	BPlusLowerLimit float64
	BUpperLimit     float64
	BLowerLimit     float64
	CUpperLimit     float64
	CLowerLimit     float64
	DUpperLimit     float64
	DLowerLimit     float64
	BusinessUnits   []BusinessUnit
}
