package domain

type ProjectStatus struct {
	ID       string
	Name     string
	Index    int
	HexColor string
}

type SearchProjectStatusesFilter struct {
	Name *string
}
