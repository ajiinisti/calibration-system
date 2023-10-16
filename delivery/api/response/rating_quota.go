package response

type RatingQuota struct {
	APlus         int
	A             int
	BPlus         int
	B             int
	C             int
	D             int
	Total         int
	ScoringMethod string
}

type RatingQuotaResponse struct {
	RatingQuotaGroups map[string]RatingQuota
}
