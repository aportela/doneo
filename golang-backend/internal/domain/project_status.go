package domain

type ProjectStatus struct {
	ID       string
	Name     string
	HexColor string
}

type SearchProjectStatusesFilter struct {
	Name *string
}
