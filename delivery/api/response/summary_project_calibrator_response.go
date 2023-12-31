package response

type SummaryProject struct {
	Summary           []*CalibratorBusinessUnit
	APlusTotalScore   int
	ATotalScore       int
	BPlusTotalScore   int
	BTotalScore       int
	CTotalScore       int
	DTotalScore       int
	APlusGuidance     int
	AGuidance         int
	BPlusGuidance     int
	BGuidance         int
	CGuidance         int
	DGuidance         int
	AverageTotalScore float64
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
	TotalCalibratedScore     float64
	UserCount                int
	AverageScore             float64
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
	SummaryData []*BUPerformanceSummarySPMO
}

type BUPerformanceSummarySPMO struct {
	BusinessUnitName    string
	BusinessUnitID      string
	MaximumTotalData    int
	ProjectPhaseSummary []*ProjectPhaseSummarySPMO
}

type ProjectPhaseSummarySPMO struct {
	ProjectPhaseID     string
	Order              int
	DataCount          int
	CalibratorSummarys []*CalibratorSummary
}

type CalibratorSummary struct {
	CalibratorName string
	CalibratorID   string
	Count          int
	Status         string
}

// type DepartmentCountSummarySPMO struct {
// 	DepartmentName   string
// 	ProjectPhaseData []*ProjectPhaseSummarySPMO
// }
