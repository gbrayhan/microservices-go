package domain

import "time"

type DateRangeFilter struct {
	Field string     `json:"field"`
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

func (sd SortDirection) IsValid() bool {
	return sd == SortAsc || sd == SortDesc
}

type DataFilters struct {
	LikeFilters      map[string][]string `json:"likeFilters"`
	Matches          map[string][]string `json:"matches"`
	DateRangeFilters []DateRangeFilter   `json:"dateRanges"`
	SortBy           []string            `json:"sortBy"`
	SortDirection    SortDirection       `json:"sortDirection"`
	Page             int                 `json:"page"`
	PageSize         int                 `json:"pageSize"`
}
