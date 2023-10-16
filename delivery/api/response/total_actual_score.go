package response

type TotalActualScore struct {
	APlus int
	A     int
	BPlus int
	B     int
	C     int
	D     int
	Total int
}

type TotalActualScoreResponse struct {
	TotalActualScoreGroups map[string]TotalActualScore
}
