package response

type SummaryProject struct {
	Summary []*CalibratorBusinessUnit
}

type CalibratorBusinessUnit struct {
	CalibratorName           string
	CalibratorID             string
	CalibratorBusinessUnit   string
	CalibratorBusinessUnitID string
	APlus                    int
	A                        int
	BPlus                    int
	B                        int
	C                        int
	D                        int
	APlusGuidance            int
	AGuidance                int
	BPlusGuidance            int
	BGuidance                int
	CGuidance                int
	DGuidance                int
	Status                   string
}
