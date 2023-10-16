package response

type SummaryProject struct {
	Summary []*CalibratorBusinessUnit
}

type CalibratorBusinessUnit struct {
	CalibratorName         string
	CalibratorBusinessUnit string
	APlus                  int
	A                      int
	BPlus                  int
	B                      int
	C                      int
	D                      int
	APlusGuidance          int
	AGuidance              int
	BPlusGuidance          int
	BGuidance              int
	CGuidance              int
	DGuidance              int
	Status                 string
}
