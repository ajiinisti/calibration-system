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

// DB Result Query
type SPMOSummaryResult struct {
	Count            int
	BusinessUnitName string
	BusinessUnitID   string
	Department       string
	CalibratorName   string
	CalibratorID     string
	ProjectPhaseID   string
	Order            int
}

// Response to User
type SummarySPMO struct {
	SummaryData []BUPerformanceSummarySPMO
}

type BUPerformanceSummarySPMO struct {
	BusinessUnitName string
	BusinessUnitID   string
	DepartmentData   []DepartmentCountSummarySPMO
}

type DepartmentCountSummarySPMO struct {
	DepartmentName   string
	ProjectPhaseData []*ProjectPhaseSummarySPMO
}

type ProjectPhaseSummarySPMO struct {
	CalibratorName string
	CalibratorID   string
	ProjectPhaseID string
	Order          int
	Count          int
	Status         string
}
