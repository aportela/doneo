package userhandler

type addRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsSuperUser bool   `json:"isSuperUser"`
}

type updateRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	IsSuperUser bool    `json:"isSuperUser"`
}

type patchRequest struct {
	DeletedAt *int64 `json:"deletedAt"`
}

type PagerRequest struct {
	Enabled     bool `json:"enabled"`
	CurrentPage int  `json:"currentPage"`
	ResultsPage int  `json:"resultsPage"`
}

type SortOrder string

const (
	SortAsc  SortOrder = "ASC"
	SortDesc SortOrder = "DESC"
)

func (o SortOrder) IsValid() bool {
	return o == SortAsc || o == SortDesc
}

type OrderRequest struct {
	Field string    `json:"field"`
	Sort  SortOrder `json:"sort"`
}

type FilterRequest struct {
}

type searchRequest struct {
	Pager  PagerRequest  `json:"pager"`
	Order  OrderRequest  `json:"order"`
	Filter FilterRequest `json:"filter"`
}

type userResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt"`
	DeletedAt   *int64 `json:"deletedAt"`
	IsSuperUser bool   `json:"isSuperUser"`
	AvatarURL   string `json:"avatarUrl"`
}

type addResponse struct {
	User userResponse `json:"user"`
}

type updateResponse struct {
	User userResponse `json:"user"`
}

type getResponse struct {
	User userResponse `json:"user"`
}

type PagerResponse struct {
	Enabled      bool `json:"enabled"`
	CurrentPage  int  `json:"currentPage"`
	ResultsPage  int  `json:"resultsPage"`
	TotalPages   int  `json:"totalPages"`
	TotalResults int  `json:"totalResults"`
}

type searchResponse struct {
	Users []userResponse `json:"users"`
	Pager PagerResponse  `json:"pager"`
}
