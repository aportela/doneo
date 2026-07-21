package handlers

type EmptyResponse struct {
}

type PagerRequest struct {
	Enabled     bool `json:"enabled"`
	CurrentPage int  `json:"currentPage"`
	ResultsPage int  `json:"resultsPage"`
}

type PagerResponse struct {
	Enabled      bool `json:"enabled"`
	CurrentPage  int  `json:"currentPage"`
	ResultsPage  int  `json:"resultsPage"`
	TotalPages   int  `json:"totalPages"`
	TotalResults int  `json:"totalResults"`
}

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "ASC"
	OrderDirectionDesc OrderDirection = "DESC"
)

func (o OrderDirection) IsValid() bool {
	return o == OrderDirectionAsc || o == OrderDirectionDesc
}

type OrderRequest struct {
	Field     string         `json:"field"`
	Direction OrderDirection `json:"direction"`
}
