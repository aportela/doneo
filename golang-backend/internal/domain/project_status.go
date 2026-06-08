package domain

type ProjectStatus struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
}

type SearchProjectStatusesFilter struct {
	Name *string
}
