package response

type Paging struct {
	Page        int
	RowsPerPage int
	TotalRows   int
	TotalPages  int
}
