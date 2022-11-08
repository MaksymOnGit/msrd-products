package models

type QueryRequest struct {
	Rows      int64  `json:"rows" validate:"required,min=5,max=30"`
	Offset    int64  `json:"offset,omitempty" validate:"min=0"`
	SortField string `json:"sortField"`
	SortOrder int    `json:"sortOrder,omitempty" validate:"oneof=-1 0 1"`
}

type QueryResponse[T any] struct {
	Result              []T   `json:"result"`
	Page                int64 `json:"page"`
	TotalPagesCount     int64 `json:"totalPagesCount"`
	TotalRecordsCount   int64 `json:"totalRecordsCount"`
	RecordsPerPageCount int64 `json:"recordsPerPageCount"`
	IsNext              bool  `json:"isNext"`
	IsPrev              bool  `json:"isPrev"`
}
